package leaven

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// linuxCondaTriple is the conda-forge sysroot directory name for goarch.
func linuxCondaTriple(goarch string) string {
	if goarch == "arm64" {
		return "aarch64-conda-linux-gnu"
	}
	return "x86_64-conda-linux-gnu"
}

// condaLinuxSysroot is prefix/<triple>/sysroot on Linux, empty elsewhere.
func condaLinuxSysroot(prefix string) string {
	if runtime.GOOS != "linux" {
		return ""
	}
	return filepath.Join(prefix, linuxCondaTriple(runtime.GOARCH), "sysroot")
}

func darwinSDK() string {
	if runtime.GOOS != "darwin" {
		return ""
	}
	if v := os.Getenv("SDKROOT"); v != "" {
		return v
	}
	out, err := exec.Command("xcrun", "--show-sdk-path").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func clangSysrootFlags() []string {
	if sdk := darwinSDK(); sdk != "" {
		return []string{"-isysroot", sdk}
	}
	return nil
}

func lldNames(goos string) []string {
	// Bare `lld` is a dispatcher. The real ELF linker is ld.lld.
	if goos == "darwin" {
		return []string{"ld64.lld"}
	}
	return []string{"ld.lld"}
}

func lldPath() string {
	for _, n := range lldNames(runtime.GOOS) {
		if p, err := exec.LookPath(n); err == nil {
			return p
		}
	}
	return ""
}

// clangNativeFlags is for linking an executable. Darwin uses ld64.lld
// so Apple ld never sees conda's -lto_library path. Linux keeps GNU ld
// (lld 22 + clang 14 sysroot looks for /lib64/libc.so.6 and misses).
func clangNativeFlags() []string {
	flags := clangSysrootFlags()
	if runtime.GOOS == "darwin" {
		if ld := lldPath(); ld != "" {
			flags = append(flags, "-fuse-ld="+ld)
		}
	}
	return flags
}
