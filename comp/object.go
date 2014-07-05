package comp

import (
	"fmt"

	"github.com/jxwr/doubi/ast"
)

type Object interface {
	Dispatch(method string, args ...Object) []Object
	Name() string
	String() string
}

/// string

type StringObject struct {
	val string
}

func NewStringObject(val string) Object {
	obj := &StringObject{val}
	return obj
}

func (self *StringObject) Name() string {
	return "string"
}

func (self *StringObject) String() string {
	return self.val
}

func (self *StringObject) Dispatch(method string, args ...Object) (results []Object) {
	switch method {
	case "__add__":
		obj := NewStringObject(self.val + args[0].String())
		results = append(results, obj)
	}
	return
}

/// integer

type IntegerObject struct {
	val int
}

func NewIntegerObject(val int) Object {
	obj := &IntegerObject{val}
	return obj
}

func (self *IntegerObject) Name() string {
	return "integer"
}

func (self *IntegerObject) String() string {
	return fmt.Sprintf("%d", self.val)
}

func (self *IntegerObject) Dispatch(method string, args ...Object) (results []Object) {
	switch method {
	case "__add__":
		val := self.val + args[0].(*IntegerObject).val
		results = append(results, NewIntegerObject(val))
	}
	return
}

/// float

type FloatObject struct {
	val float64
}

func NewFloatObject(val float64) Object {
	obj := &FloatObject{val}
	return obj
}

func (self *FloatObject) Name() string {
	return "float"
}

func (self *FloatObject) String() string {
	return fmt.Sprintf("%f", self.val)
}

func (self *FloatObject) Dispatch(method string, args ...Object) (results []Object) {
	switch method {
	case "__add__":
		var val float64
		switch arg := args[0].(type) {
		case *IntegerObject:
			val = self.val + float64(arg.val)
		case *FloatObject:
			val = self.val + arg.val
		}
		results = append(results, NewFloatObject(val))
	}
	return
}

/// array

type ArrayObject struct {
	vals []Object
}

func NewArrayObject(vals []Object) Object {
	obj := &ArrayObject{vals}
	return obj
}

func (self *ArrayObject) Name() string {
	return "array"
}

func (self *ArrayObject) String() string {
	s := "["
	ln := len(self.vals)
	for i, val := range self.vals {
		s += val.String()
		if i < ln-1 {
			s += ","
		}
	}
	s += "]"
	return s
}

func (self *ArrayObject) Dispatch(method string, args ...Object) (results []Object) {
	switch method {
	case "__add__":
		fmt.Println("__add__")
	}
	return
}

/// set

type SetObject struct {
	vals []Object
}

func NewSetObject(vals []Object) Object {
	obj := &SetObject{vals}
	return obj
}

func (self *SetObject) Name() string {
	return "set"
}

func (self *SetObject) String() string {
	s := "#["
	ln := len(self.vals)
	for i, val := range self.vals {
		s += val.String()
		if i < ln-1 {
			s += ","
		}
	}
	s += "]"
	return s
}

func (self *SetObject) Dispatch(method string, args ...Object) (results []Object) {
	switch method {
	case "__add__":
		fmt.Println("__add__")
	}
	return
}

/// function

type FuncObject struct {
	name string
	node *ast.Node
}

func NewFuncObject(name string, node *ast.Node) Object {
	obj := &FuncObject{name, node}
	return obj
}

func (self *FuncObject) Name() string {
	return "function"
}

func (self *FuncObject) String() string {
	return self.name
}

func (self *FuncObject) Dispatch(method string, args ...Object) (results []Object) {
	switch method {
	case "__call__":
		if self.node == nil {
			results = self.callBuiltin(args...)
		}
	}
	return
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

func (self *FuncObject) callBuiltin(args ...Object) (results []Object) {
	fn, ok := Builtins[self.name]
	if ok {
		fn(args...)
	}
	return
}
