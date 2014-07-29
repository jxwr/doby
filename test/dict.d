
func add(a, b) { a + b }

a = #{"name": "jiaoxiang", "age": add(1,2)}

import "fmt"

fmt.Println(a, a.name, a.age)
