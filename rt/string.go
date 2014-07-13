package rt

/// string

type StringObject struct {
	Property

	Val string
}

func NewStringObject(val string) Object {
	obj := &StringObject{Property(map[string]Object{}), val}
	obj.SetProp("length", NewBuiltinFuncObject("length", obj, nil))
	obj.SetProp("size", NewBuiltinFuncObject("size", obj, nil))
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
	switch method {
	case "__add__":
		obj := NewStringObject(self.Val + args[0].String())
		results = append(results, obj)
	case "__+=__":
		self.Val += args[0].String()
	case "__eql__":
		cmp := self.Val == args[0].String()
		results = append(results, NewBoolObject(cmp))
	case "length", "size":
		ret := NewIntegerObject(len(self.Val))
		results = append(results, ret)
	case "__get_index__":
		idx := args[0].(*IntegerObject)
		ch := string(self.Val[idx.Val])
		obj := NewStringObject(ch)
		results = append(results, obj)
	case "__slice__":
		low := 0
		high := len(self.Val)

		lo := args[0]
		if lo != nil {
			low = lo.(*IntegerObject).Val
		}
		ho := args[1]
		if ho != nil {
			high = ho.(*IntegerObject).Val
		}

		ret := NewStringObject(self.Val[low:high])
		results = append(results, ret)
	}
	return
}
