// Package shell exposes utilities for creating and managing long-running shell
// processes, while reading and writing data to them interactively
package shell

// Shell represents a shell worker
type Shell interface {
	// Import a source file into the running process
	Include(file string) error
	// Execute a command in the running process
	Exec(cmd string, defs map[string]string) (out []byte, err error)
	// Inform the caller whether a command is currently running
	Busy() bool
}
