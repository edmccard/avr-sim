package cpu

import "github.com/edmccard/avr-sim/instr"

type Cpu struct {
	reg   [32]int
	flags [8]bool
	sp    int
	pc    int
	ramp  [5]int // D,X,Y,Z,EIND
	rmask [5]int
	skip  bool
	am    instr.AddrMode
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

func (c *Cpu) Reset(sp int, pc int) {
	c.pc = pc
	c.sp = sp
	for i := range c.flags {
		c.flags[i] = false
	}
}

func (c *Cpu) Step(mem Memory, d *instr.Decoder) {
	op, op2, mnem := c.fetch(mem, d)
	d.DecodeAddr(&c.am, mnem, op, op2)
	opFuncs[mnem](c, &c.am, mem)
	if c.skip {
		c.skip = false
		op, op2, mnem = c.fetch(mem, d)
	}
}

func (c *Cpu) fetch(mem Memory, d *instr.Decoder) (instr.Opcode, instr.Opcode, instr.Mnemonic) {
	var op2 instr.Opcode
	op := instr.Opcode(mem.ReadProgram(instr.Addr(c.pc)))
	c.pcInc(1)
	mnem, ln := d.DecodeMnem(op)
	if ln == 2 {
		op2 = instr.Opcode(mem.ReadProgram(instr.Addr(c.pc)))
		c.pcInc(1)
	}
	return op, op2, mnem
}

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

func nop(cpu *Cpu, am *instr.AddrMode, mem Memory) {
}

func adc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addition(cpu, am, cpu.flags[FlagC])
}

func add(cpu *Cpu, am *instr.AddrMode, mem Memory) {
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

func adiw(cpu *Cpu, am *instr.AddrMode, mem Memory) {
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

func sub(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = subtractionNoCarry(cpu, cpu.reg[am.A1], cpu.reg[am.A2])
}

func subi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = subtractionNoCarry(cpu, cpu.reg[am.A1], int(am.A2))
}

func cp(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	subtractionNoCarry(cpu, cpu.reg[am.A1], cpu.reg[am.A2])
}

func cpi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	subtractionNoCarry(cpu, cpu.reg[am.A1], int(am.A2))
}

func neg(cpu *Cpu, am *instr.AddrMode, mem Memory) {
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

func sbc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = subtractionCarry(cpu, cpu.reg[am.A1], cpu.reg[am.A2],
		cpu.flags[FlagC])
}

func sbci(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = subtractionCarry(cpu, cpu.reg[am.A1], int(am.A2),
		cpu.flags[FlagC])
}

func cpc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
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

