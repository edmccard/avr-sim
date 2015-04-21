package cpu

import (
	"fmt"
	"github.com/edmccard/avr-sim/instr"
	"github.com/edmccard/testcase"
	"testing"
)

type idxData struct {
	pre, addr, post int
}

// Branches with X, Y, and Z as base index registers.
func brBase(tree testcase.Tree, init, exp testcase.Testable) {
	initCpu := init.(tCpuDm)
	initCpu.am.Ireg = instr.X
	tree.Run("X", initCpu, exp)
	initCpu.am.Ireg = instr.Y
	tree.Run("Y", initCpu, exp)
	initCpu.am.Ireg = instr.Z
	tree.Run("Z", initCpu, exp)
}

func setupLoad(initCpu tCpuDm, d idxData) (tCpuDm, tCpuDm) {
	initCpu, expCpu := setupLoadStore(initCpu, d)
	initCpu.dmem.SetReadData([]int{0xff})
	expCpu.SetReg(0, 0xff)
	expCpu.dmem.SetReadData([]int{0xff})
	expCpu.dmem.ReadData(instr.Addr(d.addr))
	return initCpu, expCpu
}

func setupStore(initCpu tCpuDm, d idxData) (tCpuDm, tCpuDm) {
	initCpu.SetReg(0, 0xff)
	initCpu, expCpu := setupLoadStore(initCpu, d)
	expCpu.dmem.WriteData(instr.Addr(d.addr), 0xff)
	return initCpu, expCpu
}

func setupLoadStore(initCpu tCpuDm, d idxData) (tCpuDm, tCpuDm) {
	reg := initCpu.am.Ireg.Reg()
	expCpu := initCpu
	expCpu.SetReg(reg, byte(d.post))
	expCpu.SetReg(reg+1, byte(d.post>>8))
	initCpu.SetReg(reg, byte(d.pre))
	initCpu.SetReg(reg+1, byte(d.pre>>8))
	return initCpu, expCpu
}

