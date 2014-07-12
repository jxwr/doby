package rt

type Object interface {
	Dispatch(ctx *Runtime, method string, args ...Object) []Object
	Name() string
	String() string
	HashCode() string
}

type Property map[string]Object

func (self *Property) SetProp(key string, val Object) {
	(*self)[key] = val
}

func (self *Property) GetProp(key string) Object {
	return (*self)[key]
}

func (self *Property) AccessPropMethod(method string, args ...Object) (isPropMethod bool, results []Object) {
	if method == "__get_property__" {
		idx := args[0].(*StringObject)
		results = append(results, self.GetProp(idx.Val))
		isPropMethod = true
	} else if method == "__set_property__" {
		idx := args[0].(*StringObject)
		val := args[1]
		self.SetProp(idx.Val, val)
		isPropMethod = true
	}
	return
}

type NilObject struct {
}

func (self *NilObject) Name() string {
	return "nil"
}

func (self *NilObject) HashCode() string {
	return "nil"
}

func (self *NilObject) String() string {
	return "<nil>"
}

func (self *NilObject) Dispatch(ctx *Runtime, method string, args ...Object) (results []Object) {
	return
}
