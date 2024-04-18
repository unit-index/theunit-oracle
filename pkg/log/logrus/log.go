package logrus

import (
	"github.com/sirupsen/logrus"

	"github.com/toknowwhy/theunit-oracle/pkg/log"
)

type Logger struct {
	log logrus.FieldLogger
	lvl log.Level
}

func New(logger logrus.FieldLogger) log.Logger {
	lvl := log.Debug
	if l, ok := logger.(*logrus.Logger); ok {
		switch l.Level {
		case logrus.PanicLevel:
			lvl = log.Panic
		case logrus.FatalLevel:
			lvl = log.Panic
		case logrus.ErrorLevel:
			lvl = log.Error
		case logrus.WarnLevel:
			lvl = log.Warn
		case logrus.InfoLevel:
			lvl = log.Info
		case logrus.DebugLevel:
			lvl = log.Debug
		case logrus.TraceLevel:
			lvl = log.Debug
		}
	}
	return &Logger{log: logger, lvl: lvl}
}

func (l *Logger) Level() log.Level {
	return l.lvl
}

func (l *Logger) WithField(key string, value interface{}) log.Logger {
	return &Logger{log: l.log.WithField(key, value)}
}

func (l *Logger) WithFields(fields log.Fields) log.Logger {
	return &Logger{log: l.log.WithFields(fields)}
}

func (l *Logger) WithError(err error) log.Logger {
	if fErr, ok := err.(log.ErrorWithFields); ok {
		return &Logger{log: l.log.WithFields(fErr.Fields()).WithError(err)}
	}

	return &Logger{log: l.log.WithError(err)}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.log.Panic(args...)
}
