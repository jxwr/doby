package rt

import (
	"fmt"
	"reflect"
)

/// go object wrapper

type GoObject struct {
	obj interface{}
}

func NewGoObject(obj interface{}) *GoObject {
	gobj := &GoObject{obj}
	return gobj
}

func (self *GoObject) Dispatch(ctx *Runtime, method string, args ...Object) (results []Object) {
	return
}

func (self *GoObject) Name() string {
	return "goobj"
}

func (self *GoObject) String() string {
	return fmt.Sprint(self.obj)
}

func (self *GoObject) HashCode() string {
	return fmt.Sprintf("%p", self.obj)
}

/// function

type GoFuncObject struct {
	name string
	typ  reflect.Type
	fn   interface{}
}

func NewGoFuncObject(fname string, fn interface{}) *GoFuncObject {
	gf := &GoFuncObject{fname, reflect.TypeOf(fn), fn}
	return gf
}

func (self *GoFuncObject) callGoFunc(ctx *Runtime, args ...Object) (results []Object) {
	inNum := self.typ.NumIn()
	inArgs := []reflect.Value{}

	dummy := []interface{}{}
	for i := 0; i < inNum; i++ {
		arg := args[i]
		if i == inNum-1 && self.typ.In(i) == reflect.TypeOf(dummy) {
			for j := i; j < len(args); j++ {
				arg := args[j]
				switch arg := arg.(type) {
				case *IntegerObject:
					v := reflect.ValueOf(arg.Val)
					inArgs = append(inArgs, v)
				case *FloatObject:
					v := reflect.ValueOf(arg.Val)
					inArgs = append(inArgs, v)
				case *StringObject:
					v := reflect.ValueOf(arg.Val)
					inArgs = append(inArgs, v)
				case *BoolObject:
					v := reflect.ValueOf(arg.Val)
					inArgs = append(inArgs, v)
				case *GoObject:
					if arg.obj == nil {
						var nilObj *NilObject
						inArgs = append(inArgs, reflect.ValueOf(nilObj))
					} else {
						v := reflect.ValueOf(arg.obj)
						inArgs = append(inArgs, v)
					}
				default:
					inArgs = append(inArgs, reflect.ValueOf(arg))
				}
			}
		} else {
			switch arg := arg.(type) {
			case *IntegerObject:
				v := reflect.ValueOf(arg.Val)
				t := reflect.TypeOf(arg.Val)
				if t.ConvertibleTo(self.typ.In(i)) {
					v = v.Convert(self.typ.In(i))
				}
				inArgs = append(inArgs, v)
			case *FloatObject:
				v := reflect.ValueOf(arg.Val)
				t := reflect.TypeOf(arg.Val)
				if t.ConvertibleTo(self.typ.In(i)) {
					v = v.Convert(self.typ.In(i))
				}
				inArgs = append(inArgs, v)
			case *StringObject:
				v := reflect.ValueOf(arg.Val)
				t := reflect.TypeOf(arg.Val)
				if t.ConvertibleTo(self.typ.In(i)) {
					v = v.Convert(self.typ.In(i))
				}
				inArgs = append(inArgs, v)
			case *BoolObject:
				v := reflect.ValueOf(arg.Val)
				t := reflect.TypeOf(arg.Val)
				if t.ConvertibleTo(self.typ.In(i)) {
					v = v.Convert(self.typ.In(i))
				}
				inArgs = append(inArgs, v)
			case *GoObject:
				v := reflect.ValueOf(arg.obj)
				t := reflect.TypeOf(arg.obj)
				if t.ConvertibleTo(self.typ.In(i)) {
					v = v.Convert(self.typ.In(i))
				}
				inArgs = append(inArgs, v)
			default:
				v := reflect.ValueOf(arg)
				t := reflect.TypeOf(arg)
				if t.ConvertibleTo(self.typ.In(i)) {
					v = v.Convert(self.typ.In(i))
				}
				inArgs = append(inArgs, v)
			}
		}
	}

	outVals := reflect.ValueOf(self.fn).Call(inArgs)
	for _, val := range outVals {
		switch val.Kind() {
		case reflect.Bool:
			results = append(results, NewBoolObject(val.Bool()))
		case reflect.String:
			results = append(results, NewStringObject(val.String()))
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			results = append(results, NewIntegerObject(int(val.Int())))
		case reflect.Float64, reflect.Float32:
			results = append(results, NewFloatObject(val.Float()))
		default:
			results = append(results, NewGoObject(val.Interface()))
		}
	}

	return
}

func (self *GoFuncObject) Dispatch(ctx *Runtime, method string, args ...Object) (results []Object) {
	switch method {
	case "__call__":
		results = self.callGoFunc(ctx, args...)
	}
	return
}

func (self *GoFuncObject) Name() string {
	return "gofunc"
}

func (self *GoFuncObject) String() string {
	return self.name
}

func (self *GoFuncObject) HashCode() string {
	return fmt.Sprintf("%p", self.fn)
}
