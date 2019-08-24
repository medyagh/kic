package runner

import "io"

// Cmd abstracts over running a command somewhere, this is useful for testing
type Cmd interface {
	Run() error
	// Each entry should be of the form "key=value"
	SetEnv(...string) Cmd
	SetStdin(io.Reader) Cmd
	SetStdout(io.Writer) Cmd
	SetStderr(io.Writer) Cmd
}

// Cmder abstracts over creating commands
type Cmder interface {
	// command, args..., just like os/runner.Cmd
	Command(string, ...string) Cmd
}
