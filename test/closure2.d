func add_n(n) {
     return func(a) { return a + n }
}

add_100 = add_n(100)
add_150 = add_n(150)

a = add_100(1)
b = add_150(1)

import "fmt"

fmt.Println(a, b)
