package cmd

import (
	"github.com/michilu/bazel-bin-go/cmd/copy"
	"github.com/michilu/bazel-bin-go/cmd/echo"
	"github.com/michilu/bazel-bin-go/cmd/version"
)

func addCommand() {
	copy.AddCommand(app)
	echo.AddCommand(app)
	version.AddCommand(app)
}
