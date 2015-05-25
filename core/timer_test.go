package core

import (
	"reflect"
	"testing"
)

func TestTimerGoBoom(t *testing.T) {
	witness := 0
	timer := NewTimer()
	timer.AddCounter(NewCounter(1000, func() bool {
		witness++
		return false
	}))
	timer.Tick(999)
	if witness != 0 {
		t.Error("Timer expired too early")
	}
	timer.Tick(1)
	if witness != 1 {
		t.Error("Timer did not expire")
	}
	timer.Tick(1000)
	if witness != 1 {
		t.Error("One-shot timer was not removed")
	}
}

func TestTimerRepeat(t *testing.T) {
	witness := 0
	timer := NewTimer()
	timer.AddCounter(NewCounter(1000, func() bool {
		witness++
		return true
	}))
	timer.Tick(1000)
	if witness != 1 {
		t.Error("Repeat timer didn't fire")
	}
	timer.Tick(1004)
	if witness != 2 {
		t.Error("Repeat timer didn't repeat")
	}
	timer.Tick(996)
	if witness != 3 {
		t.Error("Repeat timer fuse length ignored overlap")
	}
	timer.Tick(2000)
	if witness != 5 {
		t.Error("Repeat timer didn't fire twice in single Tick")
	}
}

func TestTimerMulti(t *testing.T) {
	timer := NewTimer()
	record := make(map[int64]string)
	timer.AddCounter(NewCounter(21, func() bool {
		record[timer.cycleCount] = "A"
		return true
	}))
	timer.AddCounter(NewCounter(50, func() bool {
		record[timer.cycleCount] = "B"
		return true
	}))
	timer.Tick(100)
	if !reflect.DeepEqual(record, map[int64]string{
		21: "A", 42: "A", 50: "B", 63: "A", 84: "A", 100: "B",
	}) {
		t.Error("Multiple timer error")
	}
}
