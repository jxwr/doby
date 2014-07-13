package rt

import (
	"fmt"
)

/// array

type ArrayObject struct {
	Property

	Vals []Object
}

func NewArrayObject(vals []Object) Object {
	obj := &ArrayObject{Property(map[string]Object{}), vals}
	obj.SetProp("Append", NewBuiltinFuncObject("Append", obj, nil))
	obj.SetProp("Length", NewBuiltinFuncObject("Length", obj, nil))

	return obj
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

func (self *ArrayObject) Append(ctx *Runtime, args ...Object) (results []Object) {
	val := args[0]
	self.Vals = append(self.Vals, val)
	return
}

func (self *ArrayObject) Length(ctx *Runtime, args ...Object) (results []Object) {
	ret := NewIntegerObject(len(self.Vals))
	results = append(results, ret)
	return
}

func (self *ArrayObject) Size(ctx *Runtime, args ...Object) (results []Object) {
	ret := NewIntegerObject(len(self.Vals))
	results = append(results, ret)
	return
}

func (self *ArrayObject) OP__add__(ctx *Runtime, args ...Object) (results []Object) {
	vals := append(self.Vals[:], args[0].(*ArrayObject).Vals...)
	ret := NewArrayObject(vals)
	results = append(results, ret)
	return
}

func (self *ArrayObject) OP__add_assign__(ctx *Runtime, args ...Object) (results []Object) {
	self.Vals = append(self.Vals, args[0].(*ArrayObject).Vals...)
	return
}

func (self *ArrayObject) OP__get_index__(ctx *Runtime, args ...Object) (results []Object) {
	idx := args[0].(*IntegerObject)
	obj := self.Vals[idx.Val]
	results = append(results, obj)
	return
}

func (self *ArrayObject) OP__set_index__(ctx *Runtime, args ...Object) (results []Object) {
	idx := args[0].(*IntegerObject)
	val := args[1]
	self.Vals[idx.Val] = val
	return
}

func (self *ArrayObject) OP__slice__(ctx *Runtime, args ...Object) (results []Object) {
	low := 0
	high := len(self.Vals)

	lo := args[0]
	if lo != nil {
		low = lo.(*IntegerObject).Val
	}
	ho := args[1]
	if ho != nil {
		high = ho.(*IntegerObject).Val
	}

	vals := self.Vals[low:high]
	ret := NewArrayObject(vals)
	results = append(results, ret)
	return
}
