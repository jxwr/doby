package rt

type Frame struct {
	Locals       []Object
	Upvals       []Object
	Parent       *Frame
	JumpTarget   int
	NeedReturn   bool
	NeedBreak    bool
	NeedContinue bool
	blockTargets []int
}

func NewFrame(numLocals, numUpvals int, parent *Frame) *Frame {
	frame := &Frame{
		make([]Object, numLocals),
		make([]Object, numUpvals),
		parent,
		-1,
		false,
		false,
		false,
		[]int{},
	}
	return frame
}

func (self *Frame) PushBlock(pc int) {
	self.blockTargets = append(self.blockTargets, pc)
}

func (self *Frame) PopBlock() (target int) {
	target = self.blockTargets[len(self.blockTargets)-1]
	self.blockTargets = self.blockTargets[:len(self.blockTargets)-1]
	return
}

func (self *Frame) BlockEndPc() int {
	if len(self.blockTargets) == 0 {
		return -1
	}
	return self.blockTargets[len(self.blockTargets)-1]
}
