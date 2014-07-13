import "fmt"

gaps = [701, 301, 132, 57, 23, 10, 4, 1]

func bubble(lst) {
	swapped = true
	for i = lst.Length() - 1; i > 0 && swapped; i-- {
		swapped = false
		for j = 0; j < i; j++ {
			if lst[j] > lst[j+1] {
				tmp = lst[j]
				lst[j] = lst[j+1]
				lst[j+1] = tmp
				fmt.Println(lst)
				swapped = true
			}
		}
	}
	return lst
}

fmt.Println(gaps)
lst = bubble(gaps)
fmt.Println(lst)
