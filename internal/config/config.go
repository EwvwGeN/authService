package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string `mapstructure:"log_level"`
	Port int `mapstructure:"port"`
	Validator Validator `mapstructure:"validator"`
	MongoConfig MongoConfig `mapstructure:"mongo"`
	TokenTTL time.Duration `mapstructure:"token_ttl"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AutomaticEnv()
	if path != "" {
		viper.AddConfigPath(path)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
        return nil, err
    }
	config.Validator.mustBeRegex()
	config.TokenTTL *= time.Second
	return &config, nil
}