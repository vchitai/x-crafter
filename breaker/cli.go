package breaker

import (
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

func dismantleLayers(source, dest string) {
	var (
		guide   = loadGuide(source)
		breaker = newBreaker(source, guide)
	)
	if err := breaker.make(source, dest); err != nil {
		log.Fatal(err)
	}
	breaker.wg.Wait()
}

func Command() *cobra.Command {
	var dest = "./layers"
	cmd := &cobra.Command{
		Use:   "break",
		Short: "Convert your prototype into recraftable golang templates",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Printf("Please provide arguments. Use --help for more help options.")
				return
			}
			var source = args[0]
			if len(args) > 1 {
				dest = args[1]
			} else {
				dest = source + "_broken"
			}
			dismantleLayers(source, filepath.Join(dest, "layers"))
			versioning(dest)
		},
	}
	return cmd
}
