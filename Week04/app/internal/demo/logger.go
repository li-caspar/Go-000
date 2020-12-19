package demo

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Logger struct {
	*logrus.Logger
	File *os.File
}

func NewLogger() *Logger {
	logger := &Logger{
		logrus.StandardLogger(),
		nil,
	}
	logrus.SetLevel(logrus.Level(viper.GetInt("log.level")))
	switch viper.GetString("log.format") {
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		logrus.SetFormatter(new(logrus.TextFormatter))
	}
	if viper.GetString("log.output") != "" {
		switch viper.GetString("log.output") {
		case "stdout":
			logrus.SetOutput(os.Stdout)
		case "stderr":
			logrus.SetOutput(os.Stderr)
		case "file":
			if name := viper.GetString("log.output_file"); name != "" {
				_ = os.MkdirAll(filepath.Dir(name), 0777)
				f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					panic(err)
				}
				logrus.SetOutput(f)
				logger.File = f
			}
		}
	}
	return logger
}

func (l *Logger) Stop() {
	if l.File != nil {
		_ = l.File.Close()
	}
}
