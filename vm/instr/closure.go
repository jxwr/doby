package instr

import (
	"fmt"
)

type ClosureProto struct {
	localVariables     map[string]int
	localOffset        int
	upvalVariables     map[string]int
	upvalOffset        int
	innerClosureProtos []*ClosureProto
	outerClosureProto  *ClosureProto
	instrs             []Instr
	args               []string
	seq                int
}

func NewClosureProto(outer *ClosureProto) *ClosureProto {
	c := &ClosureProto{
		localVariables:     map[string]int{},
		upvalVariables:     map[string]int{},
		innerClosureProtos: []*ClosureProto{},
		outerClosureProto:  outer,
		instrs:             []Instr{},
		args:               []string{},
		seq:                closure_seq,
	}
	closure_seq++
	return c
}

func (self *ClosureProto) Seq() int {
	return self.seq
}

func (self *ClosureProto) OuterClosureProto() *ClosureProto {
	return self.outerClosureProto
}

func (self *ClosureProto) LocalVariables() map[string]int {
	return self.localVariables
}

func (self *ClosureProto) NumLocalVariable() int {
	return len(self.localVariables)
}

func (self *ClosureProto) UpvalVariables() map[string]int {
	return self.upvalVariables
}

func (self *ClosureProto) NumUpvalVariable() int {
	return len(self.upvalVariables)
}

func (self *ClosureProto) Instrs() []Instr {
	return self.instrs
}

func (self *ClosureProto) Emit(instr Instr) int {
	self.instrs = append(self.instrs, instr)
	return len(self.instrs) - 1
}

func (self *ClosureProto) DumpClosureProto() {
	fmt.Println()
	fmt.Printf("CLOSURE seq:%d local:%d upval:%d\n", self.seq,
		len(self.localVariables), len(self.upvalVariables))
	for k, v := range self.localVariables {
		fmt.Printf("  .local %d %s\n", v, k)
	}
	for k, v := range self.upvalVariables {
		fmt.Printf("  .upval %d %d %d %s\n", v&0xffff, (v>>32)&0xffff, (v>>16)&0xffff, k)
	}

	fmt.Println("CODE:")
	for i, instr := range self.instrs {
		fmt.Printf("%3d: %s\n", i, instr.String())
	}
	for _, ic := range self.innerClosureProtos {
		ic.DumpClosureProto()
	}
}

func (self *ClosureProto) AddClosureProto(c *ClosureProto) int {
	self.innerClosureProtos = append(self.innerClosureProtos, c)
	return len(self.innerClosureProtos)
}

func (self *ClosureProto) AddLocalVariable(name string) (offset int) {
	exist, offset := self.LookUpLocal(name)
	if !exist {
		offset = self.localOffset
		self.localVariables[name] = offset
		self.localOffset++
	}
	return
}

func (self *ClosureProto) LookUpLocal(name string) (exist bool, offset int) {
	offset, ok := self.localVariables[name]
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
		self.upvalVariables[name] = pos
		self.upvalOffset++
	}
	return
}

func (self *ClosureProto) LookUpUpval(name string) (exist bool, offset int) {
	offset, ok := self.upvalVariables[name]
	if ok {
		exist = true
		return
	}
	return
}

func (self *ClosureProto) LookUpOuter(name string) (exist bool, depth int, offset int) {
	c := self.outerClosureProto
	depth = 1

	for c != nil {
		of, ok := c.localVariables[name]
		if ok {
			offset = of
			exist = true
			return
		}
		depth++
		c = c.outerClosureProto
	}

	depth = -1
	exist = false
	return
}

var closure_seq int = 0
