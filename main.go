package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jxwr/doubi/parser"
)

func main() {
	fi := bufio.NewReader(os.NewFile(0, "stdin"))

	for {
		var src string
		var ok bool

		fmt.Printf("> ")
		if src, ok = readGist(fi); ok {
			parser.CalcParse(&parser.Lexer{Src: src})
		} else {
			break
		}
	}
}

func readGist(fi *bufio.Reader) (string, bool) {
	s, err := fi.ReadString('~')

	if err != nil || s == "q\n" {
		return "", false
	}

	return strings.TrimSuffix(s, "~"), true
}
