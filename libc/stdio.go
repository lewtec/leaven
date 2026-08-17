package libc

import (
	"fmt"
	"os"
	"unicode"
	"unsafe"
)

func Getchar() int32 {
	var buf [1]byte
	_, err := os.Stdin.Read(buf[:])
	if err != nil {
		return -1
	}
	return int32(buf[0])
}

// Poll is poll(2). pollfd is {i32 fd; i16 events; i16 revents} (8 bytes).
// Marks every slot ready (revents=events) so rhai-run stdin checks proceed.
func Poll(fds *byte, nfds int64, timeout int32) int32 {
	_ = timeout
	if fds == nil || nfds <= 0 {
		return 0
	}
	const pollfdSize = 8
	const eventsOff = 4
	const reventsOff = 6
	base := Ptr(fds)
	for i := int64(0); i < nfds; i++ {
		p := Off(base, int(i*pollfdSize))
		ev := Load[int16](p, eventsOff)
		Store(p, reventsOff, ev)
	}
	return int32(nfds)
}

// Signal is signal(2). Returns previous disposition (pretend SIG_DFL=0).
// rhai-run ignores SIGPIPE (13) with SIG_IGN (1).
func Signal(sig int32, handler int64) int64 {
	_, _ = sig, handler
	return 0
}

// Sysconf is sysconf(3). Linux _SC_PAGESIZE is 30.
func Sysconf(name int32) int64 {
	switch name {
	case 30: // _SC_PAGESIZE
		return 4096
	default:
		return -1
	}
}

// PthreadSelf is pthread_self. Single-threaded: fixed non-zero id.
func PthreadSelf() int64 { return 1 }

// PthreadGetattrNp fills a dummy attr (stack base/size via getstack).
func PthreadGetattrNp(thread int64, attr *byte) int32 {
	_, _ = thread, attr
	return 0
}

// PthreadAttrGetstack reports a synthetic stack for stack-overflow guards.
func PthreadAttrGetstack(attr, stackaddr *byte, stacksize *byte) int32 {
	_ = attr
	if stackaddr != nil {
		Store(Ptr(stackaddr), 0, dummyStackAddr())
	}
	if stacksize != nil {
		Store[uint64](Ptr(stacksize), 0, dummyStackSize)
	}
	return 0
}

// PthreadAttrDestroy is a no-op for the dummy attr.
func PthreadAttrDestroy(attr *byte) int32 {
	_ = attr
	return 0
}

const dummyStackSize = 8 << 20 // 8 MiB

var dummyStack [dummyStackSize]byte

func dummyStackAddr() unsafe.Pointer { return unsafe.Pointer(&dummyStack[0]) }

// PthreadGetStackaddrNp is Darwin pthread_get_stackaddr_np.
func PthreadGetStackaddrNp(thread int64) unsafe.Pointer {
	_ = thread
	return dummyStackAddr()
}

// PthreadGetStacksizeNp is Darwin pthread_get_stacksize_np.
func PthreadGetStacksizeNp(thread int64) int64 {
	_ = thread
	return dummyStackSize
}

// Single-threaded: mutex/attr calls succeed and do nothing.
func PthreadMutexattrInit(attr *byte) int32    { _ = attr; return 0 }
func PthreadMutexattrDestroy(attr *byte) int32 { _ = attr; return 0 }
func PthreadMutexattrSettype(attr *byte, typ int32) int32 {
	_, _ = attr, typ
	return 0
}
func PthreadMutexInit(m, attr *byte) int32 { _, _ = m, attr; return 0 }
func PthreadMutexDestroy(m *byte) int32    { _ = m; return 0 }
func PthreadMutexLock(m *byte) int32       { _ = m; return 0 }
func PthreadMutexUnlock(m *byte) int32     { _ = m; return 0 }
func PthreadMutexTrylock(m *byte) int32    { _ = m; return 0 }

// Sigaction is sigaction(2). No-op success for rhai-run startup.
func Sigaction(signum int32, act, oldact *byte) int32 {
	_, _, _ = signum, act, oldact
	return 0
}

func Putc(c int32, stream *os.File) int32 {
	_, err := stream.Write([]byte{byte(c)})
	if err != nil {
		return -1
	}
	return c
}

func Putchar(c int32) int32 {
	return Putc(c, os.Stdout)
}

// Fprintf is C fprintf(stream, format, ...).
func Fprintf(stream *os.File, format *byte, args ...any) int32 {
	return printfTo(stream, format, args)
}

// Fputs is C fputs(s, stream).
func Fputs(s *byte, stream *os.File) int32 {
	if stream == nil {
		return -1
	}
	_, err := stream.Write(Bytes(s, int(Strlen(s))))
	if err != nil {
		return -1
	}
	return 0
}

// Fputc is C fputc(c, stream).
func Fputc(c int32, stream *os.File) int32 {
	return Putc(c, stream)
}

func cbool(ok bool) int32 {
	if ok {
		return 1
	}
	return 0
}

// Iswspace is C iswspace(c) from <wctype.h>.
func Iswspace(c int32) int32 {
	return cbool(unicode.IsSpace(rune(uint32(c))))
}

