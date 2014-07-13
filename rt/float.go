package rt

import (
	"fmt"
)

/// float

type FloatObject struct {
	Property

	Val float64
}

func NewFloatObject(val float64) Object {
	obj := &FloatObject{Property(map[string]Object{}), val}
	return obj
}

func (self *FloatObject) HashCode() string {
	return self.String()
}

func (self *FloatObject) Name() string {
	return "float"
}

func (self *FloatObject) String() string {
	return fmt.Sprintf("%f", self.Val)
}

func (self *FloatObject) Abs(rt *Runtime, args ...Object) (results []Object) {
	val := self.Val
	if val < 0 {
		val = 0 - val
	}
	results = append(results, NewFloatObject(val))
	return
}

// +=
func (self *FloatObject) OP__add_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__+=__", args[0])
	return
}

// -=
func (self *FloatObject) OP__sub_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__-=__", args[0])
	return
}

// *=
func (self *FloatObject) OP__mul_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__*=__", args[0])
	return
}

// /=
func (self *FloatObject) OP__quo_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__/=__", args[0])
	return
}

// +
func (self *FloatObject) OP__add__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__add__", args[0])
	return
}

// -
func (self *FloatObject) OP__sub__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__sub__", args[0])
	return
}

// *
func (self *FloatObject) OP__mul__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__mul__", args[0])
	return
}

// /
func (self *FloatObject) OP__quo__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__quo__", args[0])
	return
}

// ==
func (self *FloatObject) OP__eql__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__eql__", args[0])
	return
}

// <
func (self *FloatObject) OP__lss__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__lss__", args[0])
	return
}

// >
func (self *FloatObject) OP__gtr__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__gtr__", args[0])
	return
}

// !=
func (self *FloatObject) OP__neq__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__neq__", args[0])
	return
}

// >=
func (self *FloatObject) OP__geq__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__geq__", args[0])
	return
}

// <=
func (self *FloatObject) OP__leq__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__leq__", args[0])
	return
}

func (self *FloatObject) assign(method string, obj Object) (results []Object) {
	var val float64

	switch arg := obj.(type) {
	case *IntegerObject:
		val = float64(arg.Val)
	case *FloatObject:
		val = arg.Val
	}

	switch method {
	case "__+=__":
		self.Val += val
		results = append(results, NewFloatObject(val))
	case "__-=__":
		self.Val -= val
		results = append(results, NewFloatObject(val))
	case "__*=__":
		self.Val *= val
		results = append(results, NewFloatObject(val))
	case "__/=__":
		self.Val /= val
		results = append(results, NewFloatObject(val))
	}
	return
}

func (self *FloatObject) binary(method string, obj Object) (results []Object) {
	var val float64

	switch arg := obj.(type) {
	case *IntegerObject:
		val = float64(arg.Val)
	case *FloatObject:
		val = arg.Val
	}

	switch method {
	case "__add__":
		val = self.Val + val
		results = append(results, NewFloatObject(val))
	case "__sub__":
		val = self.Val - val
		results = append(results, NewFloatObject(val))
	case "__mul__":
		val = self.Val * val
		results = append(results, NewFloatObject(val))
	case "__quo__":
		val = self.Val / val
		results = append(results, NewFloatObject(val))
	case "__eql__":
		cmp := self.Val == val
		results = append(results, NewBoolObject(cmp))
	case "__lss__":
		cmp := self.Val < val
		results = append(results, NewBoolObject(cmp))
	case "__gtr__":
		cmp := self.Val > val
		results = append(results, NewBoolObject(cmp))
	case "__leq__":
		cmp := self.Val <= val
		results = append(results, NewBoolObject(cmp))
	case "__geq__":
		cmp := self.Val >= val
		results = append(results, NewBoolObject(cmp))
	case "__neq__":
		cmp := self.Val != val
		results = append(results, NewBoolObject(cmp))
	}
	return
}
