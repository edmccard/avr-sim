package cpu

import "github.com/edmccard/avr-sim/instr"

type Cpu struct {
	R     [32]int
	flags [8]bool
	SP    int
	PC    int
}

type Flag int

const (
	FlagC Flag = iota
	FlagZ
	FlagN
	FlagV
	FlagS
	FlagH
	FlagT
	FlagI
)

func (c *Cpu) GetFlag(f Flag) bool {
	return c.flags[f]
}

func (c *Cpu) SetFlag(f Flag, b bool) {
	c.flags[f] = b
}

func (c *Cpu) SregFromByte(b byte) {
	c.flags[FlagC] = (b & 0x01) != 0
	c.flags[FlagZ] = (b & 0x02) != 0
	c.flags[FlagN] = (b & 0x04) != 0
	c.flags[FlagV] = (b & 0x08) != 0
	c.flags[FlagS] = (b & 0x10) != 0
	c.flags[FlagH] = (b & 0x20) != 0
	c.flags[FlagT] = (b & 0x40) != 0
	c.flags[FlagI] = (b & 0x80) != 0
}

func (c *Cpu) ByteFromSreg() (b byte) {
	if c.flags[FlagC] {
		b |= 0x01
	}
	if c.flags[FlagZ] {
		b |= 0x02
	}
	if c.flags[FlagN] {
		b |= 0x04
	}
	if c.flags[FlagV] {
		b |= 0x08
	}
	if c.flags[FlagS] {
		b |= 0x10
	}
	if c.flags[FlagH] {
		b |= 0x20
	}
	if c.flags[FlagT] {
		b |= 0x40
	}
	if c.flags[FlagI] {
		b |= 0x80
	}
	return
}

type OpFunc func(*Cpu, *instr.AddrMode, Memory)

func Adc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addition(cpu, am, cpu.flags[FlagC])
}

func Add(cpu *Cpu, am *instr.AddrMode, mem Memory) {
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
	cpu.flags[FlagH] = hres > 0xf

	res := d + r + c
	cpu.flags[FlagC] = res > 0xff
	cpu.flags[FlagV] = (((d ^ r) & 0x80) == 0) && (((d ^ res) & 0x80) != 0)

	res &= 0xff
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]

	cpu.R[am.A1] = res
}

func Adiw(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := cpu.R[am.A1] | (cpu.R[am.A1+1] << 8)
	k := int(am.A2)

	res := (d + k) & 0xffff
	hr := (res & 0x8000) != 0
	hd := (d & 0x8000) != 0
	cpu.flags[FlagC] = !hr && hd
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = hr
	cpu.flags[FlagV] = !hd && hr
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]

	cpu.R[am.A1] = res & 0xff
	cpu.R[am.A1+1] = res >> 8
}

func Sub(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = subtractionNoCarry(cpu, cpu.R[am.A1], cpu.R[am.A2])
}

func Subi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = subtractionNoCarry(cpu, cpu.R[am.A1], int(am.A2))
}

func Cp(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	subtractionNoCarry(cpu, cpu.R[am.A1], cpu.R[am.A2])
}

func Cpi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	subtractionNoCarry(cpu, cpu.R[am.A1], int(am.A2))
}

func Neg(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = subtractionNoCarry(cpu, 0, cpu.R[am.A1])
}

func subtractionNoCarry(cpu *Cpu, d, r int) int {
	r = ^r

	hres := (d & 0xf) + (r & 0xf) + 1
	cpu.flags[FlagH] = hres <= 0xf

	res := d + r + 1
	cpu.flags[FlagC] = res < 0
	cpu.flags[FlagV] = (((d ^ r) & 0x80) == 0) && (((d ^ res) & 0x80) != 0)

	res &= 0xff
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]

	return res
}

func Sbc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = subtractionCarry(cpu, cpu.R[am.A1], cpu.R[am.A2],
		cpu.flags[FlagC])
}

func Sbci(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = subtractionCarry(cpu, cpu.R[am.A1], int(am.A2),
		cpu.flags[FlagC])
}

func Cpc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	subtractionCarry(cpu, cpu.R[am.A1], cpu.R[am.A2], cpu.flags[FlagC])
}

func subtractionCarry(cpu *Cpu, d, r int, carry bool) int {
	r = ^r
	c := 1
	if carry {
		c = 0
	}

	hres := (d & 0xf) + (r & 0xf) + c
	cpu.flags[FlagH] = hres <= 0xf

	res := d + r + c
	cpu.flags[FlagC] = res < 0
	cpu.flags[FlagV] = (((d ^ r) & 0x80) == 0) && (((d ^ res) & 0x80) != 0)

	res &= 0xff
	if res != 0 {
		cpu.flags[FlagZ] = false
	}
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]

	return res
}

