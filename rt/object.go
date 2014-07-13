package rt

import (
	"fmt"
	"reflect"
	"strings"
)

type Object interface {
	Name() string
	String() string
	HashCode() string
	SetProp(key string, val Object)
	GetProp(key string) Object
}

func Invoke(rt *Runtime, obj Object, method string, args ...Object) (results []Object) {
	if strings.HasPrefix(method, "__") {
		if method == "__get_property__" {
			idx := args[0].(*StringObject)
			results = append(results, obj.GetProp(idx.Val))
			return
		} else if method == "__set_property__" {
			idx := args[0].(*StringObject)
			val := args[1]
			obj.SetProp(idx.Val, val)
			return
		} else {
			method = "OP" + method
		}
	}

	theMethod := reflect.ValueOf(obj).MethodByName(method)
	if !theMethod.IsValid() {
		fmt.Printf("Error: Unknown Method %s for %s\n", method, obj)
	}

	theArgs := []reflect.Value{reflect.ValueOf(rt)}

	if args != nil {
		for _, arg := range args {
			theArgs = append(theArgs, reflect.ValueOf(arg))
		}
	}

	vals := theMethod.Call(theArgs)
	results = vals[0].Interface().([]Object)
	return
}

type Property map[string]Object

func (self *Property) SetProp(key string, val Object) {
	(*self)[key] = val
}

func (self *Property) GetProp(key string) Object {
	return (*self)[key]
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
