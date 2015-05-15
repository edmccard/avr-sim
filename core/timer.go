package core

type Timer struct {
	cycleCount int64
}

func (t *Timer) Tick(cycles int64) {
	t.cycleCount += cycles
}

func (t *Timer) GetCount() int64 {
	return t.cycleCount
}
