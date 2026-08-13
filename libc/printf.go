package libc

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

// asUnsigned reinterprets a signed integer arg as unsigned for C-style
// %x/%X/%o/%u formatting (Go's fmt would otherwise print a leading '-').
func asUnsigned(a any) any {
	switch v := a.(type) {
	case int8:
		return uint8(v)
	case int16:
		return uint16(v)
	case int32:
		return uint32(v)
	case int64:
		return uint64(v)
	case int:
		return uint(v)
	default:
		return a
	}
}

// parseCFormatSpec reads a C %…verb starting at format[pct] ('%').
// verb=='%' with empty flags is a literal %%. ok is false if the string
// ends mid-spec.
func parseCFormatSpec(format []byte, pct int) (flags string, verb byte, after int, ok bool) {
	i := pct + 1
	if i < len(format) && format[i] == '%' {
		return "", '%', i + 1, true
	}
	for i < len(format) {
		switch format[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '#', '-', '.', ',', ' ', 'h', 'j', 'l':
			i++
			continue
		}
		break
	}
	if i >= len(format) {
		return "", 0, i, false
	}
	flags = string(format[pct:i])
	verb = format[i]
	flags = strings.ReplaceAll(flags, "h", "")
	flags = strings.ReplaceAll(flags, "j", "")
	flags = strings.ReplaceAll(flags, "l", "")
	if j := strings.Index(flags, "#0"); j >= 0 && verb == 'x' {
		k := j + 2
		for k < len(flags) && '0' <= flags[k] && flags[k] <= '9' {
			k++
		}
		n, err := strconv.Atoi(flags[j+2 : k])
		if err != nil {
			n = 0
		}
		flags = flags[:j+2] + fmt.Sprint(n-2) + flags[k:]
	}
	return flags, verb, i + 1, true
}

// fixPrintfFormat converts a printf format string from C-style to Go-style,
// and makes needed changes to the other arguments as well.
//
// It is based on the function of the same name in github.com/andybalholm/c2go.
func fixPrintfFormat(f *byte, args []any) string {
	format := Bytes(f, int(Strlen(f)))
	var buf strings.Builder
	narg := 0
	start := 0
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		buf.Write(format[start:i])
		flags, verb, after, ok := parseCFormatSpec(format, i)
		if !ok {
			return string(format)
		}
		if verb == '%' && flags == "" {
			buf.WriteString("%%")
			start = after
			i = after - 1
			continue
		}
		start = after
		i = after - 1

		switch verb {
		default:
			buf.WriteString("%")
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 's':
			if narg < len(args) {
				switch a := args[narg].(type) {
				case *byte:
					args[narg] = Bytes(a, int(Strlen(a)))
				}
			}
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 'f', 'e', 'g', 'c', 'p':
			// usual meanings
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 'x', 'X', 'o', 'd', 'b', 'i':
			if verb == 'i' {
				verb = 'd'
			}
			// C %x/%X/%o print unsigned; Go's fmt signs int32/int64. Match C.
			if (verb == 'x' || verb == 'X' || verb == 'o') && narg < len(args) {
				args[narg] = asUnsigned(args[narg])
			}
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 'u':
			buf.WriteString(flags)
			buf.WriteString("d")
			if narg < len(args) {
				args[narg] = asUnsigned(args[narg])
			}
		}

		narg++
	}
	buf.Write(format[start:])

	return buf.String()
}

func printfTo(w io.Writer, format *byte, args []any) int32 {
	f := fixPrintfFormat(format, args)
	n, err := fmt.Fprintf(w, f, args...)
	if err != nil {
		return -1
	}
	return int32(n)
}

func Printf(format *byte, args ...any) int32 {
	return printfTo(os.Stdout, format, args)
}

func Puts(s *byte) int32 {
	n, err := fmt.Printf("%s\n", Bytes(s, int(Strlen(s))))
	if err != nil {
		return -1
	}
	return int32(n)
}

