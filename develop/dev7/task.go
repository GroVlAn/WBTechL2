package main

import (
	"fmt"
	"sync"
	"time"
)

// merge - функция для слияния n каналов
func merge(channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	// новый канал, в который будут записаны значения
	merged := make(chan interface{})

	// создаём столько горутин, сколько у нас пришло каналов
	for _, chIt := range channels {
		// обновляем счётчик WaitGroup - количество горутин, чьё завершение нужно ожидать
		wg.Add(1)

		go func(ch <-chan interface{}) {
			// в самом конце закрываем канал
			defer wg.Done()
			for val := range ch {
				// записываем в результирующий канал значения
				merged <- val
			}
		}(chIt)
	}

	// в отдельной горутине ожидаем выполнения всех горутин и закрываем результирующий канал
	// для того, чтобы наша функция была не блокирующей
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-merge(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %v", time.Since(start))

}
