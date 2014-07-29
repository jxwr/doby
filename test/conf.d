import "fmt"

a = true
if a && false {
	fmt.Println("true:", a)
} else if a {
	fmt.Println("false:", a)
}


