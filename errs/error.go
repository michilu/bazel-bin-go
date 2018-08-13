package errs

// https://middlemost.com/failure-is-your-domain/

import (
	"bytes"
	"fmt"

	"github.com/michilu/bazel-bin-go/log"
)

// Error defines a standard application error.
type Error struct {
	// Machine-readable error code.
	Code string

	// Human-readable message.
	Message string

	// Logical operation and nested error.
	Op  string
	Err error
}

// Error returns the string representation of the error message.
func (e *Error) Error() string {
	const op = "errs.Error.Error()"
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if e.Op != "" {
		_, err := fmt.Fprintf(&buf, "%s: ", e.Op)
		if err != nil {
			log.Logger().Error().
				Str("op", op).
				Err(&Error{Op: op, Err: err}).
				Msg("error")
		}
	}

	// If wrapping an error, print its Error() message.
	// Otherwise print the error code & message.
	if e.Err != nil {
		_, err := buf.WriteString(e.Err.Error())
		if err != nil {
			log.Logger().Error().
				Str("op", op).
				Err(&Error{Op: op, Err: err}).
				Msg("error")
		}
	} else {
		if e.Code != "" {
			_, err := fmt.Fprintf(&buf, "<%s> ", e.Code)
			if err != nil {
				log.Logger().Error().
					Str("op", op).
					Err(&Error{Op: op, Err: err}).
					Msg("error")
			}
		}
		_, err := buf.WriteString(e.Message)
		if err != nil {
			log.Logger().Error().
				Str("op", op).
				Err(&Error{Op: op, Err: err}).
				Msg("error")
		}
	}
	return buf.String()
}
