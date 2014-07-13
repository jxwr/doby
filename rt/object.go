package rt

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Object interface {
	Name() string
	String() string
	HashCode() string
	SetProp(key string, val Object)
	GetProp(key string) Object
	ToString(*Runtime, ...Object) []Object
}

func Invoke(rt *Runtime, obj Object, method string, args ...Object) (results []Object) {
	if strings.HasPrefix(method, "__") {
		if method == "__get_property__" {
			idx := args[0].(*StringObject)
			val := obj.GetProp(idx.Val)
			fnobj, ok := val.(*FuncObject)
			if ok {
				fnobj.Obj = obj
			}
			results = append(results, val)
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
	if theMethod.IsValid() {
		// doubi object methods
		theArgs := []reflect.Value{reflect.ValueOf(rt)}
		if args != nil {
			for _, arg := range args {
				if arg != nil {
					theArgs = append(theArgs, reflect.ValueOf(arg))
				}
			}
		}
		vals := theMethod.Call(theArgs)
		results = vals[0].Interface().([]Object)
		return
	} else {
		// go object methods
		gobj, ok := obj.(*GoObject)
		if ok {
			theMethod = reflect.ValueOf(gobj.obj).MethodByName(method)
			if !theMethod.IsValid() {
				goto err
			}
			theArgs := []reflect.Value{}
			for _, arg := range args {
				theArgs = append(theArgs, ObjectToValue(arg))
			}
			rets := theMethod.Call(theArgs)
			for _, ret := range rets {
				results = append(results, rt.NewGoObject(ret.Interface()))
			}
		} else {
			goto err
		}
	}
	return

err:
	fmt.Printf("Error: Unknown Method %s for %s\n", method, obj.String())
	os.Exit(1)
	return
}

type Property struct {
	Slots  map[string]Object
	Parent *Property
}

func MakeProperty(slots map[string]Object, parent *Property) Property {
	return Property{slots, parent}
}

func EmptyProperty() Property {
	return Property{nil, nil}
}

func (self *Property) SetProp(key string, val Object) {
	if self.Slots == nil {
		self.Slots = map[string]Object{}
	}
	self.Slots[key] = val
}

func (self *Property) GetProp(key string) Object {
	s := self
	for {
		val, ok := s.Slots[key]
		if !ok {
			if s.Parent != nil {
				s = self.Parent
				continue
			} else {
				panic(fmt.Sprintf("Error: no property named %s\n", key))
			}
		}
		return val
	}
	return nil
}

type NilObject struct {
	Property
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

func (self *NilObject) ToString(rt *Runtime, args ...Object) []Object {
	return []Object{rt.NewStringObject(self.String())}
}
