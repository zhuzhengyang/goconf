package goconf

import (
	"fmt"
	"strings"
)

// errorFormatFunc is a function callback that is called by Error to
// turn the list of errors into a string.
type errorFormatFunc func([]error) string

func errArrayDotFormatFunc(es []error) string {
	var errstr []string = make([]string, 0)
	for i := 0; i < len(es); i++ {
		errstr = append(errstr, es[i].Error())
	}
	return strings.Join(errstr, ",")
}

// Error is an error type to track multiple errors. This is used to
// accumulate errors in cases and return them as a single "error".
type errArray struct {
	Errors      []error
	ErrorFormat errorFormatFunc
}

func (e errArray) Error() string {
	fn := e.ErrorFormat
	if fn == nil {
		fn = errArrayDotFormatFunc
	}
	return fn(e.Errors)
}

func (e *errArray) Push(err error) {
	if err == nil {
		return
	}
	e.Errors = append(e.Errors, err)
}

// ErrorOrNil returns an error interface if this Error represents
// a list of errors, or returns nil if the list of errors is empty. This
// function is useful at the end of accumulation to make sure that the value
// returned represents the existence of errors.
func (e *errArray) Err() error {
	if e == nil {
		return nil
	}
	if len(e.Errors) == 0 {
		return nil
	}
	return e
}

func (e *errArray) GoString() string {
	return fmt.Sprintf("*%#v", *e)
}

// WrappedErrors returns the list of errors that this Error is wrapping.
// It is an implementation of the errwrap.Wrapper interface so that
// multierror.Error can be used with that library.
//
// This method is not safe to be called concurrently and is no different
// than accessing the Errors field directly. It is implemented only to
// satisfy the errwrap.Wrapper interface.
func (e *errArray) WrappedErrors() []error {
	return e.Errors
}
