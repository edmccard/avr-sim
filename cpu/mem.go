package cpu

import "github.com/edmccard/avr-sim/instr"

type Memory interface {
	ReadProgramWord(instr.Addr) (uint16, error)
	WriteProgramWord(instr.Addr, uint16) error
	ReadDataByte(instr.Addr) (byte, error)
	WriteDataByte(instr.Addr, byte) error
}
