package cpu

import "github.com/edmccard/avr-sim/instr"

type Cpu struct {
	reg   [32]int
	flags [8]bool
	sp    int
	pc    int
	ramp  [5]int // D,X,Y,Z,EIND
	rmask [5]int
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

type Ramp int

const (
	RampD Ramp = iota
	RampX
	RampY
	RampZ
	Eind
)

func (c *Cpu) GetReg(r instr.Addr) byte {
	return byte(c.reg[r])
}

func (c *Cpu) SetReg(r instr.Addr, val byte) {
	c.reg[r] = int(val)
}

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

func (c *Cpu) SetRamp(reg Ramp, val byte) {
	c.ramp[reg] = (int(val) << 16) & c.rmask[reg]
}

func (c *Cpu) GetRamp(reg Ramp) byte {
	return byte((c.ramp[reg] & c.rmask[reg]) >> 16)
}

func (c *Cpu) setRmask(reg Ramp, mask byte) {
	c.rmask[reg] = (int(mask) << 16)
}

func (c *Cpu) indirect(ireg instr.IndexReg, q instr.Addr) instr.Addr {
	base := ireg.Base()
	mode := ireg.Mode()
	r := 24 + base*2
	addr := c.ramp[base] | c.reg[r] | (c.reg[r+1] << 8)
	switch mode {
	case instr.NoMode:
		addr = (addr + int(q)) & (c.rmask[base] | 0xffff)
	case instr.PreDec:
		addr = (addr - 1) & (c.rmask[base] | 0xffff)
		c.reg[r] = addr & 0xff
		c.reg[r+1] = (addr >> 8) & 0xff
		c.ramp[base] = addr >> 16
	case instr.PostInc:
		a2 := (addr + 1) & (c.rmask[base] | 0xffff)
		c.reg[r] = a2 & 0xff
		c.reg[r+1] = (a2 >> 8) & 0xff
		c.ramp[base] = a2 >> 16
	}
	return instr.Addr(addr)
}

func (c *Cpu) spInc(offset int) {
	c.sp = (c.sp + offset) & 0xffff
}

func (c *Cpu) pcInc(offset int) {
	c.pc = (c.pc + offset) & (c.rmask[Eind] | 0xffff)
}

type OpFunc func(*Cpu, *instr.AddrMode, Memory)

func Nop(cpu *Cpu, am *instr.AddrMode, mem Memory) {
}

func Adc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addition(cpu, am, cpu.flags[FlagC])
}

func Add(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addition(cpu, am, false)
}

func addition(cpu *Cpu, am *instr.AddrMode, carry bool) {
	d := cpu.reg[am.A1]
	r := cpu.reg[am.A2]
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

	cpu.reg[am.A1] = res
}

func Adiw(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := cpu.reg[am.A1] | (cpu.reg[am.A1+1] << 8)
	k := int(am.A2)

	res := (d + k) & 0xffff
	hr := (res & 0x8000) != 0
	hd := (d & 0x8000) != 0
	cpu.flags[FlagC] = !hr && hd
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = hr
	cpu.flags[FlagV] = !hd && hr
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]

	cpu.reg[am.A1] = res & 0xff
	cpu.reg[am.A1+1] = res >> 8
}

func Sub(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = subtractionNoCarry(cpu, cpu.reg[am.A1], cpu.reg[am.A2])
}

func Subi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = subtractionNoCarry(cpu, cpu.reg[am.A1], int(am.A2))
}

func Cp(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	subtractionNoCarry(cpu, cpu.reg[am.A1], cpu.reg[am.A2])
}

func Cpi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	subtractionNoCarry(cpu, cpu.reg[am.A1], int(am.A2))
}

func Neg(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = subtractionNoCarry(cpu, 0, cpu.reg[am.A1])
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
	cpu.reg[am.A1] = subtractionCarry(cpu, cpu.reg[am.A1], cpu.reg[am.A2],
		cpu.flags[FlagC])
}

func Sbci(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = subtractionCarry(cpu, cpu.reg[am.A1], int(am.A2),
		cpu.flags[FlagC])
}

