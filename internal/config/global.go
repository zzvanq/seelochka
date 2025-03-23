package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env       string     `mapstructure:"env"`
	DbPath    string     `mapstructure:"db_path"`
	SentryDSN string     `mapstructure:"sentry_dsn"`
	Http      HttpServer `mapstructure:"http_server"`
}

type HttpServer struct {
	Address         string        `mapstructure:"address"`
	Timeout         time.Duration `mapstructure:"timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

func MustLoad() *Config {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	v := viper.New()
	v.AddConfigPath("./configs")
	v.SetConfigName("conf")
	v.SetConfigType("yaml")
	err = v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	viper.MergeConfigMap(v.AllSettings())

	var c Config

	err = viper.UnmarshalExact(&c)
	if err != nil {
		log.Fatal(err)
	}

	return &c
}
