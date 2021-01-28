package config

import (
	"os"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string
	LogLevel    string
	Release     string
	ApiListener string
	RunStatus   string
	WorkerCount int
}

func (c Config) Validate() error {
	return v.ValidateStruct(&c,
		v.Field(&c.LogLevel, v.Min(0), v.Max(5)),
		v.Field(&c.RunStatus, v.Required),
	)
}

type SentryConfig struct {
	Dsn string
}

func InitConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Error(err)
	}
	viper.AutomaticEnv()
	c := new(Config)
	c.RunStatus = "INIT"
	c.ServiceName = "web_page_analyzer"
	c.LogLevel = viper.GetString("LOG_LEVEL")
	c.Release = viper.GetString("RELEASE")
	if c.Release == "" {
		c.Release = viper.GetString("VCS_REF")
	}
	c.ApiListener = viper.GetString("API_LISTENER")
	c.WorkerCount = viper.GetInt("WORKER_COUNT")
	if err := c.Validate(); err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
	return c
}
