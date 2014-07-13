package rt

import (
	"fmt"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/env"
)

/// function

type FuncObject struct {
	Property

	name string
	Decl *ast.FuncDeclExpr

	IsBuiltin bool
	Obj       Object
	E         *env.Env
}

func NewFuncObject(name string, decl *ast.FuncDeclExpr, e *env.Env) Object {
	obj := &FuncObject{Property(map[string]Object{}), name, decl, false, nil, e}
	return obj
}

func NewBuiltinFuncObject(name string, recv Object, e *env.Env) *FuncObject {
	obj := &FuncObject{Property(map[string]Object{}), name, nil, true, recv, e}
	return obj
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

func (self *FuncObject) OP__call__(rt *Runtime, args ...Object) (results []Object) {
	results = Invoke(rt, self.Obj, self.name, args...)
	return
}
