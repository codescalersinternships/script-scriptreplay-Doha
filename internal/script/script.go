package script

import (
	"io"
	"os"
	"os/exec"

	"golang.org/x/term"
)

func Script(filename string) error {
	c := exec.Command("bash")
	pty, err := PtyFork(c)
	if err != nil {
		return err
	}
	defer pty.Close()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	outputWriter := io.MultiWriter(os.Stdout, logFile)

	done := make(chan error)

	go func() {
		_, err := io.Copy(outputWriter, pty)
		done <- err
	}()

	go func() {
		_, err = io.Copy(pty, os.Stdin)
		done <- err
	}()

	c.Wait()
	return err
}
