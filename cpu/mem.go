package cpu

import "github.com/edmccard/avr-sim/instr"

type Memory interface {
	ReadData(instr.Addr) (byte, error)
	WriteData(instr.Addr, byte) error
	// TODO: Flash memory interface
}
