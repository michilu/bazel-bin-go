package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/michilu/bazel-bin-go/meta"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(meta.Get())
		},
	})
}
