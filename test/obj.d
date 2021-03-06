
import "fmt"

10.Abs()
10.Times(func(i){fmt.Println(i)})

a = 10
a++

a += 10 // 21
a *= 2  // 42
a /= 3  // 14
a -= 5  // 9

if a * 3 == 81 / 3 {
	fmt.Println(true)
}

// list
list = [] + [1,2,3,4] + [3,4,5] + [a] + []
fmt.Println(list[1:] + [2])

fmt.Println("list length", list.Length())
list.Push(6)
fmt.Println(list)

list.Each(func(elem){fmt.Print(elem, ",")})
fmt.Println()
addone = list.Map(func(elem){ 
   fmt.Println(elem)
   return 1 
})
fmt.Println(addone)

// str
str = "1n2n3n5"
subs = str[1:4]
fmt.Println(subs)

if str == "1n2n3n5" {
	fmt.Println(true, str[1])
} else {
	fmt.Println(false)
}

// dict
a = #{str: 100, 100: 200}
fmt.Println(a, a[str], a[100])

// bool

a = true       // true
b = a && false // false
c = a || b     // true
d = !c||a        // false
fmt.Println(a,b,c,d)

// float 
a = 0.1
b = -0.2

c = a + b
fmt.Println(c)
