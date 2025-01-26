// cmd/mini-eth/mini-eth.go
package main

import (
	cli "github.com/edumar111/my-geth-edu/cli"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mini-eth",
		Short: "Cliente mini Ethereum",
	}

	// Subcomandos
	rootCmd.AddCommand(cli.InitCmd())
	rootCmd.AddCommand(cli.RunCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
