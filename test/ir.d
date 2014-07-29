
a = 1 + 2 * 4
b = a + 100
c = "hello" + " world"
d = 1.02
e = a + d

f = [1,2,3,4]

//func add(a, b) {
//	c = a + b
//	d = c + 1
//	c
//	return
//	c + 1
//}

//m = add(1, 2)

func add_n(x) {
	return func(y) { return x + y }
}

add88 = add_n(88)
n = add88(100)
n
