package main

import (
	"github.com/alecthomas/kong"
)

var CLI struct {
	TestModel TestModelCmd `cmd:"" help:"Test a model"`
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run(&CLI)
	ctx.FatalIfErrorf(err)
}
