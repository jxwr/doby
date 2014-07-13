package rt

import (
	"fmt"
)

/// bool

type BoolObject struct {
	Property

	Val bool
}

func NewBoolObject(val bool) Object {
	obj := &BoolObject{Property(map[string]Object{}), val}
	return obj
}

func (self *BoolObject) Name() string {
	return "bool"
}

func (self *BoolObject) HashCode() string {
	return self.String()
}

func (self *BoolObject) String() string {
	return fmt.Sprintf("%v", self.Val)
}

func (self *BoolObject) Dispatch(ctx *Runtime, method string, args ...Object) (results []Object) {
	val := args[0].(*BoolObject).Val
	switch method {
	case "__land__":
		val = self.Val && val
	case "__lor__":
		val = self.Val || val
	case "__not__":
		val = !self.Val
	}

	results = append(results, NewBoolObject(val))
	return
}
