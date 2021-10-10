package errors

import (
	stderrors "errors"
	"fmt"
	"io"
)

func Is(err error, target error) bool {
	return stderrors.Is(err, target)
}

func As(err error, target interface{}) bool {
	//goland:noinspection GoErrorsAs
	return stderrors.As(err, target)
}

func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

type stringError struct {
	msg string
}

func New(msg string) error {
	return &stringError{msg: msg}
}

func Newf(format string, v ...interface{}) error {
	return &stringError{msg: fmt.Sprintf(format, v...)}
}

func (e *stringError) Error() string {
	return e.msg
}

type wrapError struct {
	msg string
	err error
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return &wrapError{
		msg: msg,
		err: err,
	}
}

func Wrapf(err error, msg string, v ...interface{}) error {
	if err == nil {
		return nil
	}

	return &wrapError{
		msg: fmt.Sprintf(msg, v...),
		err: err,
	}
}

func (e *wrapError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.err)
}

func (e *wrapError) Unwrap() error {
	return e.err
}

//goland:noinspection GoUnhandledErrorResult
func (e *wrapError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.msg)
			fmt.Fprintf(s, ": %+v\n", e.err)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, e.Error())
	}
}
