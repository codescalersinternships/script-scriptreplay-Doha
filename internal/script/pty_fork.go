package script

import (
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

type Winsize struct {
	Rows uint16
	Cols uint16
	X    uint16
	Y    uint16
}

func PtyFork(c *exec.Cmd) (*os.File, error) {
	pty, tty, err := PtyMasterOpen()
	if err != nil {
		return nil, err
	}
	defer tty.Close()


	sz, err := getSizeFull(os.Stdin)
	if err != nil {
		return nil, err
	}

	if err := setSize(pty, sz); err != nil {
		_ = pty.Close()
		return nil, err
	}

	c.Stdout, c.Stderr, c.Stdin = tty, tty, tty
	c.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}

	if err := c.Start(); err != nil {
		_ = pty.Close()
		return nil, err
	}
	return pty, nil
}

func setSize(f *os.File, ws *Winsize) error {
	return ioctl(f.Fd(), syscall.TIOCSWINSZ, uintptr(unsafe.Pointer(ws)))
}

func getSizeFull(f *os.File) (*Winsize, error) {
	var ws Winsize
	err := ioctl(f.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws)))
	if err != nil {
		return nil, err
	}
	return &ws, nil
}
