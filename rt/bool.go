package rt

import (
	"fmt"
)

/// bool

type BoolObject struct {
	Property

	Val bool
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

func (self *BoolObject) ToString(rt *Runtime, args ...Object) []Object {
	return []Object{rt.NewStringObject(self.String())}
}

func (self *BoolObject) OP__land__(rt *Runtime, args ...Object) (results []Object) {
	val := args[0].(*BoolObject).Val
	val = self.Val && val
	results = append(results, rt.NewBoolObject(val))
	return
}

func (self *BoolObject) OP__lor__(rt *Runtime, args ...Object) (results []Object) {
	val := args[0].(*BoolObject).Val
	val = self.Val || val
	results = append(results, rt.NewBoolObject(val))
	return
}

func (self *BoolObject) OP__not__(rt *Runtime, args ...Object) (results []Object) {
	val := !self.Val
	results = append(results, rt.NewBoolObject(val))
	return
}
