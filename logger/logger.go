package logger

import (
	"os"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func InitLogger() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	SetLogLevel()
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func SetLogLevel() {
	switch config.Config.Logging.Level {
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	case "ERROR":
		log.SetLevel(logrus.ErrorLevel)
	case "WARN":
		log.SetLevel(logrus.WarnLevel)
	}
}
