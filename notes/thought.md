### Some Thought

CodeGen

```
e = 100

func add(a, b) {
     c = a
     d = c + b + e
     return d
}

nd = add(100, 200)

add sym
code:
push_args 0
store_local 0
push_args 0
push operator:add
push_args 1
send 3
store_local 1
push_local 1
return 1

push 100
push 100
push add
call
store 0

100.times(func(i){
   print(i)
})

push func
push times
push 100
send 3

Clsesh = getClassOfObject(100)
method = cls.getMethod("times")
method.call(args)

push list
push len
call
```

1. class-based or object-based
2. how to add new builtin datetypes and their methods
3. return
4. closure
5. builtins

### feature
1. concurrent
2. go reflection
3. ruby-like object extention
4. pattern matching
5. more like go
