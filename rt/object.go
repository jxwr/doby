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
	SetProp(obj Object, val Object)
	GetProp(obj Object) Object
	ToString(*Runtime, ...Object) []Object
}

func Invoke(rt *Runtime, obj Object, method string, args ...Object) (results []Object) {
	isBuiltin := false
	if strings.HasPrefix(method, "__") {
		if method == "__get_property__" {
			// builtin function
			val := obj.GetProp(args[0])
			fnobj, ok := val.(*FuncObject)
			if ok {
				fnobj.SetRecv(obj)
			}
			results = append(results, val)
			return
		} else if method == "__set_property__" {
			val := args[1]
			obj.SetProp(args[0], val)
			return
		} else {
			method = "OP" + method
		}
		isBuiltin = true
	}

	theMethod := reflect.ValueOf(obj).MethodByName(method)
	if theMethod.IsValid() || isBuiltin {
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
			methodType := theMethod.Type()
			theArgs := []reflect.Value{}

			if methodType.NumIn() > 0 {
				i := 0
				for ; i < methodType.NumIn()-1; i++ {
					reqTyp := methodType.In(i)
					theArgs = append(theArgs, ObjectToValue(args[i], reqTyp))
				}

				for ; i < len(args); i++ {
					var reqTyp reflect.Type
					if i < methodType.NumIn() && methodType.In(i).Kind() != reflect.Slice {
						reqTyp = methodType.In(i)
					}
					theArgs = append(theArgs, ObjectToValue(args[i], reqTyp))
				}
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

type Slot struct {
	Key Object
	Val Object
}

type Property struct {
	Slots  map[string]Slot
	Parent *Property
}

func MakeProperty(slots map[string]Slot, parent *Property) Property {
	return Property{slots, parent}
}

func EmptyProperty() Property {
	return Property{nil, nil}
}

func (self *Property) SetProp(obj Object, val Object) {
	if self.Slots == nil {
		self.Slots = map[string]Slot{}
	}
	hash := obj.HashCode()
	self.Slots[hash] = Slot{obj, val}
}

func (self *Property) GetProp(obj Object) Object {
	s := self
	for {
		hash := obj.HashCode()
		slot, ok := s.Slots[hash]
		if !ok {
			if s.Parent != nil {
				s = self.Parent
				continue
			} else {
				panic(fmt.Sprintf("Error: no property %v\n", obj))
			}
		}
		return slot.Val
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
