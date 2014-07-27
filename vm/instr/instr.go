package instr

import (
	"fmt"
)

type InstrType int

const (
	PUSH_NIL InstrType = iota
	PUSH_TRUE
	PUSH_FALSE
	PUSH_INT
	PUSH_STRING
	PUSH_FLOAT
	LOAD_LOCAL
	LOAD_UPVAL
	SET_LOCAL
	SET_UPVAL
	SEND_METHOD
	NEW_ARRAY
	NEW_DICT
	NEW_SET
	LABEL
	JUMP
	JUMP_IF_FALSE
	IMPORT
	PUSH_MODULE
	PUSH_CLOSURE
	RAISE_RETURN
	RAISE_BREAK
	RAISE_CONTINUE
)

var TypName = map[InstrType]string{
	PUSH_NIL:       "PUSH_NIL",
	PUSH_TRUE:      "PUSH_TRUE",
	PUSH_FALSE:     "PUSH_FALSE",
	PUSH_INT:       "PUSH_INT",
	PUSH_STRING:    "PUSH_STRING",
	PUSH_FLOAT:     "PUSH_FLOAT",
	LOAD_LOCAL:     "LOAD_LOCAL",
	LOAD_UPVAL:     "LOAD_UPVAL",
	SET_LOCAL:      "SET_LOCAL",
	SET_UPVAL:      "SET_UPVAL",
	SEND_METHOD:    "SEND_METHOD",
	NEW_ARRAY:      "NEW_ARRAY",
	NEW_DICT:       "NEW_DICT",
	NEW_SET:        "NEW_SET",
	LABEL:          "LABEL",
	JUMP:           "JUMP",
	JUMP_IF_FALSE:  "JUMP_IF_FALSE",
	IMPORT:         "IMPORT",
	PUSH_MODULE:    "PUSH_MODULE",
	PUSH_CLOSURE:   "PUSH_CLOSURE",
	RAISE_RETURN:   "RAISE_RETURN",
	RAISE_BREAK:    "RAISE_BREAK",
	RAISE_CONTINUE: "RAISE_CONTINUE",
}

type Instr interface {
	Type() InstrType
	String() string
}

type PushNilInstr struct {
	Typ InstrType
}

func PushNil() *PushNilInstr {
	instr := &PushNilInstr{PUSH_NIL}
	return instr
}

type PushTrueInstr struct {
	Typ InstrType
}

func PushTrue() *PushTrueInstr {
	instr := &PushTrueInstr{PUSH_TRUE}
	return instr
}

type PushFalseInstr struct {
	Typ InstrType
}

func PushFalse() *PushFalseInstr {
	instr := &PushFalseInstr{PUSH_FALSE}
	return instr
}

type PushIntInstr struct {
	Typ InstrType
	Val int
}

func PushInt(val int) *PushIntInstr {
	instr := &PushIntInstr{PUSH_INT, val}
	return instr
}

type PushFloatInstr struct {
	Typ InstrType
	Val float64
}

func PushFloat(val float64) *PushFloatInstr {
	instr := &PushFloatInstr{PUSH_FLOAT, val}
	return instr
}

type PushStringInstr struct {
	Typ InstrType
	Val string
}

func PushString(val string) *PushStringInstr {
	instr := &PushStringInstr{PUSH_STRING, val}
	return instr
}

type LoadLocalInstr struct {
	Typ    InstrType
	Offset int
}

func LoadLocal(offset int) *LoadLocalInstr {
	instr := &LoadLocalInstr{LOAD_LOCAL, offset}
	return instr
}

type LoadUpvalInstr struct {
	Typ    InstrType
	Offset int
}

func LoadUpval(offset int) *LoadUpvalInstr {
	instr := &LoadUpvalInstr{LOAD_UPVAL, offset}
	return instr
}

type SetLocalInstr struct {
	Typ    InstrType
	Offset int
}

func SetLocal(offset int) *SetLocalInstr {
	instr := &SetLocalInstr{SET_LOCAL, offset}
	return instr
}

type SetUpvalInstr struct {
	Typ    InstrType
	Offset int
}

func SetUpval(offset int) *SetUpvalInstr {
	instr := &SetUpvalInstr{SET_UPVAL, offset}
	return instr
}

type SendMethodInstr struct {
	Typ    InstrType
	Method string
	Offset int
}

func SendMethod(method string, offset int) *SendMethodInstr {
	instr := &SendMethodInstr{SEND_METHOD, method, offset}
	return instr
}

type NewArrayInstr struct {
	Typ InstrType
	Num int
}

func NewArray(num int) *NewArrayInstr {
	instr := &NewArrayInstr{NEW_ARRAY, num}
	return instr
}

type NewDictInstr struct {
	Typ InstrType
	Num int
}

func NewDict(num int) *NewDictInstr {
	instr := &NewDictInstr{NEW_DICT, num}
	return instr
}

type NewSetInstr struct {
	Typ InstrType
	Num int
}

func NewSet(num int) *NewSetInstr {
	instr := &NewSetInstr{NEW_SET, num}
	return instr
}

type LabelInstr struct {
	Typ   InstrType
	Label string
}

func Label(label string) *LabelInstr {
	instr := &LabelInstr{LABEL, label}
	return instr
}

type JumpInstr struct {
	Typ InstrType
	Pos int
}

func Jump(pos int) *JumpInstr {
	instr := &JumpInstr{JUMP, pos}
	return instr
}

type JumpIfFalseInstr struct {
	Typ InstrType
	Pos int
}

