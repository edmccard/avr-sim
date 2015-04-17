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
	return tc.Cpu == o.Cpu
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
		if this.reg[i] != that.reg[i] {
			thisR = append(thisR, fmt.Sprintf("%d=%02x", i, this.reg[i]))
			thatR = append(thatR, fmt.Sprintf("%d=%02x", i, that.reg[i]))
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
	if this.sp != that.sp {
		thisS = append(thisS, fmt.Sprintf("sp=%04x", this.sp))
		thatS = append(thatS, fmt.Sprintf("sp=%04x", that.sp))
	}
	if this.pc != that.pc {
		thisS = append(thisS, fmt.Sprintf("pc=%04x", this.pc))
		thatS = append(thatS, fmt.Sprintf("pc=%04x", that.pc))
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

type tDataMem struct {
	readCount  int
	readAddrs  [1]instr.Addr
	writeCount int
	writeAddrs [1]instr.Addr
	writeVals  [1]byte
}

func (tdm *tDataMem) ReadData(addr instr.Addr) byte {
	if tdm.readCount < 1 {
		tdm.readAddrs[0] = addr
	}
	tdm.readCount += 1
	return 0xff
}

func (tdm *tDataMem) WriteData(addr instr.Addr, val byte) {
	if tdm.writeCount < 1 {
		tdm.writeAddrs[0] = addr
		tdm.writeVals[0] = val
	}
	tdm.writeCount += 1
}

func (this *tDataMem) diff(that *tDataMem) string {
	// Assumes this != that; assumes 1 read xor 1 write;
	// assumes this is expected and that is actual
	switch {
	case this.readCount < that.readCount:
		return "MEM: too many reads"
	case this.readCount > that.readCount:
		return "MEM: too few reads"
	case this.writeCount < that.writeCount:
		return "MEM: too many writes"
	case this.writeCount > that.writeCount:
		return "MEM: too few writes"
	}
	if this.readCount > 0 {
		if this.readAddrs[0] != that.readAddrs[0] {
			return fmt.Sprintf("MEM: expected read %04x got %04x",
				this.readAddrs[0], that.readAddrs[0])
		}
	}
	if this.writeCount > 0 {
		if this.writeAddrs[0] != that.writeAddrs[0] {
			return fmt.Sprintf("MEM: expected write at %04x got %04x",
				this.writeAddrs[0], that.writeAddrs[0])
		}
		if this.writeVals[0] != that.writeVals[0] {
			return fmt.Sprintf("MEM: expected write of %02x got %02x",
				this.writeVals[0], that.writeVals[0])
		}
	}
	return ""
}

type tCpuDm struct {
	tCpu
	dmem tDataMem
}

func (tc tCpuDm) Equals(other testcase.Testable) bool {
	o := other.(tCpuDm)
	return tc.tCpu == o.tCpu && tc.dmem == o.dmem
}

func (tc tCpuDm) Diff(other testcase.Testable) interface{} {
	o := other.(tCpuDm)
	cDiff := tc.tCpu.Diff(o.tCpu)
	mDiff := fmt.Sprintf("%s", tc.dmem.diff(&o.dmem))
	if mDiff != "" {
		if cDiff != nil {
			return fmt.Sprintf("%s", cDiff) + "\n" + mDiff
		} else {
			return "\n" + mDiff
		}
	} else {
		return cDiff
	}
}
