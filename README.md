## Are you a DOUBI sent by the Monkey King?

A script language on top of golang runtime with golang syntax and ruby semantic, it could use golang packages if they were exported to the Doubi envrionment.

> DO NOT USE, IT'S EXPERIMENTAL

## Progress

* [DONE] Hacking IR generation
* [TODO] Clean code, remove eval ast tree related code, refactor

## DataType

### Array 
```go
import "fmt"

a = [1,2,3,4,5,6]

a.Each(func(e){fmt.Print(e)})
fmt.Println()

fn = func(e){return e*2 + 1}
b = a.Map(fn)
fmt.Println(b)

c = b.Select(func(e){return e > 10})
fmt.Println(c)

a.Push(7,8,9)
fmt.Println(a)
a.Pop(3)
fmt.Println(a)

e = a.Take(4)
fmt.Println(e)

f = e.Drop(2)
fmt.Println(f, e)
```

Outputs:
```
123456
[3,5,7,9,11,13]
[11,13]
[1,2,3,4,5,6,7,8,9]
[1,2,3,4,5,6]
[1,2,3,4]
[1,2] [3,4]
```

### Dict
```go
dict = #{"name": "jiaoxiang", "height": 180}
dict.birth = 1987
```

### Set
```
set = #[1,1,1,1,2]
```

### String

### Integer

### Float

## Examples:

### Quicksort

```go
import "fmt"

func filter(list, fn) {
     sublist = []
     for _, elem = range list {
         if fn(elem) {
            sublist.Push(elem)
         }
     }
     return sublist
}

func qsort(list) {
    if list.Length() <= 1 {
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
fmt.Println(a)

```
> [1,2,3,4,5,6,7,8,100,199,200,299,2229] 

### BubbleSort

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
				swapped = true
			}
		}
	}
	return lst
}

fmt.Println(gaps)
lst = bubble(gaps)
fmt.Println(lst)
```

#### Outputs:
```go
 =============>  test/bubble.d  <=============
 [701,301,132,57,23,10,4,1] 
 [1,4,10,23,57,132,301,701] 
```

### SUDOKU

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
	['5','3','.','.','7','.','9','.','.'],
	['6','.','.','1','9','5','.','.','.'],
	['.','9','8','.','.','.','.','6','.'],
	['8','.','.','.','6','.','.','.','3'],
	['4','.','6','8','.','3','7','.','1'],
	['7','.','.','.','2','.','.','.','6'],
	['.','6','1','.','.','.','2','8','.'],
	['.','.','.','4','1','9','.','.','5'],
	['3','.','5','.','8','.','.','7','9']
]

func showBoard(board) {
	fmt.Println("-----------------------------------")
	for i, line = range board {
		fmt.Println(line)
	}
}

solveSudoku(board)
```

#### Outputs:
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

### Use go modules
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

This shows how to import go modules into doubi enviroment (rt/runtime.go)
```go
rt.RegisterFunctions("fmt", []interface{}{
	fmt.Errorf,
	fmt.Println, fmt.Print, fmt.Printf,
	fmt.Fprint, fmt.Fprint, fmt.Fprintln, fmt.Fscan, fmt.Fscanf, fmt.Fscanln,
	fmt.Scan, fmt.Scanf, fmt.Scanln,
	fmt.Sscan, fmt.Sscanf, fmt.Sscanln,
	fmt.Sprint, fmt.Sprintf, fmt.Sprintln,
})

rt.RegisterFunctions("log", []interface{}{
	log.Fatal, log.Fatalf, log.Fatalln, log.Flags, log.Panic, log.Panicf, log.Panicln,
	log.Print, log.Printf, log.Println, log.SetFlags, log.SetOutput, log.SetPrefix,
})

rt.RegisterFunctions("os", []interface{}{
	os.Chdir, os.Chmod, os.Chown, os.Exit, os.Getpid, os.Hostname,
})

rt.RegisterFunctions("time", []interface{}{
	time.Sleep, time.Now, time.Unix,
})

rt.RegisterFunctions("math/rand", []interface{}{
	rand.New, rand.NewSource,
	rand.Float64, rand.ExpFloat64, rand.Float32, rand.Int,
	rand.Int31, rand.Int31n, rand.Int63, rand.Int63n, rand.Intn,
	rand.NormFloat64, rand.Perm, rand.Seed, rand.Uint32,
})
```

### math/rand Example

Code from https://gobyexample.com/random-numbers

No ':=' and you should call main() manually

```
import "fmt"
import "math/rand"

func main() {
	fmt.Print(rand.Intn(100), ",")
	fmt.Print(rand.Intn(100))
	fmt.Println()

	fmt.Println(rand.Float64())

	fmt.Print((rand.Float64()*5)+5, ",")
	fmt.Print((rand.Float64() * 5) + 5)
	fmt.Println()
	// To make the pseudorandom generator deterministic,
	// give it a well-known seed.
	s1 = rand.NewSource(42)
	r1 = rand.New(s1)

	// Call the resulting `rand.Source` just like the
	// functions on the `rand` package.
	fmt.Print(r1.Intn(100), ",")
	fmt.Print(r1.Intn(100))
	fmt.Println()

	// If you seed a source with the same number, it
	// produces the same sequence of random numbers.
	s2 = rand.NewSource(42)
	r2 = rand.New(s2)
	fmt.Print(r2.Intn(100), ",")
	fmt.Print(r2.Intn(100))
	fmt.Println()
}

main()
```

### A little ruby like syntax

```go
import "fmt"

10.times(func(i){ fmt.Println(i) })

hundred = (-100).abs()

fmt.Println(hundred)
```

### Object-based language ( No object model defined yet)

```
list = ["hello", "world"]
list.name = func() {
  return list[0] + " " + list[1]
}

fmt.Println(list.name())
```
> hello world 

### Dict

```go
import "fmt"

person = #{
  "name": "jiaoxiang",
  "age": 28,
  "summary": func(obj) {
     fmt.Println(obj["name"] + ":" + obj["age"])
  }
}

person.weight = 125
fmt.Println(person)
person.summary(person)
```
> {name:jiaoxiang,age:28,summary:#<closure>,weight:125}

> jiaoxiang:28

### Error Report

```
File: "test/play.d", Line 233, Col 2
 229) 
 230) // fibornacci
 231) count = 0
 232) func fib(n) {
*233) 	coun++ 
      	^
 234)     if n < 2 {
 235)         return n
 236)     }
 237)     return fib(n-2) + fib(n-1)
Error: 'coun' not Found
```

## Notes

> See notes and tests
