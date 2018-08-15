package cmd

import (
	"os"
	"path/filepath"

	valid "github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"

	"github.com/michilu/bazel-bin-go/bus"
	"github.com/michilu/bazel-bin-go/errs"
	"github.com/michilu/bazel-bin-go/log"
	"github.com/michilu/bazel-bin-go/meta"
)

var (
	app *cobra.Command
)

func init() {
	const op = "cmd.init"
	app = &cobra.Command{
		Use:   meta.Name(),
		Short: "A command-line tool that copies the Go files from the bazel-bin directory to anywhere.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunE(cmd, args, flag)
		},
	}
	initFlag()
	cobra.OnInitialize(initialize)
	for _, n := range cmds {
		c, err := n()
		if err != nil {
			log.Logger().Fatal().
				Err(&errs.Error{Op: op, Err: err}).
				Msg("error")
		}
		app.AddCommand(c)
	}
}

func preRunE(cmd *cobra.Command, args []string, f *flags) error {
	const op = "cmd.preRunE"
	ok, err := valid.ValidateStruct(&opt{f.config})
	if err != nil {
		return &errs.Error{Op: op, Err: err}
	}
	if !ok {
		return &errs.Error{Op: op, Code: codes.InvalidArgument, Message: "invalid arguments"}
	}
	return nil
}

func initialize() {
	const op = "cmd.initialize"
	f := flag

	if f.debug {
		err := log.SetLevel("debug")
		if err != nil {
			log.Logger().Fatal().
				Err(&errs.Error{Op: op, Err: err}).
				Msg("error")
		}
	}

	switch f.config {
	case "":
		viper.AddConfigPath(filepath.Dir(os.Args[0]))
		viper.SetConfigName(meta.Name())
	default:
		viper.SetConfigFile(f.config)
	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	switch err.(type) {
	case nil,
		viper.ConfigFileNotFoundError:
	default:
		log.Logger().Fatal().
			Err(&errs.Error{Op: op, Err: err}).
			Msg("error")
	}

	log.Logger().Debug().
		Str("op", op).
		Str("config", viper.ConfigFileUsed()).
		Msg("using config file")

	debugFlag()
	setSem()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the app.
func Execute() {
	err := app.Execute()
	if err != nil {
		log.Logger().Fatal().
			Err(err).
			Msg("error")
	}
	bus.Wait()
}
