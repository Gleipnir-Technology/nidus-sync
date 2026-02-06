package report

import (
	"fmt"
)

type ErrorWithCode struct {
	code    string
	err     error
	message string
}

func (e *ErrorWithCode) Code() string {
	return e.code
}
func (e *ErrorWithCode) Error() string {
	return e.message
}

func newErrorWithCode(code string, format string, args ...any) *ErrorWithCode {
	if len(args) > 0 {
		return &ErrorWithCode{
			err:  fmt.Errorf(format, args...),
			code: code,
		}
	} else {
		return &ErrorWithCode{
			code:    code,
			err:     nil,
			message: format,
		}
	}
}
