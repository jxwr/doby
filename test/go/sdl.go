package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"

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
	f, err := os.Create("prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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
