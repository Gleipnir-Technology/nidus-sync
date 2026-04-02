package http

import (
	"fmt"
	"net/http"
)

type ErrorWithStatus struct {
	Message string
	Status  int
}

func (e *ErrorWithStatus) Error() string {
	return e.Message
}
func NewBadRequest(mesg_format string, args ...any) *ErrorWithStatus {
	return NewErrorStatus(http.StatusBadRequest, mesg_format, args...)
}
func NewError(mesg_format string, args ...any) *ErrorWithStatus {
	return NewErrorStatus(http.StatusInternalServerError, mesg_format, args...)
}
func NewErrorMaybe(mesg_format string, err error, args ...any) *ErrorWithStatus {
	if err == nil {
		return nil
	}
	allArgs := append([]any{err}, args...)
	return NewErrorStatus(http.StatusInternalServerError, mesg_format, allArgs...)
}
func NewErrorStatus(status int, mesg_format string, args ...any) *ErrorWithStatus {
	w := fmt.Errorf(mesg_format, args...)
	return &ErrorWithStatus{
		Message: w.Error(),
		Status:  status,
	}
}
func NewForbidden(mesg_format string, args ...any) *ErrorWithStatus {
	return NewErrorStatus(http.StatusForbidden, mesg_format, args...)
}
