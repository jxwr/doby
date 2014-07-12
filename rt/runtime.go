package rt

import (
	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/env"
)

type Runtime struct {
	Visitor ast.Visitor
	Env     *env.Env
}

func RegisterGlobals(env *env.Env) {
	env.Put("fmt", nil)
	env.Put("rand", nil)
}

func NewRuntime(visitor ast.Visitor) *Runtime {
	env := env.NewEnv(nil)

	RegisterGlobals(env)

	rt := &Runtime{visitor, env}
	return rt
}
