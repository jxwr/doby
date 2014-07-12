### Are you a DOUBI sent by the Monkey King?

#### Examples:

* use go module (developing, module level method for now)
```go
import "fmt"
import "os"
import "time"
import "math/rand"

/// math/rand
obj = rand.Float64()
fmt.Println(obj)
fmt.Println(rand.Float64)

/// time
Nanosecond = 1
Microsecond          = 1000 * Nanosecond
Millisecond          = 1000 * Microsecond
Second               = 1000 * Millisecond
Minute               = 60 * Second
Hour                 = 60 * Minute

time.Sleep(200 * Millisecond)
fmt.Println(time.Now())

/// os
hostname, err = os.Hostname()
fmt.Println(hostname, err)
os.Exit(1)

fmt.Println("nerver reach")
```

this is how to import go modules into doubi enviroment (rt/runtime.go)
```go
	env.Put("fmt", NewDictObject(funcMap([]interface{}{
		fmt.Errorf,
		fmt.Println, fmt.Print, fmt.Printf,
		fmt.Fprint, fmt.Fprint, fmt.Fprintln, fmt.Fscan, fmt.Fscanf, fmt.Fscanln,
		fmt.Scan, fmt.Scanf, fmt.Scanln,
		fmt.Sscan, fmt.Sscanf, fmt.Sscanln,
		fmt.Sprint, fmt.Sprintf, fmt.Sprintln,
	})))

	env.Put("log", NewDictObject(funcMap([]interface{}{
		log.Fatal, log.Fatalf, log.Fatalln, log.Flags, log.Panic, log.Panicf, log.Panicln,
		log.Print, log.Printf, log.Println, log.SetFlags, log.SetOutput, log.SetPrefix,
	})))

	env.Put("os", NewDictObject(funcMap([]interface{}{
		os.Chdir, os.Chmod, os.Chown, os.Exit, os.Getpid, os.Hostname,
	})))

	env.Put("time", NewDictObject(funcMap([]interface{}{
		time.Sleep, time.Now, time.Unix,
	})))

	env.Put("math/rand", NewDictObject(funcMap([]interface{}{
		rand.Float64, rand.ExpFloat64, rand.Float32, rand.Int,
		rand.Int31, rand.Int31n, rand.Int63, rand.Int63n, rand.Intn,
		rand.NormFloat64, rand.Perm, rand.Seed, rand.Uint32,
	})))
```

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

* bubblesort
```go
gaps = [701, 301, 132, 57, 23, 10, 4, 1]

func bubble(lst) {
	swapped = true
	for i = lst.length() - 1; i > 0 && swapped; i-- {
		swapped = false
		for j = 0; j < i; j++ {
			if lst[j] > lst[j+1] {
				tmp = lst[j]
				lst[j] = lst[j+1]
				lst[j+1] = tmp
				print(lst, "\n")
				swapped = true
			}
		}
	}
	return lst
}

print(gaps, "\n")
lst = bubble(gaps)

print(lst, "\n")
```

> =============>  test/bubble.d  <=============

> [701,301,132,57,23,10,4,1] 

> [301,701,132,57,23,10,4,1] 

> [301,132,701,57,23,10,4,1] 

> [301,132,57,701,23,10,4,1] 

> [301,132,57,23,701,10,4,1] 

> [301,132,57,23,10,701,4,1] 

> [301,132,57,23,10,4,701,1] 

> [301,132,57,23,10,4,1,701] 

> [132,301,57,23,10,4,1,701] 

> [132,57,301,23,10,4,1,701] 

> [132,57,23,301,10,4,1,701] 

> [132,57,23,10,301,4,1,701] 

> [132,57,23,10,4,301,1,701] 

> [132,57,23,10,4,1,301,701] 

> [57,132,23,10,4,1,301,701] 

> [57,23,132,10,4,1,301,701] 

> [57,23,10,132,4,1,301,701] 

> [57,23,10,4,132,1,301,701] 

> [57,23,10,4,1,132,301,701] 

> [23,57,10,4,1,132,301,701] 

> [23,10,57,4,1,132,301,701] 

> [23,10,4,57,1,132,301,701] 

> [23,10,4,1,57,132,301,701] 

> [10,23,4,1,57,132,301,701] 

> [10,4,23,1,57,132,301,701] 

> [10,4,1,23,57,132,301,701] 

> [4,10,1,23,57,132,301,701] 

> [4,1,10,23,57,132,301,701] 

> [1,4,10,23,57,132,301,701] 

> [1,4,10,23,57,132,301,701] 

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

* sudoku

```go
import "fmt"
import "os"

func isValid(board, x, y, c) {
    for i = 0; i < 9; i++ {
        if board[x][i] == c {
            return false
        }
        if board[i][y] == c {
            return false
        }
    }


    for i = 3*(x/3); i < 3*(x/3+1); i++ {
        for j = 3*(y/3); j < 3*(y/3+1); j++ {
            if board[i][j] == c {
                return false
            }
        }
    }
    return true
}

count = 0

func solveSudoku(board) {
	showBoard(board)
	count++
	fmt.Println(count)

    for i = 0; i < board.length(); i++ {
        for j = 0; j < board[i].length(); j++ {
            if board[i][j] == "." {
                for k = 0; k < 9; k++ {
                    c = "" + (k+1)
                    if isValid(board, i, j, c) {
                        board[i][j] = c
                        if solveSudoku(board) {
                            return true
                        }
                        board[i][j] = "."
                    }
                }
                return false
            }
        }
    }
    return true
}

board = [
	["5","3",".",".","7",".","9",".","."],
	["6",".",".","1","9","5",".",".","."],
	[".","9","8",".",".",".",".","6","."],
	["8",".",".",".","6",".",".",".","3"],
	["4",".","6","8",".","3","7",".","1"],
	["7",".",".",".","2",".",".",".","6"],
	[".","6","1",".",".",".","2","8","."],
	[".",".",".","4","1","9",".",".","5"],
	["3",".","5",".","8",".",".","7","9"]
]

func showBoard(board) {
	fmt.Println("-----------------------------------")
	for i, line = range board {
		fmt.Println(line)
	}
}

solveSudoku(board)
```

Results:
```go
...
-----------------------------------
[5,3,4,6,7,8,9,1,2]
[6,7,2,1,9,5,3,4,8]
[1,9,8,3,4,2,5,6,7]
[8,5,9,7,6,1,4,2,3]
[4,2,6,8,5,3,7,9,1]
[7,1,3,9,2,4,8,5,6]
[9,6,1,5,3,7,2,8,4]
[2,8,7,4,1,9,6,3,5]
[3,4,5,2,8,6,1,7,9]
730
```

#### Notes

> See notes and tests

DO NOT USE, IT'S EXPERIMENTAL
