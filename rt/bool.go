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

func (self *BoolObject) OP__land__(rt *Runtime, args ...Object) (results []Object) {
	val := args[0].(*BoolObject).Val
	val = self.Val && val
	results = append(results, NewBoolObject(val))
	return
}

func (self *BoolObject) OP__lor__(rt *Runtime, args ...Object) (results []Object) {
	val := args[0].(*BoolObject).Val
	val = self.Val || val
	results = append(results, NewBoolObject(val))
	return
}

func (self *BoolObject) OP__not__(rt *Runtime, args ...Object) (results []Object) {
	val := !self.Val
	results = append(results, NewBoolObject(val))
	return
}
