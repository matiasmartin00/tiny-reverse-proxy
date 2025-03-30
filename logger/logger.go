package logger

import (
	"os"
	"sync"

	"github.com/matiasmartin00/tiny-reverse-proxy/config"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	setLogLevel()
}

type logger struct {
	config config.Config
	log    *logrus.Logger
}

var instance Logger
var once sync.Once

func getLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	return log
}

func (l *logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *logger) setLogLevel() {
	switch l.config.GetLoggingLevel() {
	case "INFO":
		l.log.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		l.log.SetLevel(logrus.DebugLevel)
	case "ERROR":
		l.log.SetLevel(logrus.ErrorLevel)
	case "WARN":
		l.log.SetLevel(logrus.WarnLevel)
	}
}

func NewLogger(cfg config.Config) Logger {
	once.Do(func() {
		log := getLogger()
		instance = &logger{
			config: cfg,
			log:    log,
		}
		instance.setLogLevel()
	})
	return instance
}

func GetLogger() Logger {
	return instance
}
