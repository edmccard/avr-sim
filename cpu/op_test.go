package cpu

import (
	"fmt"
	"github.com/edmccard/testcase"
	"testing"
)

// Branches with all flags initially set/all flags initially clear.
func flagsOnOff(t testcase.Tree, init testcase.Testable) {
	initTc := init.(tCpu)
	initTc.SregFromByte(0xff)
	t.Run("SRff", initTc)
	initTc.SregFromByte(0x00)
	t.Run("SR00", initTc)
}

// Branches with carry flag set/carry flag clear
func carryOnOff(t testcase.Tree, init testcase.Testable) {
	initTc := init.(tCpu)
	initTc.FlagC = true
	t.Run("C1", initTc)
	initTc.FlagC = false
	t.Run("C0", initTc)
}

// Tests Adc with carry false, and Add with carry true/false,
// with cases for each possible status outcome.
func addNoCarry(t testcase.Tree, init testcase.Testable) {
	cases := []struct {
		status, v1, v2, res int
	}{
		{0x00, 0x00, 0x01, 0x01},
		{0x01, 0x10, 0xf1, 0x01},
		{0x02, 0x00, 0x00, 0x00},
		{0x03, 0x10, 0xf0, 0x00},
		{0x0c, 0x10, 0x70, 0x80},
		{0x14, 0x00, 0x80, 0x80},
		{0x15, 0x90, 0xf0, 0x80},
		{0x19, 0x80, 0x81, 0x01},
		{0x1b, 0x80, 0x80, 0x00},
		{0x20, 0x01, 0x0f, 0x10},
		{0x21, 0x02, 0xff, 0x01},
		{0x23, 0x01, 0xff, 0x00},
		{0x2c, 0x01, 0x7f, 0x80},
		{0x34, 0x01, 0x8f, 0x90},
		{0x35, 0x81, 0xff, 0x80},
		{0x39, 0x81, 0x8f, 0x10},
	}
	initTc := init.(tCpu)
	var initCpu tCpu
	for n, c := range cases {
		ac := arithCase{0x3f, c.status, c.v1, c.v2, c.res}

		initCpu, t.Exp = initTc.setD5R5Case(ac)
		initCpu.FlagC = false
		Adc(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("Adc C0[%d]", n), initCpu)

		initCpu, t.Exp = initTc.setD5R5Case(ac)
		initCpu.FlagC = false
		Add(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("Add C0[%d]", n), initCpu)

		initCpu, t.Exp = initTc.setD5R5Case(ac)
		initCpu.FlagC = true
		Add(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("Add C1[%d]", n), initCpu)
	}
}

// Tests Adc with carry true, with cases for each possible status outcome.
func addCarry(t testcase.Tree, init testcase.Testable) {
	cases := []struct {
		status, v1, v2, res int
	}{
		{0x00, 0x00, 0x00, 0x01},
		{0x01, 0x10, 0xf0, 0x01},
		{0x0c, 0x10, 0x70, 0x81},
		{0x14, 0x00, 0x80, 0x81},
		{0x15, 0x90, 0xf0, 0x81},
		{0x19, 0x80, 0x80, 0x01},
		{0x20, 0x00, 0x0f, 0x10},
		{0x21, 0x01, 0xff, 0x01},
		{0x23, 0x00, 0xff, 0x00},
		{0x2c, 0x00, 0x7f, 0x80},
		{0x34, 0x00, 0x8f, 0x90},
		{0x35, 0x80, 0xff, 0x80},
		{0x39, 0x80, 0x8f, 0x10},
	}
	initTc := init.(tCpu)
	var initCpu tCpu
	for n, c := range cases {
		ac := arithCase{0x3f, c.status, c.v1, c.v2, c.res}
		initCpu, t.Exp = initTc.setD5R5Case(ac)
		initCpu.FlagC = true
		Adc(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("Adc C0[%d]", n), initCpu)
	}
}

func TestAddition(t *testing.T) {
	testcase.NewTree(t, "+", flagsOnOff, addNoCarry).Run("", tCpu{})
	testcase.NewTree(t, "+", flagsOnOff, addCarry).Run("", tCpu{})
}

var subCases = [][]struct {
	status, v1, v2, res int
}{
	{
		{0x00, 0x01, 0x00, 0x01},
		{0x01, 0x00, 0x90, 0x70},
		{0x02, 0x00, 0x00, 0x00},
		{0x0d, 0x00, 0x80, 0x80},
		{0x14, 0x80, 0x00, 0x80},
		{0x15, 0x00, 0x10, 0xf0},
		{0x18, 0x80, 0x10, 0x70},
		{0x20, 0x10, 0x01, 0x0f},
		{0x21, 0x00, 0x81, 0x7f},
		{0x2d, 0x10, 0x81, 0x8f},
		{0x34, 0x90, 0x01, 0x8f},
		{0x35, 0x00, 0x01, 0xff},
		{0x38, 0x80, 0x01, 0x7f}},
	{
		{0x00, 0x02, 0x00, 0x01},
		{0x01, 0x01, 0x90, 0x70},
		{0x02, 0x01, 0x00, 0x00},
		{0x0d, 0x01, 0x80, 0x80},
		{0x14, 0x81, 0x00, 0x80},
		{0x15, 0x01, 0x10, 0xf0},
		{0x18, 0x81, 0x10, 0x70},
		{0x20, 0x10, 0x00, 0x0f},
		{0x21, 0x00, 0x80, 0x7f},
		{0x22, 0x10, 0x0f, 0x00},
		{0x23, 0x00, 0xff, 0x00},
		{0x2d, 0x10, 0x80, 0x8f},
		{0x34, 0x90, 0x00, 0x8f},
		{0x35, 0x00, 0x00, 0xff},
		{0x38, 0x80, 0x00, 0x7f},
		{0x3a, 0x80, 0x7f, 0x00},
	},
}

