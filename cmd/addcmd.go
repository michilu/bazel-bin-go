package cmd

import (
	echo "github.com/michilu/bazel-bin-go/cmd/echo"
)

func addCommand() {
	echo.AddCommand(app)
}
