package rt

import (
	"fmt"
)

/// set

type SetObject struct {
	Property

	Vals []Object
}

func (self *SetObject) Name() string {
	return "set"
}

func (self *SetObject) HashCode() string {
	return fmt.Sprintf("%p", self)
}

func (self *SetObject) String() string {
	s := "#["
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

func (self *SetObject) ToString(rt *Runtime, args ...Object) []Object {
	return []Object{rt.NewStringObject(self.String())}
}
