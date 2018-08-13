package cmd

import (
	"github.com/spf13/cobra"

	"github.com/michilu/bazel-bin-go/cmd/copy"
	"github.com/michilu/bazel-bin-go/cmd/echo"
	"github.com/michilu/bazel-bin-go/cmd/version"
)

var (
	cmds = []func() (*cobra.Command, error){
		copy.New,
		echo.New,
		version.New,
	}
)
