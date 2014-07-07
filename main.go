package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/comp"
	"github.com/jxwr/doubi/parser"
)

func EvalStmt(stmt *ast.Stmt) {
	fmt.Printf("%#v\n", *stmt)
}

func Eval(stmts []ast.Stmt) {
	pretty := &comp.PrettyPrinter{false, 0, true}
	attr := &comp.Attr{false, comp.NewEnv(nil)}
	eval := &comp.Eval{false, comp.NewEnv(nil), comp.NewStack(),
		false, 0, false, false}

	if false {
		for _, stmt := range stmts {
			stmt.Accept(pretty)
		}
	}

	for _, stmt := range stmts {
		stmt.Accept(attr)
	}

	for _, stmt := range stmts {
		stmt.Accept(eval)
	}

}

func runTest(filename string) {
	var contents []byte
	var err error

	fmt.Println("=============> ", filename, " <=============")

	contents, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	parser.ProgramAst = nil
	parser.DoubiParse(parser.NewLexer(string(contents)))
	Eval(parser.ProgramAst)
}

func repl() {
	fi := bufio.NewReader(os.NewFile(0, "stdin"))
	for {
		var src string
		var ok bool

		fmt.Printf("> ")
		if src, ok = readGist(fi); ok {
			parser.DoubiParse(&parser.Lexer{Src: src, Line: 1, Pos: 0})
			Eval(parser.ProgramAst)
		} else {
			break
		}
	}
}

func readGist(fi *bufio.Reader) (string, bool) {
	s, err := fi.ReadString('\n')

	if err != nil || s == "q\n" || s == "exit\n" {
		return "", false
	}

	return s, true
}

var input string

func init() {
	flag.StringVar(&input, "i", "", "input file")
}

func main() {
	flag.Parse()

	if input != "" {
		runTest(input)
	} else {
		//runTest("test/vars.d")
		//runTest("test/datatype.d")
		//runTest("test/func.d")
		//runTest("test/cond.d")
		//runTest("test/quicksort.d")
		runTest("test/play.d")
	}
}
