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

var Builtins = map[string]func(args ...Object) []Object{
	"print": func(args ...Object) (results []Object) {
		ifs := []interface{}{}
		for _, arg := range args {
			ifs = append(ifs, arg)
		}
		fmt.Print(ifs...)
		return
	},
}

func (self *FuncObject) OP__call__(rt *Runtime, args ...Object) (results []Object) {
	if self.Decl == nil && self.Obj == nil {
		fn, ok := Builtins[self.name]
		if ok {
			results = fn(args...)
		}
	} else {
		results = Invoke(rt, self.Obj, self.name, args...)
	}
	return
}
