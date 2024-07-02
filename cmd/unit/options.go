package main

import (
	"fmt"
	"strings"

	logrusFlag "github.com/toknowwhy/theunit-oracle/pkg/log/logrus/flag"

	"github.com/toknowwhy/theunit-oracle/internal/gofer/marshal"
)

// These are the command options that can be set by CLI flags.
type options struct {
	LogVerbosity   string
	LogFormat      logrusFlag.FormatTypeValue
	ConfigFilePath string
	Format         formatTypeValue
	Config         Config
	NoRPC          bool
	Version        string
}

var formatMap = map[marshal.FormatType]string{
	marshal.Plain:  "plain",
	marshal.Trace:  "trace",
	marshal.JSON:   "json",
	marshal.NDJSON: "ndjson",
}

// formatTypeValue is a wrapper for the FormatType to allow implement
// the flag.Value and spf13.pflag.Value interfaces.
type formatTypeValue struct {
	format marshal.FormatType
}

// Will return the default value if none is set
// and will fail if the `format` is set to an unsupported value for some reason.
func (v *formatTypeValue) String() string {
	if v != nil {
		return formatMap[v.format]
	}
	return formatMap[marshal.Plain]
}

func (v *formatTypeValue) Set(s string) error {
	s = strings.ToLower(s)

	for ct, st := range formatMap {
		if s == st {
			v.format = ct
			return nil
		}
	}

	return fmt.Errorf("unsupported format")
}

func (v *formatTypeValue) Type() string {
	return "plain|trace|json|ndjson"
}
