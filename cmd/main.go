package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/vchitai/x-crafter/breaker"
	"github.com/vchitai/x-crafter/builder"
	"github.com/vchitai/x-crafter/parser"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "x-crafter",
		Short: fmt.Sprintf("%s is used to quickly make a go code prototype quickly become reproducible", "X-Crafter"),
	}
	rootCmd.AddCommand(breaker.Command())
	rootCmd.AddCommand(parser.Command())
	rootCmd.AddCommand(builder.Command())
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
