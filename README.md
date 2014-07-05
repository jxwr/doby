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
```

* Object-based_language

```
list = ["hello", "world"]
list.name = func() {
  return list[0] + " " + list[1]
}

println(list.name())
```

#### Notes

> See notes and tests

DO NOT USE, IT'S EXPERIMENTAL
