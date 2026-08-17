//go:build windows

package libc

import "os"

func fcntlFile(f *os.File, cmd, arg int) int32 {
	_, _, _ = f, cmd, arg
	setErrno(38) // ENOSYS
	return -1
}
