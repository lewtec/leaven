package leaven

import (
	"context"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// Command transpiles LLVM IR to Go. Fill Input (and usually Output), then call Run.
type Command struct {
	// Package is the generated Go package name. Empty means "main".
	Package string
	// Name is an optional source label for parse errors (file path or "<stdin>").
	Name string
	// Input is LLVM IR assembly.
	Input io.Reader
	// Output receives generated Go source.
	Output io.Writer
}

// Run reads LLVM IR from Input and writes Go to Output.
// ctx is checked before parse and before compile.
func (c *Command) Run(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	pkg := c.Package
	if pkg == "" {
		pkg = "main"
	}
	if err := validatePackageName(pkg); err != nil {
		return err
	}
	if c.Input == nil {
		return errInputRequired
	}
	if c.Output == nil {
		return errOutputRequired
	}

	m, err := parseIR(c.Name, c.Input)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return Compile(c.Output, m, pkg)
}

// validatePackageName checks that name is a legal Go package identifier.
func validatePackageName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: empty package name", errInvalidPackage)
	}
	r, size := utf8.DecodeRuneInString(name)
	if r == utf8.RuneError && size == 1 {
		return fmt.Errorf("%w: %q", errInvalidPackage, name)
	}
	if !unicode.IsLetter(r) && r != '_' {
		return fmt.Errorf("%w: %q", errInvalidPackage, name)
	}
	for _, r := range name[size:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return fmt.Errorf("%w: %q", errInvalidPackage, name)
		}
	}
	return nil
}
