package rt

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/jxwr/doby/env"
	"github.com/jxwr/doby/vm/instr"
)

type Runtime struct {
	Env   *env.Env
	Stack *Stack
	Nil   Object
	True  Object
	False Object

	Runner ClosureRunner

	tmpString  *StringObject
	tmpInteger *IntegerObject

	goTypeMap map[string]*Property

	integerProperties Property
	floatProperties   Property
	stringProperties  Property
	arrayProperties   Property
	dictProperties    Property
	setProperties     Property
	boolProperties    Property
	nilProperties     Property
	funcProperties    Property
	gofuncProperties  Property
	goobjProperties   Property
}

func NewRuntime() *Runtime {
	env := env.NewEnv(nil)

	rt := &Runtime{Env: env, Stack: NewStack()}

	rt.tmpString = rt.NewStringObject("")
	rt.Nil = &NilObject{}
	rt.goTypeMap = map[string]*Property{}

	rt.registerGlobals(env)
	rt.initBuiltinObjectProperties()
	rt.True = rt.NewBoolObject(true)
	rt.False = rt.NewBoolObject(false)

	return rt
}

func (self *Runtime) CallFuncObj(fnobj *ClosureObject, args ...Object) {
	self.MarkN(-len(args))
	for _, arg := range args {
		self.Push(arg)
	}
	self.Runner.RunClosure(fnobj)
}

func (self *Runtime) NewIntegerObject(val int) *IntegerObject {
	obj := &IntegerObject{MakeProperty(nil, &self.integerProperties), val}
	return obj
}

func (self *Runtime) NewStringObject(val string) *StringObject {
	obj := &StringObject{MakeProperty(nil, &self.stringProperties), val}
	return obj
}

func (self *Runtime) NewFloatObject(val float64) *FloatObject {
	obj := &FloatObject{MakeProperty(nil, &self.floatProperties), val}
	return obj
}

func (self *Runtime) NewGoFuncObject(fname string, fn interface{}) *GoFuncObject {
	gf := &GoFuncObject{MakeProperty(nil, &self.gofuncProperties), fname, reflect.TypeOf(fn), fn}
	return gf
}

func (self *Runtime) NewGoObject(obj interface{}) *GoObject {
	gobj := &GoObject{MakeProperty(nil, &self.goobjProperties), obj}
	val := reflect.ValueOf(obj)

	if obj != nil && reflect.Indirect(val).IsValid() && val.Kind() > reflect.Invalid && val.Kind() <= reflect.UnsafePointer {
		key := reflect.Indirect(val).Type().PkgPath() + "::" + reflect.Indirect(val).Type().String()
		_, ok := self.goTypeMap[key]
		if !ok {
			prop := MakeProperty(nil, &self.goobjProperties)
			self.addObjectProperties(obj, &prop)
			self.goTypeMap[key] = &prop
		}
		gobj = &GoObject{MakeProperty(nil, self.goTypeMap[key]), obj}
	}
	return gobj
}

func (self *Runtime) NewClosureObject(proto *instr.ClosureProto,
	frame *Frame) *ClosureObject {
	obj := &ClosureObject{MakeProperty(nil, &self.funcProperties), proto, frame}
	return obj
}

func (self *Runtime) NewBuiltinFuncObject(name string) *FuncObject {
	obj := &FuncObject{MakeProperty(nil, &self.funcProperties), name, nil}
	return obj
}

func (self *Runtime) NewDictObject(fields map[string]Slot) Object {
	obj := &DictObject{MakeProperty(fields, &self.dictProperties), nil}
	return obj
}

func (self *Runtime) NewArrayObject(vals []Object) Object {
	obj := &ArrayObject{MakeProperty(nil, &self.arrayProperties), vals}
	return obj
}

func (self *Runtime) NewSetObject(vals []Object) Object {
	obj := &SetObject{MakeProperty(nil, &self.setProperties), vals}
	return obj
}

func (self *Runtime) NewBoolObject(val bool) Object {
	obj := &BoolObject{MakeProperty(nil, &self.boolProperties), val}
	return obj
}

func (self *Runtime) NewNilObject(vals []Object) Object {
	return self.Nil
}

