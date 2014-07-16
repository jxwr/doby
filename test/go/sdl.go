package main

import (
	"flag"

	"github.com/jxwr/doubi/runner"
	"github.com/veandco/go-sdl2/sdl"
)

var input string

func init() {
	flag.StringVar(&input, "i", "", "input file")
}

func NewRect() sdl.Rect {
	return sdl.Rect{1, 2, 3, 4}
}

func main() {
	flag.Parse()

	r := runner.NewRunner()

	r.RegisterFunctions("sdl", []interface{}{NewRect})

	if input != "" {
		r.Run(input)
	} else {
		r.Run("test/play.d")
	}
}
