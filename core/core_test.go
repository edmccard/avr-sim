package core

import (
	"fmt"
	it "github.com/edmccard/avr-sim/instr"
	"reflect"
	"testing"
)

type key int

const (
	action key = iota
	addr
	bit
	cyc
	disp
	dstreg
	dstval
	ireg
	mulval
	mval
	opcds
	pc
	port
	ptr
	pval
	ramp
	savepc
	sp
	srcreg
	srcval
	status
	imm key = 32
)

type stack map[int]int

type flash map[int]int

type flags map[Flag]bool

type pair [2]int

type cdata map[key]interface{}

func (this cdata) merge(that cdata) cdata {
	merged := make(map[key]interface{})
	for k, v := range this {
		merged[k] = v
	}
	for k, v := range that {
		if prev, ok := merged[k]; ok && (k == status) {
			x := make(flags)
			for k2, v2 := range prev.(flags) {
				x[k2] = v2
			}
			for k2, v2 := range v.(flags) {
				x[k2] = v2
			}
			merged[k] = x
		} else {
			merged[k] = v
		}
	}
	return merged
}

func (data cdata) musthave(k key) interface{} {
	if val, ok := data[k]; ok {
		return val
	}
	panic(fmt.Sprintf("missing case data %s", k))
}

type branch []tcase

type branches []branch

type tcase struct {
	tag  string
	init cdata
	exp  cdata
	mnem it.Mnemonic
}

func (this tcase) merge(that tcase) tcase {
	return tcase{
		tag:  this.tag + " " + that.tag,
		init: this.init.merge(that.init),
		exp:  this.exp.merge(that.exp),
		mnem: this.mnem,
	}
}

func (tc tcase) run(t *testing.T) {
	init := newsystem()
	init.apply(tc.init)
	exp := newsystem()
	exp.apply(tc.init.merge(tc.exp))
	if _, ok := tc.init[opcds]; ok {
		insts := tc.init[opcds].(flash)
		var ipc int
		if val, ok := tc.init[pc]; ok {
			ipc = val.(int)
		}
		for k, v := range insts {
			init.mem.prog[Addr(k+ipc)] = uint16(v)
			exp.mem.prog[Addr(k+ipc)] = uint16(v)
		}
		init.cpu.Step(&init.mem, &decoder)
	} else {
		opFuncs[tc.mnem](&init.cpu, &init.cpu.ops, &init.mem)
	}
	if _, ok := tc.init[cyc]; !ok {
		exp.cpu.cycles = init.cpu.cycles
	}
	if !init.equals(&exp) {
		t.Error(tc.tag)
		fmt.Println("INIT:", init)
		fmt.Println("EXP: ", exp)
	}
}

type tmem struct {
	data map[Addr]byte
	prog map[Addr]uint16
}

func newtmem() tmem {
	return tmem{data: make(map[Addr]byte), prog: make(map[Addr]uint16)}
}

func (this tmem) equals(that tmem) bool {
	return reflect.DeepEqual(this.data, that.data) &&
		reflect.DeepEqual(this.prog, that.prog)
}

func (m *tmem) ReadData(addr Addr) byte {
	if val, ok := m.data[addr]; ok {
		return val
	}
	return 0x9e
}

func (m *tmem) WriteData(addr Addr, val byte) {
	m.data[addr] = val
}

func (m *tmem) ReadProgram(addr Addr) uint16 {
	if val, ok := m.prog[addr]; ok {
		return val
	}
	return 0
}

func (m *tmem) LoadProgram(addr Addr) byte {
	if val, ok := m.prog[addr>>1]; ok {
		return byte(val >> ((uint(addr) & 0x1) * 8))
	}
	return 0
}

type system struct {
	cpu Cpu
	mem tmem
}

func newsystem() system {
	return system{cpu: Cpu{}, mem: newtmem()}
}

func (s *system) apply(data cdata) {
	if statval, ok := data[status]; ok {
		for k, v := range statval.(flags) {
			s.cpu.flags[k] = v
		}
	}
	if spval, ok := data[sp]; ok {
		s.cpu.sp = spval.(int)
	}
	if dsp, ok := data[disp]; ok {
		s.cpu.ops.Off = dsp.(int)
	}

	if val, ok := data[pc]; ok {
		// jumps, calls, returns, branches
		s.applyoffset(val.(int), data)
	} else if val, ok := data[bit]; ok {
		// bset/bclr, bld/bst, sbi/cbi
		s.applybitop(val.(int), data)
	} else if val, ok := data[ireg]; ok {
		// indirect loads/stores/atomics
		s.applyindirect(val.(it.IndexReg), data)
	} else if val, ok := data[addr]; ok {
		// direct loads/stores, push/pop, in/out
		s.applydirect(Addr(val.(int)), data)
	} else {
		s.applyarith(data)
	}
}

func (s *system) applyoffset(offset int, data cdata) {
	if eindval, ok := data[ramp]; ok {
		// eixxx
		s.setramp(Eind, eindval.(int))
	}
	if val, ok := data[ptr]; ok {
		// ixxx/eixxx
		s.setindex(it.Z.Reg(), val.(int))
	}
	if stk, ok := data[savepc]; ok {
		// calls/returns
		for a, v := range stk.(stack) {
			s.mem.WriteData(Addr(a), byte(v))
		}
	}
	if bitnum, ok := data[bit]; ok {
		// branches
		s.cpu.ops.Src = bitnum.(int)
	}
	s.cpu.pc = offset
}

