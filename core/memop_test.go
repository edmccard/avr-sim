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
