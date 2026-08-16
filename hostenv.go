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
	if goos == "darwin" {
		return []string{"ld64.lld", "lld"}
	}
	return []string{"lld", "ld.lld"}
}

func lldPath() string {
	for _, n := range lldNames(runtime.GOOS) {
		if p, err := exec.LookPath(n); err == nil {
			return p
		}
	}
	return ""
}

// clangNativeFlags is for linking an executable. Prefer lld so we do
// not hit Apple ld's libLTO.dylib basename check.
func clangNativeFlags() []string {
	flags := clangSysrootFlags()
	if ld := lldPath(); ld != "" {
		flags = append(flags, "-fuse-ld="+ld)
	}
	return flags
}
