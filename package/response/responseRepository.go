package response

import (
	"fmt"
)

func PrepareError(operation string, err error) error {
	return fmt.Errorf("repository: failed prepare %s statement: %w", operation, err)
}

func ExecError(operation string, err error) error {
	return fmt.Errorf("repository: failed execute %s statement: %w", operation, err)
}

func Error(message string, err error) error {
	return fmt.Errorf("repository: %s: %w", message, err)
}
