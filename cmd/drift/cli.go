package main

import (
	"fmt"
	"io"
	"os"
)

type Verbosity int

const (
	SilentLevel Verbosity = iota
	InfoLevel
	DebugLevel
)

type CLI struct {
	stdout io.Writer
	stderr io.Writer

	verbosity Verbosity
}

func (cli *CLI) SetVerbosity(v Verbosity) {
	cli.verbosity = v
}

func (cli CLI) fwritef(w io.Writer, level Verbosity, format string, args ...interface{}) (n int, err error) {
	if cli.verbosity < level {
		return
	}
	return fmt.Fprintf(w, format+"\n", args...)
}

func (cli CLI) Exitf(code int, format string, args ...interface{}) {
	cli.fwritef(cli.stderr, InfoLevel, format, args...)
	os.Exit(code)
}

func (cli CLI) Infof(format string, args ...interface{}) (n int, err error) {
	return cli.fwritef(cli.stderr, InfoLevel, format, args...)
}

func (cli CLI) Debugf(format string, args ...interface{}) (n int, err error) {
	return cli.fwritef(cli.stderr, DebugLevel, format, args...)
}

func (cli CLI) Printf(format string, args ...interface{}) (n int, err error) {
	return cli.fwritef(cli.stdout, SilentLevel, format, args...)
}
