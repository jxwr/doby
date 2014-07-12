
import "fmt"
import "os"
import "log"
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
fmt.Println(hostname)
fmt.Println(err)
fmt.Println(hostname, err)
fmt.Println("err", err)

fmt.Println("nerver reach")

/// log
log.Println("hello", "world")

list = [1,2,3,4]
list[0], list[1] = list[1], list[0]
name = "jiao"
age = 28
fmt.Println(list, ["hello", "world"], #{"name": name, "age": age})
