### bytecode

```
a = 1 + 2

func add(a, b) {
	e = 100
	c = a + b + e
	return c
}

d = 1
e = 2
func add(a, b) {
	c = a + b + d + e
	return c
}
```

output
```
=============>  test/upval.d  <=============

CLOSURE seq:0 local:4 upval:0
  .local 0 a
  .local 1 add
  .local 2 d
  .local 3 e
CODE:
  0: push_int 2
  1: push_int 1
  2: send_stack :__add__ 1
  3: set_local 0
  4: push_closure 1
  5: set_local 1
  6: push_int 1
  7: set_local 2
  8: push_int 2
  9: set_local 3
 10: push_closure 2
 11: set_local 1

CLOSURE seq:1 local:4 upval:0
  .local 0 a
  .local 1 b
  .local 2 e
  .local 3 c
CODE:
  0: set_local 1
  1: set_local 0
  2: push_int 100
  3: set_local 2
  4: load_local 2
  5: load_local 1
  6: load_local 0
  7: send_stack :__add__ 1
  8: send_stack :__add__ 1
  9: set_local 3
 10: load_local 3
 11: raise_return 1

CLOSURE seq:2 local:3 upval:2
  .local 0 a
  .local 1 b
  .local 2 c
  .upval 0 1 3 e
  .upval 1 1 2 d
CODE:
  0: set_local 1
  1: set_local 0
  2: load_upval 0
  3: load_upval 1
  4: load_local 1
  5: load_local 0
  6: send_stack :__add__ 1
  7: send_stack :__add__ 1
  8: send_stack :__add__ 1
  9: set_local 2
 10: load_local 2
 11: raise_return 1
```