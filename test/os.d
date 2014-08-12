import "os"
import "fmt"

fmt.Println(os.Args)

if os.Args[1] == "hello" {
	fmt.Println("world")
}
