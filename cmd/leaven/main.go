package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/lewtec/leaven"
	"github.com/lewtec/lewkit/x/cmd"
)

type args struct {
	Package cmd.StringArg `long:"package" short:"p" help:"Go package name for generated code" default:"main"`
	Input   cmd.StringArg `long:"input" short:"i" help:"LLVM IR file; omit or - for stdin"`
	file    string
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
	app, input, err := parseCLI(argv)
	if err != nil {
		return err
	}
	app.Args.file = input
	return app.Run(ctx)
}

func parseCLI(argv []string) (cmd.App[args], string, error) {
	app, err := cmd.Parse[cmd.App[args]](argv...)
	if err == nil {
		return app, "", nil
	}
	if !errors.Is(err, cmd.ErrUnknownCommand) {
		return app, "", err
	}
	name := unknownCommandName(err)
	if name == "" {
		return app, "", err
	}
	app, err = cmd.Parse[cmd.App[args]](dropToken(argv, name)...)
	if err != nil {
		return app, "", err
	}
	return app, name, nil
}

func unknownCommandName(err error) string {
	s, ok := strings.CutPrefix(err.Error(), "unknown command: ")
	if !ok {
		return ""
	}
	return s
}

func dropToken(argv []string, tok string) []string {
	for i, a := range argv {
		if a == tok {
			out := make([]string, 0, len(argv)-1)
			out = append(out, argv[:i]...)
			return append(out, argv[i+1:]...)
		}
	}
	return argv
}

func (a *args) path() string {
	if v := a.Input.Value(); v != "" {
		return v
	}
	return a.file
}

func (a *args) Run(ctx context.Context) error {
	name, in, out, closer, err := openIO(a.path())
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
