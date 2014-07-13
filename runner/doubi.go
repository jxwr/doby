// test code runner

package runner

import (
	"fmt"
	"io/ioutil"

	"github.com/jxwr/doubi/comp"
	"github.com/jxwr/doubi/env"
	"github.com/jxwr/doubi/parser"
	"github.com/jxwr/doubi/rt"
)

type Runner struct {
	pretty  *comp.PrettyPrinter
	attr    *comp.Attr
	eval    *comp.Eval
	runtime *rt.Runtime
}

func NewRunner() *Runner {
	pretty := &comp.PrettyPrinter{false, 0, true}
	attr := &comp.Attr{false, env.NewEnv(nil), nil}
	eval := comp.NewEvaluater()

	runtime := rt.NewRuntime(eval)
	eval.SetRuntime(runtime)

	runner := &Runner{pretty, attr, eval, runtime}
	return runner
}

func (self *Runner) RegisterFunctions(name string, fns []interface{}) {
	self.runtime.RegisterFunctions(name, fns)
}

func (self *Runner) Run(filename string) {
	var contents []byte
	var err error

	fmt.Println("=============> ", filename, " <=============")

	contents, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	parser.ProgramAst = nil
	parser.DoubiParse(parser.NewLexer(string(contents)))

	for _, stmt := range parser.ProgramAst {
		stmt.Accept(self.attr)
	}

	for _, stmt := range parser.ProgramAst {
		stmt.Accept(self.eval)
	}
}
