package config

type GRPCConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	RgistrationLink string `mapstructure:"reg_link"`
}