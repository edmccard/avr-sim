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

var dispCases = []struct {
	disp instr.Addr
	idxData
}{
	{0x00, idxData{0x0800, 0x0800, 0x0800}},
	{0x01, idxData{0x0800, 0x0801, 0x0800}},
	{0x3f, idxData{0xffc1, 0x0000, 0xffc1}},
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
	expCpu.SetReg(0, 0xff)
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
	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		for _, indCase := range indirectCases {
			start := init.(tCpuDm)
			ireg := start.am.Ireg.WithMode(indCase.mode)
			start.am.Ireg = ireg
			for n, c := range indCase.data {
				tag := fmt.Sprintf(" %s [%d]", ireg, n)
				initCpu, expCpu := setupLoad(start, c)
				Ld(&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
				tree.Run("Ld"+tag, initCpu, expCpu)
				initCpu, expCpu = setupStore(start, c)
				St(&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
				tree.Run("St"+tag, initCpu, expCpu)
			}
		}
	}
	testcase.NewTree(t, "<->", brBase, run).Start(tCpuDm{})
}

func TestLddStd(t *testing.T) {
	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		start := init.(tCpuDm)
		for n, c := range dispCases {
			start.am.A2 = c.disp
			tag := fmt.Sprintf(" [%d]", n)
			initCpu, expCpu := setupLoad(start, c.idxData)
			Ldd(&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
			tree.Run("Ldd"+tag, initCpu, expCpu)
			initCpu, expCpu = setupStore(start, c.idxData)
			Std(&initCpu.Cpu, &initCpu.am, &initCpu.dmem)
			tree.Run("Std"+tag, initCpu, expCpu)
		}
	}
	testcase.NewTree(t, "<->", brBase, run).Start(tCpuDm{})
}
