// test code runner

package runner

import (
	"fmt"
	"io/ioutil"

	"github.com/jxwr/doby/comp"
	"github.com/jxwr/doby/parser"
	"github.com/jxwr/doby/rt"
	"github.com/jxwr/doby/vm"
)

type Runner struct {
	pretty  *comp.PrettyPrinter
	attr    *comp.Attr
	irb     *comp.IRBuilder
	runtime *rt.Runtime

	dumpInstrs bool
	printStack bool
}

func NewRunner() *Runner {
	pretty := comp.NewPrettyPrinter()
	attr := comp.NewAttr()
	irb := comp.NewIRBuilder()
	runtime := rt.NewRuntime()

	runner := &Runner{pretty, attr, irb, runtime, false, false}
	return runner
}

func (self *Runner) SetDumpInstrs(dumpInstrs bool) {
	self.dumpInstrs = dumpInstrs
}

func (self *Runner) SetPrintStack(printStack bool) {
	self.printStack = printStack
}

func (self *Runner) RegisterFunctions(name string, fns []interface{}) {
	self.runtime.RegisterFunctions(name, fns)
}

func (self *Runner) RegisterVars(name string, vars map[string]interface{}) {
	self.runtime.RegisterVars(name, vars)
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
	lexer := parser.NewLexer(filename, string(contents))
	self.irb.SetLexer(lexer)
	parser.DobyParse(lexer)

	for _, stmt := range parser.ProgramAst {
		stmt.Accept(self.attr)
	}

	// IR generation
	irb := self.irb
	for _, stmt := range parser.ProgramAst {
		stmt.Accept(irb)
	}

	if self.dumpInstrs {
		irb.RootClosure().DumpClosureProto()
	}

	fmt.Println("===========================")

	// run IRs in the vm
	vm := vm.NewVM(irb.RootClosure(), irb.ClosureTable(), self.runtime)
	self.runtime.Runner = vm
	vm.Run()

	if self.printStack {
		self.runtime.Stack.Print()
	}
}
