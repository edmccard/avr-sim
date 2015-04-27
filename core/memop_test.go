package core

import (
	"fmt"
	"github.com/edmccard/avr-sim/instr"
	"github.com/edmccard/testcase"
	"testing"
)

type idxData struct {
	pre, addr, post int
}

// Branches with X, Y, and Z as base index registers in Off.
func brBase(tree testcase.Tree, init, exp testcase.Testable) {
	initCpu := init.(tCpuDm)
	initCpu.ops.Off = int(instr.X)
	tree.Run("X", initCpu, exp)
	initCpu.ops.Off = int(instr.Y)
	tree.Run("Y", initCpu, exp)
	initCpu.ops.Off = int(instr.Z)
	tree.Run("Z", initCpu, exp)
}

func setupLoad(ireg instr.IndexReg, initCpu tCpuDm, d idxData) (tCpuDm,
	tCpuDm) {

	initCpu.ops.Src = int(ireg)
	reg := ireg.Reg()
	expCpu := initCpu
	expCpu.SetReg(reg, byte(d.post))
	expCpu.SetReg(reg+1, byte(d.post>>8))
	initCpu.SetReg(reg, byte(d.pre))
	initCpu.SetReg(reg+1, byte(d.pre>>8))
	initCpu.dmem.SetReadData([]int{0xff})
	expCpu.SetReg(0, 0xff)
	expCpu.dmem.SetReadData([]int{0xff})
	expCpu.dmem.ReadData(Addr(d.addr))
	return initCpu, expCpu
}

func setupStore(ireg instr.IndexReg, initCpu tCpuDm, d idxData) (tCpuDm,
	tCpuDm) {

	initCpu.ops.Dst = int(ireg)
	initCpu.SetReg(0, 0xff)
	reg := ireg.Reg()
	expCpu := initCpu
	expCpu.SetReg(reg, byte(d.post))
	expCpu.SetReg(reg+1, byte(d.post>>8))
	initCpu.SetReg(reg, byte(d.pre))
	initCpu.SetReg(reg+1, byte(d.pre>>8))
	expCpu.dmem.WriteData(Addr(d.addr), 0xff)
	return initCpu, expCpu
}

func TestLdSt(t *testing.T) {
	// TODO: RAMP
	var indirectCases = []struct {
		action instr.IndexAction
		data   []idxData
	}{
		{instr.NoAction, []idxData{
			{0x0800, 0x0800, 0x0800}}},
		{instr.PostInc, []idxData{
			{0x0800, 0x0800, 0x0801},
			{0xffff, 0xffff, 0x0000}}},
		{instr.PreDec, []idxData{
			{0x0800, 0x07ff, 0x07ff},
			{0x0000, 0xffff, 0xffff}}},
	}

	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		for _, indCase := range indirectCases {
			start := init.(tCpuDm)
			ireg := instr.IndexReg(start.ops.Off)
			ireg = ireg.WithAction(indCase.action)
			for n, c := range indCase.data {
				tag := fmt.Sprintf(" %s [%d]", ireg, n)
				initCpu, expCpu := setupLoad(ireg, start, c)
				ld(&initCpu.Cpu, &initCpu.ops, &initCpu.dmem)
				tree.Run("Ld"+tag, initCpu, expCpu)
				initCpu, expCpu = setupStore(ireg, start, c)
				st(&initCpu.Cpu, &initCpu.ops, &initCpu.dmem)
				tree.Run("St"+tag, initCpu, expCpu)
			}
		}
	}
	testcase.NewTree(t, "<->", brBase, run).Start(tCpuDm{})
}

func TestLddStd(t *testing.T) {
	var dispCases = []struct {
		disp int
		idxData
	}{
		{0x00, idxData{0x0800, 0x0800, 0x0800}},
		{0x01, idxData{0x0800, 0x0801, 0x0800}},
		{0x3f, idxData{0xffc1, 0x0000, 0xffc1}},
	}

	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		start := init.(tCpuDm)
		for n, c := range dispCases {
			ireg := instr.IndexReg(start.ops.Off)
			start.ops.Dst = 0
			start.ops.Off = c.disp
			tag := fmt.Sprintf(" [%d]", n)
			initCpu, expCpu := setupLoad(ireg, start, c.idxData)
			ldd(&initCpu.Cpu, &initCpu.ops, &initCpu.dmem)
			tree.Run("Ldd"+tag, initCpu, expCpu)
			start.ops.Src = 0
			start.ops.Off = c.disp
			initCpu, expCpu = setupStore(ireg, start, c.idxData)
			std(&initCpu.Cpu, &initCpu.ops, &initCpu.dmem)
			tree.Run("Std"+tag, initCpu, expCpu)
		}
	}
	testcase.NewTree(t, "<->", brBase, run).Start(tCpuDm{})
}

