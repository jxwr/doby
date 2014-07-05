func println(str) {
     print(str, "\n")
}

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
        println(list)
        return list
    }

    pivot = list[0]

    left = filter(list, func (x) {
        if x <= pivot {
            return true
        } else {
            return false
        }
    })

    right = filter(list, func (x) { return x > pivot })

    return qsort(left) + pivot + qsort(right)
}

l = [3,34,1,445,14]
a = qsort(l)
print(a, "\n")
