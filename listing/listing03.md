Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main
 
import (
    "fmt"
    "os"
)
 
func Foo() error {
    var err *os.PathError = nil
    return err
}
 
func main() {
    err := Foo()
    fmt.Println(err)
    fmt.Println(err == nil)
}
```

вывод:
```bash
<nil>
false
```

Под капотом интерфейс представляет структуру с сылкой на тип и ссылкой на объект, а 
так же таблицу методов. Пустой интерфейс не имеет таблицы методов.
