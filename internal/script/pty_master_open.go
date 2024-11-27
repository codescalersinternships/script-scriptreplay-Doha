package script

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

func PtyMasterOpen() (*os.File, *os.File, error) {
	pty, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}

	if err := unlockpt(pty); err != nil {
		return nil, nil, err
	}

	ptsName, err := ptsname(pty)
	if err != nil {
		return nil, nil, err
	}

	tty, err := os.OpenFile(ptsName, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	
	return pty, tty, nil
}

func unlockpt(f *os.File) error {
	var u int
	return ioctl(f.Fd(), syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&u)))
}

func ptsname(f *os.File) (string, error) {
	var n int
	err := ioctl(f.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n)))
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(n), nil
}

func ioctl(fd, cmd, ptr uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, ptr)
	if errno != 0 {
		return errno
	}
	return nil
}
