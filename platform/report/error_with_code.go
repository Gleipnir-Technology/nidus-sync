package report

import (
	"fmt"

	"github.com/rs/zerolog/log"
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

func newInternalError(err error, format string, args ...any) *ErrorWithCode {
	log.Error().Err(err).Str("format", format).Msg("internal server error")
	return newErrorWithCode("internal-error", format, args...)
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
