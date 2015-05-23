package atmega8

import (
	"fmt"
	"io"
	"time"

	"github.com/edmccard/avr-sim/core"
	"github.com/edmccard/avr-sim/instr"
)

type System struct {
	Cpu     *core.Cpu
	Decoder *instr.Decoder
	Memory  *Mem
	Timer   *core.Timer
}

func NewSystem() *System {
	set := instr.NewSetEnhanced8k()
	set[instr.Jmp] = false
	set[instr.Call] = false
	decoder := instr.NewDecoder(set)
	cpu := &core.Cpu{}
	return &System{
		Cpu:     cpu,
		Decoder: &decoder,
		Memory:  NewMem(cpu),
		Timer:   core.NewTimer(),
	}
}

func (sys *System) LoadProgHex(data io.Reader) {
	sys.Memory.LoadHex(data)
}

func (sys *System) Step() uint {
	elapsed := sys.Cpu.Step(sys.Memory, sys.Decoder)
	sys.Timer.Tick(int64(elapsed))
	return elapsed
}

func (sys *System) Go(hertz, slicePerSec int, onSlice SliceFunc) chan struct{} {
	cycPerSlice := uint(hertz / slicePerSec)
	quit := make(chan struct{})
	ticker := time.NewTicker(time.Second / time.Duration(slicePerSec))

	go func() {
		cycles := uint(0)
		for {
			select {
			case <-ticker.C:
				for cycles < cycPerSlice {
					cycles += sys.Step()
				}
				cycles -= cycPerSlice
				err := onSlice()
				if err != nil {
					fmt.Println("ERROR:", err)
					break
				}
			case <-quit:
				break
			}
		}
		ticker.Stop()
	}()

	return quit
}

type SliceFunc func() error
