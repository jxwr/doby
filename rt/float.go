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

func (self *FloatObject) Dispatch(ctx *Runtime, method string, args ...Object) (results []Object) {
	var is bool
	if is, results = self.AccessPropMethod(method, args...); is {
		return
	}

	var val float64

	switch arg := args[0].(type) {
	case *IntegerObject:
		val = float64(arg.Val)
	case *FloatObject:
		val = arg.Val
	}

	switch method {
	case "__+=__":
		self.Val += val
		return
	case "__-=__":
		self.Val -= val
		return
	case "__*=__":
		self.Val *= val
		return
	case "__/=__":
		self.Val /= val
		return
	case "__add__":
		val = self.Val + val
	case "__sub__":
		val = self.Val - val
	case "__mul__":
		val = self.Val * val
	case "__quo__":
		val = self.Val / val
	case "__eql__":
		cmp := self.Val == val
		results = append(results, NewBoolObject(cmp))
		return
	case "__lss__":
		cmp := self.Val < val
		results = append(results, NewBoolObject(cmp))
		return
	case "__gtr__":
		cmp := self.Val > val
		results = append(results, NewBoolObject(cmp))
		return
	case "__leq__":
		cmp := self.Val <= val
		results = append(results, NewBoolObject(cmp))
		return
	case "__geq__":
		cmp := self.Val >= val
		results = append(results, NewBoolObject(cmp))
		return
	case "__neq__":
		cmp := self.Val != val
		results = append(results, NewBoolObject(cmp))
		return
	}
	results = append(results, NewFloatObject(val))
	return
}
