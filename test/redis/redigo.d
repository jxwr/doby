import "fmt"
import "redis"

c, err = redis.Dial("tcp", ":6379")

if err == nil {
	fmt.Println("conn succ")
}

fmt.Println(c, err)

reply, err = c.Do("SET", "a", "100")
fmt.Println("Set:", reply, err)

reply, err = c.Do("GET", "a")
fmt.Println("Get:", reply, err)

str, err = redis.String(reply, err)
fmt.Println(str, err)

fmt.Println(nil)

c.Close()
