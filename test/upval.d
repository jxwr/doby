a = 1 + 2

func add1(a, b) {
	e = 100
	c = a + b + e
	return c
}

d = 1
e = 2
func add2(a, b) {
	c = a + b + d + e
	return c
}

import "fmt"

fmt.Println(add2(10, 20))
