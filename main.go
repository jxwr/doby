package main

import (
	"flag"

	"github.com/jxwr/doby/runner"
)

var input string
var dumpInstrs bool
var printStack bool

func init() {
	flag.StringVar(&input, "f", "", "input file")
	flag.BoolVar(&dumpInstrs, "i", true, "dump instrs")
	flag.BoolVar(&printStack, "s", false, "print stack")
}

func main() {
	flag.Parse()

	r := runner.NewRunner()

	r.SetDumpInstrs(dumpInstrs)
	r.SetPrintStack(printStack)

	if input != "" {
		r.Run(input)
	} else {
		r.Run("test/play.d")
	}
}
