### Are you a DOUBI sent by the Monkey King?

#### Examples:

* quicksort

```go
func filter(list, fn) {
     sublist = []
     for _, elem = range list {
         if fn(elem) {
            sublist.append(elem)
         }
     }
     return sublist
}

func qsort(list) {
    if list.length() <= 1 {
        return list
    }

    pivot = list[0]
    list = list[1:]

    left = filter(list, func (x) { return x <= pivot })
    right = filter(list, func (x) { return x > pivot })

    return qsort(left) + [pivot] + qsort(right)
}

lst = [200,299,199,3,4,1,2,7,8,5,6,100,2229]
a = qsort(lst)
print(a, "\n")

```
> [1,2,3,4,5,6,7,8,100,199,200,299,2229] 

* A little ruby like syntax

```go
func println(str) {
     print(str, "\n")
}

10.times(func(i){ println(i) })

hundred = -100.abs()
println(hundred)
```
>0 
>1 
>2 
>3 
>4 
>5 
>6 
>7 
>8 
>9 
>100

* Object-based_language

```
list = ["hello", "world"]
list.name = func() {
  return list[0] + " " + list[1]
}

println(list.name())
```
> hello world 

* Dict

```go
person = #{
  "name": "jiaoxiang",
  "age": 28,
  "summary": func(obj) {
     println(obj["name"] + ":" + obj["age"])
  }
}

person.weight = 125
println(person)
person.summary(person)
```
> {name:jiaoxiang,age:28,summary:#<closure>,weight:125}

> jiaoxiang:28

* Error Report

```
=============>  test/play.d  <=============
Syntax Error: Line:68 Col:15 NEARLINES:
  64)     print(i, "=", v, "\n")
  65)     return true
  66) }
  67) 
* 68) list = [1,2,3,4]]]]]]]
  69) list.append(5)
  70) println(list)
  71) 
  72) cl = [1,2,3,4] + [5,6,7,8]
  73) println(cl)
```

* misc
```go
func testLoop() {
	for i < 10 {
		i++
    
		if i == 2 {
			n = 0
			for n < 5 {
				n++
				if n == 3 {
					break
					print("break")
				}
				print("[", n, "]")	
			} 
			continue
		}

		if i == 9 {
			println("quit")
			return i
			print("never reach here")
		}
		print(i, "")
	}
}

n = testLoop()
println("return:"+n)
```
> 1 [ 1 ] [ 2 ] 3 4 5 6 7 8 quit

> return: 9

#### Notes

> See notes and tests

DO NOT USE, IT'S EXPERIMENTAL