func ObjectToValue(obj Object, typ reflect.Type) reflect.Value {
	var v reflect.Value
	switch obj := obj.(type) {
	case *IntegerObject:
		if typ == nil {
			v = reflect.ValueOf(obj.Val).Convert(reflect.TypeOf(1))
		} else {
			v = reflect.ValueOf(obj.Val).Convert(typ)
		}
	case *FloatObject:
		if typ == nil {
			v = reflect.ValueOf(obj.Val).Convert(reflect.TypeOf(0.1))
		} else {
			v = reflect.ValueOf(obj.Val).Convert(typ)
		}
	case *StringObject:
		v = reflect.ValueOf(obj.Val)
	case *BoolObject:
		v = reflect.ValueOf(obj.Val)
	case *GoObject:
		if obj.obj == nil {
			var nilObj *NilObject
			v = reflect.ValueOf(nilObj)
		} else {
			v = reflect.ValueOf(obj.obj)
		}
	default:
		v = reflect.ValueOf(obj)
	}
	return v
}

/// init object methods

func (self *Runtime) addObjectProperties(obj interface{}, prop *Property) {
	typ := reflect.TypeOf(obj)
	// if obj is a struct, add all fileds to its property
	val := reflect.ValueOf(obj)
	if typ.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			ch := typ.Field(i).Name[0]
			if ch >= 'A' && ch <= 'Z' {
				self.tmpString.Val = typ.Field(i).Name
				prop.SetProp(self.tmpString, self.NewGoObject(val.Field(i).Interface()))
			}
		}
	}

	// add methods of the type
	numMethods := typ.NumMethod()

	to_s, ok := typ.MethodByName("ToString")
	if ok {
		for i := 0; i < numMethods; i++ {
			m := typ.Method(i)
			if m.Type == to_s.Type {
				fn := self.NewBuiltinFuncObject(m.Name)
				self.tmpString.Val = m.Name
				prop.SetProp(self.tmpString, fn)
			}
		}
	} else {
		for i := 0; i < numMethods; i++ {
			m := typ.Method(i)
			fn := self.NewBuiltinFuncObject(m.Name)
			self.tmpString.Val = m.Name
			prop.SetProp(self.tmpString, fn)
		}
	}
}

func (self *Runtime) initBuiltinObjectProperties() {
	intObj := self.NewIntegerObject(0)
	self.addObjectProperties(intObj, &self.integerProperties)

	floatObj := self.NewFloatObject(0)
	self.addObjectProperties(floatObj, &self.floatProperties)

	stringObj := self.NewStringObject("")
	self.addObjectProperties(stringObj, &self.stringProperties)

	arrayObj := self.NewArrayObject(nil)
	self.addObjectProperties(arrayObj, &self.arrayProperties)

	dictObj := self.NewDictObject(nil)
	self.addObjectProperties(dictObj, &self.dictProperties)

	setObj := self.NewSetObject(nil)
	self.addObjectProperties(setObj, &self.setProperties)

	boolObj := self.NewBoolObject(false)
	self.addObjectProperties(boolObj, &self.boolProperties)

	gofuncObj := self.NewGoFuncObject("init", nil)
	self.addObjectProperties(gofuncObj, &self.gofuncProperties)

	goObj := self.NewGoObject(nil)
	self.addObjectProperties(goObj, &self.goobjProperties)

	self.addObjectProperties(self.Nil, &self.nilProperties)
}

/// register

// take care of slice object value
func (self *Runtime) GoValueToObject(obj interface{}) Object {
	val := reflect.ValueOf(obj)
	kind := val.Kind()

	switch kind {
	case reflect.Slice, reflect.Array:
		elems := []Object{}
		for i := 0; i < val.Len(); i++ {
			elems = append(elems, self.GoValueToObject(val.Index(i).Interface()))
		}
		return self.NewArrayObject(elems)
	case reflect.String:
		return self.NewStringObject(obj.(string))
	case reflect.Int, reflect.Int64:
		return self.NewIntegerObject(obj.(int))
	case reflect.Float32, reflect.Float64:
		return self.NewFloatObject(obj.(float64))
	case reflect.Bool:
		if obj.(bool) == true {
			return self.True
		} else {
			return self.False
		}
	default:
		return self.NewGoObject(obj)
	}
}

