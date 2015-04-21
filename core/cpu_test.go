package cpu

import (
	"fmt"
	"github.com/edmccard/avr-sim/instr"
	"github.com/edmccard/testcase"
	"strings"
)

// tCpu wraps Cpu to implement Testable, and to simplify setting up
// initial state for tests.
// TODO: a "ControllableCpu" interface (but with a better name)
//       so the branch functions work with either tCpu or tCpuDm
type tCpu struct {
	Cpu
}

func (tc tCpu) Equals(other testcase.Testable) bool {
	o := other.(tCpu)
	return tc.reg == o.reg && tc.flags == o.flags && tc.sp == o.sp &&
		tc.pc == o.pc && tc.ramp == o.ramp && tc.rmask == o.rmask
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

// tDataMem implements the Memory interface for testing. Setup of
// expected state is done by calling the Read/Write methods in the
// same order (with the same addresses and written values) as they
// should be called by an instruction; for read tests, use SetReadData
// on both the initial and expected states.
type tDataMem struct {
	readVals      []int
	readAttempts  string
	writeAttempts string
}

func (tdm *tDataMem) ReadData(addr instr.Addr) byte {
	var val byte = 0x00
	if len(tdm.readVals) != 0 {
		val = byte(tdm.readVals[0])
		tdm.readVals = tdm.readVals[1:]
	}
	tdm.readAttempts += fmt.Sprintf("%02x<-%04x ", val, addr)
	return val
}

// The data stored by SetReadData is used for both data and program
// reads.
func (tdm *tDataMem) SetReadData(vals []int) {
	for _, v := range vals {
		tdm.readVals = append(tdm.readVals, v)
	}
}

func (tdm *tDataMem) WriteData(addr instr.Addr, val byte) {
	tdm.writeAttempts += fmt.Sprintf("%04x->%02x ", addr, val)
}

func (tdm *tDataMem) ReadProgram(addr instr.Addr) uint16 {
	var val uint16 = 0x0000
	if len(tdm.readVals) != 0 {
		val = uint16(tdm.readVals[0])
		tdm.readVals = tdm.readVals[1:]
	}
	tdm.readAttempts += fmt.Sprintf("%04x<-%04x ", val, addr)
	return val
}

func (this *tDataMem) equals(that *tDataMem) bool {
	return this.readAttempts == that.readAttempts &&
		this.writeAttempts == that.writeAttempts
}

func (this *tDataMem) diff(that *tDataMem) string {
	thisLine, thatLine := "", ""
	if this.readAttempts != that.readAttempts {
		thisLine += this.readAttempts
		thatLine += that.readAttempts
	}
	if this.writeAttempts != that.writeAttempts {
		if thisLine != "" {
			thisLine += " "
			thatLine += " "
		}
		thisLine += this.writeAttempts
		thatLine += that.writeAttempts
	}
	if thisLine == "" && thatLine == "" {
		return ""
	}
	return "MEM: " + thisLine + "\n" + "MEM: " + thatLine + "\n"
}

type tCpuDm struct {
	tCpu
	dmem tDataMem
}

func (tc tCpuDm) Equals(other testcase.Testable) bool {
	o := other.(tCpuDm)
	return tc.tCpu.Equals(o.tCpu) && tc.dmem.equals(&o.dmem)
}

func (tc tCpuDm) Diff(other testcase.Testable) interface{} {
	o := other.(tCpuDm)
	cDiff := tc.tCpu.Diff(o.tCpu)
	mDiff := tc.dmem.diff(&o.dmem)
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

type tIoMem struct {
	data [96]byte
	bad  bool
}

func (im *tIoMem) ReadData(addr instr.Addr) byte {
	if addr < 0x20 || addr >= 0x60 {
		im.bad = true
		return 0
	}
	return im.data[addr]
}

func (im *tIoMem) WriteData(addr instr.Addr, val byte) {
	if addr < 0x20 || addr >= 0x60 {
		im.bad = true
		return
	}
	im.data[addr] = val
}

func (im *tIoMem) ReadProgram(addr instr.Addr) uint16 {
	return 0
}

func (this *tIoMem) diff(that *tIoMem) string {
	if this.bad || that.bad {
		return "MEM: out of range access"
	}
	var mThis, mThat string
	var thisM, thatM []string
	for i := 0; i < 32; i++ {
		if this.data[i] != that.data[i] {
			thisM = append(thisM, fmt.Sprintf("%d=%02x", i, this.data[i]))
			thatM = append(thatM, fmt.Sprintf("%d=%02x", i, that.data[i]))
		}
	}
	if thisM != nil {
		mThis = "MEM[" + strings.Join(thisM, ",") + "]"
		mThat = "MEM[" + strings.Join(thatM, ",") + "]"
		return mThis + "\n" + mThat
	}
	return ""
}

type tCpuIm struct {
	tCpu
	imem tIoMem
}

func (tc tCpuIm) Equals(other testcase.Testable) bool {
	o := other.(tCpuIm)
	return tc.tCpu.Equals(o.tCpu) && tc.imem == o.imem
}

func (tc tCpuIm) Diff(other testcase.Testable) interface{} {
	o := other.(tCpuIm)
	cDiff := tc.tCpu.Diff(o.tCpu)
	mDiff := tc.imem.diff(&o.imem)
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
