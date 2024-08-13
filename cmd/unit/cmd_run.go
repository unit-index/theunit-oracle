package main

import (
	"context"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func NewFeedCmd(opts *options) *cobra.Command {
	return &cobra.Command{
		Use:     "run",
		Aliases: []string{"run"},
		Args:    cobra.MinimumNArgs(0),
		Short:   "Return prices for given TOKEN",
		Long:    `Return prices for given TOKEN.`,
		RunE: func(c *cobra.Command, args []string) (err error) {
			srv, err := PrepareUnitServerServices(context.Background(), opts)
			if err != nil {
				return err
			}
			//defer func() {
			//	if err != nil {
			//		exitCode = 1
			//		_ = srv.Marshaller.Write(os.Stderr, err)
			//	}
			//	_ = srv.Marshaller.Flush()
			//	// Set err to nil because error was already handled by marshaller.
			//	err = nil
			//}()
			if err = srv.Start(); err != nil {
				return err
			}
			defer srv.CancelAndWait()
			o := make(chan os.Signal, 1)
			signal.Notify(o, os.Interrupt, syscall.SIGTERM)
			<-o
			return
		},
	}
}
