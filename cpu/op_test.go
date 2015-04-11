package cpu

import (
	"github.com/edmccard/testcase"
	"testing"
)

// Branches with all flags initially set/all flags initially clear.
func flagsOnOff(t testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	initTc.SregFromByte(0xff)
	t.Run("SRff", initTc, exp)
	initTc.SregFromByte(0x00)
	t.Run("SR00", initTc, exp)
}

// Branches with carry flag set/carry flag clear
func carryOnOff(t testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	initTc.FlagC = true
	t.Run("C1", initTc, exp)
	initTc.FlagC = false
	t.Run("C0", initTc, exp)
}

var addCases = [][]struct {
	status, v1, v2, res int
}{
	{
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
	},
	{
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
	},
}

func addIgnoreCarry(t testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	for n, c := range addCases[0] {
		ac := arithCase{0x3f, c.status, c.v1, c.v2, c.res, n}
		ac.testD5R5(t, initTc, Add, "(ign.) Add")
	}
}

func addRespectCarry(t testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	cIdx := 0
	if initTc.FlagC {
		cIdx = 1
	}

	for n, c := range addCases[cIdx] {
		ac := arithCase{0x3f, c.status, c.v1, c.v2, c.res, n}
		ac.testD5R5(t, initTc, Adc, "(resp.) Adc")
	}
}

func TestAddition(t *testing.T) {
	testcase.NewTree(t, "+",
		flagsOnOff, carryOnOff, addIgnoreCarry).Start(tCpu{})
	testcase.NewTree(t, "+",
		flagsOnOff, carryOnOff, addRespectCarry).Start(tCpu{})

	adiw := func(t testcase.Tree, init, exp testcase.Testable) {
		var cases = []struct {
			status, v1, v2, res int
		}{
			{0x00, 0x0000, 0x01, 0x0001},
			{0x01, 0xffc3, 0x3e, 0x0001},
			{0x02, 0x0000, 0x00, 0x0000},
			{0x03, 0xffc2, 0x3e, 0x0000},
			{0x0c, 0x7fc2, 0x3e, 0x8000},
			{0x14, 0x8000, 0x00, 0x8000},
		}
		initTc := init.(tCpu)
		for n, c := range cases {
			ac := arithCase{0x1f, c.status, c.v1, c.v2, c.res, n}
			ac.testDDK6(t, initTc, Adiw, "Adiw")
		}
	}
	testcase.NewTree(t, "+", flagsOnOff, carryOnOff, adiw).Start(tCpu{})
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
func subIgnoreCarry(t testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	for n, c := range subCases[0] {
		acsub := arithCase{0x3f, c.status, c.v1, c.v2, c.res, n}
		accp := arithCase{0x3f, c.status, c.v1, c.v2, c.v1, n}
		acsub.testD5R5(t, initTc, Sub, "(ign.) Sub")
		acsub.testD4K8(t, initTc, Subi, "(ign.) Subi")
		accp.testD5R5(t, initTc, Cp, "(ign.) Cp")
		accp.testD4K8(t, initTc, Cpi, "(ign.) Cpi")
	}
}

// Tests Sbc, Cpc, Sbci with cases for each possible status outcome.
func subRespectCarry(t testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	cIdx := 0
	if initTc.FlagC {
		cIdx = 1
	}

	for n, c := range subCases[cIdx] {
		acsub := arithCase{0x3f, c.status, c.v1, c.v2, c.res, n}
		accp := arithCase{0x3f, c.status, c.v1, c.v2, c.v1, n}
		if c.res == 0 {
			acsub.mask = 0x3d
			accp.mask = 0x3d
		}
		acsub.testD5R5(t, initTc, Sbc, "(resp.) Sbc")
		accp.testD5R5(t, initTc, Cpc, "(resp.) Cpc")
		acsub.testD4K8(t, initTc, Sbci, "(resp.) Sbci")
	}
}

func TestSubtraction(t *testing.T) {
	testcase.NewTree(t, "-",
		flagsOnOff, carryOnOff, subIgnoreCarry).Start(tCpu{})
	testcase.NewTree(t, "-",
		flagsOnOff, carryOnOff, subRespectCarry).Start(tCpu{})

	sbiw := func(t testcase.Tree, init, exp testcase.Testable) {
		var cases = []struct {
			status, v1, v2, res int
		}{
			{0x00, 0x0001, 0x00, 0x0001},
			{0x02, 0x0000, 0x00, 0x0000},
			{0x14, 0x8000, 0x00, 0x8000},
			{0x15, 0x0000, 0x01, 0xffff},
			{0x18, 0x8000, 0x01, 0x7fff},
		}
		initTc := init.(tCpu)
		for n, c := range cases {
			ac := arithCase{0x1f, c.status, c.v1, c.v2, c.res, n}
			ac.testDDK6(t, initTc, Sbiw, "Sbiw")
		}
	}
	testcase.NewTree(t, "-", flagsOnOff, carryOnOff, sbiw).Start(tCpu{})
}

func andAndi(t testcase.Tree, init, exp testcase.Testable) {
	var andCases = []struct {
		status, v1, v2, res int
	}{
		{0x00, 0x01, 0x01, 0x01},
		{0x02, 0xaa, 0x55, 0x00},
		{0x14, 0x80, 0x81, 0x80},
	}
	initTc := init.(tCpu)
	for n, c := range andCases {
		ac := arithCase{0x1e, c.status, c.v1, c.v2, c.res, n}
		ac.testD5R5(t, initTc, And, "And")
		ac.testD4K8(t, initTc, Andi, "Andi")
	}
}

func orOri(t testcase.Tree, init, exp testcase.Testable) {
	var orCases = []struct {
		status, v1, v2, res int
	}{
		{0x00, 0x01, 0x03, 0x03},
		{0x02, 0x00, 0x00, 0x00},
		{0x14, 0x80, 0x01, 0x81},
	}

	initTc := init.(tCpu)
	for n, c := range orCases {
		ac := arithCase{0x1e, c.status, c.v1, c.v2, c.res, n}
		ac.testD5R5(t, initTc, Or, "Or")
		ac.testD4K8(t, initTc, Ori, "Ori")
	}
}

func eorEor(t testcase.Tree, init, exp testcase.Testable) {
	var eorCases = []struct {
		status, v1, v2, res int
	}{
		{0x00, 0x01, 0x03, 0x02},
		{0x02, 0xaa, 0xaa, 0x00},
		{0x14, 0xaa, 0x55, 0xff},
	}
	initTc := init.(tCpu)
	for n, c := range eorCases {
		ac := arithCase{0x1e, c.status, c.v1, c.v2, c.res, n}
		ac.testD5R5(t, initTc, Eor, "Eor")
	}
}

func TestBoolean(t *testing.T) {
	testcase.NewTree(t, "&", flagsOnOff, andAndi).Start(tCpu{})
	testcase.NewTree(t, "|", flagsOnOff, orOri).Start(tCpu{})
	testcase.NewTree(t, "^", flagsOnOff, eorEor).Start(tCpu{})
}
