package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	valid "github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"

	"github.com/michilu/bazel-bin-go/bazel"
	"github.com/michilu/bazel-bin-go/bus"
	"github.com/michilu/bazel-bin-go/errs"
	"github.com/michilu/bazel-bin-go/log"
	"github.com/michilu/bazel-bin-go/semaphore"
)

var (
	wgCopy sync.WaitGroup
)

type (
	optCopy struct {
		From string `valid:"filepath"`
		To   string `valid:"filepath"`
	}
)

func init() {
	app.AddCommand(newCopy())
}

func newCopy() *cobra.Command {
	const o = "cmd.newCopy"
	var (
		f string
		t string
	)
	c := &cobra.Command{
		Use:   "copy",
		Short: "copy the Go files",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunECopy(cmd, args, f, t)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runCopy(cmd, args, f, t)
		},
	}
	c.Flags().StringVarP(&f, "from", "f", "", "a copy source directory")
	viper.BindPFlag("from", c.Flags().Lookup("from"))
	c.Flags().StringVarP(&t, "to", "t", "", "a copy destoribute directory")
	viper.BindPFlag("to", c.Flags().Lookup("to"))
	return c
}

func preRunECopy(cmd *cobra.Command, args []string, f string, t string) error {
	const op = "cmd.copy.preRunECopy"
	ok, err := valid.ValidateStruct(&optCopy{})
	if err != nil {
		return &errs.Error{Op: op, Err: err}
	}
	if !ok {
		return &errs.Error{Op: op, Code: codes.InvalidArgument.String(), Message: "invalid arguments"}
	}
	for _, s := range []string{f, t} {
		if s == "" {
			continue
		}
		i, err := os.Stat(s)
		if err != nil {
			return &errs.Error{Op: op, Err: err}
		}
		if !i.IsDir() {
			return &errs.Error{Op: op, Code: codes.InvalidArgument.String(), Message: fmt.Sprintf("must be a directory: %s", s)}
		}
	}
	return nil
}

func runCopy(cmd *cobra.Command, args []string, f string, t string) {
	const op = "cmd.copy.runCopy"

	log.Debug().
		Str("op", op).
		Str("from", f).
		Str("to", t).
		Msg("copy files")

	bus.Subscribe("copy", copyFile)
	defer wgCopy.Wait()
	defer bus.Unsubscribe("copy", copyFile)

	err := filepath.Walk(f, func(p string, i os.FileInfo, err error) error {
		const op = "filepath.Walk"

		log.Debug().
			Str("op", op).
			Str("path", p).
			Msg("start")

		if err != nil {
			return &errs.Error{Op: op, Err: err}
		}
		if i.IsDir() {
			log.Debug().
				Str("op", op).
				Str("path", p).
				Msg("skip directory")
			return nil
		}
		_ = semaphore.Acquire(nil, 1)
		bus.Publish("copy", p, i, f, t)
		wgCopy.Add(1)
		return nil
	})
	if err != nil {
		log.Logger().Fatal().
			Err(&errs.Error{Op: op, Err: err}).
			Msg("error")
	}
}

// https://stackoverflow.com/a/9739903/1085087
func copyFile(p string, i os.FileInfo, f string, t string) error {
	const op = "cmd.copy.copyFile"
	defer wgCopy.Done()
	defer func() { semaphore.Release(1) }()

	log.Debug().
		Str("op", op).
		Str("path", p).
		Msg("copy a file")

	fi, err := os.Open(p)
	if err != nil {
		log.Logger().Warn().
			Err(&errs.Error{Op: op, Err: err}).
			Msg("error")
		return nil
	}
	defer func() {
		const op = "input.Close"
		if err := fi.Close(); err != nil {
			log.Logger().Warn().
				Err(&errs.Error{Op: op, Err: err}).
				Msg("error")
		}
	}()

	log.Debug().
		Str("op", op).
		Str("path", p).
		Msg("opened a source file")

	err = bazel.Query()
	if err != nil {
		log.Logger().Warn().
			Err(&errs.Error{Op: op, Err: err}).
			Msg("error")
		return nil
	}

	fo, err := os.Open(t)
	if err != nil {
		log.Logger().Warn().
			Err(&errs.Error{Op: op, Err: err}).
			Msg("error")
		return nil
	}
	defer func() {
		const op = "output.Close"
		if err := fo.Close(); err != nil {
			log.Logger().Warn().
				Err(&errs.Error{Op: op, Err: err}).
				Msg("error")
		}
	}()

	log.Debug().
		Str("op", op).
		Str("path", t).
		Msg("opened a destoribute file")

	r := bufio.NewReader(fi)
	w := bufio.NewWriter(fo)
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			log.Logger().Warn().
				Err(&errs.Error{Op: op, Err: err}).
				Msg("error")
			break
		}
		if n == 0 {
			break
		}
		if _, err := w.Write(buf[:n]); err != nil {
			log.Logger().Warn().
				Err(&errs.Error{Op: op, Err: err}).
				Msg("error")
			break
		}
	}
	if err = w.Flush(); err != nil {
		log.Logger().Warn().
			Err(&errs.Error{Op: op, Err: err}).
			Msg("error")
		return nil
	}
	return nil
}
