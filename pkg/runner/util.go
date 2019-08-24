package runner

import (
	"bufio"
	"bytes"
	"io"
	"os"
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

// RunLoggingOutputOnFail runs the cmd, logging error output if Run returns an error
func RunLoggingOutputOnFail(c Cmd) (string, error) {
	var buff bytes.Buffer
	c.SetStdout(&buff)
	c.SetStderr(&buff)
	err := c.Run()
	scanner := bufio.NewScanner(&buff)
	return string(scanner.Bytes()), err
}

// RunWithStdoutReader runs cmd with stdout piped to readerFunc
func RunWithStdoutReader(c Cmd, readerFunc func(io.Reader) error) error {
	pr, pw, err := os.Pipe()
	if err != nil {
		return err
	}
	defer pw.Close()
	defer pr.Close()
	c.SetStdout(pw)

	errChan := make(chan error, 1)
	go func() {
		errChan <- readerFunc(pr)
		pr.Close()
	}()

	err = c.Run()
	if err != nil {
		return err
	}
	err2 := <-errChan
	if err2 != nil {
		return err2
	}
	return nil
}
