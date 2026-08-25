//go:build unix

package libc

import (
	"errors"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func fcntlFile(f *os.File, cmd, arg int) int32 {
	conn, err := f.SyscallConn()
	if err != nil {
		setErrno(5) // EIO
		return -1
	}
	var r int
	var ferr error
	if err := conn.Control(func(fd uintptr) {
		r, ferr = unix.FcntlInt(fd, cmd, arg)
	}); err != nil {
		setErrno(5)
		return -1
	}
	if ferr != nil {
		var errno syscall.Errno
		if errors.As(ferr, &errno) {
			setErrno(int32(errno))
		} else {
			setErrno(5)
		}
		return -1
	}
	return int32(r)
}
