import "fmt"

for j = 0; j < 20; j++ {
	for i = 0; i < 10; i++{
		if i > 5 {
			fmt.Println(i*10)
			if i > 7 {
				break
			}
		} else {
//			fmt.Println(i)
		}
		if j > 10 {
			break
		}
	}
//	fmt.Println(j)
	if j > 10 {
		break
	}
}

for i = 0; i < 10; i++ {
	if i == 6 {
		continue
	}
	fmt.Println(i)
}

//i=0
//cond:
//i<10
//jump_if_false end
//block:
//fmt.Println(i)
//post:
//i++
//jump cond
//end:
