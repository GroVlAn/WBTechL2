package main

import (
	"context"
	"fmt"
	"github.com/beevik/ntp"
	"os/signal"
	"syscall"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

type Clock struct {
	currentTime time.Time
}

func NewClock() *Clock {
	return &Clock{}
}

func (c *Clock) Update() error {
	var err error
	c.currentTime, err = ntp.Time("0.beevik-ntp.pool.ntp.org")

	return err
}

func (c *Clock) GetFormattedTime() string {
	return c.currentTime.Format(timeFormat)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	_ = cancel
	clock := NewClock()

	go func() {
		for {
			err := clock.Update()
			if err == nil {
				time.Sleep(1 * time.Second)
				fmt.Printf("\r%s", clock.GetFormattedTime())
			}
		}
	}()

	<-ctx.Done()
}
