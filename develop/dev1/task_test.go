package main

import (
	"testing"
	"time"
)

func TestClock(t *testing.T) {
	clock := NewClock()
	_ = clock.Update()
	timeNow := time.Now().Format(timeFormat)

	if timeNow != clock.GetFormattedTime() {
		t.Error("it's not correct now time")
	}
}
