package builder

import (
	"os"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var guidePath string
	cmd := &cobra.Command{
		Use:   "Build [SOURCE] [DESTINATION]",
		Short: "Using your broken pile to once again craft your thing again",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				source      = args[0]
				destination = args[1]
			)
			bdr, err := Create(
				WithSourceFS(os.DirFS(source), "."),
				WithGuidePath(guidePath),
			)
			if err != nil {
				cmd.PrintErr(err)
			}
			if err := bdr.Execute(destination); err != nil {
				cmd.PrintErr(err)
			}
		},
	}
	cmd.Flags().StringVarP(&guidePath, "guide", "g", "", "The Build guide to rebuild project")
	return cmd
}
