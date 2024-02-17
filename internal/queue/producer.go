package queue

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/EwvwGeN/authService/internal/config"
	"github.com/EwvwGeN/authService/internal/domain/models"
	"github.com/goombaio/namegenerator"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	log *slog.Logger
	cfg config.RabbitMQConfig
	msgCh chan *messageWithErrCh
	closer func() error 
}

func StartProducer(ctx context.Context, lg *slog.Logger, cfg config.RabbitMQConfig) (p *Producer){
	var (
		err error
		closer = make(chan struct{})
		errCh = make(chan error)
		outErrCh = make(chan error)
		msgCh = make(chan *messageWithErrCh)
		closing bool = false
	)

	lg = lg.With(slog.String("op", "rabbitMQ"))

	p = &Producer{
		log: lg,
		cfg: cfg,
		msgCh: msgCh,
	}
	
	seed := time.Now().UTC().UnixNano()
    ng := namegenerator.NewNameGenerator(seed)

	pubWG := new(sync.WaitGroup)

	for i:= 0; i < p.cfg.ProducerConfig.Count; i++ {
		pubWG.Add(1)
		go p.startPublisher(ctx, pubWG, ng.Generate(), errCh, &closing)
	}
	
	go func() {
		for {
			select {
				case <- closer: {
					closing = true
					pubWG.Wait()
					outErrCh <- err
					return
				}
				case err = <-errCh: {
					if err != nil {
						p.log.Error("error occurred", slog.String("error", err.Error()))
					}
					if p.cfg.ProducerConfig.RestartOnErr {
						p.log.Info("publisher restart attempt")
						pubWG.Add(1)
						go p.startPublisher(ctx, pubWG, ng.Generate(), errCh, &closing)
						continue
					}
					if p.cfg.ProducerConfig.CancelOnError {
						p.log.Info("closing all publishers")
						closing = true
						pubWG.Wait()
						outErrCh <- fmt.Errorf("pulbishers closing after error: %w", err)
						return
					}
				}
			}
		}
	}()

	p.closer = func() error {
		closer <- struct{}{}
		return <- outErrCh
	}
	return p
}

func (p *Producer) SendMsg(ctx context.Context, msg *models.Message) error {
	if p.msgCh == nil {
		return ErrMsgChanNil
	}

	wrappedMsg := &messageWithErrCh{
		msg: msg,
		errCh: make(chan error),
	}

	p.msgCh <- wrappedMsg
	for {
		select {
		case err := <- wrappedMsg.errCh:
			return err
		case <-time.After(p.cfg.ProducerConfig.SendAwaitTime):
			return fmt.Errorf("exceeded waiting time")
		}
	}
}

