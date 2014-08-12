import "bufio"
import "io/ioutil"
import "fmt"
import "os"

err = ioutil.WriteFile("dat", "hello", 420)
fmt.Println(err)

f, err = os.Open("dat")
fmt.Println(f, err)
