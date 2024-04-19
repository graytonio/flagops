package destinations

import (
	"fmt"
)

var _ Output = &ConsoleOutput{}

type ConsoleOutput struct {
}

func newConsoleOutput() (*ConsoleOutput, error) {
	return &ConsoleOutput{
	}, nil
}

// Init implements Output.
func (c *ConsoleOutput) Init() error {
	return nil
}

// ExecuteFile implements Output.
func (c *ConsoleOutput) ExecuteFile(path string, content []byte) error {
	fmt.Print(string(content))
	return nil
}

// Finalize implements Output.
func (c *ConsoleOutput) Finalize() error {
	return nil
}