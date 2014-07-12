package rt

/// string

type StringObject struct {
	Property

	Val string
}

func NewStringObject(val string) Object {
	obj := &StringObject{Property(map[string]Object{}), val}
	return obj
}

func (self *StringObject) Name() string {
	return "string"
}

func (self *StringObject) HashCode() string {
	return self.String()
}

func (self *StringObject) String() string {
	return self.Val
}

func (self *StringObject) Dispatch(ctx *Runtime, method string, args ...Object) (results []Object) {
	var is bool
	if is, results = self.AccessPropMethod(method, args...); is {
		return
	}

	switch method {
	case "__add__":
		obj := NewStringObject(self.Val + args[0].String())
		results = append(results, obj)
	case "__+=__":
		self.Val += args[0].String()
	}
	return
}
