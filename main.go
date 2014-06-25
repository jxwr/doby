package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jxwr/doubi/parser"
)

func main() {
	fi := bufio.NewReader(os.NewFile(0, "stdin"))

	for {
		var line string
		var ok bool

		fmt.Printf("> ")
		if line, ok = readline(fi); ok {
			parser.CalcParse(&parser.Lexer{Src: line})
		} else {
			break
		}
	}
}

func readline(fi *bufio.Reader) (string, bool) {
	s, err := fi.ReadString('\n')

	if err != nil || s == "q\n" {
		return "", false
	}

	return s, true
}
