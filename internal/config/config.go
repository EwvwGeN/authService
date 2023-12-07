package config

import "github.com/spf13/viper"

type Config struct {
	Test string `mapstructure:"Test"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AutomaticEnv()
	if path != "" {
		viper.AddConfigPath(path)
		viper.SetConfigName("")
		viper.SetConfigType("env")
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
	return &config, nil
}