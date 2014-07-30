package rt

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/env"
	"github.com/jxwr/doubi/vm/instr"
)

type Runtime struct {
	Visitor ast.Visitor
	Env     *env.Env
	Stack   *Stack
	Nil     Object
	True    Object
	False   Object

	TmpString *StringObject
	Runner    ClosureRunner

	NeedReturn   bool
	LoopDepth    int
	NeedBreak    bool
	NeedContinue bool

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

type Frame struct {
	Locals       []Object
	Upvals       []Object
	Parent       *Frame
	JumpTarget   int
	NeedReturn   bool
	NeedBreak    bool
	NeedContinue bool
	blockTargets []int
}

func NewFrame(numLocals, numUpvals int, parent *Frame) *Frame {
	frame := &Frame{
		make([]Object, numLocals),
		make([]Object, numUpvals),
		parent,
		-1,
		false,
		false,
		false,
		[]int{},
	}
	return frame
}

func (self *Frame) PushBlock(pc int) {
	self.blockTargets = append(self.blockTargets, pc)
}

func (self *Frame) PopBlock() (target int) {
	target = self.blockTargets[len(self.blockTargets)-1]
	self.blockTargets = self.blockTargets[:len(self.blockTargets)-1]
	return
}

func (self *Frame) BlockEndPc() int {
	if len(self.blockTargets) == 0 {
		return -1
	}
	return self.blockTargets[len(self.blockTargets)-1]
}

func NewRuntime(visitor ast.Visitor) *Runtime {
	env := env.NewEnv(nil)

	rt := &Runtime{Visitor: visitor, Env: env, Stack: NewStack()}
	rt.TmpString = rt.NewStringObject("")

	rt.registerGlobals(env)

	rt.Nil = &NilObject{}
	rt.goTypeMap = map[string]*Property{}
	rt.initBuiltinObjectProperties()
	rt.True = rt.NewBoolObject(true)
	rt.False = rt.NewBoolObject(false)

	return rt
}

func (self *Runtime) CallFuncObj(fnobj *ClosureObject, args ...Object) {
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

	if obj != nil && val.Kind() > reflect.Invalid && val.Kind() <= reflect.UnsafePointer {
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

func (self *Runtime) NewFuncObject(name string, decl *ast.FuncDeclExpr, e *env.Env) Object {
	obj := &FuncObject{MakeProperty(nil, &self.funcProperties), name, decl, false, nil, e}
	return obj
}

func (self *Runtime) NewBuiltinFuncObject(name string, recv Object, e *env.Env) *FuncObject {
	obj := &FuncObject{MakeProperty(nil, &self.funcProperties), name, nil, true, recv, e}
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
		v = reflect.ValueOf(obj.Val).Convert(typ)
	case *FloatObject:
		v = reflect.ValueOf(obj.Val).Convert(typ)
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
				self.TmpString.Val = typ.Field(i).Name
				prop.SetProp(self.TmpString, self.NewGoObject(val.Field(i).Interface()))
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
				fn := self.NewBuiltinFuncObject(m.Name, nil, nil)
				self.TmpString.Val = m.Name
				prop.SetProp(self.TmpString, fn)
			}
		}
	} else {
		for i := 0; i < numMethods; i++ {
			m := typ.Method(i)
			fn := self.NewBuiltinFuncObject(m.Name, nil, nil)
			self.TmpString.Val = m.Name
			prop.SetProp(self.TmpString, fn)
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

	funcObj := self.NewFuncObject("init", nil, nil)
	self.addObjectProperties(funcObj, &self.funcProperties)

	gofuncObj := self.NewGoFuncObject("init", nil)
	self.addObjectProperties(gofuncObj, &self.gofuncProperties)

	goObj := self.NewGoObject(nil)
	self.addObjectProperties(goObj, &self.goobjProperties)

	self.addObjectProperties(self.Nil, &self.nilProperties)
}

/// register

func (self *Runtime) RegisterVars(name string, vars map[string]interface{}) {
	dict, _ := self.Env.LookUp(name)

	m := map[string]Slot{}
	for k, v := range vars {
		self.TmpString.Val = k
		m[self.TmpString.HashCode()] = Slot{self.TmpString, self.NewGoObject(v)}
	}

	if dict == nil {
		dict = self.NewDictObject(m)
	} else {
		for k, v := range m {
			self.TmpString.Val = k
			dict.(*DictObject).SetProp(self.TmpString, v.Val)
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
		self.TmpString.Val = xs[len(xs)-1]
		m[self.TmpString.HashCode()] = Slot{self.TmpString, self.NewGoFuncObject(name, v)}
	}

	if dict == nil {
		dict = self.NewDictObject(m)
	} else {
		for k, v := range m {
			self.TmpString.Val = k
			dict.(*DictObject).SetProp(self.TmpString, v.Val)
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
		os.Chdir, os.Chmod, os.Chown, os.Exit, os.Getpid, os.Hostname,
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

func (self *Runtime) Rewind() {
	self.Stack.Rewind()
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
