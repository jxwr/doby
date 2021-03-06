import "fmt"

10.Times(func(i){
  fmt.Println(i)
})

func add_n(n) {
     return func(a) { return a + n }
}

add_100 = add_n(100)
add_150 = add_n(150)
c = add_100(20)
d = add_150(20)
fmt.Println(c, d)

func add(a, b) {
	c = a + b
	return c
}

d = 1
func add(a, b) {
	c = a + b + d
	return c
}

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
}

f0()