type stackCase struct {
	spPre, spPost, addr int
	regPre, regPost     byte
}

func TestPushPop(t *testing.T) {
	pushCases := []stackCase{
		{0x0000, 0xffff, 0x0000, 0xff, 0xff},
		{0x2000, 0x1fff, 0x2000, 0xff, 0xff},
	}
	popCases := []stackCase{
		{0xffff, 0x0000, 0x0000, 0x00, 0xff},
		{0x2000, 0x2001, 0x2001, 0x00, 0xff},
	}
	run := func(cases []stackCase, op opFunc, tag string) testcase.Branch {
		return func(tree testcase.Tree, init, exp testcase.Testable) {
			for n, c := range cases {
				initCpu := init.(tCpuDm)
				initCpu.sp = c.spPre
				initCpu.SetReg(0, c.regPre)
				expCpu := init.(tCpuDm)
				expCpu.sp = c.spPost
				expCpu.SetReg(0, c.regPost)
				if tag == "Push" {
					expCpu.dmem.WriteData(Addr(c.addr), c.regPost)
				} else {
					initCpu.dmem.SetReadData([]int{0xff})
					expCpu.dmem.SetReadData([]int{0xff})
					expCpu.dmem.ReadData(Addr(c.addr))
				}
				op(&initCpu.Cpu, &initCpu.ops, &initCpu.dmem)
				tree.Run(fmt.Sprintf("%s [%d]", tag, n), initCpu, expCpu)
			}
		}
	}
	testcase.NewTree(t, "STK", run(pushCases, push, "Push")).Start(tCpuDm{})
	testcase.NewTree(t, "STK", run(popCases, pop, "Pop")).Start(tCpuDm{})
}

func TestCallRcall(t *testing.T) {
	var cases = []struct {
		mnem                     instr.Mnemonic
		rmask, pcPre, a1, pcPost int
		spPre, spPost            int
	}{
		{instr.Call, 0x00, 0x0000, 0x12345, 0x2345, 0x1fff, 0x1ffd},
		{instr.Call, 0x3f, 0x0000, 0x12345, 0x12345, 0x1fff, 0x1ffc},
		{instr.Call, 0x00, 0x0000, 0x2345, 0x2345, 0x1fff, 0x1ffd},
		{instr.Call, 0x3f, 0x0000, 0x2345, 0x2345, 0x1fff, 0x1ffc},
		{instr.Rcall, 0x00, 0x0000, -1, 0xffff, 0x1fff, 0x1ffd},
		{instr.Rcall, 0x3f, 0x0000, -1, 0x3fffff, 0x1fff, 0x1ffc},
		{instr.Rcall, 0x00, 0xffff, 1, 0x0000, 0x1fff, 0x1ffd},
		{instr.Rcall, 0x3f, 0xffff, 1, 0x10000, 0x1fff, 0x1ffc},
		{instr.Rcall, 0x00, 0x1000, 0x07ff, 0x17ff, 0x1fff, 0x1ffd},
		{instr.Rcall, 0x3f, 0x1000, 0x07ff, 0x17ff, 0x1fff, 0x1ffc},
	}
	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		for n, c := range cases {
			initCpu := init.(tCpuDm)
			initCpu.ops.Off = c.a1
			initCpu.pc = c.pcPre
			initCpu.sp = c.spPre
			initCpu.rmask[Eind] = c.rmask << 16
			expCpu := initCpu
			expCpu.pc = c.pcPost
			expCpu.sp = c.spPost
			expCpu.dmem.WriteData(Addr(c.spPre), byte(c.pcPre))
			expCpu.dmem.WriteData(Addr(c.spPre-1), byte(c.pcPre>>8))
			if c.rmask != 0 {
				expCpu.dmem.WriteData(Addr(c.spPre-2), byte(c.pcPre>>16))
			}
			opFuncs[c.mnem](&initCpu.Cpu, &initCpu.ops, &initCpu.dmem)
			tree.Run(fmt.Sprintf("%s [%d]", c.mnem, n), initCpu, expCpu)
		}
	}
	testcase.NewTree(t, "JMP", run).Start(tCpuDm{})
}

