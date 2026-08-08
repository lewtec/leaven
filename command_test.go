package leaven

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
)

func TestCommandRun(t *testing.T) {
	var out bytes.Buffer
	c := &Command{Name: "anon_dot.ll", Input: bytes.NewReader(testdataIR(t)), Output: &out}
	if err := c.Run(t.Context()); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(out.Bytes(), []byte("package main")) {
		t.Fatalf("generated source missing package clause:\n%s", out.Bytes())
	}
}

func TestCommandRunCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	c := &Command{Input: bytes.NewReader(nil), Output: io.Discard}
	if err := c.Run(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}
}

func TestCommandRunInvalidPackage(t *testing.T) {
	c := &Command{Package: "123bad", Input: bytes.NewReader(nil), Output: io.Discard}
	if err := c.Run(t.Context()); !errors.Is(err, errInvalidPackage) {
		t.Fatalf("Run() error = %v, want errInvalidPackage", err)
	}
}

func TestCommandRunRequiresInput(t *testing.T) {
	c := &Command{Output: io.Discard}
	if err := c.Run(t.Context()); !errors.Is(err, errInputRequired) {
		t.Fatalf("Run() error = %v, want errInputRequired", err)
	}
}

func TestCommandRunRequiresOutput(t *testing.T) {
	c := &Command{Input: bytes.NewReader(nil)}
	if err := c.Run(t.Context()); !errors.Is(err, errOutputRequired) {
		t.Fatalf("Run() error = %v, want errOutputRequired", err)
	}
}
