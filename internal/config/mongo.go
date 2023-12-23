package config

type MongoConfig struct {
	ConectionFormat string `mapstructure:"db_con_format"`
	Host            string `mapstructure:"db_host"`
	Port            string `mapstructure:"db_port"`
	User            string `mapstructure:"db_user"`
	Password        string `mapstructure:"db_pass"`
	AuthSourse      string `mapstructure:"db_auth_source"`
	Database        string `mapstructure:"db_name"`
	UserCollection  string `mapstructure:"db_col_user"`
	AppCollection   string `mapstructure:"db_col_app"`
}