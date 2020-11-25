package errors

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"sync"
	"unicode/utf8"
)

var (
	ShowStackTraceOnError bool
	MaxStackTraceSize     = 2048
)

type Error interface {
	Message() string
	StackTrace() string
	Cause() error
	Error() string
	ErrorWithStackTrace() string
}

type BaseError struct {
	message    string
	stackTrace string
	cause      error
}

func New(msg string) *BaseError {
	return _new(nil, msg, 3)
}

func Errorf(format string, args ...interface{}) *BaseError {
	return _new(nil, fmt.Sprintf(format, args...), 3)
}

func _new(cause error, msg string, skipStack int) *BaseError {
	err := BaseError{
		message: msg,
		cause:   cause,
	}
	err.stackTrace = StackTrace(skipStack, MaxStackTraceSize)
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
	var msg string
	if e.message != "" {
		msg = e.message
	} else if e.cause != nil {
		msg = e.cause.Error()
	}
	if ShowStackTraceOnError {
		return fmt.Sprintf("message: %s\nstacktrace:\n%s", msg, e.ErrorWithStackTrace())
	} else {
		return msg
	}
}

func (e *BaseError) ErrorWithStackTrace() string {
	var buf bytes.Buffer
	buf.WriteString(e.message)
	buf.WriteString("\n")
	buf.WriteString(e.stackTrace)
	if e.cause != nil {
		buf.WriteString("  caused by \n")
		buf.WriteString(e.cause.Error())
		buf.WriteString("\n")
	}
	return buf.String()
}

func Wrap(e error, msg string) error {
	if e == nil {
		return nil
	}
	return _new(e, msg, 3)
}

func Wrapf(e error, msg string, args ...interface{}) error {
	if e == nil {
		return nil
	}
	return _new(e, fmt.Sprintf(msg, args...), 3)
}

func WrapOr(e error) error {
	if e == nil {
		return nil
	}
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

func Find(e error, f func(e error) bool) error {
	for {
		if f(e) {
			return e
		}
		err, ok := e.(Error)
		if !ok {
			return nil
		}
		cause := err.Cause()
		if cause == nil {
			return nil
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

func (me MultiError) HasAdditional() interface{} {
	return me
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

// SyncMultiError describes synchronized MultiError.
// Note: wrapped MultiError is not synchronized if you write directly.
type SyncMultiError struct {
	MultiError
	sync.Mutex
}

// Append append an error.
func (sm *SyncMultiError) Append(err error) {
	sm.Lock()
	defer sm.Unlock()
	sm.MultiError = append(sm.MultiError, err)
}

// Len returns length of errors.
func (sm *SyncMultiError) Len() int {
	return len(sm.MultiError)
}
