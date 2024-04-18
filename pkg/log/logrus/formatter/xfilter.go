package formatter

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// XFilterFormatter removes all fields with the "x-" prefix. This will allow
// adding more data fields to logs without making the CLI output to messy.
type XFilterFormatter struct {
	Formatter logrus.Formatter
}

func (f *XFilterFormatter) Format(e *logrus.Entry) ([]byte, error) {
	data := logrus.Fields{}
	for k, v := range e.Data {
		if !strings.HasPrefix(k, "x-") {
			data[k] = v
		}
	}
	e.Data = data
	return f.Formatter.Format(e)
}
