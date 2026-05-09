package utils

import (
	"time"

	"github.com/spf13/viper"
)


const (
	KeyAccessToken = "access_token"
	KeyRefreshToken = "refresh_token"
)



type Config struct {
	DB_URI               string        `mapstructure:"DB_URI"`
	Port                 int           `mapstructure:"PORT"`
	AppEnv               string        `mapstructure:"APP_ENV"`
	AppDomain            string        `mapstructure:"APP_DOMAIN"`
	SecretKey            string        `mapstructure:"SECRET_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
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
	return
}
