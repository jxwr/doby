import "fmt"

a = true
if a && false {
	fmt.Println("true:", a)
} else if a {
	fmt.Println("false:", a)
}


b = [1,2,3,4]
fmt.Println(b[3])
b[3] = 10
fmt.Println(b + [4,4,3])
b.Push(11)
fmt.Println(b)