func Cpc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	subtractionCarry(cpu, cpu.reg[am.A1], cpu.reg[am.A2], cpu.flags[FlagC])
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
	d := cpu.reg[am.A1] | (cpu.reg[am.A1+1] << 8)
	k := int(am.A2)

	res := (d - k) & 0xffff
	hr := (res & 0x8000) != 0
	hd := (d & 0x8000) != 0
	cpu.flags[FlagC] = hr && !hd
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = hr
	cpu.flags[FlagV] = hd && !hr
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]

	cpu.reg[am.A1] = res & 0xff
	cpu.reg[am.A1+1] = res >> 8
}

func And(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = boolAnd(cpu, cpu.reg[am.A1], cpu.reg[am.A2])
}

func Andi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = boolAnd(cpu, cpu.reg[am.A1], int(am.A2))
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
	cpu.reg[am.A1] = boolOr(cpu, cpu.reg[am.A1], cpu.reg[am.A2])
}

func Ori(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = boolOr(cpu, cpu.reg[am.A1], int(am.A2))
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
	res := cpu.reg[am.A1] ^ cpu.reg[am.A2]
	cpu.flags[FlagV] = false
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.reg[am.A1] = res
}

func Mul(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := cpu.reg[am.A1] * cpu.reg[am.A2]
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func Fmul(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := cpu.reg[am.A1] * cpu.reg[am.A2]
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func Muls(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.reg[am.A1])
	r := int8(cpu.reg[am.A2])
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func Fmuls(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.reg[am.A1])
	r := int8(cpu.reg[am.A2])
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func Mulsu(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.reg[am.A1])
	r := cpu.reg[am.A2]
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func Fmulsu(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.reg[am.A1])
	r := cpu.reg[am.A2]
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func Mov(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = cpu.reg[am.A2]
}

func Movw(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = cpu.reg[am.A2]
	cpu.reg[am.A1+1] = cpu.reg[am.A2+1]
}

func Ldi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = int(am.A2)
}

func Com(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := ^cpu.reg[am.A1] & 0xff
	cpu.flags[FlagC] = true
	cpu.flags[FlagV] = false
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.reg[am.A1] = res
}

func Swap(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.reg[am.A1]
	cpu.reg[am.A1] = ((val & 0xf) << 4) | ((val & 0xf0) >> 4)
}

func Dec(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := (cpu.reg[am.A1] - 1) & 0xff
	cpu.flags[FlagV] = res == 0x7f
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.reg[am.A1] = res
}

func Inc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := (cpu.reg[am.A1] + 1) & 0xff
	cpu.flags[FlagV] = res == 0x80
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.reg[am.A1] = res
}

func Asr(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.reg[am.A1]
	res := (val >> 1) | (val & 0x80)
	cpu.flags[FlagC] = (val & 0x1) != 0
	cpu.flags[FlagN] = (val & 0x80) != 0
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagV] = cpu.flags[FlagN] != cpu.flags[FlagC]
	cpu.flags[FlagS] = cpu.flags[FlagN] != cpu.flags[FlagV]
	cpu.reg[am.A1] = res
}

func Lsr(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.reg[am.A1]
	res := val >> 1
	cpu.flags[FlagC] = (val & 0x1) != 0
	cpu.flags[FlagN] = false
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagV] = cpu.flags[FlagC]
	cpu.flags[FlagS] = cpu.flags[FlagV]
	cpu.reg[am.A1] = res
}

func Ror(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.reg[am.A1]
	res := val >> 1
	if cpu.flags[FlagC] {
		res |= 0x80
	}
	cpu.flags[FlagN] = cpu.flags[FlagC]
	cpu.flags[FlagC] = (val & 0x1) != 0
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagV] = cpu.flags[FlagN] != cpu.flags[FlagC]
	cpu.flags[FlagS] = cpu.flags[FlagN] != cpu.flags[FlagV]
	cpu.reg[am.A1] = res
}

func Bset(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.flags[am.A1] = true
}

