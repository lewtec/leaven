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
	"github.com/lewtec/lewkit/x/profile"
	"github.com/lewtec/lewkit/x/release"
)

type args struct {
	verbose    cmd.Count     `short:"v" long:"verbose" help:"log verbosity" default:"0"`
	profileDir cmd.StringArg `long:"profile-dir" help:"write pprof profiles here"`
	help       cmd.Flag      `short:"h" long:"help" help:"show help"`
	version    cmd.Flag      `long:"version" help:"print version"`
	Package    cmd.StringArg `long:"package" short:"p" help:"Go package name for generated code" default:"main"`
	Input      cmd.StringArg `help:"LLVM IR file; omit or - for stdin"`
}

func (args) Description() string {
	return "Transpile LLVM IR to Go.\n\nWith no file (or -), read LLVM IR from stdin and write Go to stdout."
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, argv []string) error {
	a, err := cmd.Parse[args](argv...)
	if err != nil {
		return err
	}
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
	if err := a.setup(ctx); err != nil {
		return err
	}
	return a.Run(ctx)
}

func (a args) setup(ctx context.Context) error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo - slog.Level(4*a.verbose.Value()),
	})))
	if dir := a.profileDir.Value(); dir == "" {
		return nil
	}
	p := profile.NewProfile(a.profileDir.Value())
	go func() {
		if err := p.Run(ctx); err != nil {
			slog.Error(err.Error())
		}
	}()
	return nil
}

func (a *args) Run(ctx context.Context) error {
	name, in, out, closer, err := openIO(a.Input.Value())
	if err != nil {
		return err
	}
	defer closer()
	c := &leaven.Command{
		Package: a.Package.Value(),
		Name:    name,
		Input:   in,
		Output:  out,
	}
	return c.Run(ctx)
}

// openIO uses stdin/stdout when path is empty or "-"; otherwise path → path with .go suffix.
func openIO(path string) (name string, in io.Reader, out io.Writer, closer func(), err error) {
	if path == "" || path == "-" {
		return "<stdin>", os.Stdin, os.Stdout, func() {}, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return "", nil, nil, nil, err
	}
	outPath := strings.TrimSuffix(path, ".ll") + ".go"
	of, err := os.Create(outPath)
	if err != nil {
		f.Close()
		return "", nil, nil, nil, err
	}
	return path, f, of, func() {
		f.Close()
		of.Close()
	}, nil
}
