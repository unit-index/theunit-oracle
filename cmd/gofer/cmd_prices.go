package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
)

func NewPricesCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:     "prices [PAIR...]",
		Aliases: []string{"price"},
		Args:    cobra.MinimumNArgs(0),
		Short:   "Return prices for given PAIRs",
		Long:    `Return prices for given PAIRs.`,
		RunE: func(c *cobra.Command, args []string) (err error) {
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

			prices, err := srv.Gofer.Prices(pairs...)

			if err != nil {
				return err
			}

			for _, p := range prices {
				//fmt.Println(p)
				if mErr := srv.Marshaller.Write(os.Stdout, p); mErr != nil {
					_ = srv.Marshaller.Write(os.Stderr, mErr)
				}
			}

			// If any pair was returned with an error, then we should return a non-zero status code.
			for _, p := range prices {
				if p.Error != "" {
					exitCode = 1
					break
				}
			}

			return
		},
	}
}
