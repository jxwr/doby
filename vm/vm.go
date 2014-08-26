package vm

import (
	"fmt"
	"strings"

	"github.com/jxwr/doby/rt"
	"github.com/jxwr/doby/vm/instr"
)

type VM struct {
	cc      *instr.ClosureProto
	cs      map[int]*instr.ClosureProto
	mods    map[string]*rt.DictObject
	frame   *rt.Frame
	runtime *rt.Runtime
}

func NewVM(c *instr.ClosureProto, cs map[int]*instr.ClosureProto, runtime *rt.Runtime) *VM {
	vm := &VM{cc: c, cs: cs, runtime: runtime, mods: map[string]*rt.DictObject{}}
	return vm
}

func (self *VM) Run() {
	obj := self.runtime.NewClosureObject(self.cc, nil)
	self.RunClosure(obj)
}

func (self *VM) RunClosure(obj *rt.ClosureObject) {
	c := obj.Proto
	f := self.frame

	self.frame = rt.NewFrame(c.NumLocalVariable(), c.NumUpvalVariable(), obj.Frame)
	instrs := c.Instrs()
	for i := 0; i < len(instrs); i++ {
		instrs[i].Accept(self)
		if self.frame.JumpTarget > 0 {
			i = self.frame.JumpTarget - 1
			self.frame.JumpTarget = -1
		}
		if self.frame.NeedBreak {
			pc := self.frame.BlockEndPc()
			if pc < 0 {
				panic("wrong break stmt")
			}
			i = pc - 1
			self.frame.NeedBreak = false
		}
		if self.frame.NeedReturn {
			self.frame.NeedReturn = false
			break
		}
	}
	self.frame = f
}

func (self *VM) VisitPushClosure(ir *instr.PushClosureInstr) {
	obj := self.runtime.NewClosureObject(self.cs[ir.Seq], self.frame)
	self.runtime.Push(obj)
}

func (self *VM) VisitPushNil(ir *instr.PushNilInstr) {
	self.runtime.Push(self.runtime.Nil)
}

func (self *VM) VisitPushTrue(ir *instr.PushTrueInstr) {
	self.runtime.Push(self.runtime.True)
}

func (self *VM) VisitPushFalse(ir *instr.PushFalseInstr) {
	self.runtime.Push(self.runtime.False)
}

func (self *VM) VisitPushInt(ir *instr.PushIntInstr) {
	obj := self.runtime.NewIntegerObject(ir.Val)
	self.runtime.Push(obj)
}

func (self *VM) VisitPushFloat(ir *instr.PushFloatInstr) {
	obj := self.runtime.NewFloatObject(ir.Val)
	self.runtime.Push(obj)
}

func (self *VM) VisitPushString(ir *instr.PushStringInstr) {
	obj := self.runtime.NewStringObject(ir.Val)
	self.runtime.Push(obj)
}

func (self *VM) VisitLoadLocal(ir *instr.LoadLocalInstr) {
	obj := self.frame.Locals[ir.Offset]
	self.runtime.Push(obj)
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
	self.runtime.Push(obj)
}

func (self *VM) VisitSetLocal(ir *instr.SetLocalInstr) {
	obj := self.runtime.Pop()
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
	f.Locals[remoteOffset] = self.runtime.Pop()
}

