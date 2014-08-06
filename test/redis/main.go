package main

import (
	"flag"

	"github.com/garyburd/redigo/redis"
	"github.com/jxwr/doubi/runner"
)

var input string
var dumpInstrs bool
var printStack bool

func init() {
	flag.StringVar(&input, "f", "", "input file")
	flag.BoolVar(&dumpInstrs, "i", false, "dump instrs")
	flag.BoolVar(&printStack, "s", false, "print stack")
}

func main() {
	flag.Parse()

	r := runner.NewRunner()

	r.SetDumpInstrs(dumpInstrs)
	r.SetPrintStack(printStack)

	r.RegisterFunctions("redis", []interface{}{
		redis.Bool, redis.Bytes, redis.Float64, redis.Int, redis.Int64,
		redis.String, redis.Strings, redis.Uint64, redis.Values,
		redis.Dial, redis.DialTimeout, redis.NewConn, redis.NewLoggingConn,
		redis.NewPool, redis.NewScript,
	})

	if input != "" {
		r.Run(input)
	}
}
