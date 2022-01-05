package lang

import "fmt"

type Error struct {
	Code   ErrorCode
	Format string
	Args   []interface{}
}

var _ error = (*Error)(nil)

func NewError(code ErrorCode, format string, args ...interface{}) *Error {
	return &Error{Code: code, Format: format, Args: args}
}

func (err *Error) Error() string {
	return fmt.Sprintf("runtime error: %v: %v", err.Code, fmt.Sprintf(err.Format, err.Args...))
}

type ErrorCode int

const (
	Unknown ErrorCode = iota
	NotFound
	InvalidType
	InvalidArgument
)

func (code ErrorCode) String() string {
	switch code {
	case Unknown:
		return "Unknown"
	case NotFound:
		return "NotFound"
	case InvalidType:
		return "InvalidType"
	case InvalidArgument:
		return "InvalidArgument"
	default:
		return fmt.Sprintf("%T(%d)", code, int(code))
	}
}
