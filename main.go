package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func generate(source, dest string) {
	var (
		convention = loadConvention(source)
		parser     = newParser(convention)
	)

	if err := parser.parse(source, dest); err != nil {
		log.Fatal(err)
	}
	parser.wg.Wait()
}

func maker() *cobra.Command {
	var dest string
	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert your prototype into reproducible golang templates",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Printf("Dunno")
				return
			}
			generate(args[0], dest)
		},
	}
	cmd.Flags().StringVarP(&dest, "destination", "d", "./templates", "Template destination")
	return cmd
}

func main() {
	rootCmd := &cobra.Command{
		Use:   fmt.Sprintf("%s", "go-template-maker"),
		Short: fmt.Sprintf("%s is used to quickly make a go code prototype quickly become reproducible", "GoTemplateMaker"),
	}
	rootCmd.AddCommand(maker())
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
