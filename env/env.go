package env

type Env struct {
	Outer *Env
	tab   *Symtab
}

func NewEnv(outer *Env) *Env {
	e := &Env{outer, NewSymtab()}
	return e
}

func (e *Env) Put(name string, node interface{}) {
	e.tab.Put(name, node)
}

func (e *Env) LookUp(name string) (interface{}, *Env) {
	env := e
	for env != nil {
		ne := env.tab.LookUp(name)
		if ne != nil {
			return ne, env
		} else {
			env = env.Outer
		}
	}
	return nil, nil
}
