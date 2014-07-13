package rt

/// string

type StringObject struct {
	Property

	Val string
}

func NewStringObject(val string) Object {
	obj := &StringObject{Property(map[string]Object{}), val}
	obj.SetProp("Length", NewBuiltinFuncObject("Length", obj, nil))
	obj.SetProp("Size", NewBuiltinFuncObject("Size", obj, nil))
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

/// methods

func (self *StringObject) Length(ctx *Runtime, args ...Object) (results []Object) {
	ret := NewIntegerObject(len(self.Val))
	results = append(results, ret)
	return
}

func (self *StringObject) Size(ctx *Runtime, args ...Object) (results []Object) {
	ret := NewIntegerObject(len(self.Val))
	results = append(results, ret)
	return
}

/// operators

func (self *StringObject) OP__add__(rt *Runtime, args ...Object) (results []Object) {
	obj := NewStringObject(self.Val + args[0].String())
	results = append(results, obj)
	return
}

func (self *StringObject) OP__add_assign__(rt *Runtime, args ...Object) (results []Object) {
	self.Val += args[0].String()
	return
}

func (self *StringObject) OP__eql__(rt *Runtime, args ...Object) (results []Object) {
	cmp := self.Val == args[0].String()
	results = append(results, NewBoolObject(cmp))
	return
}

func (self *StringObject) OP__neq__(rt *Runtime, args ...Object) (results []Object) {
	cmp := self.Val != args[0].String()
	results = append(results, NewBoolObject(cmp))
	return
}

func (self *StringObject) OP__get_index__(rt *Runtime, args ...Object) (results []Object) {
	idx := args[0].(*IntegerObject)
	ch := string(self.Val[idx.Val])
	obj := NewStringObject(ch)
	results = append(results, obj)
	return
}

func (self *StringObject) OP__slice__(rt *Runtime, args ...Object) (results []Object) {
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
	return
}
