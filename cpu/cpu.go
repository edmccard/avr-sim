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

type OpFunc func(*Cpu, *instr.AddrMode)

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

func Adiw(cpu *Cpu, am *instr.AddrMode) {
	d := cpu.R[am.A1] | (cpu.R[am.A1+1] << 8)
	k := int(am.A2)

	res := (d + k) & 0xffff
	hr := (res & 0x8000) != 0
	hd := (d & 0x8000) != 0
	cpu.FlagC = !hr && hd
	cpu.FlagZ = res == 0
	cpu.FlagN = hr
	cpu.FlagV = !hd && hr
	cpu.FlagS = cpu.FlagV != cpu.FlagN

	cpu.R[am.A1] = res & 0xff
	cpu.R[am.A1+1] = res >> 8
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

func Neg(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = subtractionNoCarry(cpu, 0, cpu.R[am.A1])
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

func Sbiw(cpu *Cpu, am *instr.AddrMode) {
	d := cpu.R[am.A1] | (cpu.R[am.A1+1] << 8)
	k := int(am.A2)

	res := (d - k) & 0xffff
	hr := (res & 0x8000) != 0
	hd := (d & 0x8000) != 0
	cpu.FlagC = hr && !hd
	cpu.FlagZ = res == 0
	cpu.FlagN = hr
	cpu.FlagV = hd && !hr
	cpu.FlagS = cpu.FlagV != cpu.FlagN

	cpu.R[am.A1] = res & 0xff
	cpu.R[am.A1+1] = res >> 8
}

func And(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = boolAnd(cpu, cpu.R[am.A1], cpu.R[am.A2])
}

func Andi(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = boolAnd(cpu, cpu.R[am.A1], int(am.A2))
}

func boolAnd(cpu *Cpu, d, r int) int {
	res := d & r
	cpu.FlagV = false
	cpu.FlagN = res >= 0x80
	cpu.FlagS = cpu.FlagN
	cpu.FlagZ = res == 0
	return res
}

func Or(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = boolOr(cpu, cpu.R[am.A1], cpu.R[am.A2])
}

func Ori(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = boolOr(cpu, cpu.R[am.A1], int(am.A2))
}

func boolOr(cpu *Cpu, d, r int) int {
	res := d | r
	cpu.FlagV = false
	cpu.FlagN = res >= 0x80
	cpu.FlagS = cpu.FlagN
	cpu.FlagZ = res == 0
	return res
}

func Eor(cpu *Cpu, am *instr.AddrMode) {
	res := cpu.R[am.A1] ^ cpu.R[am.A2]
	cpu.FlagV = false
	cpu.FlagN = res >= 0x80
	cpu.FlagS = cpu.FlagN
	cpu.FlagZ = res == 0
	cpu.R[am.A1] = res
}

func Mul(cpu *Cpu, am *instr.AddrMode) {
	res := cpu.R[am.A1] * cpu.R[am.A2]
	cpu.FlagC = res >= 0x8000
	cpu.FlagZ = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Fmul(cpu *Cpu, am *instr.AddrMode) {
	res := cpu.R[am.A1] * cpu.R[am.A2]
	cpu.FlagC = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.FlagZ = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Muls(cpu *Cpu, am *instr.AddrMode) {
	d := int8(cpu.R[am.A1])
	r := int8(cpu.R[am.A2])
	res := (int(d) * int(r)) & 0xffff
	cpu.FlagC = res >= 0x8000
	cpu.FlagZ = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Fmuls(cpu *Cpu, am *instr.AddrMode) {
	d := int8(cpu.R[am.A1])
	r := int8(cpu.R[am.A2])
	res := (int(d) * int(r)) & 0xffff
	cpu.FlagC = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.FlagZ = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Mulsu(cpu *Cpu, am *instr.AddrMode) {
	d := int8(cpu.R[am.A1])
	r := cpu.R[am.A2]
	res := (int(d) * int(r)) & 0xffff
	cpu.FlagC = res >= 0x8000
	cpu.FlagZ = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Fmulsu(cpu *Cpu, am *instr.AddrMode) {
	d := int8(cpu.R[am.A1])
	r := cpu.R[am.A2]
	res := (int(d) * int(r)) & 0xffff
	cpu.FlagC = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.FlagZ = res == 0
	cpu.R[0] = res & 0xff
	cpu.R[1] = res >> 8
}

func Mov(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = cpu.R[am.A2]
}

func Movw(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = cpu.R[am.A2]
	cpu.R[am.A1+1] = cpu.R[am.A2+1]
}

func Ldi(cpu *Cpu, am *instr.AddrMode) {
	cpu.R[am.A1] = int(am.A2)
}

func Com(cpu *Cpu, am *instr.AddrMode) {
	res := ^cpu.R[am.A1] & 0xff
	cpu.FlagC = true
	cpu.FlagV = false
	cpu.FlagZ = res == 0
	cpu.FlagN = res >= 0x80
	cpu.FlagS = cpu.FlagN
	cpu.R[am.A1] = res
}

func Swap(cpu *Cpu, am *instr.AddrMode) {
	val := cpu.R[am.A1]
	cpu.R[am.A1] = ((val & 0xf) << 4) | ((val & 0xf0) >> 4)
}

func Dec(cpu *Cpu, am *instr.AddrMode) {
	res := (cpu.R[am.A1] - 1) & 0xff
	cpu.FlagV = res == 0x7f
	cpu.FlagN = res >= 0x80
	cpu.FlagS = cpu.FlagV != cpu.FlagN
	cpu.FlagZ = res == 0
	cpu.R[am.A1] = res
}

func Inc(cpu *Cpu, am *instr.AddrMode) {
	res := (cpu.R[am.A1] + 1) & 0xff
	cpu.FlagV = res == 0x80
	cpu.FlagN = res >= 0x80
	cpu.FlagS = cpu.FlagV != cpu.FlagN
	cpu.FlagZ = res == 0
	cpu.R[am.A1] = res
}

func Asr(cpu *Cpu, am *instr.AddrMode) {
	val := cpu.R[am.A1]
	res := (val >> 1) | (val & 0x80)
	cpu.FlagC = (val & 0x1) != 0
	cpu.FlagN = (val & 0x80) != 0
	cpu.FlagZ = res == 0
	cpu.FlagV = cpu.FlagN != cpu.FlagC
	cpu.FlagS = cpu.FlagN != cpu.FlagV
	cpu.R[am.A1] = res
}

func Lsr(cpu *Cpu, am *instr.AddrMode) {
	val := cpu.R[am.A1]
	res := val >> 1
	cpu.FlagC = (val & 0x1) != 0
	cpu.FlagN = false
	cpu.FlagZ = res == 0
	cpu.FlagV = cpu.FlagC
	cpu.FlagS = cpu.FlagV
	cpu.R[am.A1] = res
}

func Ror(cpu *Cpu, am *instr.AddrMode) {
	val := cpu.R[am.A1]
	res := val >> 1
	if cpu.FlagC {
		res |= 0x80
	}
	cpu.FlagC = (val & 0x1) != 0
	cpu.FlagN = false
	cpu.FlagZ = res == 0
	cpu.FlagV = cpu.FlagC
	cpu.FlagS = cpu.FlagV
	cpu.R[am.A1] = res
}
