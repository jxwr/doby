package rt

import (
	"fmt"
)

/// integer

type IntegerObject struct {
	Property

	Val int
}

func NewIntegerObject(val int) Object {
	obj := &IntegerObject{Property(map[string]Object{}), val}
	obj.SetProp("Times", NewBuiltinFuncObject("Times", obj, nil))
	obj.SetProp("Abs", NewBuiltinFuncObject("Abs", obj, nil))
	return obj
}

func (self *IntegerObject) Name() string {
	return "integer"
}

func (self *IntegerObject) HashCode() string {
	return self.String()
}

func (self *IntegerObject) String() string {
	return fmt.Sprintf("%d", self.Val)
}

func (self *IntegerObject) Times(rt *Runtime, args ...Object) (results []Object) {
	fnobj := args[0].(*FuncObject)
	fnDecl := fnobj.Decl
	for i := 0; i < self.Val; i++ {
		fnobj.E.Put(fnDecl.Args[0].Name, NewIntegerObject(i))
		fnDecl.Body.Accept(rt.Visitor)
	}
	return
}

func (self *IntegerObject) Abs(rt *Runtime, args ...Object) (results []Object) {
	val := self.Val
	if val < 0 {
		val = 0 - val
	}
	results = append(results, NewIntegerObject(val))
	return
}

func (self *IntegerObject) OP__inc__(rt *Runtime, args ...Object) (results []Object) {
	self.Val++
	return
}

func (self *IntegerObject) OP__dec__(rt *Runtime, args ...Object) (results []Object) {
	self.Val--
	return
}

// +=
func (self *IntegerObject) OP__add_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__+=__", args[0])
	return
}

// -=
func (self *IntegerObject) OP__sub_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__-=__", args[0])
	return
}

// *=
func (self *IntegerObject) OP__mul_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__*=__", args[0])
	return
}

// /=
func (self *IntegerObject) OP__quo_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__/=__", args[0])
	return
}

// %=
func (self *IntegerObject) OP__rem_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__%=__", args[0])
	return
}

// &=
func (self *IntegerObject) OP__and_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__&=__", args[0])
	return
}

// |=
func (self *IntegerObject) OP__or_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__|=__", args[0])
	return
}

// ^=
func (self *IntegerObject) OP__xor_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__^=__", args[0])
	return
}

// <<=
func (self *IntegerObject) OP__shl_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__<<=__", args[0])
	return
}

// >>=
func (self *IntegerObject) OP__shr_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__>>=__", args[0])
	return
}

// &^=
func (self *IntegerObject) OP__and_not_assign__(rt *Runtime, args ...Object) (results []Object) {
	results = self.assign("__&^=__", args[0])
	return
}

// +
func (self *IntegerObject) OP__add__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__add__", args[0])
	return
}

// -
func (self *IntegerObject) OP__sub__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__sub__", args[0])
	return
}

// *
func (self *IntegerObject) OP__mul__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__mul__", args[0])
	return
}

// /
func (self *IntegerObject) OP__quo__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__quo__", args[0])
	return
}

// %
func (self *IntegerObject) OP__rem__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__rem__", args[0])
	return
}

// &
func (self *IntegerObject) OP__and__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__and__", args[0])
	return
}

// |
func (self *IntegerObject) OP__or__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__or__", args[0])
	return
}

// ^
func (self *IntegerObject) OP__xor__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__xor__", args[0])
	return
}

// <<
func (self *IntegerObject) OP__shl__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__shl__", args[0])
	return
}

// >>
func (self *IntegerObject) OP__shr__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__shr__", args[0])
	return
}

// &^
func (self *IntegerObject) OP__and_not__(rt *Runtime, args ...Object) (results []Object) {
	results = self.binary("__and_not__", args[0])
	return
}

// ==
func (self *IntegerObject) OP__eql__(rt *Runtime, args ...Object) (results []Object) {
	results = self.logic("__eql__", args[0])
	return
}

