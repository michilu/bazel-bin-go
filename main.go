package main

import (
	"github.com/spf13/cobra"

	"v/cmd"

	"github.com/michilu/bazel-bin-go/cmd/copy"
	"github.com/michilu/bazel-bin-go/cmd/echo"
	"github.com/michilu/bazel-bin-go/cmd/version"
)

func main() {
	cmd.AddCommand([]func() (*cobra.Command, error){
		copy.New,
		echo.New,
		version.New,
	})
	cmd.Execute()
}
