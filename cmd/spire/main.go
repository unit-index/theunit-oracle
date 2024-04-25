package main

import (
	_ "embed"
	"os"
)

func main() {
	opts := options{Version: "1"}
	rootCmd := NewRootCommand(&opts)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
