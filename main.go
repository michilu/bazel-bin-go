package main

import (
	"github.com/spf13/cobra"

	"v/cmd"

	"github.com/michilu/bazel-bin-go/cmds/copy"
	"github.com/michilu/bazel-bin-go/cmds/echo"
	"github.com/michilu/bazel-bin-go/cmds/version"
)

func main() {
	cmd.AddCommand([]func() (*cobra.Command, error){
		copy.New,
		echo.New,
		version.New,
	})
	cmd.Execute()
}
