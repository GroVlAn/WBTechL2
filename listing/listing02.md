Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и порядок их вызовов.

```go
package main
 
import (
    "fmt"
)
 
func test() (x int) {
    defer func() {
        x++
    }()
    x = 1
    return
}
 
 
func anotherTest() int {
    var x int
    defer func() {
        x++
    }()
    x = 1
    return x
}
 
 
func main() {
    fmt.Println(test())
    fmt.Println(anotherTest())
}
```

вывод будет следующим
```bash
2
1
```

Так как последний defer отработает первым
