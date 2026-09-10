package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/lewtec/leaven"
	"github.com/lewtec/lewkit/x/cmd"
)

type args struct {
	Package cmd.StringArg `long:"package" short:"p" help:"Go package name for generated code" default:"main"`
	Input   cmd.StringArg `help:"LLVM IR file; omit or - for stdin"`
}

func (args) Description() string {
	return "Transpile LLVM IR to Go.\n\nWith no file (or -), read LLVM IR from stdin and write Go to stdout."
}

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	app, err := cmd.Parse[cmd.App[args]](os.Args[1:]...)
	if err != nil {
		return err
	}
	return app.Run(ctx)
}

func (a *args) Run(ctx context.Context) error {
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
	out, err := os.Create(strings.TrimSuffix(path, ".ll") + ".go")
	if err != nil {
		return err
	}
	defer out.Close()
	c.Name, c.Input, c.Output = path, f, out
	return c.Run(ctx)
}
