package main

import (
	"flag"

	"github.com/jxwr/doubi/runner"
)

var input string

func init() {
	flag.StringVar(&input, "f", "", "input file")
}

func main() {
	flag.Parse()

	r := runner.NewRunner()

	if input != "" {
		r.Run(input)
	} else {
		r.Run("test/play.d")
	}
}
