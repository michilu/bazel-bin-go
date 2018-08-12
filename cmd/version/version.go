package version

import (
	"github.com/spf13/cobra"

	"github.com/michilu/bazel-bin-go/meta"
)

func AddCommand(c *cobra.Command) {
	c.AddCommand(newCmd())
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