package main

import (
	"context"
	"fmt"
	"github.com/beevik/ntp"
	"os/signal"
	"syscall"
	"time"
)

// Формат к которому будем приводить полученное время
const (
	timeFormat = "2006-01-02 15:04:05"
)

// Clock - структура хранящая текущее время
type Clock struct {
	currentTime time.Time
}

func NewClock() *Clock {
	return &Clock{}
}

// Update - метод запрашивающий текущее время через пакер ntp
func (c *Clock) Update() error {
	var err error
	c.currentTime, err = ntp.Time("0.beevik-ntp.pool.ntp.org")

	return err
}

// FormattedTime - метод для форматирования текудего времени
func (c *Clock) FormattedTime() string {
	return c.currentTime.Format(timeFormat)
}

/*
PrintTime - функция для вывода текущего времени в консоль
ctx - конктекст для остановки функции, если пользователь вышел из программы
*/
func (c *Clock) PrintTime(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("close")
			return
		default:
			err := c.Update()
			if err == nil {
				fmt.Printf("\r%s", c.FormattedTime())
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	_ = cancel
	clock := NewClock()

	go clock.PrintTime(ctx)

	<-ctx.Done()
}
