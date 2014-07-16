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

func MakeRect(x, y, h, w int32) sdl.Rect {
	return sdl.Rect{x, y, h, w}
}

func NewRect(x, y, h, w int32) *sdl.Rect {
	return &sdl.Rect{x, y, h, w}
}

func main() {
	flag.Parse()

	r := runner.NewRunner()

	r.RegisterFunctions("sdl", []interface{}{
		MakeRect, NewRect, sdl.CreateWindow, sdl.Delay, sdl.PollEvent,
	})

	r.RegisterVars("sdl", map[string]interface{}{
		"WINDOWPOS_UNDEFINED": sdl.WINDOWPOS_UNDEFINED,
		"WINDOW_SHOWN":        sdl.WINDOW_SHOWN,
	})

	if input != "" {
		r.Run(input)
	} else {
		r.Run("test/play.d")
	}
}
