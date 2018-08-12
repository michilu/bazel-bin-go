package cmd

import (
	"github.com/michilu/bazel-bin-go/cmd/echo"
	"github.com/michilu/bazel-bin-go/cmd/version"
)

func addCommand() {
	echo.AddCommand(app)
	version.AddCommand(app)
}
