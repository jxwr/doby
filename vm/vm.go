package vm

import (
	_ "fmt"

	"github.com/jxwr/doubi/comp"
	"github.com/jxwr/doubi/vm/instr"
)

type VM struct {
	C *comp.ClosureProto
}

func (self *VM) Run() {
	for _, instr := range self.C.Instrs {
		instr.Accept(self)
	}
}

func (self *VM) VisitPushNil(ir *instr.PushNilInstr)             {}
func (self *VM) VisitPushTrue(ir *instr.PushTrueInstr)           {}
func (self *VM) VisitPushFalse(ir *instr.PushFalseInstr)         {}
func (self *VM) VisitPushInt(ir *instr.PushIntInstr)             {}
func (self *VM) VisitPushFloat(ir *instr.PushFloatInstr)         {}
func (self *VM) VisitPushString(ir *instr.PushStringInstr)       {}
func (self *VM) VisitLoadLocal(ir *instr.LoadLocalInstr)         {}
func (self *VM) VisitLoadUpval(ir *instr.LoadUpvalInstr)         {}
func (self *VM) VisitSetLocal(ir *instr.SetLocalInstr)           {}
func (self *VM) VisitSetUpval(ir *instr.SetUpvalInstr)           {}
func (self *VM) VisitSendMethod(ir *instr.SendMethodInstr)       {}
func (self *VM) VisitNewArray(ir *instr.NewArrayInstr)           {}
func (self *VM) VisitNewDict(ir *instr.NewDictInstr)             {}
func (self *VM) VisitNewSet(ir *instr.NewSetInstr)               {}
func (self *VM) VisitLabel(ir *instr.LabelInstr)                 {}
func (self *VM) VisitJump(ir *instr.JumpInstr)                   {}
func (self *VM) VisitJumpIfFalse(ir *instr.JumpIfFalseInstr)     {}
func (self *VM) VisitImport(ir *instr.ImportInstr)               {}
func (self *VM) VisitPushModule(ir *instr.PushModuleInstr)       {}
func (self *VM) VisitPushClosure(ir *instr.PushClosureInstr)     {}
func (self *VM) VisitRaiseReturn(ir *instr.RaiseReturnInstr)     {}
func (self *VM) VisitRaiseBreak(ir *instr.RaiseBreakInstr)       {}
func (self *VM) VisitRaiseContinue(ir *instr.RaiseContinueInstr) {}
