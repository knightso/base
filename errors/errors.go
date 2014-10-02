package errors

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"unicode/utf8"
)

type Error interface {
	Message() string
	StackTrace() string
	Cause() error
	Error() string
}

type BaseError struct {
	message    string
	stackTrace string
	cause      error
}

func New(msg string) *BaseError {
	return _new(nil, msg, 3)
}

func _new(cause error, msg string, skipStack int) *BaseError {
	err := BaseError{
		message: msg,
		cause:   cause,
	}
	err.stackTrace = StackTrace(skipStack, 2048)
	return &err
}

func (e *BaseError) Message() string {
	return e.message
}

func (e *BaseError) StackTrace() string {
	return e.stackTrace
}

func (e *BaseError) Cause() error {
	return e.cause
}

func (e *BaseError) Error() string {
	var buf bytes.Buffer
	buf.WriteString(e.message)
	buf.WriteString("\n")
	buf.WriteString(e.stackTrace)
	if e.cause != nil {
		buf.WriteString("  caused by\n")
		buf.WriteString(e.cause.Error())
		buf.WriteString("\n")
	}
	return buf.String()
}

func Wrap(e error, msg string) *BaseError {
	return _new(e, msg, 3)
}

func Wrapf(e error, msg string, args ...interface{}) *BaseError {
	return _new(e, fmt.Sprintf(msg, args...), 3)
}

func WrapOr(e error) *BaseError {
	if ee, ok := e.(*BaseError); ok {
		return ee
	}
	return _new(e, "", 3)
}

func Root(e error) error {
	for {
		err, ok := e.(Error)
		if !ok {
			return e
		}
		cause := err.Cause()
		if cause == nil {
			return err
		}
		e = cause
	}
}

type MultiError []error

func (me MultiError) Error() string {
	var buf bytes.Buffer
	for _, err := range me {
		buf.WriteString(err.Error())
		buf.WriteString("\n")
	}
	return buf.String()
}

func StackTrace(skip, maxBytes int) string {
	// this func is debug purpose and ignores errors

	buf := make([]byte, maxBytes)
	n := runtime.Stack(buf, false)
	var gotall bool
	if n < len(buf) {
		buf = buf[:n]
		gotall = true
	} else {
		for !utf8.Valid(buf) || len(buf) == 0 {
			buf = buf[:len(buf)-1]
		}
	}

	var w bytes.Buffer

	writeOrSkip := func(buf []byte, w io.Writer, line int) {
		if line == 1 || line > 1+skip*2 {
			w.Write(buf)
		}
	}

	line := 1
	for {
		lf := bytes.IndexByte(buf, '\n')
		if lf < 0 {
			writeOrSkip(buf, &w, line)
			break
		}
		writeOrSkip(buf[:lf+1], &w, line)
		buf = buf[lf+1:]
		line++
	}

	if !gotall {
		w.WriteString("\n        ... (omitted)")
	}
	w.WriteString("\n")

	return w.String()
}
