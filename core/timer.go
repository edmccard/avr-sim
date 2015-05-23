package core

import "math"

type Timer struct {
	cycleCount int64
	fuse       int64
	counters   *Counter
}

func NewTimer() *Timer {
	return &Timer{fuse: math.MaxInt64}
}

func (t *Timer) Tick(cycles int64) {
	t.cycleCount += cycles
	for cycles >= t.fuse {
		cycles -= t.fuse
		ctr := t.counters
		if ctr.fire() {
			t.AddCounter(ctr)
		} else {
			ctr = ctr.next
			t.fuse = ctr.rem
			t.counters = ctr
		}
	}
	t.fuse -= cycles
}

func (t *Timer) GetCount() int64 {
	return t.cycleCount
}

func (t *Timer) AddCounter(ctr *Counter) {
	if ctr.len == 0 {
		panic("zero-length counter")
	}
	ctr.rem = ctr.len
	ctr.end = ctr.len + t.cycleCount
	t.insertCounter(ctr)
}

func (t *Timer) insertCounter(ctr *Counter) {
	var prev *Counter
	next := t.counters
	for ; next != nil; next = next.next {
		if ctr.end < next.end {
			break
		}
		prev = next
	}
	if prev == nil {
		t.fuse = ctr.rem
		ctr.next = t.counters
		t.counters = ctr
	} else {
		ctr.next = next
		prev.next = ctr
	}
}

type Counter struct {
	next *Counter
	end  int64
	rem  int64
	len  int64
	fire func() bool
}

func NewCounter(ln int64, action func() bool) *Counter {
	return &Counter{len: ln, fire: action}
}
