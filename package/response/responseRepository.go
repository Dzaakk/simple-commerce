package response

import (
	"fmt"
)

func Error(message string, err error) error {
	return fmt.Errorf("repository: %s: %w", message, err)
}
