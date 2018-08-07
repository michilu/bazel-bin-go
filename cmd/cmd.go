package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/michilu/bazel-bin-go/errs"
	"github.com/michilu/bazel-bin-go/log"
	"github.com/michilu/bazel-bin-go/meta"
)

var (
	app *cobra.Command
	f   *flag
)

type (
	flag struct {
		debug bool

		config string
	}
)

func init() {
	app = &cobra.Command{
		Use:   meta.Name(),
		Short: "A command-line tool that copies the Go files from the bazel-bin directory to anywhere.",
	}
	f = &flag{}
	app.PersistentFlags().BoolVar(&f.debug, "debug", false, "debug mode")
	app.PersistentFlags().StringVar(&f.config, "config", "", fmt.Sprintf("config file (default is %s.yaml)", meta.Name()))
	cobra.OnInitialize(initialize)
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
	if err != nil {
		log.Debug().
			Err(&errs.Error{Op: op, Err: err}).
			Msg("error")
		return
	}
	log.Debug().
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
}