func TestLdSt(t *testing.T) {
	// TODO: RAMP
	var indirectCases = []struct {
		mode instr.IndexMode
		data []idxData
	}{
		{instr.NoMode, []idxData{
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
			ireg := start.am.Ireg.WithMode(indCase.mode)
			start.am.Ireg = ireg
			for n, c := range indCase.data {
				tag := fmt.Sprintf(" %s [%d]", ireg, n)
				initCpu, expCpu := setupLoad(start, c)
				ld(&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
				tree.Run("Ld"+tag, initCpu, expCpu)
				initCpu, expCpu = setupStore(start, c)
				st(&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
				tree.Run("St"+tag, initCpu, expCpu)
			}
		}
	}
	testcase.NewTree(t, "<->", brBase, run).Start(tCpuDm{})
}

func TestLddStd(t *testing.T) {
	var dispCases = []struct {
		disp instr.Addr
		idxData
	}{
		{0x00, idxData{0x0800, 0x0800, 0x0800}},
		{0x01, idxData{0x0800, 0x0801, 0x0800}},
		{0x3f, idxData{0xffc1, 0x0000, 0xffc1}},
	}

	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		start := init.(tCpuDm)
		for n, c := range dispCases {
			start.am.A2 = c.disp
			tag := fmt.Sprintf(" [%d]", n)
			initCpu, expCpu := setupLoad(start, c.idxData)
			ldd(&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
			tree.Run("Ldd"+tag, initCpu, expCpu)
			initCpu, expCpu = setupStore(start, c.idxData)
			std(&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
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
	run := func(cases []stackCase, op OpFunc, tag string) testcase.Branch {
		return func(tree testcase.Tree, init, exp testcase.Testable) {
			for n, c := range cases {
				initCpu := init.(tCpuDm)
				initCpu.sp = c.spPre
				initCpu.SetReg(0, c.regPre)
				expCpu := init.(tCpuDm)
				expCpu.sp = c.spPost
				expCpu.SetReg(0, c.regPost)
				if tag == "Push" {
					expCpu.dmem.WriteData(instr.Addr(c.addr), c.regPost)
				} else {
					initCpu.dmem.SetReadData([]int{0xff})
					expCpu.dmem.SetReadData([]int{0xff})
					expCpu.dmem.ReadData(instr.Addr(c.addr))
				}
				op(&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
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
			initCpu.am.A1 = instr.Addr(c.a1)
			initCpu.pc = c.pcPre
			initCpu.sp = c.spPre
			initCpu.rmask[Eind] = c.rmask << 16
			expCpu := initCpu
			expCpu.pc = c.pcPost
			expCpu.sp = c.spPost
			expCpu.dmem.WriteData(instr.Addr(c.spPre), byte(c.pcPre))
			expCpu.dmem.WriteData(instr.Addr(c.spPre-1), byte(c.pcPre>>8))
			if c.rmask != 0 {
				expCpu.dmem.WriteData(instr.Addr(c.spPre-2), byte(c.pcPre>>16))
			}
			opFuncs[c.mnem](&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
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
			expCpu.dmem.WriteData(instr.Addr(c.spPre), byte(c.pcPre))
			expCpu.dmem.WriteData(instr.Addr(c.spPre-1), byte(c.pcPre>>8))
			if c.eind != 0 {
				expCpu.dmem.WriteData(instr.Addr(c.spPre-2), byte(c.pcPre>>16))
			}
			opFuncs[c.mnem](&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
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
				expCpu.dmem.ReadData(instr.Addr(c.spPost - 2))
			} else {
				initCpu.dmem.SetReadData([]int{
					c.pc >> 8, c.pc,
				})
				expCpu.dmem.SetReadData([]int{
					c.pc >> 8, c.pc,
				})
			}
			expCpu.dmem.ReadData(instr.Addr(c.spPost - 1))
			expCpu.dmem.ReadData(instr.Addr(c.spPost))
			if c.mnem == instr.Reti {
				expCpu.SetFlag(FlagI, true)
			}
			opFuncs[c.mnem](&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
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
			expCpu.dmem.ReadProgram(instr.Addr(c.pcPre))
			skip := (c.mnem == "Sbrs" && c.reg0 != 0) ||
				(c.mnem == "Sbrc" && c.reg0 == 0)
			if skip {
				expCpu.dmem.ReadProgram(instr.Addr(c.pcPre + 1))
				if c.opNext == 0x940e {
					expCpu.dmem.ReadProgram(instr.Addr(c.pcPre + 2))
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
			expCpu.dmem.ReadProgram(instr.Addr(c.pcPre))
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
			initCpu.am.A1 = instr.Addr(c.port)
			initCpu.am.A2 = instr.Addr(c.reg)
			expCpu := initCpu
			expCpu.SetReg(instr.Addr(c.reg), byte(c.res))
			opFuncs[c.mnem](&initCpu.Cpu, &initCpu.am, &initCpu.imem)
			tree.Run(fmt.Sprintf("%s [%d]", c.mnem, n), initCpu, expCpu)
		}
	}
	testcase.NewTree(t, "IOP", run).Start(tCpuIm{})
}

var setXmega = instr.NewSetXmega()

type ldiMem struct{}

func (m *ldiMem) ReadData(addr instr.Addr) byte {
	return 0
}

func (m *ldiMem) WriteData(addr instr.Addr, val byte) {}

func (m *ldiMem) ReadProgram(addr instr.Addr) uint16 {
	return 0xe000
}

func BenchmarkStepLdi(b *testing.B) {
	var xxxam = instr.AddrMode{}
	var xxxc = Cpu{am: xxxam}
	var xxxd = instr.NewDecoder(setXmega)
	var xxxm = ldiMem{}

	for n := 0; n < b.N; n++ {
		xxxc.Step(&xxxm, &xxxd)
	}
}
