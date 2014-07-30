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

func (self *ArrayObject) Push(rt *Runtime, args ...Object) (results []Object) {
	for _, arg := range args {
		self.Vals = append(self.Vals, arg)
	}
	results = append(results, self)
	return
}

func (self *ArrayObject) Pop(rt *Runtime, args ...Object) (results []Object) {
	n := 1
	if len(args) == 1 {
		n = args[0].(*IntegerObject).Val
		if n < 0 || n > len(self.Vals) {
			panic("pop array out of range")
		}
	}
	self.Vals = self.Vals[:len(self.Vals)-n]
	results = append(results, self)
	return
}

func (self *ArrayObject) Take(rt *Runtime, args ...Object) (results []Object) {
	var n int
	if len(args) == 1 {
		n = args[0].(*IntegerObject).Val
		if n < 0 || n > len(self.Vals) {
			panic("take array out of range")
		}
	} else {
		rt.Fatalf("array::Take need one integer argumenet, %d given", len(args))
	}
	obj := rt.NewArrayObject(self.Vals[:n])
	results = append(results, obj)
	return
}

func (self *ArrayObject) Drop(rt *Runtime, args ...Object) (results []Object) {
	var n int
	if len(args) == 1 {
		n = args[0].(*IntegerObject).Val
		if n < 0 || n > len(self.Vals) {
			panic("drop array out of range")
		}
	} else {
		rt.Fatalf("array::Drop need one integer argumenet, %d given", len(args))
	}
	obj := rt.NewArrayObject(self.Vals[:n])
	self.Vals = self.Vals[n:]
	results = append(results, obj)
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
	fnobj := args[0].(*ClosureObject)
	for i := 0; i < len(self.Vals); i++ {
		rt.CallFuncObj(fnobj, self.Vals[i])
	}
	return
}

func (self *ArrayObject) Map(rt *Runtime, args ...Object) (results []Object) {
	fnobj := args[0].(*ClosureObject)
	arr := []Object{}
	for i := 0; i < len(self.Vals); i++ {
		rt.CallFuncObj(fnobj, self.Vals[i])
		arr = append(arr, rt.Pop())
	}
	obj := rt.NewArrayObject(arr)
	results = append(results, obj)
	return
}

func (self *ArrayObject) Select(rt *Runtime, args ...Object) (results []Object) {
	fnobj := args[0].(*ClosureObject)
	arr := []Object{}
	for i := 0; i < len(self.Vals); i++ {
		rt.CallFuncObj(fnobj, self.Vals[i])
		if rt.Pop().(*BoolObject).Val {
			arr = append(arr, self.Vals[i])
		}
	}
	obj := rt.NewArrayObject(arr)
	results = append(results, obj)
	return
}

func (self *ArrayObject) OP__iter__(rt *Runtime, args ...Object) (results []Object) {
	idx := args[0].(*IntegerObject)
	if idx.Val < len(self.Vals) {
		obj := self.Vals[idx.Val]
		results = append(results, args[0], obj, rt.True)
	} else {
		results = append(results, rt.False)
	}
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
	lo := args[0]
	ho := args[1]

	v, ok := lo.(*IntegerObject)
	if ok {
		low = v.Val
	}

	v, ok = ho.(*IntegerObject)
	if ok {
		high = ho.(*IntegerObject).Val
	}

	vals := self.Vals[low:high]
	ret := rt.NewArrayObject(vals)
	results = append(results, ret)
	return
}
