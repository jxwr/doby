package symtab

import (
	"github.com/jxwr/doubi/ast"
)

type Symtab struct {
	tab map[string]ast.Node
}

func NewSymtab() *Symtab {
	st := &Symtab{map[string]ast.Node{}}
	return st
}

func (t *Symtab) Put(name string, node ast.Node) {
	t.tab[name] = node
}

func (t *Symtab) LookUp(name string) ast.Node {
	val, ok := t.tab[name]
	if ok {
		return val
	}
	return nil
}
