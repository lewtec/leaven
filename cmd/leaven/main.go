package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/lewtec/leaven"
	"github.com/lewtec/lewkit/x/cmd"
	"github.com/lewtec/lewkit/x/io/atomic"
	"github.com/lewtec/lewkit/x/release"
)

type args struct {
	help    cmd.Flag      `short:"h" long:"help" help:"show help"`
	version cmd.Flag      `long:"version" help:"print version"`
	Package cmd.StringArg `long:"package" short:"p" help:"Go package name for generated code" default:"main"`
	Input   cmd.StringArg `help:"LLVM IR file; omit or - for stdin"`
}

func (args) Description() string {
	return "Transpile LLVM IR to Go.\n\nWith no file (or -), read LLVM IR from stdin and write Go to stdout."
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	a, err := cmd.Parse[args](os.Args[1:]...)
	if err != nil {
		return err
	}
	return a.Run(ctx)
}

func (a *args) Run(ctx context.Context) error {
	switch {
	case a.help.Value():
		text, err := cmd.Usage[args](filepath.Base(os.Args[0]))
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(os.Stdout, text)
		return err
	case a.version.Value():
		_, err := fmt.Fprintln(os.Stdout, release.Version())
		return err
	}

	c := &leaven.Command{Package: a.Package.Value()}
	path := a.Input.Value()
	if path == "" || path == "-" {
		c.Name, c.Input, c.Output = "<stdin>", os.Stdin, os.Stdout
		return c.Run(ctx)
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	c.Name, c.Input = path, f
	return atomic.WriteFileFunction(strings.TrimSuffix(path, ".ll")+".go", func(w io.Writer) error {
		c.Output = w
		return c.Run(ctx)
	})
}
