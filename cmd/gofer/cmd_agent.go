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
		Args:  cobra.NoArgs,
		Short: "Start an RPC server",
		Long:  `Start an RPC server.`,
		RunE: func(_ *cobra.Command, args []string) error {
			srv, err := PrepareGoferAgentService(context.Background(), opts)
			if err != nil {
				return err
			}
			if err = srv.Start(); err != nil {
				return err
			}
			defer srv.CancelAndWait()

			// Wait for the interrupt signal:
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c

			return nil
		},
	}
}