// Tests Sub, Subi, Cp, Cpi with cases for each possible status outcome.
func subIgnoreCarry(t testcase.Tree, init testcase.Testable) {
	initTc := init.(tCpu)
	var initCpu tCpu
	for n, c := range subCases[0] {
		acsub := arithCase{0x3f, c.status, c.v1, c.v2, c.res}
		accp := arithCase{0x3f, c.status, c.v1, c.v2, c.v1}

		initCpu, t.Exp = initTc.setD5R5Case(acsub)
		Sub(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("(ig.) Sub [%d]", n), initCpu)

		initCpu, t.Exp = initTc.setD4K8Case(acsub)
		Subi(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("(ig.) Subi [%d]", n), initCpu)

		initCpu, t.Exp = initTc.setD5R5Case(accp)
		Cp(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("(ig.) Cp [%d]", n), initCpu)

		initCpu, t.Exp = initTc.setD4K8Case(accp)
		Cpi(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("(ig.) Cpi [%d]", n), initCpu)
	}
}

// Tests Sbc, Cpc, Sbci with cases for each possible status outcome.
func subRespectCarry(t testcase.Tree, init testcase.Testable) {
	initTc := init.(tCpu)
	var initCpu tCpu
	cIdx := 0
	if initTc.FlagC {
		cIdx = 1
	}

	for n, c := range subCases[cIdx] {
		acsub := arithCase{0x3f, c.status, c.v1, c.v2, c.res}
		accp := arithCase{0x3f, c.status, c.v1, c.v2, c.v1}

		if c.res == 0 {
			acsub.mask = 0x3d
			accp.mask = 0x3d
		}

		initCpu, t.Exp = initTc.setD5R5Case(acsub)
		Sbc(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("(resp.) Sbc [%d]", n), initCpu)

		initCpu, t.Exp = initTc.setD5R5Case(accp)
		Cpc(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("(resp.) Cpc [%d]", n), initCpu)

		initCpu, t.Exp = initTc.setD4K8Case(acsub)
		Sbci(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("(resp.) Sbci [%d]", n), initCpu)
	}
}

func TestSubtraction(t *testing.T) {
	testcase.NewTree(t, "-",
		flagsOnOff, carryOnOff, subIgnoreCarry).Run("", tCpu{})
	testcase.NewTree(t, "-",
		flagsOnOff, carryOnOff, subRespectCarry).Run("", tCpu{})
}

func andAndi(t testcase.Tree, init testcase.Testable) {
	var andCases = []struct {
		status, v1, v2, res int
	}{
		{0x00, 0x01, 0x01, 0x01},
		{0x02, 0xaa, 0x55, 0x00},
		{0x14, 0x80, 0x81, 0x80},
	}
	initTc := init.(tCpu)
	var initCpu tCpu

	for n, c := range andCases {
		ac := arithCase{0x1e, c.status, c.v1, c.v2, c.res}

		initCpu, t.Exp = initTc.setD5R5Case(ac)
		And(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("And [%d]", n), initCpu)

		initCpu, t.Exp = initTc.setD4K8Case(ac)
		Andi(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("Andi [%d]", n), initCpu)
	}
}

func orOri(t testcase.Tree, init testcase.Testable) {
	var orCases = []struct {
		status, v1, v2, res int
	}{
		{0x00, 0x01, 0x03, 0x03},
		{0x02, 0x00, 0x00, 0x00},
		{0x14, 0x80, 0x01, 0x81},
	}

	initTc := init.(tCpu)
	var initCpu tCpu
	for n, c := range orCases {
		ac := arithCase{0x1e, c.status, c.v1, c.v2, c.res}

		initCpu, t.Exp = initTc.setD5R5Case(ac)
		Or(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("Or [%d]", n), initCpu)

		initCpu, t.Exp = initTc.setD4K8Case(ac)
		Ori(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("Ori [%d]", n), initCpu)
	}
}

func eorEor(t testcase.Tree, init testcase.Testable) {
	var eorCases = []struct {
		status, v1, v2, res int
	}{
		{0x00, 0x01, 0x03, 0x02},
		{0x02, 0xaa, 0xaa, 0x00},
		{0x14, 0xaa, 0x55, 0xff},
	}
	initTc := init.(tCpu)
	var initCpu tCpu

	for n, c := range eorCases {
		ac := arithCase{0x1e, c.status, c.v1, c.v2, c.res}

		initCpu, t.Exp = initTc.setD5R5Case(ac)
		Eor(&(initCpu.Cpu), &(initCpu.am))
		t.Run(fmt.Sprintf("Eor [%d]", n), initCpu)
	}
}

func TestBoolean(t *testing.T) {
	testcase.NewTree(t, "&", flagsOnOff, andAndi).Run("", tCpu{})
	testcase.NewTree(t, "|", flagsOnOff, orOri).Run("", tCpu{})
	testcase.NewTree(t, "^", flagsOnOff, eorEor).Run("", tCpu{})
}
