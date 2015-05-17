package atmega8

import (
	"io"

	"github.com/edmccard/avr-sim/core"
	"github.com/edmccard/ihex"
)

const (
	FlashWords = 0x4000
	SramBytes  = 0x460
	PortCount  = 0x60
)

type Mem struct {
	prog     []uint16
	data     []byte
	inports  []core.MemRead
	outports []core.MemWrite
}

func NewMem(cpu *core.Cpu) *Mem {
	mem := &Mem{
		prog:     make([]uint16, FlashWords),
		data:     make([]byte, SramBytes),
		inports:  make([]core.MemRead, PortCount),
		outports: make([]core.MemWrite, PortCount),
	}

	for i := 0; i < 32; i++ {
		mem.inports[i] = cpu.MemReadReg
		mem.outports[i] = cpu.MemWriteReg
	}

	memread := func(addr core.Addr) byte {
		return mem.data[addr]
	}
	memwrite := func(addr core.Addr, val byte) {
		mem.data[addr] = val
	}
	for i := 32; i < PortCount; i++ {
		mem.inports[i] = memread
		mem.outports[i] = memwrite
	}

	mem.inports[0x5d] = cpu.MemReadSPL
	mem.outports[0x5d] = cpu.MemWriteSPL
	mem.inports[0x5e] = cpu.MemReadSPH
	mem.outports[0x5e] = cpu.MemWriteSPH
	mem.inports[0x5f] = cpu.MemReadSreg
	mem.outports[0x5f] = cpu.MemWriteSreg

	return mem
}

func (mem *Mem) SetWriter(addr core.Addr, f core.MemWrite) {
	mem.outports[addr] = f
}

func (mem *Mem) SetReader(addr core.Addr, f core.MemRead) {
	mem.inports[addr] = f
}

func (mem *Mem) SetRW(addr core.Addr, r core.MemRead, w core.MemWrite) {
	mem.SetReader(addr, r)
	mem.SetWriter(addr, w)
}

func (mem *Mem) ReadData(addr core.Addr) byte {
	addr %= SramBytes
	if addr < PortCount {
		return mem.inports[addr](addr)
	}
	return mem.data[addr]
}

func (mem *Mem) WriteData(addr core.Addr, val byte) {
	addr %= SramBytes
	if addr < PortCount {
		mem.outports[addr](addr, val)
	} else {
		mem.data[addr] = val
	}
}

func (mem *Mem) LoadProgram(addr core.Addr) byte {
	shift := (uint(addr) & 0x1) * 8
	addr = (addr >> 1) & (FlashWords - 1)
	return byte(mem.prog[addr] >> shift)
}

func (mem *Mem) ReadProgram(addr core.Addr) uint16 {
	addr &= (FlashWords - 1)
	return mem.prog[addr]
}

func (mem *Mem) LoadHex(data io.Reader) {
	loadRecord := func(rec ihex.Record) {
		addr := rec.Address >> 1
		for i := 0; i < len(rec.Bytes); i += 2 {
			val := uint16(rec.Bytes[i]) | (uint16(rec.Bytes[i+1]) << 8)
			mem.prog[addr] = val
			addr++
		}
	}

	parser := ihex.NewParser(data)
	for parser.Parse() {
		loadRecord(parser.Data())
	}
	if parser.Err() != nil {
		panic("bad hex data")
	}
}
