
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

print(list0, "\n", list1, "\n", empty_dict, "\n", dict, "\n", set, "\n")

// 

a = 100 + 3 * 123 
b = a + 2

print("a", a, "\n")
print("b", b, "\n")

c = [1,2,3,4]

print(c, "\n")
print(100 + c[3], "\n")
print(c[0] + c[1] + c[2] + c[3], "\n")

func println(str) {
     print(str, "\n")
}

println("fun decl")

func add(a, b) {
     return a + b
}

println(add(1, 200))

c[0] = 1000
println(c[0])

c.len = func() {
         return 2 * (c[0] + 1)
}

println(c.len())
c[0] = 132
println(c.len())

if 2 > 1 {
  println("true")
} else {
  println("false")
}

if false {
  println("true")
} else {
  println("false")
}

a = 500
a++
a++
a.b = 998
a.b++
a.b++
println(a)
println(a.b)

for i = 0; i < 3; i++ {
    print(i,"")
}

for i, v = range c {
    print(i, "=", v, "\n")
    return true
}

list = [1,2,3,4]
list.append(5)
println(list)

cl = [1,2,3,4] + [5,6,7,8]
println(cl)
println(cl[1:8])

base = 99
print("abs:", (-100).abs(), "\n")
11.times(func(i){
  if i % 2 == 0 { 
    print(i, "") 
  }
})

list = ["hello", "world"]
list.name = func() {
  return list[0] + " " + list[1]
}

println(list.name())

func printA() {
     for i = 0; i< 1000; i {
     	 print("A")
     }
}

person = #{
  "name": "jiaoxiang",
  "age": 28,
  "summary": func(obj) {
    println(obj["name"] + ":" + obj["age"])
  }
}

person.weight = 125
println(person)

func nnn(obj) {
    println(obj)
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
					print("break")
				}
				print("[", n, "]")	
			} 
			continue
		}

		if i == 9 {
			println("quit")
			return i
			print("never reach here")
		}
		print(i, "")
	}
}

n = testLoop()
println("return:"+n)

// 1 [ 1 ] [ 2 ] 3 4 5 6 7 8 quit
// return: 9

// custom print function
println(220)

// closure

func test0(n) {
     return n
}

func test1() {
     n = 100

     a = test0(n)

     print("a", a, "\n")

     b = test0(n)
     
     c = a + b
     print("c", c, "\n")
}

func test2(n) {
     a = n
     print("test2\n")
     return a + 1
}

func test3() {
     print("test3\n")
     m = test2(123)
     print(m, "\n")
}

test1()
bb = test2(888)
print(bb, "\n")
test3()

// fibornacci
func fib(n) {
    if n < 2 {
        return n
    }
    return fib(n-2) + fib(n-1)
}

for i = 0; i < 10; i++ {
    print(fib(i), "")
}
print("\n")

// list reverse

func reverse(lst) {
  if lst.length() < 1 {
    return lst   
  }
  return reverse(lst[1:]) + [lst[0]]
}

list = [1,2,3,4,5,6,7]
rlist = reverse(list)
println(rlist)

// test object
a = #{}
a.size = 89
a.name = "tang"
a.hello = func() { println("hello") }

a.hello()
println(a)

// xxx_assign
a = 33
a += 11
println(a)

b = [1,2,3]
b += [4,5,6]
println(b)

c = 1.23
c += 2
println(c)

// switch

a = 11

switch a {
  case a < 10:
     print("< 10")
  case a > 100:
     print("> 10")
  case 10:
     print(10)
  default:
     print("default\n")
}

//10000.times(func(i){print(i,"")})
print("\n")

func(yes,no){print(yes,no,"\n")}(true,false)

for index, value = range ["a", "b", "c", "d"] {
    print(index, value, "\n")
}

for index, value = range #{"a": 1, "b": 2, "c": 3, "d": 4} {
    print(index, value, "\n")
}
