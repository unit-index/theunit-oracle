package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/toknowwhy/theunit-oracle/pkg/unit"
	"os"
)

func NewSupplyCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:     "supply [Token ...]",
		Aliases: []string{"prices"},
		Args:    cobra.MinimumNArgs(0),
		Short:   "Return supply for given TOKEN",
		Long:    `Return supply for given TOKEN.`,
		RunE: func(c *cobra.Command, args []string) (err error) {
			srv, err := PrepareUnitClientServices(context.Background(), opts)
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

			fmt.Println(args)

			tokens, err := unit.NewTokens(args...)
			if err != nil {
				return err
			}

			fmt.Println(tokens)

			//supply, err := srv.Unit.TokenTotalSupply(tokens)
			//if err != nil {
			//	return err
			//}
			//fmt.Println(supply)
			//supply, err := srv.Gofer.TokenTotalSupply(tokens)

			//fmt.Println(supply)
			//for _, p := range supply {
			//	if mErr := srv.Marshaller.Write(os.Stdout, p); mErr != nil {
			//		_ = srv.Marshaller.Write(os.Stderr, mErr)
			//	}
			//}
			//
			//// If any pair was returned with an error, then we should return a non-zero status code.
			//for _, p := range prices {
			//	if p.Error != "" {
			//		exitCode = 1
			//		break
			//	}
			//}

			return
		},
	}
}
