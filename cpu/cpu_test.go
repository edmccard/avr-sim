package cpu

import (
	"fmt"
	"github.com/edmccard/avr-sim/instr"
	"github.com/edmccard/testcase"
	"reflect"
	"strings"
)

// tCpu wraps Cpu to implement Testable, and to simplify setting up
// initial state for tests.
type tCpu struct {
	Cpu
	am   instr.AddrMode
	dmem tDataMem
}

func (tc tCpu) Equals(other testcase.Testable) bool {
	o := other.(tCpu)
	return tc.Cpu == o.Cpu && tc.am == o.am &&
		(&tc.dmem).equals(&o.dmem)
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

type tDataMem struct {
	readAddrs  []instr.Addr
	writeAddrs []instr.Addr
	writeVals  []byte
}

func (tdm *tDataMem) ReadData(addr instr.Addr) (byte, error) {
	tdm.readAddrs = append(tdm.readAddrs, addr)
	return 0xff, nil
}

func (tdm *tDataMem) WriteData(addr instr.Addr, val byte) error {
	tdm.writeAddrs = append(tdm.writeAddrs, addr)
	tdm.writeVals = append(tdm.writeVals, val)
	return nil
}

func (this *tDataMem) equals(that *tDataMem) bool {
	return reflect.DeepEqual(this, that)
}

func (this *tDataMem) diff(that *tDataMem) string {
	// Assumes this != that; assumes 1 read xor 1 write;
	// assumes this is expected and that is actual
	switch {
	case len(this.readAddrs) < len(that.readAddrs):
		return "MEM: too many reads"
	case len(this.readAddrs) > len(that.readAddrs):
		return "MEM: too few reads"
	case len(this.writeAddrs) < len(that.writeAddrs):
		return "MEM: too many writes"
	case len(this.writeAddrs) > len(that.writeAddrs):
		return "MEM: too few writes"
	}
	if len(this.readAddrs) > 0 {
		if this.readAddrs[0] != that.readAddrs[0] {
			return fmt.Sprintf("MEM: expected read %04x got %04x",
				this.readAddrs[0], that.readAddrs[0])
		}
	}
	if len(this.writeAddrs) > 0 {
		if this.writeAddrs[0] != that.writeAddrs[1] {
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