func (s *system) applyindirect(base it.IndexReg, data cdata) {
	action := data.musthave(action).(it.IndexAction)
	indexreg := base.WithAction(action)
	if reg, ok := data[srcreg]; ok {
		s.cpu.reg[reg.(int)] = data.musthave(srcval).(int)
		s.cpu.ops.Src = reg.(int)
		s.cpu.ops.Dst = int(indexreg)
	} else {
		reg := data.musthave(dstreg).(int)
		s.cpu.ops.Src = int(indexreg)
		s.cpu.ops.Dst = reg
		if dval, ok := data[dstval]; ok {
			s.cpu.reg[reg] = dval.(int)
		}
	}
	iptr := data.musthave(ptr).(int)
	s.setindex(indexreg.Reg(), iptr)
	if rmp, ok := data[ramp]; ok {
		s.setramp(Ramp(base), rmp.(int))
	}
	if memval, ok := data[mval]; ok {
		if maddr, ok := data[addr]; ok {
			s.mem.WriteData(Addr(maddr.(int)), byte(memval.(int)))
		} else {
			s.mem.WriteData(Addr(iptr), byte(memval.(int)))
		}
	} else if progval, ok := data[pval]; ok {
		maddr := data.musthave(addr).(int)
		s.mem.prog[Addr(maddr)] = uint16(progval.(int))
	}
}

func (s *system) applydirect(maddr Addr, data cdata) {
	ioport, hasioport := data[port]
	if !hasioport {
		s.cpu.ops.Off = int(maddr)
	}
	if memval, ok := data[mval]; ok {
		s.mem.WriteData(maddr, byte(memval.(int)))
	}
	if reg, ok := data[srcreg]; ok {
		s.cpu.reg[reg.(int)] = data.musthave(srcval).(int)
		s.cpu.ops.Src = reg.(int)
		s.cpu.ops.Dst = reg.(int)
		if hasioport {
			s.cpu.ops.Dst = ioport.(int)
		}
	} else {
		reg := data.musthave(dstreg).(int)
		s.cpu.ops.Src = reg
		s.cpu.ops.Dst = reg
		if hasioport {
			s.cpu.ops.Src = ioport.(int)
		}
		if dval, ok := data[dstval]; ok {
			s.cpu.reg[reg] = dval.(int)
		}
	}
}

func (s *system) applybitop(b int, data cdata) {
	if ioport, ok := data[port]; ok {
		// cbi/sbi
		memval := data.musthave(mval).(int)
		maddr := data.musthave(addr).(int)
		s.mem.WriteData(Addr(maddr), byte(memval))
		s.cpu.ops.Dst = ioport.(int)
		s.cpu.ops.Src = ioport.(int)
		s.cpu.ops.Off = b
	} else if reg, ok := data[srcreg]; ok {
		sval := data.musthave(srcval).(int)
		s.cpu.reg[reg.(int)] = sval
		s.cpu.ops.Src = reg.(int)
		s.cpu.ops.Dst = reg.(int)
		s.cpu.ops.Off = b
	} else {
		// bset/bclr
		s.cpu.ops.Src = b
		s.cpu.ops.Dst = b
	}
}

func (s *system) applyarith(data cdata) {
	if reg, ok := data[srcreg]; ok {
		sval := data.musthave(srcval)
		switch sval := sval.(type) {
		case int:
			switch reg := reg.(type) {
			case key:
				s.cpu.ops.Src = sval
			case int:
				s.cpu.ops.Src = reg
				s.cpu.reg[reg] = sval
			}
		case pair:
			s.cpu.ops.Src = reg.(pair)[1]
			s.cpu.reg[reg.(pair)[0]] = sval[0]
			s.cpu.reg[reg.(pair)[1]] = sval[1]
		default:
			panic("bad srcval")
		}
	}
	reg := data.musthave(dstreg)
	val := data.musthave(dstval)
	switch val := val.(type) {
	case int:
		s.cpu.ops.Dst = reg.(int)
		s.cpu.reg[reg.(int)] = val
	case pair:
		s.cpu.ops.Dst = reg.(pair)[1]
		s.cpu.reg[reg.(pair)[0]] = val[0]
		s.cpu.reg[reg.(pair)[1]] = val[1]
	default:
		panic("bad dstval")
	}
	if val, ok := data[mulval]; ok {
		s.cpu.reg[0] = val.(int) & 0xff
		s.cpu.reg[1] = val.(int) >> 8
	}
}

func (s *system) setramp(base Ramp, val int) {
	s.cpu.setRmask(base, 0x3f)
	s.cpu.SetRamp(base, byte(val))
}

func (s *system) setindex(reg, val int) {
	s.cpu.reg[reg] = val & 0xff
	s.cpu.reg[reg+1] = val >> 8
}

func (this *system) equals(that *system) bool {
	return this.cpu == that.cpu && this.mem.equals(that.mem)
}

type casetree struct {
	t        *testing.T
	branches branches
}

func (tree casetree) run(builder tcase) {
	if len(tree.branches) == 0 {
		builder.run(tree.t)
	} else {
		next := casetree{t: tree.t, branches: tree.branches[1:]}
		for _, tc := range tree.branches[0] {
			next.run(builder.merge(tc))
		}
	}
}

var decoder = it.NewDecoder(it.NewSetXmega())
