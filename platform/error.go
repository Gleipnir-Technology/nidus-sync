package platform

import (
	"fmt"
)

type ErrorNotFound struct {
	message string
}

func (e ErrorNotFound) Error() string { return fmt.Sprintf("not found: %s", e.message) }
func newNotFound(format string, m ...any) error {
	return &ErrorNotFound{message: fmt.Sprintf(format, m...)}
}
