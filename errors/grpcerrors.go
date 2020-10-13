package errors

import (
	"fmt"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorLevel int32

const (
	Unknown  ErrorLevel = 0
	Critical ErrorLevel = 1
	Error    ErrorLevel = 2
	Warning  ErrorLevel = 3
	Info     ErrorLevel = 4
	Debug    ErrorLevel = 5
)

var errorLevelNames = map[ErrorLevel]string{
	Unknown:  "UNKNOWN",
	Critical: "CRITICAL",
	Error:    "ERROR",
	Warning:  "WARNING",
	Info:     "INFO",
	Debug:    "DEBUG",
}

func (l ErrorLevel) String() string {
	if l, ok := errorLevelNames[l]; ok {
		return l
	}
	return fmt.Sprintf("invalid ErrorLevel %d", l)
}
func GetErrorLevel(level string) (ErrorLevel, error) {
	for k, v := range errorLevelNames {
		if v == level {
			return k, nil
		}
	}
	return 0, Errorf("invalid ErrorLevel name %s", level)
}

type GrpcError struct {
	origin error
	level  ErrorLevel
	stack  *stack
	status *status.Status
}

func NewError(code codes.Code, message string, args ...interface{}) *GrpcError {
	s := status.Newf(code, message, args...)
	return &GrpcError{status: s, origin: nil, level: Error, stack: callers()}
}

func NewWarning(code codes.Code, message string, args ...interface{}) *GrpcError {
	s := status.Newf(code, message, args...)
	return &GrpcError{status: s, origin: nil, level: Warning, stack: callers()}
}

func NewErrorFrom(origin error, code codes.Code, message string, args ...interface{}) *GrpcError {
	s := status.Newf(code, message, args...)
	return &GrpcError{status: s, origin: origin, level: Error, stack: callers()}
}

func NewWarningFrom(origin error, code codes.Code, message string, args ...interface{}) *GrpcError {
	s := status.Newf(code, message, args...)
	return &GrpcError{status: s, origin: origin, level: Warning, stack: callers()}
}

func NewGrpcError(origin error, code codes.Code, level ErrorLevel, message string, args ...interface{}) *GrpcError {
	s := status.Newf(code, message, args...)
	return &GrpcError{status: s, origin: origin, level: level, stack: callers()}
}

func GetCode(err error) codes.Code {
	if err == nil {
		return codes.Unknown
	}
	c := Cause(err)
	if e, ok := c.(*GrpcError); ok {
		return e.Status().Code()
	} else {
		return codes.Internal
	}
}

func GetLevel(err error) ErrorLevel {
	if err == nil {
		return Unknown
	}
	c := Cause(err)

	if e, ok := c.(*GrpcError); ok {
		return e.level
	}
	return Critical
}

func (e *GrpcError) Error() string {
	return fmt.Sprintf("[%v][%s] %v", e.level, e.status.Code(), e.status.Message())
}
func (e *GrpcError) String() string {
	return e.Error()
}

func (e *GrpcError) Status() *status.Status {
	return e.status
}

func (e *GrpcError) Origin() error {
	return e.origin
}

func (e *GrpcError) ErrorLevel() ErrorLevel {
	return e.level
}

func (e *GrpcError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			must(io.WriteString(s, e.Error()))
			e.stack.Format(s, verb)
			must(io.WriteString(s, "\n"))
			if e.origin != nil {
				must(io.WriteString(s, "... Caused by ... \n"))
				must(fmt.Fprintf(s, "%+v\n", e.origin))
			}
			return
		}
		fallthrough
	case 's':
		must(io.WriteString(s, e.Error()))
	case 'q':
		must(fmt.Fprintf(s, "%q", e.Error()))
	}
}

func must(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
