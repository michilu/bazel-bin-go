package cmd

import (
	"fmt"
	"os"

	valid "github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"

	"github.com/michilu/bazel-bin-go/bus"
	"github.com/michilu/bazel-bin-go/errs"
	"github.com/michilu/bazel-bin-go/log"
)

type (
	optEcho struct {
		F string `valid:"filepath"`
	}
)

func init() {
	app.AddCommand(newEcho())
}

func newEcho() *cobra.Command {
	const op = "cmd.newEcho"
	var (
		f string
	)
	c := &cobra.Command{
		Use:   "echo",
		Short: "echo",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunEEcho(cmd, args, f)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runEcho(cmd, args, f)
		},
	}
	c.Flags().StringVarP(&f, "file", "f", "", "filepath")
	viper.BindPFlag("file", c.Flags().Lookup("file"))
	return c
}

func preRunEEcho(cmd *cobra.Command, args []string, f string) error {
	const op = "cmd.echo.preRunEEcho"
	ok, err := valid.ValidateStruct(&optEcho{})
	if err != nil {
		return &errs.Error{Op: op, Err: err}
	}
	if !ok {
		return &errs.Error{Op: op, Code: codes.InvalidArgument.String(), Message: "invalid arguments"}
	}
	for _, s := range []string{f} {
		if s == "" {
			continue
		}
		i, err := os.Stat(s)
		if err != nil {
			return &errs.Error{Op: op, Err: err}
		}
		if i.IsDir() {
			return &errs.Error{Op: op, Code: codes.InvalidArgument.String(), Message: fmt.Sprintf("must be a file: %s", s)}
		}
	}
	return nil
}

func runEcho(cmd *cobra.Command, args []string, f string) {
	const op = "cmd.echo.runEcho"

	log.Debug().
		Str("op", op).
		Str("f", f).
		Msg("echo a file")

	bus.Subscribe("echo", echo)
	bus.Publish("echo", f)
}

func echo(s string) {
	const op = "cmd.echo.echo"

	defer bus.Unsubscribe("echo", echo)

	log.Debug().
		Str("op", op).
		Str("s", s).
		Msg("echo a file")

	log.Logger().Info().
		Msg(s)
}