// <
func (self *IntegerObject) OP__lss__(rt *Runtime, args ...Object) (results []Object) {
	results = self.logic("__lss__", args[0])
	return
}

// >
func (self *IntegerObject) OP__gtr__(rt *Runtime, args ...Object) (results []Object) {
	results = self.logic("__gtr__", args[0])
	return
}

// !=
func (self *IntegerObject) OP__neq__(rt *Runtime, args ...Object) (results []Object) {
	results = self.logic("__neq__", args[0])
	return
}

// >=
func (self *IntegerObject) OP__geq__(rt *Runtime, args ...Object) (results []Object) {
	results = self.logic("__geq__", args[0])
	return
}

// <=
func (self *IntegerObject) OP__leq__(rt *Runtime, args ...Object) (results []Object) {
	results = self.logic("__leq__", args[0])
	return
}

func (self *IntegerObject) assign(method string, obj Object) (results []Object) {
	var val float64

	switch arg := obj.(type) {
	case *IntegerObject:
		val = float64(arg.Val)
	case *FloatObject:
		val = arg.Val
	}

	switch method {
	case "__+=__":
		self.Val += int(val)
	case "__-=__":
		self.Val -= int(val)
	case "__*=__":
		self.Val *= int(val)
	case "__/=__":
		self.Val /= int(val)
	case "__%=__":
		self.Val %= int(val)
	case "__|=__":
		self.Val |= int(val)
	case "__&=__":
		self.Val &= int(val)
	case "__^=__":
		self.Val ^= int(val)
	case "__<<=__":
		self.Val <<= uint(val)
	case "__>>=__":
		self.Val >>= uint(val)
	case "__&^=___":
		self.Val &^= int(val)
	}
	return
}

func (self *IntegerObject) binary(method string, obj Object) (results []Object) {
	var val float64

	isFloat := false
	switch arg := obj.(type) {
	case *IntegerObject:
		val = float64(arg.Val)
	case *FloatObject:
		isFloat = true
		val = arg.Val
	}

	switch method {
	// binop
	case "__add__":
		val = float64(self.Val) + val
	case "__sub__":
		val = float64(self.Val) - val
	case "__mul__":
		val = float64(self.Val) * val
	case "__quo__":
		val = float64(self.Val) / val
	case "__rem__":
		val = float64(self.Val % int(val))
	case "__and__":
		val = float64(self.Val & int(val))
	case "__or__":
		val = float64(self.Val | int(val))
	case "__xor__":
		val = float64(self.Val ^ int(val))
	case "__shl__":
		val = float64(uint(self.Val) << uint(val))
	case "__shr__":
		val = float64(uint(self.Val) >> uint(val))
	case "__and_not__":
		val = float64(uint(self.Val) &^ uint(val))
	}

	if isFloat {
		results = append(results, NewFloatObject(val))
	} else {
		results = append(results, NewIntegerObject(int(val)))
	}
	return
}

func (self *IntegerObject) logic(method string, obj Object) (results []Object) {
	var val float64

	switch arg := obj.(type) {
	case *IntegerObject:
		val = float64(arg.Val)
	case *FloatObject:
		val = arg.Val
	}

	switch method {
	case "__eql__":
		cmp := float64(self.Val) == val
		results = append(results, NewBoolObject(cmp))
	case "__lss__":
		cmp := float64(self.Val) < val
		results = append(results, NewBoolObject(cmp))
	case "__gtr__":
		cmp := float64(self.Val) > val
		results = append(results, NewBoolObject(cmp))
	case "__leq__":
		cmp := float64(self.Val) <= val
		results = append(results, NewBoolObject(cmp))
	case "__geq__":
		cmp := float64(self.Val) >= val
		results = append(results, NewBoolObject(cmp))
	case "__neq__":
		cmp := float64(self.Val) != val
		results = append(results, NewBoolObject(cmp))
	}
	return
}
