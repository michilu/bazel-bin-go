package errs

// https://middlemost.com/failure-is-your-domain/

import (
	"bytes"
	"fmt"
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
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if e.Op != "" {
		fmt.Fprintf(&buf, "%s: ", e.Op)
	}

	// If wrapping an error, print its Error() message.
	// Otherwise print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			fmt.Fprintf(&buf, "<%s> ", e.Code)
		}
		buf.WriteString(e.Message)
	}
	return buf.String()
}
