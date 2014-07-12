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
	obj.SetProp("append", NewBuiltinFuncObject("append", obj, nil))
	obj.SetProp("length", NewBuiltinFuncObject("length", obj, nil))

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

func (self *ArrayObject) Dispatch(ctx *Runtime, method string, args ...Object) (results []Object) {
	var is bool
	if is, results = self.AccessPropMethod(method, args...); is {
		return
	}

	switch method {
	case "__add__":
		vals := append(self.Vals[:], args[0].(*ArrayObject).Vals...)
		ret := NewArrayObject(vals)
		results = append(results, ret)
	case "__+=__":
		self.Vals = append(self.Vals, args[0].(*ArrayObject).Vals...)
	case "__get_index__":
		idx := args[0].(*IntegerObject)
		obj := self.Vals[idx.Val]
		results = append(results, obj)
	case "__set_index__":
		idx := args[0].(*IntegerObject)
		val := args[1]
		self.Vals[idx.Val] = val
	case "__slice__":
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
	case "append":
		val := args[0]
		self.Vals = append(self.Vals, val)
	case "length":
		ret := NewIntegerObject(len(self.Vals))
		results = append(results, ret)
	}
	return
}
