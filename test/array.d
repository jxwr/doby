import "fmt"

a = [1,2,3,4,5,6]

a.Each(func(e){fmt.Print(e)})
fmt.Println()

fn = func(e){return e*2 + 1}
b = a.Map(fn)
fmt.Println(b)

c = b.Select(func(e){return e > 10})
fmt.Println(c)

a.Push(7,8,9)
fmt.Println(a)
a.Pop(3)
fmt.Println(a)

e = a.Take(4)
fmt.Println(e)

f = e.Drop(2)
fmt.Println(f, e)
