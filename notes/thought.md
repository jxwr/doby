### Some Thought

CodeGen

```
func add(a, b) {
     c = a + b
     return c
}

d = add(100, 200)

add sym
code:
push b
push operator:add
push a
send 3

push 100
push 100
push add
call

100.times(func(i){
   print(i)
})

push func
push times
push 100
send 3

cls = getClassOfObject(100)
method = cls.getMethod("times")
method.call(args)

push list
push len
call
```