package cpu

import (
	"fmt"
	"github.com/edmccard/avr-sim/instr"
	"github.com/edmccard/testcase"
	"strings"
)

// tCpu wraps Cpu to implement Testable, and to simplify setting up
// initial state for tests.
type tCpu struct {
	Cpu
	am instr.AddrMode
}

func (tc tCpu) Equals(other testcase.Testable) bool {
	o := other.(tCpu)
	return o == tc
}

func (tc tCpu) Diff(other testcase.Testable) interface{} {
	this, that := tc, other.(tCpu)

	rThis, rThat := this.rDiff(that)
	sThis, sThat := this.sDiff(that)
	if rThis == "" && sThis == "" {
		return nil
	}

	thisLine, thatLine := "", ""
	if rThis != "" {
		thisLine += rThis
		thatLine += rThat
		if sThis != "" {
			thisLine += " | "
			thatLine += " | "
		}
	}
	if sThis != "" {
		thisLine += sThis
		thatLine += sThat
	}

	return "\n" + thisLine + "\n" + thatLine + "\n"
}

// Diff for registers 0-31
func (this tCpu) rDiff(that tCpu) (rThis, rThat string) {
	var thisR, thatR []string
	for i := 0; i < 32; i++ {
		if this.R[i] != that.R[i] {
			thisR = append(thisR, fmt.Sprintf("%d=%02x", i, this.R[i]))
			thatR = append(thatR, fmt.Sprintf("%d=%02x", i, that.R[i]))
		}
	}
	if thisR != nil {
		rThis = "R[" + strings.Join(thisR, ",") + "]"
		rThat = "R[" + strings.Join(thatR, ",") + "]"
	}
	return
}

// Diff for status, stack pointer, and program counter
func (this tCpu) sDiff(that tCpu) (sThis, sThat string) {
	var thisS, thatS []string
	thisSreg, thatSreg := this.ByteFromSreg(), that.ByteFromSreg()
	if thisSreg != thatSreg {
		thisS = append(thisS, fmt.Sprintf("S=%02x", thisSreg))
		thatS = append(thatS, fmt.Sprintf("S=%02x", thatSreg))
	}
	if this.SP != that.SP {
		thisS = append(thisS, fmt.Sprintf("SP=%04x", this.SP))
		thatS = append(thatS, fmt.Sprintf("SP=%04x", that.SP))
	}
	if this.PC != that.PC {
		thisS = append(thisS, fmt.Sprintf("PC=%04x", this.PC))
		thatS = append(thatS, fmt.Sprintf("PC=%04x", that.PC))
	}
	sThis = strings.Join(thisS, " ")
	sThat = strings.Join(thatS, " ")
	return
}

// Helper for setting expected status
func (tc *tCpu) setStatus(expStatus, mask byte) {
	setFlags := mask & expStatus
	clrFlags := ^mask | expStatus
	status := tc.ByteFromSreg()
	status &= clrFlags
	status |= setFlags
	tc.SregFromByte(status)
}

type arithData struct {
	status, v1, v2, res int
}

type arithCase struct {
	t    testcase.Tree
	init tCpu
	mask int
	arithData
	n int
}

func (ac arithCase) testD5R5(op OpFunc, tag string) {
	initCpu := ac.init
	d := initCpu.am.A1
	r := initCpu.am.A2
	if d == r && ac.v1 != ac.v2 {
		return
	}
	initCpu.R[r] = ac.v2
	initCpu.R[d] = ac.v1
	expCpu := ac.init
	expCpu.R[r] = ac.v2
	expCpu.R[d] = ac.res
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	op(&(initCpu.Cpu), &(initCpu.am), nil)
	ac.t.Run(fmt.Sprintf("%s [%d]", tag, ac.n), initCpu, expCpu)
}

func (ac arithCase) testMul(op OpFunc, tag string) {
	initCpu := ac.init
	d := initCpu.am.A1
	r := initCpu.am.A2
	if d == r && ac.v1 != ac.v2 {
		return
	}
	initCpu.R[r] = ac.v2
	initCpu.R[d] = ac.v1
	expCpu := ac.init
	expCpu.R[r] = ac.v2
	expCpu.R[d] = ac.v1
	expCpu.R[0] = ac.res & 0xff
	expCpu.R[1] = ac.res >> 8
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	op(&(initCpu.Cpu), &(initCpu.am), nil)
	ac.t.Run(fmt.Sprintf("%s [%d]", tag, ac.n), initCpu, expCpu)
}

func (ac arithCase) testD4K8(op OpFunc, tag string) {
	expCpu := ac.init
	expCpu.R[16] = ac.res
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	initCpu := ac.init
	initCpu.R[16] = ac.v1
	initCpu.am = instr.AddrMode{16, instr.Addr(ac.v2), instr.NoIndex}
	op(&(initCpu.Cpu), &(initCpu.am), nil)
	ac.t.Run(fmt.Sprintf("%s [%d]", tag, ac.n), initCpu, expCpu)
}

func (ac arithCase) testDDK6(op OpFunc, tag string) {
	expCpu := ac.init
	expCpu.R[24] = ac.res & 0xff
	expCpu.R[25] = ac.res >> 8
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	initCpu := ac.init
	initCpu.R[24] = ac.v1 & 0xff
	initCpu.R[25] = ac.v1 >> 8
	initCpu.am = instr.AddrMode{24, instr.Addr(ac.v2), instr.NoIndex}
	op(&(initCpu.Cpu), &(initCpu.am), nil)
	ac.t.Run(fmt.Sprintf("%s [%d]", tag, ac.n), initCpu, expCpu)
}

func (ac arithCase) testMovw() {
	initCpu := ac.init
	d := initCpu.am.A1
	r := initCpu.am.A2
	if d == r && ac.v1 != ac.v2 {
		return
	}
	initCpu.R[r] = ac.v2 & 0xff
	initCpu.R[r+1] = ac.v2 >> 8
	initCpu.R[d] = ac.v1 & 0xff
	initCpu.R[d+1] = ac.v1 >> 8
	expCpu := ac.init
	expCpu.R[r] = ac.v2 & 0xff
	expCpu.R[r+1] = ac.v2 >> 8
	expCpu.R[d] = ac.res & 0xff
	expCpu.R[d+1] = ac.res >> 8
	Movw(&(initCpu.Cpu), &(initCpu.am), nil)
	ac.t.Run(fmt.Sprintf("Movw [%d]", ac.n), initCpu, expCpu)
}

type branchCase struct {
	status, offset, pre, post int
}

func (bc branchCase) testBranch(t testcase.Tree, init tCpu, op OpFunc,
	tag string) {

	bit := init.am.A1
	init.setStatus(byte(bc.status), byte(1<<uint(bit)))
	status := init.ByteFromSreg()
	exp := init
	exp.PC = bc.post
	init.am.A2 = instr.Addr(bc.offset)
	init.PC = bc.pre
	op(&(init.Cpu), &(init.am), nil)
	t.Run(fmt.Sprintf("%s(%d) %02x", tag, bit, status), init, exp)
}
