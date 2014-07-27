package vm

import (
	"fmt"

	"github.com/jxwr/doubi/rt"
	"github.com/jxwr/doubi/vm/instr"
)

type Frame struct {
	Locals     []rt.Object
	Upvals     []rt.Object
	Parent     *Frame
	NeedReturn bool
}

func NewFrame(numLocals, numUpvals int, parent *Frame) *Frame {
	frame := &Frame{
		make([]rt.Object, numLocals),
		make([]rt.Object, numUpvals),
		parent,
		false,
	}
	return frame
}

type VM struct {
	C     *instr.ClosureProto
	CS    map[int]*instr.ClosureProto
	RT    *rt.Runtime
	frame *Frame
}

func NewVM(c *instr.ClosureProto, cs map[int]*instr.ClosureProto, runtime *rt.Runtime) *VM {
	vm := &VM{C: c, CS: cs, RT: runtime}
	return vm
}

func (self *VM) Run() {
	self.RunClosure(self.C)
}

func (self *VM) RunClosure(c *instr.ClosureProto) {
	self.frame = NewFrame(len(c.LocalVariables), len(c.UpvalVariables), self.frame)
	for _, instr := range c.Instrs {
		if self.frame.NeedReturn {
			break
		}
		instr.Accept(self)
	}
	self.frame = self.frame.Parent
}

func (self *VM) VisitPushNil(ir *instr.PushNilInstr) {
	self.RT.Push(self.RT.Nil)
}

func (self *VM) VisitPushTrue(ir *instr.PushTrueInstr) {
	self.RT.Push(self.RT.True)
}

func (self *VM) VisitPushFalse(ir *instr.PushFalseInstr) {
	self.RT.Push(self.RT.False)
}

func (self *VM) VisitPushInt(ir *instr.PushIntInstr) {
	obj := self.RT.NewIntegerObject(ir.Val)
	self.RT.Push(obj)
}

func (self *VM) VisitPushFloat(ir *instr.PushFloatInstr) {
	obj := self.RT.NewFloatObject(ir.Val)
	self.RT.Push(obj)
}

func (self *VM) VisitPushString(ir *instr.PushStringInstr) {
	obj := self.RT.NewStringObject(ir.Val)
	self.RT.Push(obj)
}

func (self *VM) VisitLoadLocal(ir *instr.LoadLocalInstr) {
	obj := self.frame.Locals[ir.Offset]
	self.RT.Push(obj)
}

func (self *VM) VisitLoadUpval(ir *instr.LoadUpvalInstr) {}

func (self *VM) VisitSetLocal(ir *instr.SetLocalInstr) {
	obj := self.RT.Pop()
	self.frame.Locals[ir.Offset] = obj
	fmt.Println(obj)
}

func (self *VM) VisitSetUpval(ir *instr.SetUpvalInstr) {}

func (self *VM) VisitSendMethod(ir *instr.SendMethodInstr) {
	obj := self.RT.Pop()

	if ir.Method == "__call__" {
		switch v := obj.(type) {
		case *rt.ClosureObject:
			self.RunClosure(v.Proto)
		case *rt.GoFuncObject:
			args := make([]rt.Object, ir.Num)
			for i := ir.Num - 1; i >= 0; i-- {
				args[i] = self.RT.Pop()
			}
			rets := v.CallGoFunc(self.RT, args...)
			for _, ret := range rets {
				self.RT.Push(ret)
			}
		}
	} else {
		args := make([]rt.Object, ir.Num)
		for i := ir.Num - 1; i >= 0; i-- {
			args[i] = self.RT.Pop()
		}
		rets := rt.Invoke(self.RT, obj, ir.Method, args...)
		for _, ret := range rets {
			self.RT.Push(ret)
		}
	}
}

func (self *VM) VisitNewArray(ir *instr.NewArrayInstr) {
	elems := make([]rt.Object, ir.Num)
	for i := ir.Num - 1; i >= 0; i-- {
		elems[i] = self.RT.Pop()
	}
	obj := self.RT.NewArrayObject(elems)
	self.RT.Push(obj)
}

func (self *VM) VisitNewDict(ir *instr.NewDictInstr)         {}
func (self *VM) VisitNewSet(ir *instr.NewSetInstr)           {}
func (self *VM) VisitLabel(ir *instr.LabelInstr)             {}
func (self *VM) VisitJump(ir *instr.JumpInstr)               {}
func (self *VM) VisitJumpIfFalse(ir *instr.JumpIfFalseInstr) {}
func (self *VM) VisitImport(ir *instr.ImportInstr)           {}
func (self *VM) VisitPushModule(ir *instr.PushModuleInstr)   {}

func (self *VM) VisitPushClosure(ir *instr.PushClosureInstr) {
	obj := self.RT.NewClosureObject(self.CS[ir.Seq])
	self.RT.Push(obj)
}

func (self *VM) VisitRaiseReturn(ir *instr.RaiseReturnInstr) {
	self.frame.NeedReturn = true
}
func (self *VM) VisitRaiseBreak(ir *instr.RaiseBreakInstr)       {}
func (self *VM) VisitRaiseContinue(ir *instr.RaiseContinueInstr) {}
