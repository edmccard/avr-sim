package cpu

import "github.com/edmccard/avr-sim/instr"

type Cpu struct {
	R     [32]int
	FlagI bool
	FlagT bool
	FlagH bool
	FlagS bool
	FlagV bool
	FlagN bool
	FlagZ bool
	FlagC bool
	SP    int
	PC    int
}

func (c *Cpu) SregFromByte(b byte) {
	c.FlagC = (b & 0x01) != 0
	c.FlagZ = (b & 0x02) != 0
	c.FlagN = (b & 0x04) != 0
	c.FlagV = (b & 0x08) != 0
	c.FlagS = (b & 0x10) != 0
	c.FlagH = (b & 0x20) != 0
	c.FlagT = (b & 0x40) != 0
	c.FlagI = (b & 0x80) != 0
}

func (c *Cpu) ByteFromSreg() (b byte) {
	if c.FlagC {
		b |= 0x01
	}
	if c.FlagZ {
		b |= 0x02
	}
	if c.FlagN {
		b |= 0x04
	}
	if c.FlagV {
		b |= 0x08
	}
	if c.FlagS {
		b |= 0x10
	}
	if c.FlagH {
		b |= 0x20
	}
	if c.FlagT {
		b |= 0x40
	}
	if c.FlagI {
		b |= 0x80
	}
	return
}

func Adc(cpu *Cpu, am *instr.AddrMode) {
	addition(cpu, am, cpu.FlagC)
}

func Add(cpu *Cpu, am *instr.AddrMode) {
	addition(cpu, am, false)
}

func addition(cpu *Cpu, am *instr.AddrMode, carry bool) {
	d := cpu.R[am.A1]
	r := cpu.R[am.A2]
	c := 0
	if carry {
		c = 1
	}

	hres := (d & 0xf) + (r & 0xf) + c
	cpu.FlagH = hres > 0xf

	res := d + r + c
	cpu.FlagC = res > 0xff
	cpu.FlagV = (((d ^ r) & 0x80) == 0) && (((d ^ res) & 0x80) != 0)

	res &= 0xff
	cpu.FlagZ = res == 0
	cpu.FlagN = res >= 0x80
	cpu.FlagS = cpu.FlagV != cpu.FlagN

	cpu.R[am.A1] = res
}

func Sub(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = subtractionNoCarry(cpu, cpu.R[am.A1], cpu.R[am.A2])
}

func Subi(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = subtractionNoCarry(cpu, cpu.R[am.A1], int(am.A2))
}

func Cp(cpu *Cpu, am *instr.AddrMode) {
	subtractionNoCarry(cpu, cpu.R[am.A1], cpu.R[am.A2])
}

func Cpi(cpu *Cpu, am *instr.AddrMode) {
	subtractionNoCarry(cpu, cpu.R[am.A1], int(am.A2))
}

func subtractionNoCarry(cpu *Cpu, d, r int) int {
	r = ^r

	hres := (d & 0xf) + (r & 0xf) + 1
	cpu.FlagH = hres <= 0xf

	res := d + r + 1
	cpu.FlagC = res < 0
	cpu.FlagV = (((d ^ r) & 0x80) == 0) && (((d ^ res) & 0x80) != 0)

	res &= 0xff
	cpu.FlagZ = res == 0
	cpu.FlagN = res >= 0x80
	cpu.FlagS = cpu.FlagV != cpu.FlagN

	return res
}

func Sbc(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = subtractionCarry(cpu, cpu.R[am.A1], cpu.R[am.A2], cpu.FlagC)
}

func Sbci(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = subtractionCarry(cpu, cpu.R[am.A1], int(am.A2), cpu.FlagC)
}

func Cpc(cpu *Cpu, am *instr.AddrMode) {
	subtractionCarry(cpu, cpu.R[am.A1], cpu.R[am.A2], cpu.FlagC)
}

func subtractionCarry(cpu *Cpu, d, r int, carry bool) int {
	r = ^r
	c := 1
	if carry {
		c = 0
	}

	hres := (d & 0xf) + (r & 0xf) + c
	cpu.FlagH = hres <= 0xf

	res := d + r + c
	cpu.FlagC = res < 0
	cpu.FlagV = (((d ^ r) & 0x80) == 0) && (((d ^ res) & 0x80) != 0)

	res &= 0xff
	if res != 0 {
		cpu.FlagZ = false
	}
	cpu.FlagN = res >= 0x80
	cpu.FlagS = cpu.FlagV != cpu.FlagN

	return res
}