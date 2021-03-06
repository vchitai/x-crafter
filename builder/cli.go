package builder

import (
	"os"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var guidePath string
	cmd := &cobra.Command{
		Use:   "build [SOURCE] [DESTINATION]",
		Short: "Using your broken pile to once again craft your thing again",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				source      = args[0]
				destination = args[1]
			)
			bdr, err := New(
				WithGuidePath(guidePath),
				WithSourceFS(os.DirFS(source), "."),
			)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			if err := bdr.Execute(destination); err != nil {
				cmd.PrintErr(err)
			}
		},
	}
	cmd.Flags().StringVarP(&guidePath, "guide", "g", "", "The Build guide to rebuild project")
	return cmd
}
