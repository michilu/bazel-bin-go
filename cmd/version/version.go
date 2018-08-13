package version

import (
	"github.com/spf13/cobra"

	"github.com/michilu/bazel-bin-go/meta"
)

// AddCommand adds commands to given the command.
func AddCommand(cmd *cobra.Command) {
	cmd.AddCommand(newCmd())
}

func newCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Print(meta.Get())
		},
	}
}
