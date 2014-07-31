package rt

import "fmt"

/// stack

type Stack struct {
	cur  int
	mark []int
	vals []Object
}

func NewStack() *Stack {
	stack := &Stack{0, []int{}, []Object{}}
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
	self.mark = append(self.mark, self.cur)
}

func (self *Stack) MarkN(offset int) {
	self.mark = append(self.mark, self.cur+offset)
}

func (self *Stack) Rewind() {
	ln := len(self.mark)
	self.cur = self.mark[ln-1]
	self.mark = self.mark[:ln-1]
}

func (self *Stack) ShiftTopN(n, pos int) {
	top := self.cur - n
	self.cur = pos
	for i := 0; i < n; i++ {
		self.vals[self.cur] = self.vals[top+i]
		self.cur++
	}
}

func (self *Stack) PopMark() int {
	ln := len(self.mark)
	if ln == 0 {
		panic("rewind from empty mark stack, may missing return in some func")
	}
	mark := self.mark[ln-1]
	self.mark = self.mark[:ln-1]
	return mark
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
