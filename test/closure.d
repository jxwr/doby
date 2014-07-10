
10.times(func(i){
  print(i,"\n")
})

func add_n(n) {
     return func(a) { a + n }
}

add_100 = add_n(100)
add_150 = add_n(150)
c = add_100(20)
d = add_150(20)
print(c, d, "\n")

func f0() {
     c = 100
     f1 = func() {
     	return c
     }
     c = 200
     f2 = func() {
     	return c
     }
     d = f1()
     e = f2()
     print(d, e,"\n")
}

f0()

func q(n) { print(n["name"]+"\n") }

func p(str) { q(str) }

p(#{"name": "hello"})
