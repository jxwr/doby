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
    print(left,"\n")
    right = filter(list, func (x) { return x > pivot })
    print(right,"\n")
    print(pivot, "\n")

    return qsort(left) + [pivot] + qsort(right)
}

a = 100
func test() {
     a = 200
}

test()
print(a)

lst = [200,299,199,3,4,1,2,7,8,5,6,100,2229]
a = qsort(lst)
print(a, "\n")
