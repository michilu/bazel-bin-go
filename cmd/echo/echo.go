package echo

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"

	"github.com/michilu/bazel-bin-go/bus"
	"github.com/michilu/bazel-bin-go/errs"
	"github.com/michilu/bazel-bin-go/log"
)

const (
	topic = "echo"
)

type (
	flag struct {
		filepath string
	}

	opt struct {
		F string `valid:"fileexists"`
	}
)

func AddCommand(c *cobra.Command) {
	c.AddCommand(newCmd())
}

func newCmd() *cobra.Command {
	const op = "cmd.echo.new"
	f := &flag{}
	c := &cobra.Command{
		Use:   "echo",
		Short: "echo",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunE(cmd, args, f)
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, f)
		},
	}
	c.Flags().StringVarP(&f.filepath, "file", "f", "", "path to an exists file")
	viper.BindPFlag("file", c.Flags().Lookup("file"))
	return c
}

func preRunE(cmd *cobra.Command, args []string, f *flag) error {
	const op = "cmd.echo.preRunE"
	ok, err := valid.ValidateStruct(&opt{f.filepath})
	if err != nil {
		return &errs.Error{Op: op, Err: err}
	}
	if !ok {
		return &errs.Error{Op: op, Code: codes.InvalidArgument.String(), Message: "invalid arguments"}
	}
	return nil
}

func run(cmd *cobra.Command, args []string, f *flag) {
	const op = "cmd.echo.run"

	bus.Subscribe(topic, echo)
	bus.Publish(topic, f.filepath)
}

func echo(s string) {
	const op = "cmd.echo.echo"

	defer bus.Unsubscribe(topic, echo)

	log.Debug().
		Str("op", op).
		Str("s", s).
		Msg("echo a file")

	log.Logger().Info().
		Msg(s)
}