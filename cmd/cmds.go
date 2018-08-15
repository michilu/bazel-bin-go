package cmd

import (
	"github.com/spf13/cobra"

	"github.com/michilu/bazel-bin-go/cmds/copy"
	"github.com/michilu/bazel-bin-go/cmds/echo"
	"github.com/michilu/bazel-bin-go/cmds/version"
)

var (
	cmds = []func() (*cobra.Command, error){
		copy.New,
		echo.New,
		version.New,
	}
)
