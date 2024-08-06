package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDRIVER            string        `mapstructure:"DB_DRIVER"`
	DBSOURCE            string        `mapstructure:"DB_SOURCE"`
	PORT                string        `mapstructure:"PORT"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		return
	}

	return
}
