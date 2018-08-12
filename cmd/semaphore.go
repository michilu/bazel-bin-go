package cmd

import (
	sem "github.com/marusama/semaphore"
	"google.golang.org/grpc/codes"

	"github.com/michilu/bazel-bin-go/errs"
	"github.com/michilu/bazel-bin-go/log"
)

var (
	semaphore sem.Semaphore
)

func initSem() {
	const op = "cmd.initSem"

	f := flag

	if f.parallel < 1 {
		log.Logger().Fatal().
			Int("parallel", f.parallel).
			Err(&errs.Error{Op: op, Code: codes.InvalidArgument.String(), Message: "parallel must be 1 or more"}).
			Msg("error")
	}

	semaphore = sem.New(f.parallel)
}
