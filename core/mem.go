package core

type Addr int

type Memory interface {
	ReadData(Addr) byte
	WriteData(Addr, byte)
	ReadProgram(Addr) uint16
	LoadProgram(Addr) byte
}

type MemRead func(Addr) byte
type MemWrite func(Addr, byte)