// fixScanfFormat converts a scanf format string from C-style to Go-style,
// and makes needed changes to the other arguments as well.
func fixScanfFormat(f *byte, args []any) string {
	format := Bytes(f, int(Strlen(f)))
	var buf strings.Builder
	narg := 0
	start := 0
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		buf.Write(format[start:i])
		flags, verb, after, ok := parseCFormatSpec(format, i)
		if !ok {
			return string(format)
		}
		if verb == '%' && flags == "" {
			buf.WriteString("%%")
			start = after
			i = after - 1
			continue
		}
		start = after
		i = after - 1

		switch verb {
		default:
			buf.WriteString("%")
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 's':
			if narg < len(args) {
				switch a := args[narg].(type) {
				case *byte:
					args[narg] = &stringScanner{a}
				}
			}
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 'f', 'e', 'g', 'c', 'p':
			// usual meanings
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 'x', 'X', 'o', 'd', 'b', 'i':
			if verb == 'i' {
				verb = 'v'
			}
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 'u':
			buf.WriteString(flags)
			buf.WriteString("d")
			if narg >= len(args) {
				break
			}
			switch a := args[narg].(type) {
			case *int16:
				args[narg] = (*uint16)(unsafe.Pointer(a))
			case *int32:
				args[narg] = (*uint32)(unsafe.Pointer(a))
			case *int64:
				args[narg] = (*uint64)(unsafe.Pointer(a))
			}
		}

		narg++
	}
	buf.Write(format[start:])

	return buf.String()
}

type stringScanner struct {
	b *byte
}

func (s *stringScanner) Scan(state fmt.ScanState, verb rune) error {
	state.SkipSpace()
	var result []rune
	for {
		c, _, err := state.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if unicode.IsSpace(c) {
			state.UnreadRune()
			break
		}
		result = append(result, c)
	}
	str := string(result)
	b := Bytes(s.b, len(str)+1)
	copy(b, str)
	b[len(b)-1] = 0
	return nil
}

func Scanf(format *byte, args ...any) int32 {
	f := fixScanfFormat(format, args)
	n, err := fmt.Scanf(f, args...)
	if err != nil && n == 0 {
		return -1
	}
	return int32(n)
}

// Sscanf is C sscanf / glibc __isoc23_sscanf. Dest pointers from
// opaque-ptr IR arrive as *byte; write the scanned value through them.
func Sscanf(str *byte, format *byte, args ...any) int32 {
	if str == nil || format == nil {
		return -1
	}
	in := Bytes(str, int(Strlen(str)))
	pat := Bytes(format, int(Strlen(format)))
	n, eof := csscanf(in, pat, args)
	if n == 0 && eof {
		return -1
	}
	return int32(n)
}

func csscanf(in, pat []byte, args []any) (n int, eof bool) {
	si, ai := 0, 0
	for pi := 0; pi < len(pat); {
		if pat[pi] != '%' {
			if isCSpace(pat[pi]) {
				for pi < len(pat) && isCSpace(pat[pi]) {
					pi++
				}
				si = skipCSpace(in, si)
				continue
			}
			if si >= len(in) {
				return n, true
			}
			if in[si] != pat[pi] {
				return n, false
			}
			si++
			pi++
			continue
		}
		pi++
		if pi < len(pat) && pat[pi] == '%' {
			if si >= len(in) || in[si] != '%' {
				return n, si >= len(in)
			}
			si++
			pi++
			continue
		}
		suppress := false
		if pi < len(pat) && pat[pi] == '*' {
			suppress = true
			pi++
		}
		width := 0
		for pi < len(pat) && pat[pi] >= '0' && pat[pi] <= '9' {
			width = width*10 + int(pat[pi]-'0')
			pi++
		}
		long := 0
		for pi < len(pat) && (pat[pi] == 'l' || pat[pi] == 'h' || pat[pi] == 'j' || pat[pi] == 'z' || pat[pi] == 't') {
			if pat[pi] == 'l' {
				long++
			}
			pi++
		}
		if pi >= len(pat) {
			return n, false
		}
		verb := pat[pi]
		pi++
		switch verb {
		case 'd', 'i', 'u', 'x', 'X', 'o':
			si = skipCSpace(in, si)
			if si >= len(in) {
				return n, true
			}
			base := 10
			switch verb {
			case 'x', 'X':
				base = 16
			case 'o':
				base = 8
			case 'i':
				base = 0
			}
			chunk := in[si:]
			if width > 0 && width < len(chunk) {
				chunk = chunk[:width]
			}
			v, adv, ok := parseCInt(chunk, base)
			if !ok {
				return n, false
			}
			si += adv
			if !suppress {
				if ai >= len(args) {
					return n, false
				}
				storeScanInt(args[ai], v, intSize(long))
				ai++
				n++
			}
		case 's':
			si = skipCSpace(in, si)
			if si >= len(in) {
				return n, true
			}
			start := si
			lim := len(in)
			if width > 0 && si+width < lim {
				lim = si + width
			}
			for si < lim && !isCSpace(in[si]) {
				si++
			}
			if !suppress {
				if ai >= len(args) {
					return n, false
				}
				dst := destData(args[ai])
				if dst != nil {
					out := Bytes(As[byte](dst), si-start+1)
					copy(out, in[start:si])
					out[si-start] = 0
				}
				ai++
				n++
			}
		case 'c':
			if si >= len(in) {
				return n, true
			}
			cnt := width
			if cnt == 0 {
				cnt = 1
			}
			if si+cnt > len(in) {
				return n, true
			}
			if !suppress {
				if ai >= len(args) {
					return n, false
				}
				dst := destData(args[ai])
				if dst != nil {
					copy(Bytes(As[byte](dst), cnt), in[si:si+cnt])
				}
				ai++
				n++
			}
			si += cnt
		default:
			return n, false
		}
	}
	return n, false
}

func isCSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f'
}

func skipCSpace(in []byte, i int) int {
	for i < len(in) && isCSpace(in[i]) {
		i++
	}
	return i
}

func intSize(long int) int {
	if long >= 2 {
		return 8
	}
	if long == 1 {
		return 8 // LP64 unsigned long
	}
	return 4
}

func parseCInt(in []byte, base int) (uint64, int, bool) {
	if len(in) == 0 {
		return 0, 0, false
	}
	i := 0
	neg := false
	if in[0] == '+' || in[0] == '-' {
		neg = in[0] == '-'
		i++
	}
	if i >= len(in) {
		return 0, 0, false
	}
	if base == 0 {
		if in[i] == '0' && i+1 < len(in) && (in[i+1] == 'x' || in[i+1] == 'X') {
			base = 16
			i += 2
		} else if in[i] == '0' {
			base = 8
		} else {
			base = 10
		}
	} else if base == 16 && in[i] == '0' && i+1 < len(in) && (in[i+1] == 'x' || in[i+1] == 'X') {
		i += 2
	}
	if i >= len(in) {
		return 0, 0, false
	}
	var v uint64
	start := i
	for i < len(in) {
		d := hexDigit(in[i])
		if d < 0 || d >= base {
			break
		}
		v = v*uint64(base) + uint64(d)
		i++
	}
	if i == start {
		return 0, 0, false
	}
	if neg {
		v = uint64(-int64(v))
	}
	return v, i, true
}

func hexDigit(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	default:
		return -1
	}
}

func storeScanInt(a any, v uint64, nbytes int) {
	p := destData(a)
	if p == nil {
		return
	}
	switch nbytes {
	case 1:
		*(*uint8)(p) = uint8(v)
	case 2:
		*(*uint16)(p) = uint16(v)
	case 4:
		*(*uint32)(p) = uint32(v)
	default:
		*(*uint64)(p) = v
	}
}

func destData(a any) unsafe.Pointer {
	if a == nil {
		return nil
	}
	switch v := a.(type) {
	case unsafe.Pointer:
		return v
	case *byte:
		return unsafe.Pointer(v)
	case *uint64:
		return unsafe.Pointer(v)
	case *int64:
		return unsafe.Pointer(v)
	case *uint32:
		return unsafe.Pointer(v)
	case *int32:
		return unsafe.Pointer(v)
	case *uint:
		return unsafe.Pointer(v)
	case *int:
		return unsafe.Pointer(v)
	default:
		type eface struct{ typ, data unsafe.Pointer }
		return (*eface)(unsafe.Pointer(&a)).data
	}
}