func preparePub(ctx context.Context, cfg config.RabbitMQConfig, lg *slog.Logger) (*amqp.Connection, *amqp.Channel, error){
	config := amqp.Config{
		Vhost:      cfg.VirtualHost,
		Properties: amqp.NewConnectionProperties(),
	}
	config.Properties.SetClientConnectionName(cfg.ConnectionName)

	amqpURI := amqp.URI{
		Scheme: cfg.Scheme,
		Host: cfg.Host,
		Port: cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		Vhost: cfg.VirtualHost,
	}

	lg.Info("rabbitMQ dialing", slog.String("URI", amqpURI.String()))
	conn, err := amqp.DialConfig(amqpURI.String(), config)
	if err != nil {
		lg.Error(ErrStartupConnection.Error(), slog.String("error", err.Error()))
		return nil, nil, fmt.Errorf("%s: %w", ErrStartupConnection.Error(), err)
	}

	lg.Info("got Connection")
	lg.Info("getting Channel")
	channel, err := conn.Channel()
	if err != nil {
		lg.Error(ErrOpenChannel.Error(), slog.String("error", err.Error()))
		return nil, nil, fmt.Errorf("%s: %w", ErrOpenChannel.Error(), err)
	}

	lg.Info("got Channel")
	lg.Info("declaring Exchange", slog.String("exchange", cfg.ExchangerConfig.Name))
	if err := channel.ExchangeDeclare(
		cfg.ExchangerConfig.Name,
		cfg.ExchangerConfig.Type,
		true,
		false,
		false,
		true,
		nil,
	); err != nil {
		lg.Error(ErrExchangeDeclare.Error(), slog.String("error", err.Error()))
		return nil, nil, fmt.Errorf("%s: %w", ErrExchangeDeclare.Error(), err)
	}

	lg.Info("declared Exchange", slog.String("exchange", cfg.ExchangerConfig.Name))
	lg.Info("declaring Queue", slog.String("queue", cfg.QueueConfig.Name))
	queue, err := channel.QueueDeclare(
		cfg.QueueConfig.Name,
		true,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		lg.Error(ErrDeclareQueue.Error(), slog.String("error", err.Error()))
		return nil, nil, fmt.Errorf("%s: %w", ErrDeclareQueue.Error(), err)
	}

	lg.Info(
		"declared Queue",
		slog.String("queue", cfg.QueueConfig.Name),
		slog.Int("messages", queue.Messages),
		slog.Int("consumers", queue.Consumers),
	)


	lg.Info("binding to Exchange", slog.String("key", cfg.BindingConfig.Key))
	if err = channel.QueueBind(
		queue.Name,
		cfg.BindingConfig.Key,
		cfg.ExchangerConfig.Name,
		true,
		nil,
	); err != nil {
		lg.Error(ErrExchangeBind.Error(), slog.String("error", err.Error()))
		return nil, nil, fmt.Errorf("%s: %w", ErrExchangeBind.Error(), err)
	}
	return conn, channel, err
}

func (p *Producer) startPublisher(ctx context.Context, pubWG *sync.WaitGroup, name string, errCh chan error, closing *bool) {
	defer pubWG.Done()
	var (
		channel *amqp.Channel
		conn *amqp.Connection
		err error
		open bool = false
	)
	lg := p.log.With(slog.String("name", name))
	conn, channel, err = preparePub(ctx, p.cfg, lg)
	if err != nil {
		lg.Error("cant prepare publisher")
		if conn != nil {
			conn.Close()
		}
		errCh <- err
		return
	}
	open = true
	for !*closing {
		select {
			case msg, ok := <- p.msgCh: {
				if !ok {
					lg.Error(ErrReceiveMessage.Error())
					errCh <- ErrReceiveMessage
					break
				}
				if open {
					sendMsg(ctx, channel, p.cfg.ExchangerConfig.Name, p.cfg.BindingConfig.Key, msg)
				}
				if !open {
					if conn, channel, err = preparePub(ctx, p.cfg, lg); err != nil {
						lg.Error("cant reopen publisher")
						if conn != nil {
							conn.Close()
						}
						errCh <- err
						return
					}
					open = true
					sendMsg(ctx, channel, p.cfg.ExchangerConfig.Name, p.cfg.BindingConfig.Key, msg)
				}
			}
			case <-time.After(p.cfg.ProducerConfig.TimeToSleep): {
				if !open {
					if err = conn.Close(); err != nil {
						lg.Error("cant close connection", slog.String("error", err.Error()))
						errCh <- err
						break
					}
					open = false
				}
			}
		}
	}
}

func sendMsg(ctx context.Context, channel *amqp.Channel, exchanger, bindKey string, wMsg *messageWithErrCh) {
	err := channel.PublishWithContext(
		ctx,
		exchanger,
		bindKey,
		true,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{
				"To": wMsg.msg.EmailTo,
				"Subject": wMsg.msg.Subject,
			},
			ContentType:     "text/plain",
			ContentEncoding: "",
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
			Body:            wMsg.msg.Body,
		})
	wMsg.errCh <- err
}

func (p *Producer) Close() error {
	p.log.Info("closing rabbitmq producer")
	return p.closer()
}