package vm

import (
	//	"fmt"
	"strings"

	"github.com/jxwr/doubi/rt"
	"github.com/jxwr/doubi/vm/instr"
)

type VM struct {
	C     *instr.ClosureProto
	CS    map[int]*instr.ClosureProto
	RT    *rt.Runtime
	Mods  map[string]*rt.DictObject
	frame *rt.Frame
}

func NewVM(c *instr.ClosureProto, cs map[int]*instr.ClosureProto, runtime *rt.Runtime) *VM {
	vm := &VM{C: c, CS: cs, RT: runtime, Mods: map[string]*rt.DictObject{}}
	return vm
}

func (self *VM) Run() {
	obj := self.RT.NewClosureObject(self.C, nil)
	self.RunClosure(obj)
}

func (self *VM) RunClosure(obj *rt.ClosureObject) {
	c := obj.Proto
	f := self.frame
	self.frame = rt.NewFrame(len(c.LocalVariables), len(c.UpvalVariables), obj.Frame)
	for _, instr := range c.Instrs {
		instr.Accept(self)
		if self.frame.NeedReturn {
			break
		}
	}
	self.frame = f
}

func (self *VM) VisitPushClosure(ir *instr.PushClosureInstr) {
	obj := self.RT.NewClosureObject(self.CS[ir.Seq], self.frame)
	self.RT.Push(obj)
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

func (self *VM) VisitLoadUpval(ir *instr.LoadUpvalInstr) {
	depth := (ir.Offset >> 32) & 0xffff
	remoteOffset := (ir.Offset >> 16) & 0xffff

	f := self.frame
	for depth > 0 {
		f = f.Parent
		depth--
	}
	obj := f.Locals[remoteOffset]
	self.RT.Push(obj)
}

func (self *VM) VisitSetLocal(ir *instr.SetLocalInstr) {
	obj := self.RT.Pop()
	self.frame.Locals[ir.Offset] = obj
}

func (self *VM) VisitSetUpval(ir *instr.SetUpvalInstr) {
	depth := (ir.Offset >> 32) & 0xffff
	remoteOffset := (ir.Offset >> 16) & 0xffff

	f := self.frame
	for depth > 0 {
		f = f.Parent
		depth--
	}
	f.Locals[remoteOffset] = self.RT.Pop()
}

func (self *VM) VisitSendMethod(ir *instr.SendMethodInstr) {
	obj := self.RT.Pop()

	if ir.Method == "__call__" {
		switch v := obj.(type) {
		case *rt.ClosureObject:
			self.RunClosure(v)
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

func (self *VM) VisitImport(ir *instr.ImportInstr) {
	mod, _ := self.RT.Env.LookUp(ir.Path)
	xs := strings.Split(ir.Path, "/")
	self.Mods[xs[len(xs)-1]] = mod.(*rt.DictObject)
}

func (self *VM) VisitPushModule(ir *instr.PushModuleInstr) {
	self.RT.Push(self.Mods[ir.Name])
}

func (self *VM) VisitRaiseReturn(ir *instr.RaiseReturnInstr) {
	self.frame.NeedReturn = true
}

func (self *VM) VisitRaiseBreak(ir *instr.RaiseBreakInstr)       {}
func (self *VM) VisitRaiseContinue(ir *instr.RaiseContinueInstr) {}
