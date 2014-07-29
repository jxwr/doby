package instr

import (
	"fmt"
)

type ClosureProto struct {
	LocalVariables map[string]int
	localOffset    int

	UpvalVariables map[string]int
	upvalOffset    int

	InnerClosureProtos []*ClosureProto
	OuterClosureProto  *ClosureProto
	Instrs             []Instr
	Args               []string
	Seq                int
}

func NewClosureProto(outer *ClosureProto) *ClosureProto {
	c := &ClosureProto{
		LocalVariables:     map[string]int{},
		UpvalVariables:     map[string]int{},
		InnerClosureProtos: []*ClosureProto{},
		OuterClosureProto:  outer,
		Instrs:             []Instr{},
		Args:               []string{},
		Seq:                closure_seq,
	}
	closure_seq++
	return c
}

func (self *ClosureProto) Emit(instr Instr) {
	self.Instrs = append(self.Instrs, instr)
}

func (self *ClosureProto) DumpClosureProto() {
	fmt.Println()
	fmt.Printf("CLOSURE seq:%d local:%d upval:%d\n", self.Seq,
		len(self.LocalVariables), len(self.UpvalVariables))
	for k, v := range self.LocalVariables {
		fmt.Printf("  .local %d %s\n", v, k)
	}
	for k, v := range self.UpvalVariables {
		fmt.Printf("  .upval %d %d %d %s\n", v&0xffff, (v>>32)&0xffff, (v>>16)&0xffff, k)
	}

	fmt.Println("CODE:")
	for i, instr := range self.Instrs {
		fmt.Printf("%3d: %s\n", i, instr.String())
	}
	for _, ic := range self.InnerClosureProtos {
		ic.DumpClosureProto()
	}
}

func (self *ClosureProto) AddClosureProto(c *ClosureProto) int {
	self.InnerClosureProtos = append(self.InnerClosureProtos, c)
	return len(self.InnerClosureProtos)
}

func (self *ClosureProto) AddLocalVariable(name string) (offset int) {
	exist, _ := self.LookUpLocal(name)
	if !exist {
		offset = self.localOffset
		self.LocalVariables[name] = offset
		self.localOffset++
	}
	return
}

func (self *ClosureProto) LookUpLocal(name string) (exist bool, offset int) {
	offset, ok := self.LocalVariables[name]
	if ok {
		exist = true
		return
	}
	return
}

func (self *ClosureProto) AddUpvalVariable(name string, depth, remoteOffset int) (pos int) {
	exist, _ := self.LookUpLocal(name)
	if !exist {
		offset := self.upvalOffset
		pos = offset + (depth << 32) + (remoteOffset << 16)
		self.UpvalVariables[name] = pos
		self.upvalOffset++
	}
	return
}

func (self *ClosureProto) LookUpUpval(name string) (exist bool, offset int) {
	offset, ok := self.UpvalVariables[name]
	if ok {
		//		offset &= 0xffff
		exist = true
		return
	}
	return
}

func (self *ClosureProto) LookUpOuter(name string) (exist bool, depth int, offset int) {
	c := self.OuterClosureProto
	depth = 1

	for c != nil {
		of, ok := c.LocalVariables[name]
		if ok {
			offset = of
			exist = true
			return
		}
		depth++
		c = c.OuterClosureProto
	}

	depth = -1
	exist = false
	return
}

var closure_seq int = 0
