
Что выведет программа? Объяснить вывод программы.

```go
package main

func main() {
ch := make(chan int)
go func() {
for i := 0; i < 10; i++ {
ch <- i
}
}()

    for n := range ch {
        println(n)
    }
}
```

будут выведены значение с 0 по 9 включительно, после чего получим панику deadlock, так 
как канал не закрыт