func TestIcallEicall(t *testing.T) {
	var cases = []struct {
		mnem                 instr.Mnemonic
		eind, z, post        int
		pcPre, spPre, spPost int
	}{
		{instr.Icall, 0x00, 0x2000, 0x2000, 0x2345, 0x1fff, 0x1ffd},
		{instr.Icall, 0x3f, 0x2000, 0x2000, 0x2345, 0x1fff, 0x1ffc},
		{instr.Eicall, 0x00, 0x2000, 0x2000, 0x2345, 0x1fff, 0x1ffd},
		{instr.Eicall, 0x3f, 0x2000, 0x3f2000, 0x2345, 0x1fff, 0x1ffc},
	}
	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		for n, c := range cases {
			initCpu := init.(tCpuDm)
			initCpu.SetReg(30, byte(c.z))
			initCpu.SetReg(31, byte(c.z>>8))
			initCpu.setRmask(Eind, byte(c.eind))
			initCpu.SetRamp(Eind, byte(c.eind))
			initCpu.sp = c.spPre
			initCpu.pc = c.pcPre
			expCpu := initCpu
			expCpu.pc = c.post
			expCpu.sp = c.spPost
			expCpu.dmem.WriteData(Addr(c.spPre), byte(c.pcPre))
			expCpu.dmem.WriteData(Addr(c.spPre-1), byte(c.pcPre>>8))
			if c.eind != 0 {
				expCpu.dmem.WriteData(Addr(c.spPre-2), byte(c.pcPre>>16))
			}
			opFuncs[c.mnem](&initCpu.Cpu, &initCpu.ops, &initCpu.dmem)
			tree.Run(fmt.Sprintf("%s [%d]", c.mnem, n), initCpu, expCpu)
		}
	}
	testcase.NewTree(t, "JMP", run).Start(tCpuDm{})
}

func TestRetReti(t *testing.T) {
	var cases = []struct {
		mnem                     instr.Mnemonic
		emask, pc, spPre, spPost int
	}{
		{instr.Ret, 0x00, 0x2345, 0x1ffd, 0x1fff},
		{instr.Reti, 0x00, 0x2345, 0x1ffd, 0x1fff},
		{instr.Ret, 0x3f, 0x12345, 0x1ffc, 0x1fff},
		{instr.Reti, 0x3f, 0x12345, 0x1ffc, 0x1fff},
	}
	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		for n, c := range cases {
			initCpu := init.(tCpuDm)
			initCpu.sp = c.spPre
			initCpu.setRmask(Eind, byte(c.emask))
			expCpu := initCpu
			expCpu.pc = c.pc
			expCpu.sp = c.spPost
			if c.emask != 0 {
				initCpu.dmem.SetReadData([]int{
					c.pc >> 16, c.pc >> 8, c.pc,
				})
				expCpu.dmem.SetReadData([]int{
					c.pc >> 16, c.pc >> 8, c.pc,
				})
				expCpu.dmem.ReadData(Addr(c.spPost - 2))
			} else {
				initCpu.dmem.SetReadData([]int{
					c.pc >> 8, c.pc,
				})
				expCpu.dmem.SetReadData([]int{
					c.pc >> 8, c.pc,
				})
			}
			expCpu.dmem.ReadData(Addr(c.spPost - 1))
			expCpu.dmem.ReadData(Addr(c.spPost))
			if c.mnem == instr.Reti {
				expCpu.SetFlag(FlagI, true)
			}
			opFuncs[c.mnem](&initCpu.Cpu, &initCpu.ops, &initCpu.dmem)
			tree.Run(fmt.Sprintf("%s [%d]", c.mnem, n), initCpu, expCpu)
		}
	}
	testcase.NewTree(t, "JMP", run).Start(tCpuDm{})
}

