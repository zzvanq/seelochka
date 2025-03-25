package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env       string     `mapstructure:"env"`
	DbPath    string     `mapstructure:"db_path"`
	SentryDSN string     `mapstructure:"sentry_dsn"`
	Http      HttpServer `mapstructure:"http_server,squash"`
}

type HttpServer struct {
	BindHost        string `mapstructure:"bind_host"`
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Address         string
	Timeout         time.Duration `mapstructure:"timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

func (hs *HttpServer) SetAddress() {
	hs.Address = fmt.Sprintf("%s:%s", hs.Host, hs.Port)
}

func MustLoad() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	viper.AutomaticEnv()

	v := viper.New()
	configPath := viper.GetString("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("no CONFIG_PATH is set")
	}
	v.SetConfigFile(configPath)
	err = v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	viper.MergeConfigMap(v.AllSettings())

	var c Config
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	c.Http.SetAddress()

	return &c
}
