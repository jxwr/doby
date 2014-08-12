package rt

import (
	"regexp"
	"strconv"
	"strings"
)

/// string

type StringObject struct {
	Property

	Val string
}

func (self *StringObject) Name() string {
	return "string"
}

func (self *StringObject) HashCode() string {
	return self.String()
}

func (self *StringObject) String() string {
	return self.Val
}

func (self *StringObject) ToString(rt *Runtime, args ...Object) []Object {
	return []Object{rt.NewStringObject(self.String())}
}

/// methods

func (self *StringObject) Length(rt *Runtime, args ...Object) (results []Object) {
	ret := rt.NewIntegerObject(len(self.Val))
	results = append(results, ret)
	return
}

func (self *StringObject) Size(rt *Runtime, args ...Object) (results []Object) {
	ret := rt.NewIntegerObject(len(self.Val))
	results = append(results, ret)
	return
}

func (self *StringObject) Trim(rt *Runtime, args ...Object) (results []Object) {
	val := strings.TrimSpace(self.Val)
	results = append(results, rt.NewStringObject(val))
	return
}

func (self *StringObject) ParseInt(rt *Runtime, args ...Object) (results []Object) {
	v, err := strconv.ParseInt(self.Val, 10, 32)
	if err != nil {
		results = append(results, rt.Nil)
	} else {
		results = append(results, rt.NewIntegerObject(int(v)))
	}
	return
}

func (self *StringObject) ParseFloat(rt *Runtime, args ...Object) (results []Object) {
	v, err := strconv.ParseFloat(self.Val, 64)
	if err != nil {
		results = append(results, rt.Nil)
	} else {
		results = append(results, rt.NewFloatObject(v))
	}
	return
}

func (self *StringObject) Split(rt *Runtime, args ...Object) (results []Object) {
	pat := args[0].String()
	n := -1
	if len(args) == 2 {
		n = args[1].(*IntegerObject).Val
	}
	re := regexp.MustCompile(pat)
	parts := re.Split(self.Val, n)

	arr := []Object{}
	for _, part := range parts {
		arr = append(arr, rt.NewStringObject(part))
	}
	results = append(results, rt.NewArrayObject(arr))
	return
}

/// operators

func (self *StringObject) OP__add__(rt *Runtime, args ...Object) (results []Object) {
	obj := rt.NewStringObject(self.Val + args[0].String())
	results = append(results, obj)
	return
}

func (self *StringObject) OP__add_assign__(rt *Runtime, args ...Object) (results []Object) {
	self.Val += args[0].String()
	return
}

func (self *StringObject) OP__eql__(rt *Runtime, args ...Object) (results []Object) {
	cmp := self.Val == args[0].String()
	results = append(results, rt.NewBoolObject(cmp))
	return
}

func (self *StringObject) OP__neq__(rt *Runtime, args ...Object) (results []Object) {
	cmp := self.Val != args[0].String()
	results = append(results, rt.NewBoolObject(cmp))
	return
}

func (self *StringObject) OP__get_index__(rt *Runtime, args ...Object) (results []Object) {
	idx := args[0].(*IntegerObject)
	ch := string(self.Val[idx.Val])
	obj := rt.NewStringObject(ch)
	results = append(results, obj)
	return
}

func (self *StringObject) OP__slice__(rt *Runtime, args ...Object) (results []Object) {
	low := 0
	high := len(self.Val)

	lo := args[0]
	if lo != nil {
		low = lo.(*IntegerObject).Val
	}
	ho := args[1]
	if ho != nil {
		high = ho.(*IntegerObject).Val
	}

	ret := rt.NewStringObject(self.Val[low:high])
	results = append(results, ret)
	return
}
