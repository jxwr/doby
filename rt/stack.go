package rt

import "fmt"

/// stack

type Stack struct {
	cur  int
	mark int
	vals []Object
}

func NewStack() *Stack {
	stack := &Stack{0, 0, []Object{}}
	return stack
}

func (self *Stack) Push(obj Object) {
	if len(self.vals) <= self.cur {
		self.vals = append(self.vals, obj)
	} else {
		self.vals[self.cur] = obj
	}
	self.cur++
}

func (self *Stack) Pop() Object {
	if self.cur == 0 {
		panic("Pop from empty stack, maybe missing return in some func")
	}
	self.cur--
	return self.vals[self.cur]
}

func (self *Stack) Mark() {
	self.mark = self.cur
}

func (self *Stack) Rewind() {
	self.cur = self.mark
}

func (self *Stack) Print() {
	fmt.Println("~~~~~~~~~~~")
	for i := 0; i < len(self.vals); i++ {
		if i > self.cur {
			fmt.Println("-----")
		}
		fmt.Printf("%s %s\n", self.vals[i].Name(), self.vals[i].String())
	}
	fmt.Println("~~~~~~~~~~~")
}
