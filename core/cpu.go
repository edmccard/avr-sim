package core

import "github.com/edmccard/avr-sim/instr"

type Cpu struct {
	reg   [32]int
	flags [8]bool
	sp    int
	pc    int
	ramp  [5]int // D,X,Y,Z,EIND
	rmask [5]int
	skip  bool
	ops   instr.Operands
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
	d.DecodeOperands(&c.ops, mnem, op, op2)
	opFuncs[mnem](c, &c.ops, mem)
	if c.skip {
		c.skip = false
		c.fetch(mem, d)
	}
}

func (c *Cpu) fetch(mem Memory, d *instr.Decoder) (instr.Opcode, instr.Opcode,
	instr.Mnemonic) {

	var op2 instr.Opcode
	op := instr.Opcode(mem.ReadProgram(Addr(c.pc)))
	c.pcInc(1)
	mnem, ln := d.DecodeMnem(op)
	if ln == 2 {
		op2 = instr.Opcode(mem.ReadProgram(Addr(c.pc)))
		c.pcInc(1)
	}
	return op, op2, mnem
}

func (c *Cpu) GetReg(r int) byte {
	return byte(c.reg[r])
}

func (c *Cpu) SetReg(r int, val byte) {
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

func (c *Cpu) indirect(ireg instr.IndexReg, q int) int {
	base := ireg.Base()
	action := ireg.Action()
	r := 24 + base*2
	addr := c.ramp[base] | c.reg[r] | (c.reg[r+1] << 8)
	switch action {
	case instr.NoAction:
		addr = (addr + int(q)) & (c.rmask[base] | 0xffff)
	case instr.PreDec:
		addr = (addr - 1) & (c.rmask[base] | 0xffff)
		c.reg[r] = addr & 0xff
		c.reg[r+1] = (addr >> 8) & 0xff
		c.ramp[base] = addr & c.rmask[base]
	case instr.PostInc:
		a2 := (addr + 1) & (c.rmask[base] | 0xffff)
		c.reg[r] = a2 & 0xff
		c.reg[r+1] = (a2 >> 8) & 0xff
		c.ramp[base] = a2 & c.rmask[base]
	}
	return addr
}

func (c *Cpu) spInc(offset int) {
	c.sp = (c.sp + offset) & 0xffff
}

func (c *Cpu) pcInc(offset int) {
	c.pc = (c.pc + offset) & (c.rmask[Eind] | 0xffff)
}

func nop(cpu *Cpu, o *instr.Operands, mem Memory) {
}

func adc(cpu *Cpu, o *instr.Operands, mem Memory) {
	addition(cpu, o, cpu.flags[FlagC])
}

func add(cpu *Cpu, o *instr.Operands, mem Memory) {
	addition(cpu, o, false)
}

func addition(cpu *Cpu, o *instr.Operands, carry bool) {
	d := cpu.reg[o.Dst]
	r := cpu.reg[o.Src]
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

	cpu.reg[o.Dst] = res
}

func adiw(cpu *Cpu, o *instr.Operands, mem Memory) {
	d := cpu.reg[o.Dst] | (cpu.reg[o.Dst+1] << 8)
	k := o.Src

	res := (d + k) & 0xffff
	hr := (res & 0x8000) != 0
	hd := (d & 0x8000) != 0
	cpu.flags[FlagC] = !hr && hd
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = hr
	cpu.flags[FlagV] = !hd && hr
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]

	cpu.reg[o.Dst] = res & 0xff
	cpu.reg[o.Dst+1] = res >> 8
}

func sub(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = subtractionNoCarry(cpu, cpu.reg[o.Dst], cpu.reg[o.Src])
}

func subi(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = subtractionNoCarry(cpu, cpu.reg[o.Dst], o.Src)
}

func cp(cpu *Cpu, o *instr.Operands, mem Memory) {
	subtractionNoCarry(cpu, cpu.reg[o.Dst], cpu.reg[o.Src])
}

func cpi(cpu *Cpu, o *instr.Operands, mem Memory) {
	subtractionNoCarry(cpu, cpu.reg[o.Dst], o.Src)
}

