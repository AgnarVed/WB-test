package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port            int    `mapstructure:"SERVER_PORT" required:"true"`
	DBConnStr       string `mapstructure:"DB_CONN_STR" required:"true"`
	DriverName      string `mapstructure:"DB_DRIVER_NAME" required:"true"`
	ShutdownTimeout int    `mapstructure:"SHUTDOWN_TIMEOUT" required:"true"`
}

//func getTime(input string) (time.Duration, error) {
//	return time.ParseDuration(input)
//}

func NewConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, err
}
