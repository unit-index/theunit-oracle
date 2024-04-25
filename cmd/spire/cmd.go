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
	Version        string
}

func NewRootCommand(opts *options) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "spire",
		Version:       opts.Version,
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
		"config",
		"c",
		"./config.json",
		"spire config file",
	)

	rootCmd.AddCommand(
		NewAgentCmd(opts),
		NewPullCmd(opts),
		NewPushCmd(opts),
	)

	return rootCmd
}
