package cmd

import (
	"fmt"
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

	echo "github.com/michilu/bazel-bin-go/cmd/echo"
)

var (
	app *cobra.Command
	f   *flag
)

type (
	flag struct {
		config string
		debug  bool
	}

	opt struct {
		C string `valid:"fileexists"`
	}
)

func init() {
	f = &flag{}
	app = &cobra.Command{
		Use:   meta.Name(),
		Short: "A command-line tool that copies the Go files from the bazel-bin directory to anywhere.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunE(cmd, args, f)
		},
	}
	app.PersistentFlags().BoolVar(&f.debug, "debug", false, "debug mode")
	app.PersistentFlags().StringVar(&f.config, "config", "", fmt.Sprintf("config file (default is %s.yaml)", meta.Name()))
	cobra.OnInitialize(initialize)

	echo.AddCommand(app)
}

func preRunE(cmd *cobra.Command, args []string, f *flag) error {
	const op = "cmd.preRunE"
	ok, err := valid.ValidateStruct(&opt{f.config})
	if err != nil {
		return &errs.Error{Op: op, Err: err}
	}
	if !ok {
		return &errs.Error{Op: op, Code: codes.InvalidArgument.String(), Message: "invalid arguments"}
	}
	return nil
}

func initialize() {
	const op = "cmd.initialize"

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

	log.Debug().
		Str("op", op).
		Str("config", viper.ConfigFileUsed()).
		Msg("using config file")

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
