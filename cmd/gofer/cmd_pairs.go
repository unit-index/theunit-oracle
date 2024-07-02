package main

import (
	"context"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
	"os"

	"github.com/spf13/cobra"
)

func NewPairsCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:     "pairs [PAIR...]",
		Aliases: []string{"pair"},
		Args:    cobra.MinimumNArgs(0),
		Short:   "List all supported asset pairs",
		Long:    `List all supported asset pairs.`,
		RunE: func(_ *cobra.Command, args []string) (err error) {
			srv, err := PrepareGoferClientServices(context.Background(), opts)
			if err != nil {
				return err
			}
			defer func() {
				if err != nil {
					exitCode = 1
					_ = srv.Marshaller.Write(os.Stderr, err)
				}
				_ = srv.Marshaller.Flush()
				// Set err to nil because error was already handled by marshaller.
				err = nil
			}()
			if err = srv.Start(); err != nil {
				return err
			}
			defer srv.CancelAndWait()

			pairs, err := gofer.NewPairs(args...)
			if err != nil {
				return err
			}

			models, err := srv.Gofer.Models(pairs...)
			if err != nil {
				return err
			}

			for _, p := range models {
				if mErr := srv.Marshaller.Write(os.Stdout, p); mErr != nil {
					_ = srv.Marshaller.Write(os.Stderr, mErr)
				}
			}

			return
		},
	}
}
