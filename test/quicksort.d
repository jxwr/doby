
func qsort(list) {
    if list.len() <= 1 {
        return list
    }

    pivot = list[0]

    left = filter(func (x) {
        if x <= pivot {
            return list
        }
    })

    right = filter(func (x) {
        if x > pivot {
            return list
        }
    })

    return qsort(left) + pivot + qsort(right)
}
