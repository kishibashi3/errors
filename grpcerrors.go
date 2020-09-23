package errors

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type ErrorLevel int32

const (
	Critical ErrorLevel = 1
	Error    ErrorLevel = 2
	Warning  ErrorLevel = 3
	Info     ErrorLevel = 4
	Debug    ErrorLevel = 5
)

var errorLevelNames = map[ErrorLevel]string{
	Critical: "CRITICAL",
	Error:    "ERROR",
	Warning:  "WARNING",
	Info:     "INFO",
	Debug:    "DEBUG",
}

func (l ErrorLevel) String() string {
	return errorLevelNames[l]
}
func GetErrorLevel(level string) (ErrorLevel, error) {
	for k, v := range errorLevelNames {
		if v == level {
			return k, nil
		}
	}
	return 0, Errorf("illegal ErrorLevel %s", level)
}

type GrpcError struct {
	origin error
	level  ErrorLevel
	stack  *stack
	status *status.Status
}

func NewError(code codes.Code, message string, args ...interface{}) *GrpcError {
	return NewGrpcError(nil, code, Error, message, args...)
}

func NewWarning(code codes.Code, message string, args ...interface{}) *GrpcError {
	return NewGrpcError(nil, code, Warning, message, args...)
}

func NewErrorFrom(origin error, code codes.Code, message string, args ...interface{}) *GrpcError {
	return NewGrpcError(origin, code, Error, message, args...)
}

func NewWarningFrom(origin error, code codes.Code, message string, args ...interface{}) *GrpcError {
	return NewGrpcError(origin, code, Warning, message, args...)
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

func GetLevel(err error) (ErrorLevel, error) {
	if err == nil {
		return 0, New("error is null")
	}
	c := Cause(err)

	if e, ok := c.(*GrpcError); ok {
		return e.level, nil
	}
	return Critical, nil
}

func (e *GrpcError) Error() string {
	if e == nil {
		return "null pointer"
	}
	return fmt.Sprintf("[%v/%s] %v", e.level, e.status.Code(), e.status.Message())
}
func (e *GrpcError) String() string {
	if e == nil {
		return "null pointer"
	}
	return e.Error()
}

func (e *GrpcError) Status() *status.Status {
	if e == nil {
		return nil
	}
	return e.status
}

func (e *GrpcError) Origin() error {
	if e == nil {
		return nil
	}
	return e.origin
}

func (e *GrpcError) ErrorLevel() ErrorLevel {
	if e == nil {
		return 0
	}
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
