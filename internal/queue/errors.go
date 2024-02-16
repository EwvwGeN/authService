package queue

import "errors"

var (
	ErrInternal = errors.New("internal error")

	ErrStartupConnection = errors.New("failed to connect to RabbitMQ")
	ErrOpenChannel       = errors.New("failed to open a channel")
	ErrDeclareQueue      = errors.New("failed to declare a queue")
	ErrExchangeDeclare   = errors.New("failed to declare an exchange")
	ErrExchangeBind      = errors.New("failed to bind to exchange")
	ErrMsgChanNil 	 	 = errors.New("message channel doesnt exist")
	ErrReceiveMessage 	 = errors.New("cant receive message from channel")
	ErrClose     		 = errors.New("error while closing")
)