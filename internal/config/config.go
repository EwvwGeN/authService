package config

import (
	"fmt"
	p "path"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var serviceTag string = "auth_service"

type Config struct {
	LogLevel string `mapstructure:"log_level"`
	Port int `mapstructure:"port"`
	Validator Validator `mapstructure:"validator"`
	MongoConfig MongoConfig `mapstructure:"mongo"`
	TokenTTL time.Duration `mapstructure:"token_ttl"`
}

func LoadConfig(path string) (*Config, error) {
	type ServiceConfig struct {
		Cfg Config `mapstructure:"auth_service"`
	}
	v := viper.NewWithOptions()
	v.AutomaticEnv()
	v.AliasesFirstly(false)
	v.AliasesStepByStep(true)
	if path != "" {
		fileParts := strings.Split(p.Base(path), ".")
		if len(fileParts) < 2 {
			return nil, fmt.Errorf("incorrect config file: %s", path)
		}
		v.SetConfigFile(path)
		v.SetConfigType(fileParts[len(fileParts)-1])
		err := v.ReadInConfig()
		if err != nil {
			return nil, err
		}
	} else {
		v.SetConfigFile(".env")
		err := v.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}
	outOfServiceKeys := deleteKeysWithPrefix(v.AllKeys(), serviceTag)
	for _, value := range outOfServiceKeys {
		v.RegisterAlias(fmt.Sprintf("%s.%s", serviceTag, value), value)
	}
	var config ServiceConfig
	err := v.Unmarshal(&config)
	if err != nil {
        return nil, err
    }
	config.Cfg.Validator.mustBeRegex()
	return &config.Cfg, nil
}
func deleteKeysWithPrefix(keys []string, prefix string) []string {
	prefix = strings.ToLower(prefix)
	var out []string
	for _, v := range keys {
		if len(v) <= len(prefix) + 1 {
			out = append(out, v)
			continue
		}
		if v[:len(prefix)] != prefix{
			out = append(out, v)
		}
	}
	return out
}