package leaven

import "testing"

func TestLinuxCondaTriple(t *testing.T) {
	if got := linuxCondaTriple("arm64"); got != "aarch64-conda-linux-gnu" {
		t.Fatalf("arm64: %s", got)
	}
	if got := linuxCondaTriple("amd64"); got != "x86_64-conda-linux-gnu" {
		t.Fatalf("amd64: %s", got)
	}
}
