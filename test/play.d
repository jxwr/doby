

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

for i = 0; i < 100; i++ {
    println(i)
}
