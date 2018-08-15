package main

import (
	"github.com/spf13/cobra"

	"v/cmd"
	"v/errs"
	"v/log"
	"v/meta"

	"github.com/michilu/bazel-bin-go/cmd/copy"
	"github.com/michilu/bazel-bin-go/cmd/echo"
	"github.com/michilu/bazel-bin-go/cmd/version"
)

func init() {
	const op = "main.init"
	err := meta.Set(&meta.Meta{
		Build:  build,
		Hash:   hash,
		Name:   name,
		SemVer: semVer,
		Serial: serial,
	})
	if err != nil {
		log.Logger().Fatal().
			Str("op", op).
			Err(&errs.Error{Op: op, Err: err}).
			Msg("error")
	}
	cmd.Init()
}

func main() {
	cmd.AddCommand([]func() (*cobra.Command, error){
		copy.New,
		echo.New,
		version.New,
	})
	cmd.Execute()
}
