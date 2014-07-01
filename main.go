package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jxwr/doubi/eval"
	"github.com/jxwr/doubi/parser"
)

func runTest(filename string) {
	var contents []byte
	var err error

	fmt.Println("================================> ", filename)

	contents, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	parser.ProgramAst = nil
	parser.CalcParse(&parser.Lexer{Src: string(contents)})
	eval.Eval(parser.ProgramAst)
}

func repl() {
	fi := bufio.NewReader(os.NewFile(0, "stdin"))
	for {
		var src string
		var ok bool

		fmt.Printf("> ")
		if src, ok = readGist(fi); ok {
			parser.CalcParse(&parser.Lexer{Src: src})
			eval.Eval(parser.ProgramAst)
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

func main() {
	runTest("test/vars.d")
	runTest("test/datatype.d")
	runTest("test/func.d")
	runTest("test/cond.d")

	repl()
}