func sbiw(cpu *Cpu, am *instr.AddrMode, mem Memory) {
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

func and(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = boolAnd(cpu, cpu.reg[am.A1], cpu.reg[am.A2])
}

func andi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
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

func or(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = boolOr(cpu, cpu.reg[am.A1], cpu.reg[am.A2])
}

func ori(cpu *Cpu, am *instr.AddrMode, mem Memory) {
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

func eor(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := cpu.reg[am.A1] ^ cpu.reg[am.A2]
	cpu.flags[FlagV] = false
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.reg[am.A1] = res
}

func mul(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := cpu.reg[am.A1] * cpu.reg[am.A2]
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func fmul(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := cpu.reg[am.A1] * cpu.reg[am.A2]
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func muls(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.reg[am.A1])
	r := int8(cpu.reg[am.A2])
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func fmuls(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.reg[am.A1])
	r := int8(cpu.reg[am.A2])
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func mulsu(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.reg[am.A1])
	r := cpu.reg[am.A2]
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func fmulsu(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	d := int8(cpu.reg[am.A1])
	r := cpu.reg[am.A2]
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func mov(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = cpu.reg[am.A2]
}

func movw(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = cpu.reg[am.A2]
	cpu.reg[am.A1+1] = cpu.reg[am.A2+1]
}

func ldi(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = int(am.A2)
}

func com(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := ^cpu.reg[am.A1] & 0xff
	cpu.flags[FlagC] = true
	cpu.flags[FlagV] = false
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.reg[am.A1] = res
}

func swap(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.reg[am.A1]
	cpu.reg[am.A1] = ((val & 0xf) << 4) | ((val & 0xf0) >> 4)
}

func dec(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := (cpu.reg[am.A1] - 1) & 0xff
	cpu.flags[FlagV] = res == 0x7f
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.reg[am.A1] = res
}

func inc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	res := (cpu.reg[am.A1] + 1) & 0xff
	cpu.flags[FlagV] = res == 0x80
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.reg[am.A1] = res
}

func asr(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.reg[am.A1]
	res := (val >> 1) | (val & 0x80)
	cpu.flags[FlagC] = (val & 0x1) != 0
	cpu.flags[FlagN] = (val & 0x80) != 0
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagV] = cpu.flags[FlagN] != cpu.flags[FlagC]
	cpu.flags[FlagS] = cpu.flags[FlagN] != cpu.flags[FlagV]
	cpu.reg[am.A1] = res
}

func lsr(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.reg[am.A1]
	res := val >> 1
	cpu.flags[FlagC] = (val & 0x1) != 0
	cpu.flags[FlagN] = false
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagV] = cpu.flags[FlagC]
	cpu.flags[FlagS] = cpu.flags[FlagV]
	cpu.reg[am.A1] = res
}

func ror(cpu *Cpu, am *instr.AddrMode, mem Memory) {
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

func bset(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.flags[am.A1] = true
}

func bclr(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.flags[am.A1] = false
}

func bst(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	val := cpu.reg[am.A2]
	cpu.flags[FlagT] = (val & (1 << uint(am.A1))) != 0
}

func bld(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	bit := uint(am.A1)
	if cpu.flags[FlagT] {
		cpu.reg[am.A2] |= (1 << bit)
	} else {
		cpu.reg[am.A2] &= ^(1 << bit) & 0xff
	}
}

func ld(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addr := cpu.indirect(am.Ireg, 0)
	cpu.reg[am.A1] = int(mem.ReadData(addr))
}

func st(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addr := cpu.indirect(am.Ireg, 0)
	mem.WriteData(addr, byte(cpu.reg[am.A1]))
}

func ldd(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addr := cpu.indirect(am.Ireg, am.A2)
	cpu.reg[am.A1] = int(mem.ReadData(addr))
}

func std(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	addr := cpu.indirect(am.Ireg, am.A2)
	mem.WriteData(addr, byte(cpu.reg[am.A1]))
}

func lds(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.reg[am.A1] = int(mem.ReadData(instr.Addr(cpu.ramp[RampD]) | am.A2))
}

func sts(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	mem.WriteData(instr.Addr(cpu.ramp[RampD])|am.A2, byte(cpu.reg[am.A1]))
}

func push(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	mem.WriteData(instr.Addr(cpu.sp), byte(cpu.reg[am.A1]))
	cpu.spInc(-1)
}

func pop(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.spInc(1)
	cpu.reg[am.A1] = int(mem.ReadData(instr.Addr(cpu.sp)))
}

func brbs(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	if cpu.flags[am.A1] {
		cpu.pcInc(int(am.A2))
	}
}

func brbc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	if !cpu.flags[am.A1] {
		cpu.pcInc(int(am.A2))
	}
}

func eijmp(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.pc = cpu.reg[30] | (cpu.reg[31] << 8) | cpu.ramp[Eind]
}

func ijmp(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.pc = cpu.reg[30] | (cpu.reg[31] << 8)
}

func jmp(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.pc = int(am.A1) & (cpu.rmask[Eind] | 0xffff)
}

func rjmp(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	cpu.pcInc(int(am.A1))
}

func call(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	pushPC(cpu, mem)
	jmp(cpu, am, mem)
}

func eicall(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	pushPC(cpu, mem)
	eijmp(cpu, am, mem)
}

func icall(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	pushPC(cpu, mem)
	ijmp(cpu, am, mem)
}

func rcall(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	pushPC(cpu, mem)
	rjmp(cpu, am, mem)
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

func ret(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	popPC(cpu, mem)
}

func reti(cpu *Cpu, am *instr.AddrMode, mem Memory) {
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

func lac(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	mask := byte(cpu.reg[am.A1])
	addr := cpu.indirect(instr.Z, 0)
	cpu.reg[am.A1] = int(mem.ReadData(addr))
	mem.WriteData(addr, byte(cpu.reg[am.A1]) & ^mask)
}

func las(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	mask := byte(cpu.reg[am.A1])
	addr := cpu.indirect(instr.Z, 0)
	cpu.reg[am.A1] = int(mem.ReadData(addr))
	mem.WriteData(addr, byte(cpu.reg[am.A1])|mask)
}

func lat(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	mask := byte(cpu.reg[am.A1])
	addr := cpu.indirect(instr.Z, 0)
	cpu.reg[am.A1] = int(mem.ReadData(addr))
	mem.WriteData(addr, byte(cpu.reg[am.A1])^mask)
}

func xch(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	tmp := byte(cpu.reg[am.A1])
	addr := cpu.indirect(instr.Z, 0)
	cpu.reg[am.A1] = int(mem.ReadData(addr))
	mem.WriteData(addr, tmp)
}

func sbrc(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	if (cpu.reg[am.A2] & (1 << uint(am.A1))) == 0 {
		cpu.skip = true
	}
}

func sbrs(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	if (cpu.reg[am.A2] & (1 << uint(am.A1))) != 0 {
		cpu.skip = true
	}
}

func cpse(cpu *Cpu, am *instr.AddrMode, mem Memory) {
	if cpu.reg[am.A1] == cpu.reg[am.A2] {
		cpu.skip = true
	}
}

type OpFunc func(*Cpu, *instr.AddrMode, Memory)

var opFuncs = [...]OpFunc{
	nop,    // Reserved
	adc,    // Adc
	adc,    // AdcReduced
	add,    // Add
	add,    // AddReduced
	adiw,   // Adiw
	and,    // And
	and,    // AndReduced
	andi,   // Andi
	asr,    // Asr
	asr,    // AsrReduced
	bclr,   // Bclr
	bld,    // Bld
	bld,    // BldReduced
	brbc,   // Brbc
	brbs,   // Brbs
	nop,    // Break ****
	bset,   // Bset
	bst,    // Bst
	bst,    // BstReduced
	call,   // Call
	nop,    // Cbi ****
	com,    // Com
	com,    // ComReduced
	cp,     // Cp
	cp,     // CpReduced
	cpc,    // Cpc
	cpc,    // CpcReduced
	cpi,    // Cpi
	nop,    // Cpse ****
	nop,    // CpseReduced ****
	dec,    // Dec
	dec,    // DecReduced
	nop,    // Des ****
	eicall, // Eicall
	eijmp,  // Eijmp
	nop,    // Elpm ****
	nop,    // ElpmEnhanced ****
	eor,    // Eor
	eor,    // EorReduced
	fmul,   // Fmul
	fmuls,  // Fmuls
	fmulsu, // Fmulsu
	icall,  // Icall
	ijmp,   // Ijmp
	nop,    // In ****
	nop,    // InReduced ****
	inc,    // Inc
	inc,    // IncReduced
	jmp,    // Jmp
	lac,    // Lac
	las,    // Las
	lat,    // Lat
	ld,     // LdClassic
	ld,     // LdClassicReduced
	ld,     // LdMinimal
	ld,     // LdMinimalReduced
	ldd,    // Ldd
	ldi,    // Ldi
	lds,    // Lds
	nop,    // Lds16 ****
	nop,    // Lpm ****
	nop,    // LpmEnhanced ****
	lsr,    // Lsr
	lsr,    // LsrReduced
	mov,    // Mov
	mov,    // MovReduced
	movw,   // Movw
	mul,    // Mul
	muls,   // Muls
	mulsu,  // Mulsu
	neg,    // Neg
	neg,    // NegReduced
	nop,    // Nop
	or,     // Or
	or,     // OrReduced
	ori,    // Ori
	nop,    // Out ****
	nop,    // OutReduced ****
	pop,    // Pop
	pop,    // PopReduced
	push,   // Push
	push,   // PushReduced
	rcall,  // Rcall
	ret,    // Ret
	reti,   // Reti
	rjmp,   // Rjmp
	ror,    // Ror
	ror,    // RorReduced
	sbc,    // Sbc
	sbc,    // SbcReduced
	sbci,   // Sbci
	nop,    // Sbi ****
	nop,    // Sbic ****
	nop,    // Sbis ****
	sbiw,   // Sbiw
	sbrc,   // Sbrc
	sbrc,   // SbrcReduced
	sbrs,   // Sbrs
	sbrs,   // SbrsReduced
	nop,    // Sleep ****
	nop,    // Spm ****
	nop,    // SpmXmega ****
	st,     // StClassic
	st,     // StClassicReduced
	st,     // StMinimal
	st,     // StMinimalReduced
	std,    // Std
	sts,    // Sts
	nop,    // Sts16 ****
	sub,    // Sub
	sub,    // SubReduced
	subi,   // Subi
	swap,   // Swap
	swap,   // SwapReduced
	nop,    // Wdr ****
	xch,    // Xch
}
