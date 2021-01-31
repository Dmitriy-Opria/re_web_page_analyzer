package config

import (
	"os"

	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/log"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string
	LogLevel    log.Level
	Release     string
	ApiListener string
	RunStatus   string
	WorkerCount int
}

func (c Config) Validate() error {
	return v.ValidateStruct(&c,
		v.Field(&c.LogLevel, v.Required),
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
	var logLevel log.Level
	switch viper.GetInt("LOG_LEVEL") {
	case 0:
		logLevel = log.LevelPanic
	case 1:
		logLevel = log.LevelFatal
	case 2:
		logLevel = log.LevelError
	case 3:
		logLevel = log.LevelWarn
	case 4:
		logLevel = log.LevelInfo
	case 5:
		fallthrough
	default:
		logLevel = log.LevelDebug
	}

	c.LogLevel = logLevel
	c.Release = viper.GetString("RELEASE")
	if c.Release == "" {
		c.Release = viper.GetString("VCS_REF")
	}
	c.ApiListener = viper.GetString("API_LISTENER")
	c.WorkerCount = viper.GetInt("API_WORKER_COUNT")
	if err := c.Validate(); err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
	return c
}