func Bclr(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.flags[am.A1] = false
}

func Bst(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.reg[am.A2]
	cpu.flags[FlagT] = (val & (1 << uint(am.A1))) != 0
}

func Bld(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	bit := uint(am.A1)
	if cpu.flags[FlagT] {
		cpu.reg[am.A2] |= (1 << bit)
	} else {
		cpu.reg[am.A2] &= ^(1 << bit) & 0xff
	}
}

func Ld(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addr := cpu.indirect(am.Ireg, 0)
	cpu.reg[am.A1] = int(mem.ReadData(addr))
}

func St(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addr := cpu.indirect(am.Ireg, 0)
	mem.WriteData(addr, byte(cpu.reg[am.A1]))
}

func Ldd(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addr := cpu.indirect(am.Ireg, am.A2)
	cpu.reg[am.A1] = int(mem.ReadData(addr))
}

func Std(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addr := cpu.indirect(am.Ireg, am.A2)
	mem.WriteData(addr, byte(cpu.reg[am.A1]))
}

func Lds(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = int(mem.ReadData(instr.Addr(cpu.ramp[RampD]) | am.A2))
}

func Sts(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	mem.WriteData(instr.Addr(cpu.ramp[RampD])|am.A2, byte(cpu.reg[am.A1]))
}

func Push(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	mem.WriteData(instr.Addr(cpu.sp), byte(cpu.reg[am.A1]))
	cpu.spInc(-1)
}

func Pop(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.spInc(1)
	cpu.reg[am.A1] = int(mem.ReadData(instr.Addr(cpu.sp)))
}

func Brbs(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	if cpu.flags[am.A1] {
		cpu.pcInc(int(am.A2))
	}
}

func Brbc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	if !cpu.flags[am.A1] {
		cpu.pcInc(int(am.A2))
	}
}

func Eijmp(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.pc = cpu.reg[30] | (cpu.reg[31] << 8) | cpu.ramp[Eind]
}

func Ijmp(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.pc = cpu.reg[30] | (cpu.reg[31] << 8)
}

func Jmp(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.pc = int(am.A1) & (cpu.rmask[Eind] | 0xffff)
}

func Rjmp(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.pcInc(int(am.A1))
}

func Call(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	pushPC(cpu, mem)
	Jmp(cpu, am, mem)
}

func Eicall(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	pushPC(cpu, mem)
	Eijmp(cpu, am, mem)
}

func Icall(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	pushPC(cpu, mem)
	Ijmp(cpu, am, mem)
}

func Rcall(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	pushPC(cpu, mem)
	Rjmp(cpu, am, mem)
}

func pushPC(cpu *Cpu, mem Memory) {
	mem.WriteData(instr.Addr(cpu.sp), byte(cpu.pc))
	cpu.spInc(-1)
	mem.WriteData(instr.Addr(cpu.sp), byte(cpu.pc>>8))
	cpu.spInc(-1)
	if cpu.rmask[Eind] != 0 {
		mem.WriteData(instr.Addr(cpu.sp), byte(cpu.pc>>16))
		cpu.spInc(-1)
	}
}

func Ret(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	popPC(cpu, mem)
}

func Reti(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	popPC(cpu, mem)
	cpu.flags[FlagI] = true
}

func popPC(cpu *Cpu, mem Memory) {
	cpu.pc = 0
	if cpu.rmask[Eind] != 0 {
		cpu.spInc(1)
		cpu.pc |= (int(mem.ReadData(instr.Addr(cpu.sp))) << 16)
	}
	cpu.spInc(1)
	cpu.pc |= (int(mem.ReadData(instr.Addr(cpu.sp))) << 8)
	cpu.spInc(1)
	cpu.pc |= int(mem.ReadData(instr.Addr(cpu.sp)))
}

func Lac(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	// mask = reg
	// reg = (z)
	// (z) &= ^mask
}

func Las(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	// mask = reg
	// reg = (z)
	// (z) |= ^mask

}

func Lat(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	// mask = reg
	// temp = (Z)
	// reg = tmp
	// (Z) ^= mask
}
