
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
