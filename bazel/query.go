package bazel

import (
	"bytes"
	"io"
	"os/exec"
	"strings"

	"github.com/michilu/bazel-bin-go/bus"
	"github.com/michilu/bazel-bin-go/errs"
	"github.com/michilu/bazel-bin-go/log"
)

var (
	n = []byte("\n")[0]
)

func Query() error {
	const op = "bazel.Query"

	c := exec.Command("bazel", "query", "//...")

	log.Logger().Debug().
		Str("op", op).
		Str("path", c.Path).
		Strs("args", c.Args).
		Msg("create exec.Command")

	var o bytes.Buffer
	c.Stdout = &o
	err := c.Run()
	if err != nil {
		return &errs.Error{Op: op, Err: err}
	}

	log.Logger().Debug().
		Str("op", op).
		Str("path", c.Path).
		Strs("args", c.Args).
		Str("stdout", o.String()).
		Msg("ran exec.Command")

	for {
		s, err := o.ReadString(n)
		switch err {
		case io.EOF,
			nil:
		default:
			return &errs.Error{Op: op, Err: err}
		}

		ok := strings.Contains(s, "_go_")
		if !ok {
			break
		}

		bus.Publish("test", s)

		if err == io.EOF {
			break
		}
	}

	return nil
}
