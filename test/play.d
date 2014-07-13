import "fmt"

// datatype
list0 = [1,2,3,4]
list1 = [1,2,
	4,5]

empty_dict = #{}

dict = #{ 
	"name": "jxwr"
	"age": 123 
}

dict = #{ "name": "jxwr",
	"age": 123 }

dict = #{ "name": "jxwr",
	"age": 123,
}

dict = #{ 
	"name": "jxwr",
	"age": 123,
}

empty_set = #[]

a,b,c = 8,9,10
set = #[1,2,3,4,a,b,c]

set = #[1,2,3,4,
    a,b,c]

set = #[1,
	2,
	b,
	c]

set = #[    
	1, a, c
]

fmt.Print(list0, "\n", list1, "\n", empty_dict, "\n", dict, "\n", set, "\n")

// 

a = 100 + 3 * 123 
b = a + 2

fmt.Print("a", a, "\n")
fmt.Print("b", b, "\n")

c = [1,2,3,4]

fmt.Print(c, "\n")
fmt.Print(100 + c[3], "\n")
fmt.Print(c[0] + c[1] + c[2] + c[3], "\n")

fmt.Println("fun decl")

func add(a, b) {
	return a + b
}

fmt.Println(add(1, 200))

c[0] = 1000
fmt.Println(c[0])

c.len = func() {
	return 2 * (c[0] + 1)
}

fmt.Println(c.len())
c[0] = 132
fmt.Println(c.len())

if 2 > 1 {
	fmt.Println("true")
} else {
	fmt.Println("false")
}

if false {
	fmt.Println("true")
} else {
	fmt.Println("false")
}

a = 500
a++
a++
a.b = 998
a.b++
a.b++
fmt.Println(a)
fmt.Println(a.b)

for i = 0; i < 3; i++ {
    fmt.Print(i,"")
}

for i, v = range c {
    fmt.Print(i, "=", v, "\n")
    return true
}

list = [1,2,3,4]
list.Append(5)
fmt.Println(list)

cl = [1,2,3,4] + [5,6,7,8]
fmt.Println(cl)
fmt.Println(cl[1:8])

base = 99
fmt.Print("abs:", (-100).Abs(), "\n")
11.Times(func(i){
	if i % 2 == 0 { 
		fmt.Print(i, "") 
	}
})

list = ["hello", "world"]
list.name = func() {
	return list[0] + " " + list[1]
}

fmt.Println(list.name())

func printA() {
	for i = 0; i< 1000; i {
		fmt.Print("A")
	}
}

person = #{
	"name": "jiaoxiang",
	"age": 28,
	"summary": func(obj) {
		fmt.Println(obj["name"] + ":" + obj["age"])
	}
}

person.weight = 125
fmt.Println(person)

func nnn(obj) {
    fmt.Println(obj)
}

person.summary(person)
nnn(person)

//hello
i = 0

func testLoop() {
	for i < 10 {
		i++
		
		if i == 2 {
			n = 0
			for n < 5 {
				n++
				if n == 3 {
					break
					fmt.Print("break")
				}
				fmt.Print("[", n, "]")	
			} 
			continue
		}

		if i == 9 {
			fmt.Println("quit")
			return i
			fmt.Print("never reach here")
		}
		fmt.Print(i, "")
	}
}

n = testLoop()
fmt.Println("return:"+n)

// 1 [ 1 ] [ 2 ] 3 4 5 6 7 8 quit
// return: 9

// custom fmt.Print function
fmt.Println(220)

// closure

func test0(n) {
	return n
}

func test1() {
	n = 100

	a = test0(n)

	fmt.Print("a", a, "\n")

	b = test0(n)
	
	c = a + b
	fmt.Print("c", c, "\n")
}

func test2(n) {
	a = n
	fmt.Print("test2\n")
	return a + 1
}

func test3() {
	fmt.Print("test3\n")
	m = test2(123)
	fmt.Print(m, "\n")
}

test1()
bb = test2(888)
fmt.Print(bb, "\n")
test3()

// fibornacci
func fib(n) {
    if n < 2 {
        return n
    }
    return fib(n-2) + fib(n-1)
}

for i = 0; i < 10; i++ {
    fmt.Println(i, fib(i))
}
fmt.Print("\n")

// list reverse

func reverse(lst) {
	if lst.Length() < 1 {
		return lst   
	}
	return reverse(lst[1:]) + [lst[0]]
}

list = [1,2,3,4,5,6,7]
rlist = reverse(list)
fmt.Println(rlist)

// test object
a = #{}
a.size = 89
a.name = "tang"
a.hello = func() { fmt.Println("hello") }

a.hello()
fmt.Println(a)

// xxx_assign
a = 33
a += 11
fmt.Println(a)

b = [1,2,3]
b += [4,5,6]
fmt.Println(b)

c = 1.23
c += 2
fmt.Println(c)

// switch

a = 11

for i = 0; i < 20; i++ {
	fmt.Print(i)
}

switch a {
case a < 10:
	fmt.Print("< 10")
case a > 100:
	fmt.Print("> 10")
case 10:
	fmt.Print(10)
default:
	fmt.Print("default\n")
}

//10000.times(func(i){fmt.Print(i,"")})
fmt.Print("\n")

func(yes,no){fmt.Print(yes,no,"\n")}(true,false)

for index, value = range ["a", "b", "c", "d"] {
    fmt.Print(index, value, "\n")
}

for index, value = range #{"a": 1, "b": 2, "c": 3, "d": 4} {
    fmt.Print(index, value, "\n")
}
