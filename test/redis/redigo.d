import "fmt"
import "redis"

c, err = redis.Dial("tcp", ":6379")

if err == nil {
	fmt.Println("conn succ")
}

fmt.Println(c, err)

reply, err = c.Do("SET", "a", "100")
reply, err = c.Do("INFO")
info, err = redis.String(reply, err)

info.Split("\\n").Each(func(index, line) {
	xs = line.Split(":")
	if xs.Size() == 2 {
		fmt.Println(index, xs[0], xs[1])
	}
})

fmt.Println(nil)

c.Close()
