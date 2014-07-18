package instr

import (
	"fmt"
)

type Instr interface {
	String() string
}

type SendMethod struct {
	Method string
	NumArg int
}

func (self *SendMethod) String() string {
	return fmt.Sprintf("SEND_METHOD %s %d", self.Method, self.NumArg)
}

func MkSendMethod(method string, numArg int) *SendMethod {
	m := &SendMethod{method, numArg}
	return m
}

const (
	ILLEGAL = iota

	NOOP
	PUSH_NIL
	IS_NIL
	PUSH_TRUE
	PUSH_FALSE
	PUSH_INT
	GOTO
	GOTO_IF_FALSE
	GOTO_IF_TRUE
	RET
	SWAP_STACK
	DUP_TOP
	DUP_MANY
	POP
	POP_MANY
	ROTATE
	SET_LOCAL
	PUSH_LOCAL
	RAISE_BREAK
	MAKE_ARRAY
	PUSH_CONST
	SET_CONST
	SEND_METHOD
	SEND_STACK
)