func (self *VM) VisitSendMethod(ir *instr.SendMethodInstr) {
	obj := self.runtime.Pop()

	// closure object is a function defined in doby code, mark stack and rewind manually
	// gofunc object is a function defined in go lib, no way to rewind the stack here
	// func object is a function of a builtin object, mark stack by CallFuncObj(rt/runtime.go)

	if ir.Method == "__call__" {
		switch v := obj.(type) {
		case *rt.ClosureObject:
			// take care of the stack
			self.runtime.MarkN(-(ir.Num))
			self.RunClosure(v)
		case *rt.GoFuncObject:
			args := make([]rt.Object, ir.Num)
			for i := ir.Num - 1; i >= 0; i-- {
				args[i] = self.runtime.Pop()
			}
			rets := v.CallGoFunc(self.runtime, args...)
			for _, ret := range rets {
				self.runtime.Push(ret)
			}
		case *rt.FuncObject:
			// FIXME: For now, func object is the method object of the builtin object
			args := make([]rt.Object, ir.Num)
			for i := ir.Num - 1; i >= 0; i-- {
				args[i] = self.runtime.Pop()
			}
			rets := rt.Invoke(self.runtime, v, "__call__", args...)
			for _, ret := range rets {
				self.runtime.Push(ret)
			}
		default:
			fmt.Printf("%T", v)
		}
	} else {
		args := make([]rt.Object, ir.Num)
		for i := ir.Num - 1; i >= 0; i-- {
			args[i] = self.runtime.Pop()
		}
		rets := rt.Invoke(self.runtime, obj, ir.Method, args...)
		for _, ret := range rets {
			self.runtime.Push(ret)
		}
	}
}

func (self *VM) VisitNewArray(ir *instr.NewArrayInstr) {
	elems := make([]rt.Object, ir.Num)
	for i := ir.Num - 1; i >= 0; i-- {
		elems[i] = self.runtime.Pop()
	}
	obj := self.runtime.NewArrayObject(elems)
	self.runtime.Push(obj)
}

func (self *VM) VisitNewDict(ir *instr.NewDictInstr) {
	fieldMap := map[string]rt.Slot{}

	for i := 0; i < ir.Num; i++ {
		val := self.runtime.Pop()
		key := self.runtime.Pop()
		fieldMap[key.HashCode()] = rt.Slot{key, val}
	}
	obj := self.runtime.NewDictObject(fieldMap)
	self.runtime.Push(obj)
}

func (self *VM) VisitNewSet(ir *instr.NewSetInstr) {
	elems := make([]rt.Object, ir.Num)
	for i := ir.Num - 1; i >= 0; i-- {
		elems[i] = self.runtime.Pop()
	}
	obj := self.runtime.NewSetObject(elems)
	self.runtime.Push(obj)
}

func (self *VM) VisitLabel(ir *instr.LabelInstr) {}

func (self *VM) VisitJump(ir *instr.JumpInstr) {
	self.frame.JumpTarget = ir.Target
}

func (self *VM) VisitJumpIfFalse(ir *instr.JumpIfFalseInstr) {
	obj := self.runtime.Pop().(*rt.BoolObject)

	if !obj.Val {
		self.frame.JumpTarget = ir.Target
	}
}

func (self *VM) VisitImport(ir *instr.ImportInstr) {
	mod, _ := self.runtime.Env.LookUp(ir.Path)
	xs := strings.Split(ir.Path, "/")
	self.mods[xs[len(xs)-1]] = mod.(*rt.DictObject)
}

func (self *VM) VisitPushModule(ir *instr.PushModuleInstr) {
	self.runtime.Push(self.mods[ir.Name])
}

func (self *VM) VisitPushBlock(ir *instr.PushBlockInstr) {
	self.frame.PushBlock(ir.Target)
}

func (self *VM) VisitPopBlock(ir *instr.PopBlockInstr) {
	self.frame.PopBlock()
}

func (self *VM) VisitRaiseReturn(ir *instr.RaiseReturnInstr) {
	// FIXME: THIS IS VERY INEFFICIENCY
	// We need to do this explicitly to clean up the garbage values left in the stack,
	// should check out JRuby or some other interpreters to see how to resolve this
	// problem.
	self.runtime.ShiftTopN(ir.Num, self.runtime.PopMark())
	self.frame.NeedReturn = true
}

func (self *VM) VisitRaiseBreak(ir *instr.RaiseBreakInstr) {
	self.frame.NeedBreak = true
}

func (self *VM) VisitRaiseContinue(ir *instr.RaiseContinueInstr) {
	self.frame.NeedContinue = true
}
