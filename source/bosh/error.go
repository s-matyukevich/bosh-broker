package bosh

import (
	"fmt"
)

type boshError struct {
	command       string
	args          []string
	originalError error
	output        []byte
}

func (err boshError) Error() string {
	return fmt.Sprintf("Error occured during command execution:\nCommand: %s\nArgs: %v\nError: %s\nCommand output:\n%s", err.command, err.args, err.originalError, err.output)
}
