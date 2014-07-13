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
)

type Runtime struct {
	Visitor ast.Visitor
	Env     *env.Env

	integerProperties Property
}

func NewRuntime(visitor ast.Visitor) *Runtime {
	env := env.NewEnv(nil)

	rt := &Runtime{Visitor: visitor, Env: env}
	rt.RegisterGlobals(env)

	return rt
}

func (self *Runtime) NewIntegerObject(val int) *IntegerObject {
	obj := &IntegerObject{MakeProperty(nil, &self.integerProperties), val}
	return obj
}

func (self *Runtime) NewStringObject(val string) Object {
	obj := &StringObject{EmptyProperty(), val}
	return obj
}

func (self *Runtime) NewFloatObject(val float64) Object {
	obj := &FloatObject{EmptyProperty(), val}
	return obj
}

func (self *Runtime) NewGoFuncObject(fname string, fn interface{}) *GoFuncObject {
	gf := &GoFuncObject{EmptyProperty(), fname, reflect.TypeOf(fn), fn}
	return gf
}

func (self *Runtime) NewGoObject(obj interface{}) *GoObject {
	gobj := &GoObject{EmptyProperty(), obj}
	return gobj
}

func (self *Runtime) NewFuncObject(name string, decl *ast.FuncDeclExpr, e *env.Env) Object {
	obj := &FuncObject{EmptyProperty(), name, decl, false, nil, e}
	return obj
}

func (self *Runtime) NewBuiltinFuncObject(name string, recv Object, e *env.Env) *FuncObject {
	obj := &FuncObject{EmptyProperty(), name, nil, true, recv, e}
	return obj
}

func (self *Runtime) NewDictObject(fields map[string]Object) Object {
	obj := &DictObject{Property{fields, nil}}
	return obj
}

func (self *Runtime) NewBoolObject(val bool) Object {
	obj := &BoolObject{EmptyProperty(), val}
	return obj
}

func (self *Runtime) NewArrayObject(vals []Object) Object {
	obj := &ArrayObject{EmptyProperty(), vals}
	return obj
}

func (self *Runtime) NewSetObject(vals []Object) Object {
	obj := &SetObject{EmptyProperty(), vals}
	return obj
}

/// init object methods

/// register

func (self *Runtime) RegisterFunctions(name string, fns []interface{}) {
	self.Env.Put(name, self.NewDictObject(self.funcMap(fns)))
}

func (self *Runtime) RegisterGlobals(env *env.Env) {
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
		rand.Float64, rand.ExpFloat64, rand.Float32, rand.Int,
		rand.Int31, rand.Int31n, rand.Int63, rand.Int63n, rand.Intn,
		rand.NormFloat64, rand.Perm, rand.Seed, rand.Uint32,
	})
}

func (self *Runtime) funcMap(funcList []interface{}) (fm map[string]Object) {
	fm = map[string]Object{}
	for _, f := range funcList {
		fname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		xs := strings.Split(fname, ".")
		fm[xs[len(xs)-1]] = self.NewGoFuncObject(fname, f)
	}
	return
}

func WrapGoFunc(fn interface{}) {
	typ := reflect.TypeOf(fn)

	if typ.Kind() == reflect.Func {
		fmt.Println(typ.String())
		fmt.Println(typ.NumIn())
		fmt.Println(typ.NumOut())
	}
}
