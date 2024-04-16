package logger

import "github.com/sirupsen/logrus"

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

type LogrusLogger struct {
	log *logrus.Logger
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}
func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}
func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}
func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}
func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}
func (l *LogrusLogger) Debug(msg string, fields map[string]interface{}) {
	l.log.WithFields(fields).Debug(msg)
}
func (l *LogrusLogger) Info(msg string, fields map[string]interface{}) {
	l.log.WithFields(fields).Info(msg)
}
func (l *LogrusLogger) Warn(msg string, fields map[string]interface{}) {
	l.log.WithFields(fields).Warn(msg)
}
func (l *LogrusLogger) Error(msg string, fields map[string]interface{}) {
	l.log.WithFields(fields).Error(msg)
}
func (l *LogrusLogger) Fatal(msg string, fields map[string]interface{}) {
	l.log.WithFields(fields).Fatal(msg)
}

func NewLogrusLogger(logLevel string, isDev bool) *LogrusLogger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
		ForceQuote:    true,
	})

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Warn("Invalid log level: " + logLevel)
		log.Warn("Using default log level: INFO")
		log.SetLevel(logrus.InfoLevel)
	}

	log.SetLevel(level)
	if !isDev {
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	return &LogrusLogger{log: log}
}
