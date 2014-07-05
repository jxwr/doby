package symtab

type Symtab struct {
	tab map[string]interface{}
}

func NewSymtab() *Symtab {
	st := &Symtab{map[string]interface{}{}}
	return st
}

func (t *Symtab) Put(name string, obj interface{}) {
	t.tab[name] = obj
}

func (t *Symtab) LookUp(name string) interface{} {
	val, ok := t.tab[name]
	if ok {
		return val
	}
	return nil
}
