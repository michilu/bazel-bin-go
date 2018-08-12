package cmd

import (
	"fmt"
	"runtime"

	"github.com/michilu/bazel-bin-go/log"
	"github.com/michilu/bazel-bin-go/meta"
)

var (
	flag *flags

	defaultP = runtime.NumCPU()
)

type (
	flags struct {
		config   string
		debug    bool
		parallel int
	}

	opt struct {
		C string `valid:"fileexists"`
	}
)

func Flag() flags {
	return *flag
}

func initFlag() {
	flag = &flags{}
	f := flag

	app.PersistentFlags().StringVar(&f.config, "config", "", fmt.Sprintf("config file (default is %s.yaml)", meta.Name()))
	app.PersistentFlags().BoolVar(&f.debug, "debug", false, "debug mode")

	if f.debug {
		defaultP = 0
	}
	app.PersistentFlags().IntVarP(&f.parallel, "parallel", "p", defaultP, "parallel")

}

func debugFlag() {
	const op = "cmd.debugFlag"

	e := log.Debug()
	if !e.Enabled() {
		return
	}

	f := flag
	e.
		Str("op", op).
		Str("config", f.config).
		Bool("debug", f.debug).
		Int("parallel", f.parallel).
		Msg("flag")
}
