package log

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Level uint8

const (
	Panic Level = iota
	Error
	Warn
	Info
	Debug
)

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return Panic, nil
	case "error", "err":
		return Error, nil
	case "warning", "warn":
		return Warn, nil
	case "info":
		return Info, nil
	case "debug":
		return Debug, nil
	}
	return Level(0), fmt.Errorf("not a valid log level: %q", lvl)
}

func (l Level) String() string {
	switch l {
	case Panic:
		return "panic"
	case Error:
		return "error"
	case Warn:
		return "warning"
	case Info:
		return "info"
	case Debug:
	}
	return "unknown"
}

type Fields = map[string]interface{}

type ErrorWithFields interface {
	error
	Fields() Fields
}

type Logger interface {
	Level() Level

	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithError(err error) Logger

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
}

func Format(s ...interface{}) []string {
	r := make([]string, len(s))
	for i, s := range s {
		switch ts := s.(type) {
		case error:
			r[i] = ts.Error()
		default:
			rtype := reflect.TypeOf(s)
			switch rtype.Kind() {
			case reflect.Struct, reflect.Interface, reflect.Map, reflect.Slice, reflect.Array:
				j, err := json.Marshal(s)
				if err != nil {
					r[i] = err.Error()
				} else {
					r[i] = string(j)
				}
			case reflect.Ptr:
				rvalue := reflect.ValueOf(s)
				r[i] = Format(rvalue.Elem().Interface())[0]
			default:
				r[i] = fmt.Sprint(s)
			}
		}
	}
	return r
}
