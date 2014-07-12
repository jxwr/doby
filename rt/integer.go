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
	obj.SetProp("times", NewBuiltinFuncObject("times", obj, nil))
	obj.SetProp("abs", NewBuiltinFuncObject("abs", obj, nil))
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

func (self *IntegerObject) classMethods(ctx *Runtime, method string, args ...Object) (results []Object) {
	switch method {
	case "times":
		fnobj := args[0].(*FuncObject)
		fnDecl := fnobj.Decl
		for i := 0; i < self.Val; i++ {
			fnobj.E.Put(fnDecl.Args[0].Name, NewIntegerObject(i))
			fnDecl.Body.Accept(ctx.Visitor)
		}
	case "abs":
		val := self.Val
		if val < 0 {
			val = 0 - val
		}
		results = append(results, NewIntegerObject(val))
	}
	return
}

// shits
func (self *IntegerObject) Dispatch(ctx *Runtime, method string, args ...Object) (results []Object) {
	var is bool
	if is, results = self.AccessPropMethod(method, args...); is {
		return
	}

	isFloat := false
	var val float64

	if len(args) == 0 {
		if method == "__inc__" {
			self.Val++
		} else if method == "__dec__" {
			self.Val--
		} else {
			results = self.classMethods(ctx, method, args...)
			return
		}
		return
	}

	switch arg := args[0].(type) {
	case *IntegerObject:
		val = float64(arg.Val)
	case *FloatObject:
		isFloat = true
		val = arg.Val
	}

	switch method {
	// xxx_assign
	case "__+=__":
		self.Val += int(val)
		return
	case "__-=__":
		self.Val -= int(val)
		return
	case "__*=__":
		self.Val *= int(val)
		return
	case "__/=__":
		self.Val /= int(val)
		return
	case "__%=__":
		self.Val %= int(val)
		return
	case "__|=__":
		self.Val |= int(val)
		return
	case "__&=__":
		self.Val &= int(val)
		return
	case "__^=__":
		self.Val ^= int(val)
		return
	case "__<<=__":
		self.Val <<= uint(val)
		return
	case "__>>=__":
		self.Val >>= uint(val)
		return
	case "__&^___":
		self.Val &^= int(val)
		return
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
	case "__eql__":
		cmp := float64(self.Val) == val
		results = append(results, NewBoolObject(cmp))
		return
	case "__lss__":
		cmp := float64(self.Val) < val
		results = append(results, NewBoolObject(cmp))
		return
	case "__gtr__":
		cmp := float64(self.Val) > val
		results = append(results, NewBoolObject(cmp))
		return
	case "__leq__":
		cmp := float64(self.Val) <= val
		results = append(results, NewBoolObject(cmp))
		return
	case "__geq__":
		cmp := float64(self.Val) >= val
		results = append(results, NewBoolObject(cmp))
		return
	case "__neq__":
		cmp := float64(self.Val) != val
		results = append(results, NewBoolObject(cmp))
		return
	default:
		results = self.classMethods(ctx, method, args...)
		return
	}

	if isFloat {
		results = append(results, NewFloatObject(val))
	} else {
		results = append(results, NewIntegerObject(int(val)))
	}
	return
}