// Iswblank is C iswblank(c) from <wctype.h>.
func Iswblank(c int32) int32 {
	r := rune(uint32(c))
	if r == '\t' || unicode.Is(unicode.Zs, r) {
		return 1
	}
	return 0
}

// Iswalnum is C iswalnum(c) from <wctype.h>.
func Iswalnum(c int32) int32 {
	r := rune(uint32(c))
	return cbool(unicode.IsLetter(r) || unicode.IsDigit(r))
}

// Iswalpha is C iswalpha(c) from <wctype.h>.
func Iswalpha(c int32) int32 {
	return cbool(unicode.IsLetter(rune(uint32(c))))
}

// Iswdigit is C iswdigit(c) from <wctype.h>.
func Iswdigit(c int32) int32 {
	return cbool(unicode.IsDigit(rune(uint32(c))))
}

// Iswlower is C iswlower(c) from <wctype.h>.
func Iswlower(c int32) int32 {
	return cbool(unicode.IsLower(rune(uint32(c))))
}

// Iswupper is C iswupper(c) from <wctype.h>.
func Iswupper(c int32) int32 {
	return cbool(unicode.IsUpper(rune(uint32(c))))
}

// Iswcntrl is C iswcntrl(c) from <wctype.h>.
func Iswcntrl(c int32) int32 {
	return cbool(unicode.IsControl(rune(uint32(c))))
}

// Iswxdigit is C iswxdigit(c) from <wctype.h>.
func Iswxdigit(c int32) int32 {
	r := rune(uint32(c))
	if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') {
		return 1
	}
	return 0
}

// Towupper is C towupper(c) from <wctype.h>.
func Towupper(c int32) int32 {
	return int32(unicode.ToUpper(rune(uint32(c))))
}

// Towlower is C towlower(c) from <wctype.h>.
func Towlower(c int32) int32 {
	return int32(unicode.ToLower(rune(uint32(c))))
}

// Exit is C exit(status). Runs __cxa_atexit handlers first.
func Exit(status int32) {
	runCxaAtExit()
	os.Exit(int(status))
}

// Dup is C dup(fd). Pure-Go: only maps 0/1/2 to themselves; other FDs
// return -1 (tree-sitter uses this for print-dot-graph on std FDs).
func Dup(fd int32) int32 {
	switch fd {
	case 0, 1, 2:
		return fd
	}
	return -1
}

// glibc _ISprint bit on little-endian (see bits/ctype.h _ISbit(6)).
const ctypeBPrint int16 = 16384

// ctypeB is a minimal glibc-style __ctype_b table for ASCII.
// Indexed 0..127; only the print flag is populated (tree-sitter core needs it).
var ctypeB [128]int16

var (
	ctypeBPtr *int16
	ctypeBLoc **int16
)

func init() {
	for i := 1; i < 128; i++ {
		if unicode.IsPrint(rune(i)) {
			ctypeB[i] = ctypeBPrint
		}
	}
	ctypeBPtr = &ctypeB[0]
	ctypeBLoc = &ctypeBPtr
}

// CtypeBLoc is glibc __ctype_b_loc(): returns a pointer to the ctype bitmask table.
func CtypeBLoc() **int16 {
	return ctypeBLoc
}

// Snprintf is C snprintf(buf, n, format, ...).
func Snprintf(buf *byte, n int64, format *byte, args ...any) int32 {
	f := fixPrintfFormat(format, args)
	s := fmt.Sprintf(f, args...)
	if n > 0 && buf != nil {
		dst := Bytes(buf, int(n))
		copyLen := len(s)
		if int64(copyLen) >= n {
			copyLen = int(n - 1)
			if copyLen < 0 {
				copyLen = 0
			}
		}
		copy(dst[:copyLen], s)
		dst[copyLen] = 0
	}
	// C returns the number of chars that would have been written (excl. NUL).
	return int32(len(s))
}

// Vsnprintf is C vsnprintf. ap points at system va_list storage
// (__va_list_tag); llvm.va_start stores &varargs at overflow_arg_area (offset 8).
func Vsnprintf(buf *byte, n int64, format *byte, ap *byte) int32 {
	var args []any
	if ap != nil {
		vlPtr := Load[*[]interface{}](Ptr(ap), 8)
		if vlPtr != nil {
			args = append(args, (*vlPtr)...)
		}
	}
	return Snprintf(buf, n, format, args...)
}

// Fdopen is C fdopen(fd, mode). Maps 0/1/2 to stdin/stdout/stderr;
// other FDs use os.NewFile. Mode is accepted but not fully applied.
func Fdopen(fd int32, mode *byte) *os.File {
	_ = mode
	switch fd {
	case 0:
		return os.Stdin
	case 1:
		return os.Stdout
	case 2:
		return os.Stderr
	}
	return os.NewFile(uintptr(fd), fmt.Sprintf("fd%d", fd))
}

// Fclose is C fclose. Does not close stdin/stdout/stderr.
func Fclose(f *os.File) int32 {
	if f == nil {
		return -1
	}
	if f == os.Stdin || f == os.Stdout || f == os.Stderr {
		return 0
	}
	if err := f.Close(); err != nil {
		return -1
	}
	return 0
}

// Abort is C abort(): no-return process abort.
func Abort() {
	panic("abort")
}
