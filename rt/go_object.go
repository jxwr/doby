package rt

type GoFuncObject struct {
}

func (self *GoFuncObject) Dispatch(ctx *Runtime, method string, args ...Object) (results []Object) {
	return
}

func (self *GoFuncObject) Name() string {
	return "gofunc"
}

func (self *GoFuncObject) String() string {
	return ""
}

func (self *GoFuncObject) HashCode() string {
	return ""
}
