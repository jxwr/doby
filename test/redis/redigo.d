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

dict = #{}
info.Split("\\n").Each(func(index, line) {
	xs = line.Split(":")
	if xs.Size() == 2 {
		dict[xs[0]] = xs[1].Trim()
	}
})

fmt.Println(dict)

for k, v = range dict {
	fmt.Println(k, v)
}

c.Close()
