Что выведет программа? Объяснить вывод программы.

```go
package main
 
import (
    "fmt"
)
 
func main() {
    a := [5]int{76, 77, 78, 79, 80}
    var b []int = a[1:4]
    fmt.Println(b)
}
```

Выведет [77, 78, 79] так как указан диапазон с 1 включительно и 4 не включительно