package main

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/toknowwhy/theunit-oracle/pkg/transport/messages"
)

func NewPushCmd(opts *options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Args:  cobra.ExactArgs(1),
		Short: "",
		Long:  ``,
	}

	cmd.AddCommand(NewPushPriceCmd(opts))

	return cmd
}

func NewPushPriceCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "price",
		Args:  cobra.MaximumNArgs(1),
		Short: "",
		Long:  ``,
		RunE: func(_ *cobra.Command, args []string) error {
			var err error
			ctx := context.Background()
			srv, err := PrepareClientServices(ctx, opts)
			if err != nil {
				return err
			}
			if err = srv.Start(); err != nil {
				return err
			}
			defer srv.CancelAndWait()

			in := os.Stdin
			if len(args) == 1 {
				in, err = os.Open(args[0])
				if err != nil {
					return err
				}
			}

			// Read JSON and parse it:
			input, err := io.ReadAll(in)
			if err != nil {
				return err
			}

			msg := &messages.Price{}
			err = msg.Unmarshall(input)
			if err != nil {
				return err
			}

			// Send price message to RPC client:
			err = srv.Client.PublishPrice(msg)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
