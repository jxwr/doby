import "fmt"

a = [1,2,3,4]

a.Push(1)
fmt.Println(a)

b = [
  [0,1,2],
  [10,11,12]
  ]

fmt.Println(b)
b[0][1] = 99
fmt.Println(b, b[1][0])

c = [1,2,3,4]
e = c[:3]
fmt.Println(e)
