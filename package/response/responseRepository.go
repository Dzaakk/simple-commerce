package response

import (
	"errors"
	"fmt"
)

func PrepareError(operation string, err error) error {
	return fmt.Errorf("repository: failed to prepare %s statement: %w", operation, err)
}

func ExecError(operation string, err error) error {
	return fmt.Errorf("repository: failed to execute %s statement: %w", operation, err)
}

func InvalidParameter() error {
	return errors.New("repository: invalid parameter")
}

func Error(message string, err error) error {
	return fmt.Errorf("repository:  %s: %w", message, err)
}