func JumpIfFalse(pos int) *JumpIfFalseInstr {
	instr := &JumpIfFalseInstr{JUMP_IF_FALSE, pos}
	return instr
}

type ImportInstr struct {
	Typ  InstrType
	Path string
}

func Import(path string) *ImportInstr {
	instr := &ImportInstr{IMPORT, path}
	return instr
}

type PushModuleInstr struct {
	Typ  InstrType
	Name string
}

func PushModule(name string) *PushModuleInstr {
	instr := &PushModuleInstr{PUSH_MODULE, name}
	return instr
}

type PushClosureInstr struct {
	Typ InstrType
	Seq int
}

func PushClosure(seq int) *PushClosureInstr {
	instr := &PushClosureInstr{PUSH_CLOSURE, seq}
	return instr
}

type RaiseReturnInstr struct {
	Typ InstrType
	Num int
}

func RaiseReturn(num int) *RaiseReturnInstr {
	instr := &RaiseReturnInstr{RAISE_RETURN, num}
	return instr
}

type RaiseBreakInstr struct {
	Typ InstrType
}

func RaiseBreak() *RaiseBreakInstr {
	instr := &RaiseBreakInstr{RAISE_BREAK}
	return instr
}

type RaiseContinueInstr struct {
	Typ InstrType
}

func RaiseContinue() *RaiseContinueInstr {
	instr := &RaiseContinueInstr{RAISE_CONTINUE}
	return instr
}

var _t = func(args ...interface{}) string {
	s := ""
	for _, arg := range args {
		s += fmt.Sprint(arg, " ")
	}
	return s
}

func (n *PushNilInstr) String() string       { return TypName[n.Typ] }
func (n *PushTrueInstr) String() string      { return TypName[n.Typ] }
func (n *PushFalseInstr) String() string     { return TypName[n.Typ] }
func (n *PushIntInstr) String() string       { return _t(TypName[n.Typ], n.Val) }
func (n *PushFloatInstr) String() string     { return _t(TypName[n.Typ], n.Val) }
func (n *PushStringInstr) String() string    { return _t(TypName[n.Typ], n.Val) }
func (n *LoadLocalInstr) String() string     { return _t(TypName[n.Typ], n.Offset) }
func (n *LoadUpvalInstr) String() string     { return _t(TypName[n.Typ], n.Offset) }
func (n *SetLocalInstr) String() string      { return _t(TypName[n.Typ], n.Offset) }
func (n *SetUpvalInstr) String() string      { return _t(TypName[n.Typ], n.Offset) }
func (n *SendMethodInstr) String() string    { return _t(TypName[n.Typ], n.Method, n.Offset) }
func (n *NewArrayInstr) String() string      { return _t(TypName[n.Typ], n.Num) }
func (n *NewDictInstr) String() string       { return _t(TypName[n.Typ], n.Num) }
func (n *NewSetInstr) String() string        { return _t(TypName[n.Typ], n.Num) }
func (n *LabelInstr) String() string         { return _t(TypName[n.Typ], n.Label) }
func (n *JumpInstr) String() string          { return _t(TypName[n.Typ], n.Pos) }
func (n *JumpIfFalseInstr) String() string   { return _t(TypName[n.Typ], n.Pos) }
func (n *ImportInstr) String() string        { return _t(TypName[n.Typ], n.Path) }
func (n *PushModuleInstr) String() string    { return _t(TypName[n.Typ], n.Name) }
func (n *PushClosureInstr) String() string   { return _t(TypName[n.Typ], n.Seq) }
func (n *RaiseReturnInstr) String() string   { return TypName[n.Typ] }
func (n *RaiseBreakInstr) String() string    { return TypName[n.Typ] }
func (n *RaiseContinueInstr) String() string { return TypName[n.Typ] }

func (n *PushNilInstr) Type() InstrType       { return n.Typ }
func (n *PushTrueInstr) Type() InstrType      { return n.Typ }
func (n *PushFalseInstr) Type() InstrType     { return n.Typ }
func (n *PushIntInstr) Type() InstrType       { return n.Typ }
func (n *PushFloatInstr) Type() InstrType     { return n.Typ }
func (n *PushStringInstr) Type() InstrType    { return n.Typ }
func (n *LoadLocalInstr) Type() InstrType     { return n.Typ }
func (n *LoadUpvalInstr) Type() InstrType     { return n.Typ }
func (n *SetLocalInstr) Type() InstrType      { return n.Typ }
func (n *SetUpvalInstr) Type() InstrType      { return n.Typ }
func (n *SendMethodInstr) Type() InstrType    { return n.Typ }
func (n *NewArrayInstr) Type() InstrType      { return n.Typ }
func (n *NewDictInstr) Type() InstrType       { return n.Typ }
func (n *NewSetInstr) Type() InstrType        { return n.Typ }
func (n *LabelInstr) Type() InstrType         { return n.Typ }
func (n *JumpInstr) Type() InstrType          { return n.Typ }
func (n *JumpIfFalseInstr) Type() InstrType   { return n.Typ }
func (n *ImportInstr) Type() InstrType        { return n.Typ }
func (n *PushModuleInstr) Type() InstrType    { return n.Typ }
func (n *PushClosureInstr) Type() InstrType   { return n.Typ }
func (n *RaiseReturnInstr) Type() InstrType   { return n.Typ }
func (n *RaiseBreakInstr) Type() InstrType    { return n.Typ }
func (n *RaiseContinueInstr) Type() InstrType { return n.Typ }
