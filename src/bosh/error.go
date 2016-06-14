package bosh

import (
	"fmt"
)

type boshError struct {
	command       string
	originalError error
	output        []byte
}

func (err boshError) Error() string {
	return fmt.Sprintf("Error occured during command execution:Command: %s\nError: %s\nCommand output:\n%s", err.command, err.originalError, err.output)
}
