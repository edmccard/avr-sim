package cpu

type Addr int

type Memory interface {
	ReadData(Addr) byte
	WriteData(Addr, byte)
	ReadProgram(Addr) uint16
}
