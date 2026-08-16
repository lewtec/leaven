package leaven

import "testing"

func TestLldNames(t *testing.T) {
	got := lldNames("linux")
	if len(got) != 1 || got[0] != "ld.lld" {
		t.Fatalf("linux: %v", got)
	}
	got = lldNames("darwin")
	if len(got) != 1 || got[0] != "ld64.lld" {
		t.Fatalf("darwin: %v", got)
	}
}

func TestLinuxCondaTriple(t *testing.T) {
	if got := linuxCondaTriple("arm64"); got != "aarch64-conda-linux-gnu" {
		t.Fatalf("arm64: %s", got)
	}
	if got := linuxCondaTriple("amd64"); got != "x86_64-conda-linux-gnu" {
		t.Fatalf("amd64: %s", got)
	}
}
