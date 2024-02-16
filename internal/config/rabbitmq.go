package config

import "time"

type RabbitMQConfig struct {
	Scheme          string          `mapstructure:"scheme"`
	Host            string          `mapstructure:"host"`
	Port            int             `mapstructure:"port"`
	Username        string          `mapstructure:"username"`
	Password        string          `mapstructure:"password"`
	VirtualHost     string          `mapstructure:"virtual_host"`
	ConnectionName  string          `mapstructure:"connection_name"`
	ExchangerConfig ExchangerConfig `mapstructure:"exchanger"`
	BindingConfig   BindingConfig   `mapstructure:"binding"`
	QueueConfig     QueueConfig     `mapstructure:"queue"`
	ProducerConfig  ProducerConfig  `mapstructure:"producer"`
}

type ExchangerConfig struct {
	Name string `mapstructure:"name"`
	Type string `mapstructure:"type"`
}

type BindingConfig struct {
	Key string `mapstructure:"key"`
}

type QueueConfig struct {
	Name string `mapstructure:"name"`
}

type ProducerConfig struct {
	Count         int  `mapstructure:"count"`
	RestartOnErr  bool `mapstructure:"rstrt_on_error"`
	CancelOnError bool `mapstructure:"cancel_on_error"`
	TimeToSleep   time.Duration `mapstructure:"tts"`
	SendAwaitTime time.Duration `mapstructure:"send_await_time"`
}