func TestSBR(t *testing.T) {
	var cases = []struct {
		mnem                      string
		reg0                      byte
		pcPre, op, opNext, pcPost int
	}{
		{"Sbrs", 0x00, 0x1000, 0xfe00, 0x0000, 0x1001},
		{"Sbrs", 0x00, 0x1000, 0xfe00, 0x940e, 0x1001},
		{"Sbrs", 0x01, 0x1000, 0xfe00, 0x0000, 0x1002},
		{"Sbrs", 0x01, 0x1000, 0xfe00, 0x940e, 0x1003},
		{"Sbrc", 0x01, 0x1000, 0xfc00, 0x0000, 0x1001},
		{"Sbrc", 0x01, 0x1000, 0xfc00, 0x940e, 0x1001},
		{"Sbrc", 0x00, 0x1000, 0xfc00, 0x0000, 0x1002},
		{"Sbrc", 0x00, 0x1000, 0xfc00, 0x940e, 0x1003},
	}
	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		for n, c := range cases {
			initCpu := init.(tCpuDm)
			initCpu.SetReg(0, c.reg0)
			expCpu := initCpu
			initCpu.pc = c.pcPre
			expCpu.pc = c.pcPost
			initCpu.dmem.SetReadData([]int{c.op, c.opNext, 0})
			expCpu.dmem.SetReadData([]int{c.op, c.opNext, 0})
			expCpu.dmem.ReadProgram(Addr(c.pcPre))
			skip := (c.mnem == "Sbrs" && c.reg0 != 0) ||
				(c.mnem == "Sbrc" && c.reg0 == 0)
			if skip {
				expCpu.dmem.ReadProgram(Addr(c.pcPre + 1))
				if c.opNext == 0x940e {
					expCpu.dmem.ReadProgram(Addr(c.pcPre + 2))
				}
			}
			decoder := instr.NewDecoder(setXmega)
			initCpu.Step(&initCpu.dmem, &decoder)
			tree.Run(fmt.Sprintf("%s [%d]", c.mnem, n), initCpu, expCpu)
		}
	}
	testcase.NewTree(t, "SKP", run).Start(tCpuDm{})
}

func TestStepBranch(t *testing.T) {
	var cases = []struct {
		pcPre, op, pcPost int
	}{
		{0x1002, 0xf7e8, 0x1000},
		{0x1000, 0xf408, 0x1002},
	}
	decoder := instr.NewDecoder(setXmega)
	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		for n, c := range cases {
			initCpu := init.(tCpuDm)
			expCpu := initCpu
			initCpu.dmem.SetReadData([]int{c.op, 0})
			expCpu.dmem.SetReadData([]int{c.op, 0})
			expCpu.dmem.ReadProgram(Addr(c.pcPre))
			expCpu.pc = c.pcPost
			initCpu.Reset(0, c.pcPre)
			initCpu.Step(&initCpu.dmem, &decoder)
			tree.Run(fmt.Sprintf("Brcc [%d]", n), initCpu, expCpu)
		}
	}
	testcase.NewTree(t, "BRA", run).Start(tCpuDm{})
}

func TestIn(t *testing.T) {
	var cases = []struct {
		mnem                      instr.Mnemonic
		port, addr, val, reg, res int
	}{
		{instr.In, 0x00, 0x20, 0x44, 0x01, 0x44},
	}
	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		for n, c := range cases {
			initCpu := init.(tCpuIm)
			initCpu.imem.data[c.addr] = byte(c.val)
			initCpu.ops.Dst = c.reg
			initCpu.ops.Src = c.port
			expCpu := initCpu
			expCpu.SetReg(c.reg, byte(c.res))
			opFuncs[c.mnem](&initCpu.Cpu, &initCpu.ops, &initCpu.imem)
			tree.Run(fmt.Sprintf("%s [%d]", c.mnem, n), initCpu, expCpu)
		}
	}
	testcase.NewTree(t, "IOP", run).Start(tCpuIm{})
}

var setXmega = instr.NewSetXmega()

type ldiMem struct{}

func (m *ldiMem) ReadData(addr Addr) byte {
	return 0
}

func (m *ldiMem) WriteData(addr Addr, val byte) {}

func (m *ldiMem) ReadProgram(addr Addr) uint16 {
	return 0xe000
}

func (m *ldiMem) LoadProgram(addr Addr) byte {
	return 0
}

func BenchmarkLdi(b *testing.B) {
	cpu := Cpu{}
	d := instr.NewDecoder(setXmega)
	mem := ldiMem{}
	for n := 0; n < b.N; n++ {
		cpu.Step(&mem, &d)
	}
}

type tLPMmem struct{}

func (m *tLPMmem) ReadData(addr Addr) byte {
	return 0
}

func (m *tLPMmem) WriteData(addr Addr, val byte) {}

func (m *tLPMmem) ReadProgram(addr Addr) uint16 {
	return 0x95c8
}

func (m *tLPMmem) LoadProgram(addr Addr) byte {
	return 0
}

type tLPMCpu struct {
	tCpu
	mem tLPMmem
}

func (m tLPMCpu) Equals(other testcase.Testable) bool {
	o := other.(tLPMCpu)
	return m.tCpu.Equals(o.tCpu)
}

func (m tLPMCpu) Diff(other testcase.Testable) interface{} {
	o := other.(tLPMCpu)
	return m.tCpu.Diff(o.tCpu)
}
