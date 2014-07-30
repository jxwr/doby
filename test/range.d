import "fmt"

a = #{10:20,30:40}
b = [1,2,3,4]

i = 0

//for i, v = range b {
///    fmt.Println(i, v)
//}

fmt.Println("range_over_dict")
for index, value = range #{"a": 1, "b": 2, "c": 3, "d": 4} {
    fmt.Print(index, value, "\n")
}


//for i, v = range a {
//    fmt.Println(i, v)
//}


// PUSH_INT 0
// for_range_start:
// load_local expr
// SEND_METHOD OP__iter__ 1
// set_local i
// set_local v
// body_start
// ..
// jump