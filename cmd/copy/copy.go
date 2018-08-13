package copy

import (
	"bufio"
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

const (
	topic = "copy"
)

var (
	sem = semaphore.New(1)

	wg sync.WaitGroup
)

type (
	flag struct {
		from string
		to   string
	}

	opt struct {
		F string `valid:"direxists"`
		T string `valid:"direxists"`
	}
)

// AddCommand adds commands to the given command.
func AddCommand(cmd *cobra.Command) {
	cmd.AddCommand(newCmd())
}

func newCmd() *cobra.Command {
	const op = "cmd.copy.new"
	f := &flag{}
	c := &cobra.Command{
		Use:   "copy",
		Short: "copy the Go files",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return preRunE(cmd, args, f)
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, f)
		},
	}
	c.Flags().StringVarP(&f.from, "from", "f", "", "a copy source directory")
	viper.BindPFlag("from", c.Flags().Lookup("from"))
	c.Flags().StringVarP(&f.to, "to", "t", "", "a copy destoribute directory")
	viper.BindPFlag("to", c.Flags().Lookup("to"))
	return c
}

func preRunE(cmd *cobra.Command, args []string, f *flag) error {
	const op = "cmd.copy.preRunE"
	ok, err := valid.ValidateStruct(&opt{f.from, f.to})
	if err != nil {
		return &errs.Error{Op: op, Err: err}
	}
	if !ok {
		return &errs.Error{Op: op, Code: codes.InvalidArgument.String(), Message: "invalid arguments"}
	}
	return nil
}

func run(cmd *cobra.Command, args []string, f *flag) {
	const op = "cmd.copy.run"

	bus.Subscribe(topic, copyFile)
	defer wg.Wait()
	defer bus.Unsubscribe(topic, copyFile)

	err := filepath.Walk(f.from, func(p string, i os.FileInfo, err error) error {
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
		//lint:ignore SA1012 Pass a nil `context.Context` for speedup.
		_ = sem.Acquire(nil, 1)
		bus.Publish(topic, p, i, f.from, f.to)
		wg.Add(1)
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
	defer wg.Done()
	defer func() { sem.Release(1) }()

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
		if e := fi.Close(); e != nil {
			log.Logger().Warn().
				Err(&errs.Error{Op: op, Err: e}).
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
		if e := fo.Close(); e != nil {
			log.Logger().Warn().
				Err(&errs.Error{Op: op, Err: e}).
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
		n, e := r.Read(buf)
		if err != nil && e != io.EOF {
			log.Logger().Warn().
				Err(&errs.Error{Op: op, Err: e}).
				Msg("error")
			break
		}
		if n == 0 {
			break
		}
		if _, e := w.Write(buf[:n]); e != nil {
			log.Logger().Warn().
				Err(&errs.Error{Op: op, Err: e}).
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
