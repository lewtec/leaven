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
