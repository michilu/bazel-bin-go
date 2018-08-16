package main

import (
	"github.com/spf13/cobra"

	"v/cmd"

	"github.com/michilu/bazel-bin-go/cmd/copy"
	"github.com/michilu/bazel-bin-go/cmd/echo"
	"github.com/michilu/bazel-bin-go/cmd/version"
)

const (
	name   = "bazel-bin-go"
	semVer = "1.0.0-alpha"
)

var (
	ns = []func() (*cobra.Command, error){
		copy.New,
		echo.New,
		version.New,
	}
)

func main() {
	cmd.Execute()
}
