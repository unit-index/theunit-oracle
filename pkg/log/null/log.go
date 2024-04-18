package null

import (
	"fmt"

	"github.com/toknowwhy/theunit-oracle/pkg/log"
)

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (n *Logger) Level() log.Level                             { return log.Panic }
func (n *Logger) WithField(_ string, _ interface{}) log.Logger { return n }
func (n *Logger) WithFields(_ log.Fields) log.Logger           { return n }
func (n *Logger) WithError(_ error) log.Logger                 { return n }
func (n *Logger) Debugf(_ string, _ ...interface{})            {}
func (n *Logger) Infof(_ string, _ ...interface{})             {}
func (n *Logger) Warnf(_ string, _ ...interface{})             {}
func (n *Logger) Errorf(_ string, _ ...interface{})            {}
func (n *Logger) Panicf(format string, args ...interface{})    { panic(fmt.Sprintf(format, args...)) }
func (n *Logger) Debug(_ ...interface{})                       {}
func (n *Logger) Info(_ ...interface{})                        {}
func (n *Logger) Warn(_ ...interface{})                        {}
func (n *Logger) Error(_ ...interface{})                       {}
func (n *Logger) Panic(args ...interface{})                    { panic(fmt.Sprint(args...)) }
