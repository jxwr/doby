package rt

import (
	"fmt"
	"sort"
)

/// dict

type DictObject struct {
	Property
	sortedkeys []string
}

func (self *DictObject) Name() string {
	return "dict"
}

func (self *DictObject) HashCode() string {
	return fmt.Sprintf("%p", self)
}

func (self *DictObject) String() string {
	s := "#{"

	ln := len(self.Property.Slots)
	idx := 0
	for _, slot := range self.Property.Slots {
		s += slot.Key.String()
		s += ":"
		s += slot.Val.String()
		if idx < ln-1 {
			s += ","
		}
		idx++
	}
	s += "}"
	return s
}

func (self *DictObject) ToString(rt *Runtime, args ...Object) []Object {
	return []Object{rt.NewStringObject(self.String())}
}

func (self *DictObject) OP__iter__(rt *Runtime, args ...Object) (results []Object) {
	idx := args[0].(*IntegerObject)
	ln := len(self.Property.Slots)

	if idx.Val == 0 {
		hashes := make([]string, ln)
		i := 0
		for hash, _ := range self.Property.Slots {
			hashes[i] = hash
			i++
		}
		self.sortedkeys = hashes
		sort.Strings(self.sortedkeys)
	}
	if idx.Val < ln {
		hash := self.sortedkeys[idx.Val]
		slot := self.Property.Slots[hash]
		results = append(results, slot.Key, slot.Val, rt.True)
	} else {
		results = append(results, rt.False)
	}
	return
}

func (self *DictObject) OP__get_index__(rt *Runtime, args ...Object) (results []Object) {
	idx := args[0]
	results = append(results, self.GetProp(idx))
	return
}

func (self *DictObject) OP__set_index__(rt *Runtime, args ...Object) (results []Object) {
	idx := args[0]
	val := args[1]
	self.SetProp(idx, val)
	return
}
