import "fmt"

// fibornacci
count = 0
func fib(n) {
	count++ 
    if n < 2 {
        return n
    }
    return fib(n-2) + fib(n-1)
}

fmt.Println(fib(20))
fmt.Print("\n")

