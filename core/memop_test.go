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