func neg(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = subtractionNoCarry(cpu, 0, cpu.reg[o.Dst])
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

func sbc(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = subtractionCarry(cpu, cpu.reg[o.Dst], cpu.reg[o.Src],
		cpu.flags[FlagC])
}

func sbci(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = subtractionCarry(cpu, cpu.reg[o.Dst], o.Src,
		cpu.flags[FlagC])
}

func cpc(cpu *Cpu, o *instr.Operands, mem Memory) {
	subtractionCarry(cpu, cpu.reg[o.Dst], cpu.reg[o.Src], cpu.flags[FlagC])
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

func sbiw(cpu *Cpu, o *instr.Operands, mem Memory) {
	d := cpu.reg[o.Dst] | (cpu.reg[o.Dst+1] << 8)
	k := o.Src

	res := (d - k) & 0xffff
	hr := (res & 0x8000) != 0
	hd := (d & 0x8000) != 0
	cpu.flags[FlagC] = hr && !hd
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = hr
	cpu.flags[FlagV] = hd && !hr
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]

	cpu.reg[o.Dst] = res & 0xff
	cpu.reg[o.Dst+1] = res >> 8
}

func and(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = boolAnd(cpu, cpu.reg[o.Dst], cpu.reg[o.Src])
}

func andi(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = boolAnd(cpu, cpu.reg[o.Dst], o.Src)
}

func boolAnd(cpu *Cpu, d, r int) int {
	res := d & r
	cpu.flags[FlagV] = false
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	return res
}

func or(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = boolOr(cpu, cpu.reg[o.Dst], cpu.reg[o.Src])
}

func ori(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = boolOr(cpu, cpu.reg[o.Dst], o.Src)
}

func boolOr(cpu *Cpu, d, r int) int {
	res := d | r
	cpu.flags[FlagV] = false
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	return res
}

func eor(cpu *Cpu, o *instr.Operands, mem Memory) {
	res := cpu.reg[o.Dst] ^ cpu.reg[o.Src]
	cpu.flags[FlagV] = false
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.reg[o.Dst] = res
}

func mul(cpu *Cpu, o *instr.Operands, mem Memory) {
	res := cpu.reg[o.Dst] * cpu.reg[o.Src]
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func fmul(cpu *Cpu, o *instr.Operands, mem Memory) {
	res := cpu.reg[o.Dst] * cpu.reg[o.Src]
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func muls(cpu *Cpu, o *instr.Operands, mem Memory) {
	d := int8(cpu.reg[o.Dst])
	r := int8(cpu.reg[o.Src])
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func fmuls(cpu *Cpu, o *instr.Operands, mem Memory) {
	d := int8(cpu.reg[o.Dst])
	r := int8(cpu.reg[o.Src])
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func mulsu(cpu *Cpu, o *instr.Operands, mem Memory) {
	d := int8(cpu.reg[o.Dst])
	r := cpu.reg[o.Src]
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func fmulsu(cpu *Cpu, o *instr.Operands, mem Memory) {
	d := int8(cpu.reg[o.Dst])
	r := cpu.reg[o.Src]
	res := (int(d) * int(r)) & 0xffff
	cpu.flags[FlagC] = res >= 0x8000
	res = (res << 1) & 0xffff
	cpu.flags[FlagZ] = res == 0
	cpu.reg[0] = res & 0xff
	cpu.reg[1] = res >> 8
}

func mov(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = cpu.reg[o.Src]
}

func movw(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = cpu.reg[o.Src]
	cpu.reg[o.Dst+1] = cpu.reg[o.Src+1]
}

func ldi(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = o.Src
}

func com(cpu *Cpu, o *instr.Operands, mem Memory) {
	res := ^cpu.reg[o.Dst] & 0xff
	cpu.flags[FlagC] = true
	cpu.flags[FlagV] = false
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagN]
	cpu.reg[o.Dst] = res
}

func swap(cpu *Cpu, o *instr.Operands, mem Memory) {
	val := cpu.reg[o.Dst]
	cpu.reg[o.Dst] = ((val & 0xf) << 4) | ((val & 0xf0) >> 4)
}

func dec(cpu *Cpu, o *instr.Operands, mem Memory) {
	res := (cpu.reg[o.Dst] - 1) & 0xff
	cpu.flags[FlagV] = res == 0x7f
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.reg[o.Dst] = res
}

func inc(cpu *Cpu, o *instr.Operands, mem Memory) {
	res := (cpu.reg[o.Dst] + 1) & 0xff
	cpu.flags[FlagV] = res == 0x80
	cpu.flags[FlagN] = res >= 0x80
	cpu.flags[FlagS] = cpu.flags[FlagV] != cpu.flags[FlagN]
	cpu.flags[FlagZ] = res == 0
	cpu.reg[o.Dst] = res
}

func asr(cpu *Cpu, o *instr.Operands, mem Memory) {
	val := cpu.reg[o.Dst]
	res := (val >> 1) | (val & 0x80)
	cpu.flags[FlagC] = (val & 0x1) != 0
	cpu.flags[FlagN] = (val & 0x80) != 0
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagV] = cpu.flags[FlagN] != cpu.flags[FlagC]
	cpu.flags[FlagS] = cpu.flags[FlagN] != cpu.flags[FlagV]
	cpu.reg[o.Dst] = res
}

func lsr(cpu *Cpu, o *instr.Operands, mem Memory) {
	val := cpu.reg[o.Dst]
	res := val >> 1
	cpu.flags[FlagC] = (val & 0x1) != 0
	cpu.flags[FlagN] = false
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagV] = cpu.flags[FlagC]
	cpu.flags[FlagS] = cpu.flags[FlagV]
	cpu.reg[o.Dst] = res
}

func ror(cpu *Cpu, o *instr.Operands, mem Memory) {
	val := cpu.reg[o.Dst]
	res := val >> 1
	if cpu.flags[FlagC] {
		res |= 0x80
	}
	cpu.flags[FlagN] = cpu.flags[FlagC]
	cpu.flags[FlagC] = (val & 0x1) != 0
	cpu.flags[FlagZ] = res == 0
	cpu.flags[FlagV] = cpu.flags[FlagN] != cpu.flags[FlagC]
	cpu.flags[FlagS] = cpu.flags[FlagN] != cpu.flags[FlagV]
	cpu.reg[o.Dst] = res
}

func bset(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.flags[o.Dst] = true
}

func bclr(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.flags[o.Src] = false
}

func bst(cpu *Cpu, o *instr.Operands, mem Memory) {
	val := cpu.reg[o.Src]
	cpu.flags[FlagT] = (val & (1 << uint(o.Off))) != 0
}

func bld(cpu *Cpu, o *instr.Operands, mem Memory) {
	bit := uint(o.Off)
	if cpu.flags[FlagT] {
		cpu.reg[o.Dst] |= (1 << bit)
	} else {
		cpu.reg[o.Dst] &= ^(1 << bit) & 0xff
	}
}

func ld(cpu *Cpu, o *instr.Operands, mem Memory) {
	addr := Addr(cpu.indirect(instr.IndexReg(o.Src), 0))
	cpu.reg[o.Dst] = int(mem.ReadData(addr))
}

func st(cpu *Cpu, o *instr.Operands, mem Memory) {
	addr := Addr(cpu.indirect(instr.IndexReg(o.Dst), 0))
	mem.WriteData(addr, byte(cpu.reg[o.Src]))
}

func ldd(cpu *Cpu, o *instr.Operands, mem Memory) {
	addr := Addr(cpu.indirect(instr.IndexReg(o.Src), o.Off))
	cpu.reg[o.Dst] = int(mem.ReadData(addr))
}

func std(cpu *Cpu, o *instr.Operands, mem Memory) {
	addr := Addr(cpu.indirect(instr.IndexReg(o.Dst), o.Off))
	mem.WriteData(addr, byte(cpu.reg[o.Src]))
}

func lds(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = int(mem.ReadData(Addr(cpu.ramp[RampD] | o.Off)))
}

func sts(cpu *Cpu, o *instr.Operands, mem Memory) {
	mem.WriteData(Addr(cpu.ramp[RampD]|o.Off), byte(cpu.reg[o.Src]))
}

func push(cpu *Cpu, o *instr.Operands, mem Memory) {
	mem.WriteData(Addr(cpu.sp), byte(cpu.reg[o.Src]))
	cpu.spInc(-1)
}

func pop(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.spInc(1)
	cpu.reg[o.Dst] = int(mem.ReadData(Addr(cpu.sp)))
}

func brbs(cpu *Cpu, o *instr.Operands, mem Memory) {
	if cpu.flags[o.Src] {
		cpu.pcInc(o.Off)
	}
}

func brbc(cpu *Cpu, o *instr.Operands, mem Memory) {
	if !cpu.flags[o.Src] {
		cpu.pcInc(o.Off)
	}
}

func eijmp(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.pc = cpu.reg[30] | (cpu.reg[31] << 8) | cpu.ramp[Eind]
}

func ijmp(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.pc = cpu.reg[30] | (cpu.reg[31] << 8)
}

func jmp(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.pc = o.Off & (cpu.rmask[Eind] | 0xffff)
}

func rjmp(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.pcInc(o.Off)
}

func call(cpu *Cpu, o *instr.Operands, mem Memory) {
	pushPC(cpu, mem)
	jmp(cpu, o, mem)
}

func eicall(cpu *Cpu, o *instr.Operands, mem Memory) {
	pushPC(cpu, mem)
	eijmp(cpu, o, mem)
}

func icall(cpu *Cpu, o *instr.Operands, mem Memory) {
	pushPC(cpu, mem)
	ijmp(cpu, o, mem)
}

func rcall(cpu *Cpu, o *instr.Operands, mem Memory) {
	pushPC(cpu, mem)
	rjmp(cpu, o, mem)
}

func pushPC(cpu *Cpu, mem Memory) {
	mem.WriteData(Addr(cpu.sp), byte(cpu.pc))
	cpu.spInc(-1)
	mem.WriteData(Addr(cpu.sp), byte(cpu.pc>>8))
	cpu.spInc(-1)
	if cpu.rmask[Eind] != 0 {
		mem.WriteData(Addr(cpu.sp), byte(cpu.pc>>16))
		cpu.spInc(-1)
	}
}

func ret(cpu *Cpu, o *instr.Operands, mem Memory) {
	popPC(cpu, mem)
}

func reti(cpu *Cpu, o *instr.Operands, mem Memory) {
	popPC(cpu, mem)
	cpu.flags[FlagI] = true
}

func popPC(cpu *Cpu, mem Memory) {
	cpu.pc = 0
	if cpu.rmask[Eind] != 0 {
		cpu.spInc(1)
		cpu.pc |= (int(mem.ReadData(Addr(cpu.sp))) << 16)
	}
	cpu.spInc(1)
	cpu.pc |= (int(mem.ReadData(Addr(cpu.sp))) << 8)
	cpu.spInc(1)
	cpu.pc |= int(mem.ReadData(Addr(cpu.sp)))
}

func lac(cpu *Cpu, o *instr.Operands, mem Memory) {
	mask := byte(cpu.reg[o.Dst])
	addr := Addr(cpu.indirect(instr.Z, 0))
	cpu.reg[o.Dst] = int(mem.ReadData(addr))
	mem.WriteData(addr, byte(cpu.reg[o.Dst]) & ^mask)
}

func las(cpu *Cpu, o *instr.Operands, mem Memory) {
	mask := byte(cpu.reg[o.Dst])
	addr := Addr(cpu.indirect(instr.Z, 0))
	cpu.reg[o.Dst] = int(mem.ReadData(addr))
	mem.WriteData(addr, byte(cpu.reg[o.Dst])|mask)
}

func lat(cpu *Cpu, o *instr.Operands, mem Memory) {
	mask := byte(cpu.reg[o.Dst])
	addr := Addr(cpu.indirect(instr.Z, 0))
	cpu.reg[o.Dst] = int(mem.ReadData(addr))
	mem.WriteData(addr, byte(cpu.reg[o.Dst])^mask)
}

func xch(cpu *Cpu, o *instr.Operands, mem Memory) {
	tmp := byte(cpu.reg[o.Dst])
	addr := Addr(cpu.indirect(instr.Z, 0))
	cpu.reg[o.Dst] = int(mem.ReadData(addr))
	mem.WriteData(addr, tmp)
}

func sbrc(cpu *Cpu, o *instr.Operands, mem Memory) {
	if (cpu.reg[o.Src] & (1 << uint(o.Off))) == 0 {
		cpu.skip = true
	}
}

func sbrs(cpu *Cpu, o *instr.Operands, mem Memory) {
	if (cpu.reg[o.Src] & (1 << uint(o.Off))) != 0 {
		cpu.skip = true
	}
}

func cpse(cpu *Cpu, o *instr.Operands, mem Memory) {
	if cpu.reg[o.Dst] == cpu.reg[o.Src] {
		cpu.skip = true
	}
}

func sbic(cpu *Cpu, o *instr.Operands, mem Memory) {
	if (mem.ReadData(Addr(o.Src+0x20)) & (1 << uint(o.Off))) == 0 {
		cpu.skip = true
	}
}

func sbis(cpu *Cpu, o *instr.Operands, mem Memory) {
	if (mem.ReadData(Addr(o.Src+0x20)) & (1 << uint(o.Off))) != 0 {
		cpu.skip = true
	}
}

func in(cpu *Cpu, o *instr.Operands, mem Memory) {
	cpu.reg[o.Dst] = int(mem.ReadData(Addr(o.Src + 0x20)))
}

func out(cpu *Cpu, o *instr.Operands, mem Memory) {
	mem.WriteData(Addr(o.Dst+0x20), byte(cpu.reg[o.Src]))
}

func cbi(cpu *Cpu, o *instr.Operands, mem Memory) {
	addr := Addr(o.Dst + 0x20)
	val := mem.ReadData(addr) & ^(1 << uint(o.Off))
	mem.WriteData(addr, val)
}

func sbi(cpu *Cpu, o *instr.Operands, mem Memory) {
	addr := Addr(o.Dst + 0x20)
	val := mem.ReadData(addr) | ^(1 << uint(o.Off))
	mem.WriteData(addr, val)
}

func lpm(cpu *Cpu, o *instr.Operands, mem Memory) {
	o.Src = int(instr.Z)
	o.Dst = 0
	loadProgMem(cpu, o, mem, 0)
}

func lpme(cpu *Cpu, o *instr.Operands, mem Memory) {
	loadProgMem(cpu, o, mem, 0)
}

func elpm(cpu *Cpu, o *instr.Operands, mem Memory) {
	o.Src = int(instr.Z)
	o.Dst = 0
	loadProgMem(cpu, o, mem, cpu.rmask[RampZ])
}

func elpme(cpu *Cpu, o *instr.Operands, mem Memory) {
	loadProgMem(cpu, o, mem, cpu.rmask[RampZ])
}

func loadProgMem(cpu *Cpu, o *instr.Operands, mem Memory, zmask int) {
	tmpMask := cpu.rmask[RampZ]
	tmpRamp := cpu.ramp[RampZ]
	cpu.rmask[RampZ] = zmask
	addr := Addr(cpu.indirect(instr.IndexReg(o.Src), 0))
	cpu.reg[o.Dst] = int(mem.LoadProgram(addr>>1, uint(addr)&0x1))
	// cpu.reg[o.Dst] = int(val>>((uint(addr)&0x1)*8)) & 0xff
	cpu.rmask[RampZ] = tmpMask
	cpu.ramp[RampZ] = tmpRamp
}

type opFunc func(*Cpu, *instr.Operands, Memory)

var opFuncs = [...]opFunc{
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
	cpi,    // Cbi
	com,    // Com
	com,    // ComReduced
	cp,     // Cp
	cp,     // CpReduced
	cpc,    // Cpc
	cpc,    // CpcReduced
	cpi,    // Cpi
	cpse,   // Cpse
	cpse,   // CpseReduced
	dec,    // Dec
	dec,    // DecReduced
	nop,    // Des ****
	eicall, // Eicall
	eijmp,  // Eijmp
	elpm,   // Elpm
	elpme,  // ElpmEnhanced
	eor,    // Eor
	eor,    // EorReduced
	fmul,   // Fmul
	fmuls,  // Fmuls
	fmulsu, // Fmulsu
	icall,  // Icall
	ijmp,   // Ijmp
	in,     // In
	in,     // InReduced
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
	lpm,    // Lpm
	lpme,   // LpmEnhanced
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
	out,    // Out
	out,    // OutReduced
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
	sbi,    // Sbi
	sbic,   // Sbic
	sbis,   // Sbis
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
