package runner

import (
	"bufio"
	"bytes"
)

// CombinedOutputLines is like os/exec's cmd.CombinedOutput(),
// but over our Cmd interface, and instead of returning the byte buffer of
// stderr + stdout, it scans these for lines and returns a slice of output lines
func CombinedOutputLines(c Cmd) (lines []string, err error) {
	var buff bytes.Buffer
	c.SetStdout(&buff)
	c.SetStderr(&buff)
	err = c.Run()
	scanner := bufio.NewScanner(&buff)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, err
}
