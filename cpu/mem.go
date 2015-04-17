package cpu

import "github.com/edmccard/avr-sim/instr"

type Memory interface {
	ReadData(instr.Addr) byte
	WriteData(instr.Addr, byte)
}
