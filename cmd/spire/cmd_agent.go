package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func NewAgentCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "agent",
		Args:  cobra.ExactArgs(0),
		Short: "",
		Long:  ``,
		RunE: func(_ *cobra.Command, _ []string) error {
			var err error
			ctx := context.Background()
			srv, err := PrepareAgentServices(ctx, opts)
			if err != nil {
				return err
			}
			if err = srv.Start(); err != nil {
				return err
			}
			defer srv.CancelAndWait()

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c

			return nil
		},
	}
}
