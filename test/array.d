
import "fmt"

a = [1,2,3,4,5,6]

//a.Each(func(e){fmt.Print(e)})
//fmt.Println()

fn = func(e){return e*2 + 1}
b = a.Map(fn)
fmt.Println(b)
