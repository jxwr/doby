import "os"
import "fmt"

fmt.Println(os.Args)

if os.Args[1] == "hello" {
	fmt.Println("world")
}

os.Setenv("FOO", "1")
fmt.Println("FOO:", os.Getenv("FOO"))
fmt.Println("BAR:", os.Getenv("BAR"))
