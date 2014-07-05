

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
    println(i)
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
10.times(func(i){ println(i) })

println(-100.abs())

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

person.summary(person)

i = 0 * 

for i < 100 {
    println(i)
    i++
}
