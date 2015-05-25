package core

import "math"

type Timer struct {
	cycleCount int64
	fuse       int64
	counters   *Counter
}

func NewTimer() *Timer {
	timer := &Timer{}
	timer.AddCounter(NewCounter(math.MaxInt64, nil))
	return timer
}

func (t *Timer) Tick(cycles int64) {
	for cycles >= t.fuse {
		cycles -= t.fuse
		t.cycleCount += t.fuse
		ctr := t.counters
		if ctr.fire() {
			t.counters = t.counters.next
			t.AddCounter(ctr)
		} else {
			ctr = ctr.next
			t.fuse = ctr.end - t.cycleCount
			t.counters = ctr
		}
	}
	t.fuse -= cycles
	t.cycleCount += cycles
}

func (t *Timer) GetCount() int64 {
	return t.cycleCount
}

func (t *Timer) AddCounter(ctr *Counter) {
	if ctr.len == 0 {
		panic("zero-length counter")
	}
	ctr.end = ctr.len + t.cycleCount
	ctr.next = nil
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
		ctr.next = t.counters
		t.counters = ctr
	} else {
		ctr.next = next
		prev.next = ctr
	}
	t.fuse = t.counters.end - t.cycleCount
}

type Counter struct {
	next *Counter
	end  int64
	len  int64
	fire func() bool
}

func NewCounter(ln int64, action func() bool) *Counter {
	return &Counter{len: ln, fire: action}
}
