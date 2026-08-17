package libc

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"unsafe"
)

// errno for the process; __errno_location returns &errnoTLS.
var errnoTLS int32

// ErrnoLocation is __errno_location.
func ErrnoLocation() *int32 { return &errnoTLS }

func setErrno(e int32) { atomic.StoreInt32(&errnoTLS, e) }

// fd table: 0/1/2 = stdin/out/err; others are opened files.
var (
	fdMu   sync.Mutex
	fdTab  = map[int32]*os.File{0: os.Stdin, 1: os.Stdout, 2: os.Stderr}
	fdNext int32 = 3
)

func fdGet(fd int32) *os.File {
	fdMu.Lock()
	defer fdMu.Unlock()
	return fdTab[fd]
}

// Open is open/open64(path, flags[, mode]). Returns fd or -1.
func Open(path *byte, flags int32, mode ...int32) int32 {
	if path == nil {
		setErrno(14) // EFAULT
		return -1
	}
	name := GoString(path)
	// Ignore most flags; O_RDONLY=0, O_WRONLY=1, O_RDWR=2.
	how := os.O_RDONLY
	acc := flags & 3
	switch acc {
	case 1:
		how = os.O_WRONLY
	case 2:
		how = os.O_RDWR
	}
	if flags&64 != 0 { // O_CREAT
		how |= os.O_CREATE
	}
	if flags&512 != 0 { // O_TRUNC
		how |= os.O_TRUNC
	}
	if flags&1024 != 0 { // O_APPEND
		how |= os.O_APPEND
	}
	perm := os.FileMode(0o666)
	if len(mode) > 0 {
		perm = os.FileMode(mode[0] & 0o777)
	}
	f, err := os.OpenFile(name, how, perm)
	if err != nil {
		setErrno(2) // ENOENT-ish
		return -1
	}
	fdMu.Lock()
	fd := fdNext
	fdNext++
	fdTab[fd] = f
	fdMu.Unlock()
	return fd
}

// Open64 is open64 — same as Open.
func Open64(path *byte, flags int32, mode ...int32) int32 {
	return Open(path, flags, mode...)
}

// Fcntl is fcntl(2). The libc fd is looked up, then the host fcntl
// runs on that *os.File (unix.FcntlInt via SyscallConn).
func Fcntl(fd int32, cmd int32, args ...any) int32 {
	f := fdGet(fd)
	if f == nil {
		setErrno(9) // EBADF
		return -1
	}
	arg := 0
	if len(args) > 0 {
		arg = fcntlArg(args[0])
	}
	return fcntlFile(f, int(cmd), arg)
}

func fcntlArg(v any) int {
	switch x := v.(type) {
	case int:
		return x
	case int32:
		return int(x)
	case int64:
		return int(x)
	case uint32:
		return int(x)
	case uint64:
		return int(x)
	case uintptr:
		return int(x)
	case *byte:
		return int(uintptr(unsafe.Pointer(x)))
	case unsafe.Pointer:
		return int(uintptr(x))
	default:
		return 0
	}
}

// Close is close(fd).
func Close(fd int32) int32 {
	if fd < 3 {
		return 0
	}
	fdMu.Lock()
	f := fdTab[fd]
	delete(fdTab, fd)
	fdMu.Unlock()
	if f == nil {
		setErrno(9) // EBADF
		return -1
	}
	if err := f.Close(); err != nil {
		setErrno(5)
		return -1
	}
	return 0
}

// Read is read(fd, buf, n).
func Read(fd int32, buf *byte, n int64) int64 {
	f := fdGet(fd)
	if f == nil || buf == nil || n <= 0 {
		if f == nil {
			setErrno(9)
		}
		return -1
	}
	nr, err := f.Read(Bytes(buf, int(n)))
	if err != nil && nr == 0 {
		// EOF → 0
		return 0
	}
	return int64(nr)
}

// Write is write(fd, buf, n).
func Write(fd int32, buf *byte, n int64) int64 {
	f := fdGet(fd)
	if f == nil || buf == nil {
		if f == nil {
			setErrno(9)
		}
		return -1
	}
	if n <= 0 {
		return 0
	}
	nw, err := f.Write(Bytes(buf, int(n)))
	if err != nil {
		setErrno(5)
		return -1
	}
	return int64(nw)
}

// Sigaltstack is sigaltstack(2). Report existing stack (flags=0, not
// SS_DISABLE) so rustc stack_overflow::make_handler skips mmap.
// stack_t: { void *ss_sp; int ss_flags; size_t ss_size } @0,8,16.
func Sigaltstack(ss, oss *byte) int32 {
	if oss != nil {
		base := Ptr(oss)
		Store[*byte](base, 0, nil)
		Store[int32](base, 8, 0) // not SS_DISABLE
		Store[uint64](base, 16, 0)
	}
	_ = ss // install accepted
	return 0
}

// Mmap64 is mmap64. Anonymous maps use Malloc; MAP_FAILED is (void*)-1.
func Mmap64(addr *byte, length int64, prot, flags, fd int32, offset int64) *byte {
	_, _, _, _ = addr, prot, flags, offset
	if length <= 0 {
		setErrno(22)
		return (*byte)(unsafe.Pointer(^uintptr(0)))
	}
	// File maps: still allocate; content not filled (rare at startup).
	_ = fd
	p := Malloc[byte](length)
	if p == nil {
		setErrno(12) // ENOMEM
		return (*byte)(unsafe.Pointer(^uintptr(0)))
	}
	return p
}

// Mmap is mmap(2). Same allocator as Mmap64; returns unsafe.Pointer
// because Darwin rustc assigns the result to an LLVM ptr.
func Mmap(addr *byte, length int64, prot, flags, fd int32, offset int64) unsafe.Pointer {
	return unsafe.Pointer(Mmap64(addr, length, prot, flags, fd, offset))
}

