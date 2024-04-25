package main

import (
	"github.com/spf13/cobra"

	logrusFlag "github.com/toknowwhy/theunit-oracle/pkg/log/logrus/flag"
)

type options struct {
	LogVerbosity   string
	LogFormat      logrusFlag.FormatTypeValue
	ConfigFilePath string
	Config         Config
}

func NewRootCommand(opts *options) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "spire-bootstrap",
		Version:       "1",
		Short:         "",
		Long:          ``,
		SilenceErrors: false,
		SilenceUsage:  true,
	}

	rootCmd.PersistentFlags().StringVarP(
		&opts.LogVerbosity,
		"log.verbosity", "v",
		"info",
		"verbosity level",
	)
	rootCmd.PersistentFlags().Var(
		&opts.LogFormat,
		"log.format",
		"log format",
	)
	rootCmd.PersistentFlags().StringVarP(
		&opts.ConfigFilePath,
		"config", "c",
		"./config.json",
		"ghost config file",
	)

	return rootCmd
}
