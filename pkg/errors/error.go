package errors

import (
	"errors"
	"fmt"
)

func New(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return errors.New(msg)
}
