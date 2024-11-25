package script

import (
	"io/fs"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

func PtyMasterOpen() (*os.File, *os.File, error) {

	pty, err := posix_openpt(os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	defer pty.Close()

	err = unlockpt(pty)
	if err != nil {
		return nil, nil, err
	}

	ptySlave, err := ptsname(pty)
	if err != nil {
		return nil, nil, err
	}
	tty, err := os.OpenFile(ptySlave, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}

	return pty, tty, nil
}

func posix_openpt(flags int, perm fs.FileMode) (*os.File, error) {
	return os.OpenFile("/dev/ptmx", flags, perm)
}

func ptsname(f *os.File) (string, error) {
	var n int
	err := ioctl(f.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n)))
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(int(n)), nil
}

func unlockpt(f *os.File) error {
	var u int
	return ioctl(f.Fd(), syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&u))) 
}

func ioctl(fd, cmd, ptr uintptr) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, ptr)
	if err != 0 {
		return err
	}
	return nil
}
