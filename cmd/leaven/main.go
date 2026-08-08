package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/lewtec/leaven"
)

func main() {
	packageName := flag.String("package", "main", "Go package name for generated code")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: leaven [flags] [input-file.ll]\n")
		fmt.Fprintf(os.Stderr, "With no file (or -), read LLVM IR from stdin and write Go to stdout.\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	}

	name, in, out, closer, err := openIO(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer closer()

	cmd := &leaven.Command{
		Package: *packageName,
		Name:    name,
		Input:   in,
		Output:  out,
	}
	if err := cmd.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
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
