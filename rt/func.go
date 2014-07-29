package rt

import (
	"fmt"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/env"
	"github.com/jxwr/doubi/vm/instr"
)

/// closure

type ClosureObject struct {
	Property

	Proto *instr.ClosureProto
	Frame *Frame
}

func (self *ClosureObject) Name() string {
	return "closure"
}

func (self *ClosureObject) HashCode() string {
	return fmt.Sprintf("%p", self)
}

func (self *ClosureObject) String() string {
	return fmt.Sprintf("closure#%d", self.Proto.Seq)
}

func (self *ClosureObject) ToString(rt *Runtime, args ...Object) []Object {
	return []Object{rt.NewStringObject(self.String())}
}

/// function

type FuncObject struct {
	Property

	name string
	Decl *ast.FuncDeclExpr

	IsBuiltin bool
	Obj       Object
	E         *env.Env
}

func (self *FuncObject) Name() string {
	return "function"
}

func (self *FuncObject) HashCode() string {
	return fmt.Sprintf("%p", self)
}

func (self *FuncObject) String() string {
	return self.name
}

func (self *FuncObject) ToString(rt *Runtime, args ...Object) []Object {
	return []Object{rt.NewStringObject(self.String())}
}

func (self *FuncObject) OP__call__(rt *Runtime, args ...Object) (results []Object) {
	results = Invoke(rt, self.Obj, self.name, args...)
	return
}
