
true = 1
false = 0

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
    if list.len() <= 1 {
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
