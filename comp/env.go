package comp

import (
	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/symtab"
)

type Env struct {
	Outer  *Env
	Symtab *symtab.Symtab
}

func NewEnv(outer *Env) *Env {
	e := &Env{outer, symtab.NewSymtab()}
	return e
}

func (e *Env) Put(name string, node ast.Node) {
	e.Symtab.Put(name, node)
}

func (e *Env) LookUp(name string) ast.Node {
	env := e
	for env != nil {
		ne := env.Symtab.LookUp(name)
		if ne != nil {
			return ne
		} else {
			env = env.Outer
		}
	}
	return nil
}
