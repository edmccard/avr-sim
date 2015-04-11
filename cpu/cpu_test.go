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

type caseData struct {
	status, v1, v2, res int
}

type arithCase struct {
	t    testcase.Tree
	init tCpu
	mask int
	caseData
	n int
}

func (ac arithCase) testD5R5(op OpFunc, tag string) {
	expCpu := ac.init
	expCpu.R[0] = ac.res
	expCpu.R[1] = ac.v2
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	initCpu := ac.init
	initCpu.R[0] = ac.v1
	initCpu.R[1] = ac.v2
	initCpu.am = instr.AddrMode{0, 1, instr.NoIndex}
	op(&(initCpu.Cpu), &(initCpu.am))
	ac.t.Run(fmt.Sprintf("%s [%d]", tag, ac.n), initCpu, expCpu)
}

func (ac arithCase) testD4K8(op OpFunc, tag string) {
	expCpu := ac.init
	expCpu.R[16] = ac.res
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	initCpu := ac.init
	initCpu.R[16] = ac.v1
	initCpu.am = instr.AddrMode{16, instr.Addr(ac.v2), instr.NoIndex}
	op(&(initCpu.Cpu), &(initCpu.am))
	ac.t.Run(fmt.Sprintf("%s [%d]", tag, ac.n), initCpu, expCpu)
}

func (ac arithCase) testDDK6(op OpFunc, tag string) {
	expCpu := ac.init
	expCpu.R[24] = ac.res & 0xff
	expCpu.R[25] = ac.res >> 8
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	initCpu:= ac.init
	initCpu.R[24] = ac.v1 & 0xff
	initCpu.R[25] = ac.v1 >> 8
	initCpu.am = instr.AddrMode{24, instr.Addr(ac.v2), instr.NoIndex}
	op(&(initCpu.Cpu), &(initCpu.am))
	ac.t.Run(fmt.Sprintf("%s [%d]", tag, ac.n), initCpu, expCpu)
}