// Munmap is munmap(2).
func Munmap(addr *byte, length int64) int32 {
	_ = length
	if addr != nil {
		Free(addr)
	}
	return 0
}

// Mprotect is mprotect(2). Always succeeds in this model.
func Mprotect(addr *byte, length int64, prot int32) int32 {
	_, _, _ = addr, length, prot
	return 0
}

// Getauxval is getauxval(3). AT_MINSIGSTKSZ=51, AT_PAGESIZE=6.
func Getauxval(typ int64) int64 {
	switch typ {
	case 6: // AT_PAGESIZE
		return 4096
	case 51: // AT_MINSIGSTKSZ
		return 2048
	default:
		return 0
	}
}

// Gettid is gettid(2). Single-threaded: same as pid 1 stand-in.
func Gettid() int32 { return 1 }

// Getpid is getpid(2).
func Getpid() int32 { return 1 }

// Getenv is getenv(3). Always nil (no process env for fixtures).
func Getenv(name *byte) *byte {
	_ = name
	return nil
}

// Getcwd is getcwd(3).
func Getcwd(buf *byte, size int64) *byte {
	if buf == nil || size <= 0 {
		return nil
	}
	wd, err := os.Getwd()
	if err != nil || int64(len(wd)+1) > size {
		setErrno(34) // ERANGE
		return nil
	}
	dst := Bytes(buf, int(size))
	copy(dst, wd)
	dst[len(wd)] = 0
	return buf
}

// PATH_MAX for realpath resolved buffer when resolved is non-nil.
const pathMax = 4096

// Dlsym is dlsym(3). Always nil so rust getrandom uses the file backend
// (calling a Go func through a C fnptr cast is not defined).
func Dlsym(handle, name *byte) *byte {
	_, _ = handle, name
	return nil
}

// Linux x86_64 syscall numbers used by rustc std / getrandom.
const (
	sysGettid = 186
	sysFutex  = 202
	sysStatx  = 332
)

// Syscall is the libc syscall(2) shim for a few numbers rustc needs.
// Single-threaded: futex wait/wake are no-ops that succeed.
func Syscall(nr int64, args ...interface{}) int64 {
	switch nr {
	case sysGettid:
		return int64(Gettid())
	case sysFutex:
		// futex(uaddr, op, val, timeout, uaddr2, val3) — ignore, single-thread.
		return 0
	case sysStatx:
		// statx(dirfd, path, flags, mask, statxbuf)
		var path *byte
		var buf *byte
		var dirfd, flags, mask int32
		if len(args) > 0 {
			dirfd = syscallI32(args[0])
		}
		if len(args) > 1 {
			path = syscallBytePtr(args[1])
		}
		if len(args) > 2 {
			flags = syscallI32(args[2])
		}
		if len(args) > 3 {
			mask = syscallI32(args[3])
		}
		if len(args) > 4 {
			buf = syscallBytePtr(args[4])
		}
		return int64(Statx(dirfd, path, flags, mask, buf))
	default:
		setErrno(38) // ENOSYS
		return -1
	}
}

func syscallI32(v interface{}) int32 {
	switch x := v.(type) {
	case int32:
		return x
	case int64:
		return int32(x)
	case int:
		return int32(x)
	default:
		return 0
	}
}

func syscallBytePtr(v interface{}) *byte {
	switch x := v.(type) {
	case *byte:
		return x
	case unsafe.Pointer:
		return (*byte)(x)
	case nil:
		return nil
	default:
		return nil
	}
}

// Realpath is realpath(3). resolved nil → malloc; else write into pathMax buf.
func Realpath(path *byte, resolved *byte) *byte {
	if path == nil {
		setErrno(14) // EFAULT
		return nil
	}
	abs, err := filepath.Abs(GoString(path))
	if err != nil {
		setErrno(2) // ENOENT-ish
		return nil
	}
	// Clean + require the path to exist (glibc realpath does).
	abs = filepath.Clean(abs)
	if _, err := os.Stat(abs); err != nil {
		setErrno(2)
		return nil
	}
	need := len(abs) + 1
	if resolved == nil {
		buf := Malloc[byte](int64(need))
		if buf == nil {
			setErrno(12) // ENOMEM
			return nil
		}
		dst := Bytes(buf, need)
		copy(dst, abs)
		dst[len(abs)] = 0
		return buf
	}
	if need > pathMax {
		setErrno(36) // ENAMETOOLONG
		return nil
	}
	dst := Bytes(resolved, pathMax)
	copy(dst, abs)
	dst[len(abs)] = 0
	return resolved
}

// Fstat64 is fstat64 — zero-fill and succeed (size fields left 0).
func Fstat64(fd int32, st *byte) int32 {
	_ = fd
	if st != nil {
		// struct size varies; clear a generous region.
		Memset(st, 0, 256)
	}
	return 0
}

// Stat64 is stat64 — same as fstat for our needs.
func Stat64(path *byte, st *byte) int32 {
	_ = path
	return Fstat64(0, st)
}

// Lseek64 is lseek64.
func Lseek64(fd int32, offset int64, whence int32) int64 {
	f := fdGet(fd)
	if f == nil {
		setErrno(9)
		return -1
	}
	// os.File has Seek
	n, err := f.Seek(offset, int(whence))
	if err != nil {
		setErrno(22)
		return -1
	}
	return n
}

// Getrandom is getrandom(2).
func Getrandom(buf *byte, buflen int64, flags int32) int64 {
	_ = flags
	if buf == nil || buflen <= 0 {
		return 0
	}
	dst := Bytes(buf, int(buflen))
	n, err := rand.Read(dst)
	if err != nil {
		setErrno(5)
		return -1
	}
	return int64(n)
}
