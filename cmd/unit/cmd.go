package main

import (
	"github.com/spf13/cobra"
)

func NewRootCommand(opts *options) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "unit",
		Version: opts.Version,
		Short:   "Tool for providing reliable data in the blockchain ecosystem",
		Long: `
Gofer is a CLI interface for the Gofer Go Library.

It is a tool that allows for easy data retrieval from various sources
with aggregates that increase reliability in the DeFi environment.`,
		SilenceErrors: true,
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
		"config file",
	)
	rootCmd.PersistentFlags().VarP(
		&opts.Format,
		"format",
		"f",
		"output format",
	)
	rootCmd.PersistentFlags().BoolVar(
		&opts.NoRPC,
		"norpc",
		false,
		"disable the use of RPC agent",
	)

	return rootCmd
}