func Sbiw(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := cpu.R[am.A1] | (cpu.R[am.A1+1] << 8)
	k := int(am.A2)

	res := (d - k) & 0xffff
	hr := (res & 0x8000) != 0
	hd := (d & 0x8000) != 0
	cpu.flags[FlagC] = hr && !hd
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = hr
	cpu.flags[FlagV] = hd && !hr
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]

	cpu.R[am.A1] = res & 0xff
	cpu.R[am.A1+1] = res >> 8
}

func And(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = boolAnd(cpu, cpu.R[am.A1], cpu.R[am.A2])
}

func Andi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = boolAnd(cpu, cpu.R[am.A1], int(am.A2))
}

func boolAnd(cpu *Cpu, d, r int) int {
	res := d & r
	cpu.flags[FlagV] = false
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	return res
}

func Or(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = boolOr(cpu, cpu.R[am.A1], cpu.R[am.A2])
}

func Ori(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = boolOr(cpu, cpu.R[am.A1], int(am.A2))
}

func boolOr(cpu *Cpu, d, r int) int {
	res := d | r
	cpu.flags[FlagV] = false
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	return res
}

func Eor(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := cpu.R[am.A1] ^ cpu.R[am.A2]
	cpu.flags[FlagV] = false
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.R[am.A1] = res
}

func Mul(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := cpu.R[am.A1] * cpu.R[am.A2]
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Fmul(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := cpu.R[am.A1] * cpu.R[am.A2]
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Muls(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.R[am.A1])
	r := int8(cpu.R[am.A2])
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Fmuls(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.R[am.A1])
	r := int8(cpu.R[am.A2])
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Mulsu(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.R[am.A1])
	r := cpu.R[am.A2]
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Fmulsu(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.R[am.A1])
	r := cpu.R[am.A2]
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Mov(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = cpu.R[am.A2]
}

func Movw(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = cpu.R[am.A2]
	cpu.R[am.A1+1] = cpu.R[am.A2+1]
}

func Ldi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.R[am.A1] = int(am.A2)
}

func Com(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := ^cpu.R[am.A1] & 0xff
	cpu.flags[FlagC] = true
	cpu.flags[FlagV] = false
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.R[am.A1] = res
}

func Swap(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.R[am.A1]
	cpu.R[am.A1] = ((val & 0xf) << 4) | ((val & 0xf0) >> 4)
}

func Dec(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := (cpu.R[am.A1] - 1) & 0xff
	cpu.flags[FlagV] = res == 0x7f
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.R[am.A1] = res
}

func Inc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := (cpu.R[am.A1] + 1) & 0xff
	cpu.flags[FlagV] = res == 0x80
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.R[am.A1] = res
}

func Asr(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.R[am.A1]
	res := (val >> 1) | (val & 0x80)
	cpu.flags[FlagC] = (val & 0x1) != 0
	cpu.flags[FlagN] = (val & 0x80) != 0
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagV] = cpu.flags[FlagN] != cpu.flags[FlagC]
	cpu.flags[FlagS] = cpu.flags[FlagN] != cpu.flags[FlagV]
	cpu.R[am.A1] = res
}

func Lsr(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.R[am.A1]
	res := val >> 1
	cpu.flags[FlagC] = (val & 0x1) != 0
	cpu.flags[FlagN] = false
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagV] = cpu.flags[FlagC]
	cpu.flags[FlagS] = cpu.flags[FlagV]
	cpu.R[am.A1] = res
}

func Ror(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.R[am.A1]
	res := val >> 1
	if cpu.flags[FlagC] {
		res |= 0x80
	}
	cpu.flags[FlagN] = cpu.flags[FlagC]
	cpu.flags[FlagC] = (val & 0x1) != 0
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagV] = cpu.flags[FlagN] != cpu.flags[FlagC]
	cpu.flags[FlagS] = cpu.flags[FlagN] != cpu.flags[FlagV]
	cpu.R[am.A1] = res
}

func Brbs(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	if cpu.flags[am.A1] {
		cpu.PC += int(am.A2)
	}
}

func Brbc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	if !cpu.flags[am.A1] {
		cpu.PC += int(am.A2)
	}
}

func Bset(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.flags[am.A1] = true
}

func Bclr(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.flags[am.A1] = false
}

func Bst(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.R[am.A2]
	cpu.flags[FlagT] = (val & (1 << uint(am.A1))) != 0
}

func Bld(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	bit := uint(am.A1)
	if cpu.flags[FlagT] {
		cpu.R[am.A2] |= (1 << bit)
	} else {
		cpu.R[am.A2] &= ^(1 << bit) & 0xff
	}
}
