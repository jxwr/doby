package rt

import (
	"fmt"
)

/// array

type ArrayObject struct {
	Property

	Vals []Object
}

func (self *ArrayObject) Name() string {
	return "array"
}

func (self *ArrayObject) HashCode() string {
	return fmt.Sprintf("%p", self)
}

func (self *ArrayObject) String() string {
	s := "["
	ln := len(self.Vals)
	for i, val := range self.Vals {
		s += val.String()
		if i < ln-1 {
			s += ","
		}
	}
	s += "]"
	return s
}

func (self *ArrayObject) ToString(rt *Runtime, args ...Object) []Object {
	return []Object{rt.NewStringObject(self.String())}
}

func (self *ArrayObject) Append(rt *Runtime, args ...Object) (results []Object) {
	val := args[0]
	self.Vals = append(self.Vals, val)
	return
}

func (self *ArrayObject) Length(rt *Runtime, args ...Object) (results []Object) {
	ret := rt.NewIntegerObject(len(self.Vals))
	results = append(results, ret)
	return
}

func (self *ArrayObject) Size(rt *Runtime, args ...Object) (results []Object) {
	ret := rt.NewIntegerObject(len(self.Vals))
	results = append(results, ret)
	return
}

func (self *ArrayObject) Each(rt *Runtime, args ...Object) (results []Object) {
	fnobj := args[0].(*FuncObject)
	fnDecl := fnobj.Decl
	for i := 0; i < len(self.Vals); i++ {
		fnobj.E.Put(fnDecl.Args[0].Name, self.Vals[i])
		fnDecl.Body.Accept(rt.Visitor)
	}
	return
}

func (self *ArrayObject) Map(rt *Runtime, args ...Object) (results []Object) {
	fnobj := args[0].(*FuncObject)
	fnDecl := fnobj.Decl

	arr := []Object{}
	for i := 0; i < len(self.Vals); i++ {
		fnobj.E.Put(fnDecl.Args[0].Name, self.Vals[i])
		rt.NeedReturn = false
		fnDecl.Body.Accept(rt.Visitor)
		arr = append(arr, rt.Pop())
	}
	obj := rt.NewArrayObject(arr)
	results = append(results, obj)
	return
}

func (self *ArrayObject) OP__add__(rt *Runtime, args ...Object) (results []Object) {
	vals := append(self.Vals[:], args[0].(*ArrayObject).Vals...)
	ret := rt.NewArrayObject(vals)
	results = append(results, ret)
	return
}

func (self *ArrayObject) OP__add_assign__(rt *Runtime, args ...Object) (results []Object) {
	self.Vals = append(self.Vals, args[0].(*ArrayObject).Vals...)
	return
}

func (self *ArrayObject) OP__get_index__(rt *Runtime, args ...Object) (results []Object) {
	idx := args[0].(*IntegerObject)
	obj := self.Vals[idx.Val]
	results = append(results, obj)
	return
}

func (self *ArrayObject) OP__set_index__(rt *Runtime, args ...Object) (results []Object) {
	idx := args[0].(*IntegerObject)
	val := args[1]
	self.Vals[idx.Val] = val
	return
}

func (self *ArrayObject) OP__slice__(rt *Runtime, args ...Object) (results []Object) {
	low := 0
	high := len(self.Vals)

	if len(args) > 0 {
		lo := args[0]
		low = lo.(*IntegerObject).Val
	}
	if len(args) > 1 {
		ho := args[1]
		high = ho.(*IntegerObject).Val
	}

	vals := self.Vals[low:high]
	ret := rt.NewArrayObject(vals)
	results = append(results, ret)
	return
}