func (self *Runtime) RegisterVars(name string, vars map[string]interface{}) {
	dict, _ := self.Env.LookUp(name)

	m := map[string]Slot{}
	for k, v := range vars {
		self.tmpString.Val = k
		m[self.tmpString.HashCode()] = Slot{self.tmpString, self.GoValueToObject(v)}
	}

	if dict == nil {
		dict = self.NewDictObject(m)
	} else {
		for k, v := range m {
			self.tmpString.Val = k
			dict.(*DictObject).SetProp(self.tmpString, v.Val)
		}
	}
	self.Env.Put(name, dict)
}

func (self *Runtime) RegisterFunctions(name string, vars []interface{}) {
	dict, _ := self.Env.LookUp(name)

	m := map[string]Slot{}
	for _, v := range vars {
		name := runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
		xs := strings.Split(name, ".")
		self.tmpString.Val = xs[len(xs)-1]
		m[self.tmpString.HashCode()] = Slot{self.tmpString, self.NewGoFuncObject(name, v)}
	}

	if dict == nil {
		dict = self.NewDictObject(m)
	} else {
		for k, v := range m {
			self.tmpString.Val = k
			dict.(*DictObject).SetProp(self.tmpString, v.Val)
		}
	}
	self.Env.Put(name, dict)
}

func (self *Runtime) registerGlobals(env *env.Env) {
	self.RegisterFunctions("fmt", []interface{}{
		fmt.Errorf,
		fmt.Println, fmt.Print, fmt.Printf,
		fmt.Fprint, fmt.Fprint, fmt.Fprintln, fmt.Fscan, fmt.Fscanf, fmt.Fscanln,
		fmt.Scan, fmt.Scanf, fmt.Scanln,
		fmt.Sscan, fmt.Sscanf, fmt.Sscanln,
		fmt.Sprint, fmt.Sprintf, fmt.Sprintln,
	})

	self.RegisterFunctions("log", []interface{}{
		log.Fatal, log.Fatalf, log.Fatalln, log.Flags, log.Panic, log.Panicf, log.Panicln,
		log.Print, log.Printf, log.Println, log.SetFlags, log.SetOutput, log.SetPrefix,
	})

	self.RegisterFunctions("os", []interface{}{
		os.Chdir, os.Chmod, os.Chown, os.Exit, os.Getpid,
		os.Hostname, os.Environ, os.Getenv, os.Setenv,
		os.Create, os.Open,
	})

	self.RegisterVars("os", map[string]interface{}{
		"Args": os.Args[2:],
	})

	self.RegisterFunctions("time", []interface{}{
		time.Sleep, time.Now, time.Unix,
	})

	self.RegisterFunctions("math/rand", []interface{}{
		rand.New, rand.NewSource,
		rand.Float64, rand.ExpFloat64, rand.Float32, rand.Int,
		rand.Int31, rand.Int31n, rand.Int63, rand.Int63n, rand.Intn,
		rand.NormFloat64, rand.Perm, rand.Seed, rand.Uint32,
	})

	self.RegisterFunctions("io/ioutil", []interface{}{
		ioutil.WriteFile, ioutil.ReadFile, ioutil.TempDir, ioutil.TempFile,
		ioutil.ReadAll, ioutil.ReadDir, ioutil.NopCloser,
	})

	self.RegisterFunctions("bufio", []interface{}{
		bufio.NewWriter, bufio.NewReader, bufio.NewReadWriter, bufio.NewScanner,
	})
}

/// stack wrapper

func (self *Runtime) Push(obj Object) {
	self.Stack.Push(obj)
}

func (self *Runtime) Pop() Object {
	return self.Stack.Pop()
}

func (self *Runtime) Mark() {
	self.Stack.Mark()
}

func (self *Runtime) MarkN(offset int) {
	self.Stack.MarkN(offset)
}

func (self *Runtime) Rewind() {
	self.Stack.Rewind()
}

func (self *Runtime) PopMark() int {
	return self.Stack.PopMark()
}

func (self *Runtime) ShiftTopN(n, pos int) {
	self.Stack.ShiftTopN(n, pos)
}

func (self *Runtime) Fatalf(format string, a ...interface{}) {
	fmt.Printf("Runtime Error: "+format, a...)
	fmt.Println()
	os.Exit(1)
}

func WrapGoFunc(fn interface{}) {
	typ := reflect.TypeOf(fn)

	if typ.Kind() == reflect.Func {
		fmt.Println(typ.String())
		fmt.Println(typ.NumIn())
		fmt.Println(typ.NumOut())
	}
}
