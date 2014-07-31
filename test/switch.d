import "fmt"

for _, a = range [100, 200, 300] {
	switch a {
	case 100:
		fmt.Println("hello")
		fmt.Println("one hundred")
	case 200:
		fmt.Println("yes")
		fmt.Println("two hundred")
	default:
		fmt.Println(a)
	}
}
