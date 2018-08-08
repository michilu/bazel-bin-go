package log

import (
	"log"
	"os"

	"github.com/rs/zerolog"

	"github.com/michilu/bazel-bin-go/errs"
)

var (
	logger          zerolog.Logger
	timeFieldFormat string
)

func init() {
	setTimeFieldFormat()
	setDefaultLogger()
}

func setTimeFieldFormat() {
	timeFieldFormat = zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = ""
}

func setDefaultLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()
	log.SetOutput(logger)
	log.SetFlags(0)
}

func SetLevel(s string) error {
	const op = "log.SetLevel"
	l, err := zerolog.ParseLevel(s)
	if err != nil {
		return &errs.Error{Op: op, Err: err}
	}
	zerolog.SetGlobalLevel(l)
	switch l {
	case zerolog.DebugLevel:
		zerolog.TimeFieldFormat = timeFieldFormat
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().
			Timestamp().
			Logger()
		log.SetOutput(logger)
	}
	return nil
}

func Logger() *zerolog.Logger {
	return &logger
}

func Debug() *zerolog.Event {
	return logger.Debug()
}
