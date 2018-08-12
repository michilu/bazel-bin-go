package cmd

import (
	"fmt"
	"runtime"

	"github.com/michilu/bazel-bin-go/meta"
)

var (
	flag *flags
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

func initFlag() {
	flag = &flags{}
	f := flag

	app.PersistentFlags().StringVar(&f.config, "config", "", fmt.Sprintf("config file (default is %s.yaml)", meta.Name()))
	app.PersistentFlags().BoolVar(&f.debug, "debug", false, "debug mode")

	defaultP := runtime.NumCPU()
	if f.debug {
		defaultP = 0
	}
	app.PersistentFlags().IntVarP(&f.parallel, "parallel", "p", defaultP, "parallel")
}
