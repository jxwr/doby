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
}

func NewRuntime(visitor ast.Visitor) *Runtime {
	env := env.NewEnv(nil)

	RegisterGlobals(env)

	rt := &Runtime{visitor, env}
	return rt
}

func funcMap(funcList []interface{}) (fm map[string]Object) {
	fm = map[string]Object{}
	for _, f := range funcList {
		fname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		xs := strings.Split(fname, ".")
		fm[xs[len(xs)-1]] = NewGoFuncObject(fname, f)
	}
	return
}

func RegisterGlobals(env *env.Env) {
	env.Put("fmt", NewDictObject(funcMap([]interface{}{
		fmt.Errorf,
		fmt.Println, fmt.Print, fmt.Printf,
		fmt.Fprint, fmt.Fprint, fmt.Fprintln, fmt.Fscan, fmt.Fscanf, fmt.Fscanln,
		fmt.Scan, fmt.Scanf, fmt.Scanln,
		fmt.Sscan, fmt.Sscanf, fmt.Sscanln,
		fmt.Sprint, fmt.Sprintf, fmt.Sprintln,
	})))

	env.Put("log", NewDictObject(funcMap([]interface{}{
		log.Fatal, log.Fatalf, log.Fatalln, log.Flags, log.Panic, log.Panicf, log.Panicln,
		log.Print, log.Printf, log.Println, log.SetFlags, log.SetOutput, log.SetPrefix,
	})))

	env.Put("os", NewDictObject(funcMap([]interface{}{
		os.Chdir, os.Chmod, os.Chown, os.Exit, os.Getpid, os.Hostname,
	})))

	env.Put("time", NewDictObject(funcMap([]interface{}{
		time.Sleep, time.Now, time.Unix,
	})))

	env.Put("math/rand", NewDictObject(funcMap([]interface{}{
		rand.Float64, rand.ExpFloat64, rand.Float32, rand.Int,
		rand.Int31, rand.Int31n, rand.Int63, rand.Int63n, rand.Intn,
		rand.NormFloat64, rand.Perm, rand.Seed, rand.Uint32,
	})))
}

func WrapGoFunc(fn interface{}) {
	typ := reflect.TypeOf(fn)

	if typ.Kind() == reflect.Func {
		fmt.Println(typ.String())
		fmt.Println(typ.NumIn())
		fmt.Println(typ.NumOut())
	}
}
