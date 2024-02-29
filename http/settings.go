package main

import (
	"github.com/spf13/viper"
)

type SettingsHttp struct {
	AppName string `mapstructure:"app_name"`
	Addr string
	Env string
	RedisDsn string `mapstructure:"redis_dsn"`
	RedisExp int64 `mapstructure:"redis_exp"`
	PostgresDsn string `mapstructure:"postgres_dsn"`
	SentryDsn string `mapstructure:"sentry_dsn"`
	// Any other settings below
}

func NewSettingsHttp() (settings SettingsHttp, err error) {
	// Set the file name of the configurations file
	viper.SetConfigName("local.http.env")
	viper.SetConfigType("env")
	viper.SetDefault("addr", "0.0.0.0:8080")
	viper.SetDefault("app_name", "Test App !Change Name!")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&settings)
	return
}

func (s SettingsHttp) GetRedisDsn() string {
	return s.RedisDsn
}
