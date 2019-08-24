package runner

// DefaultCmder is the default commander
var DefaultCmder = &LocalCmder{}

// Command is a convience wrapper over DefaultCmder.Command
func Command(cmdStr string, args ...string) Cmd {
	return DefaultCmder.Command(cmdStr, args...)
}
