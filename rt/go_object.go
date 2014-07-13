package rt

import (
	"fmt"
	"reflect"
)

/// go object wrapper

type GoObject struct {
	Property
	obj interface{}
}

func (self *GoObject) Name() string {
	return "goobj"
}

func (self *GoObject) String() string {
	return fmt.Sprint(self.obj)
}

func (self *GoObject) ToString(rt *Runtime, args ...Object) []Object {
	return []Object{rt.NewStringObject(self.String())}
}

func (self *GoObject) HashCode() string {
	return fmt.Sprintf("%p", self.obj)
}

/// function

type GoFuncObject struct {
	Property
	name string
	typ  reflect.Type
	fn   interface{}
}

func (self *GoFuncObject) callGoFunc(rt *Runtime, args ...Object) (results []Object) {
	inNum := self.typ.NumIn()
	inArgs := []reflect.Value{}

	dummy := []interface{}{}
	for i := 0; i < inNum && i < len(args); i++ {
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
			results = append(results, rt.NewBoolObject(val.Bool()))
		case reflect.String:
			results = append(results, rt.NewStringObject(val.String()))
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			results = append(results, rt.NewIntegerObject(int(val.Int())))
		case reflect.Float64, reflect.Float32:
			results = append(results, rt.NewFloatObject(val.Float()))
		default:
			results = append(results, rt.NewGoObject(val.Interface()))
		}
	}

	return
}

func (self *GoFuncObject) OP__call__(rt *Runtime, args ...Object) (results []Object) {
	results = self.callGoFunc(rt, args...)
	return
}

func (self *GoFuncObject) Name() string {
	return "gofunc"
}

func (self *GoFuncObject) String() string {
	return self.name
}

func (self *GoFuncObject) ToString(rt *Runtime, args ...Object) []Object {
	return []Object{rt.NewStringObject(self.String())}
}

func (self *GoFuncObject) HashCode() string {
	return fmt.Sprintf("%p", self.fn)
}
