package main

import (
	"fmt"
	"os"
)

// exitCode to be returned by the application.
var exitCode = 0

func main() {
	opts := options{
		Version: "0.1.0",
	}

	rootCmd := NewRootCommand(&opts)
	rootCmd.AddCommand(
		NewPairsCmd(&opts),
		NewPricesCmd(&opts),
		NewAgentCmd(&opts),
		NewSupplyCmd(&opts),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %s\n", err)
		if exitCode == 0 {
			os.Exit(1)
		}
	}
	os.Exit(exitCode)
}
