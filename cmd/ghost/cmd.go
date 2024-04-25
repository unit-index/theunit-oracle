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
	GoferNoRPC     bool
}

func NewRootCommand(opts *options) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "ghost",
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
	rootCmd.PersistentFlags().BoolVar(
		&opts.GoferNoRPC,
		"gofer.norpc",
		false,
		"disable the use of Gofer RPC agent",
	)

	return rootCmd
}
