
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

list = [1,2,3,4]

fmt.Println("list length", list.Length())
list.Append(6)
fmt.Println(list)

str = "1n2n3n5"
subs = str[1:4]
fmt.Println(subs)

if str == "1n2n3n5" {
	fmt.Println(true, str[1])
} else {
	fmt.Println(false)
}

a = #{str: 100, 100: 200}

fmt.Println(a, a[str], a[100])
