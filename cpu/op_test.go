package cpu

import (
	"fmt"
	"github.com/edmccard/avr-sim/instr"
	"github.com/edmccard/testcase"
	"math/rand"
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

// Branches with carry flag set/carry flag clear.
func carryOnOff(t testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	initTc.FlagC = true
	t.Run("C1", initTc, exp)
	initTc.FlagC = false
	t.Run("C0", initTc, exp)
}

// Branches with (d, d+1) and (d,d) as dest/source registers.
func regD5R5(t testcase.Tree, init, exp testcase.Testable) {
	d := instr.Addr(rand.Intn(32))
	r := (d + 1) & 0x1f
	initTc := init.(tCpu)
	initTc.am = instr.AddrMode{d, r, instr.NoIndex}
	t.Run(fmt.Sprintf("r%02d,r%02d", d, r), initTc, exp)
	initTc.am = instr.AddrMode{d, d, instr.NoIndex}
	t.Run(fmt.Sprintf("r%02d,r%02d", d, d), initTc, exp)
}

// Branches with various dest/source registers for Mul instruction.
func reg5Mul(t testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	regs := []struct {
		d, r instr.Addr
	}{{0, 1}, {0, 0}, {1, 1}, {1, 2}, {2, 3}}
	for _, reg := range regs {
		initTc.am = instr.AddrMode{reg.d, reg.r, instr.NoIndex}
		t.Run(fmt.Sprintf("r%02d,r%02d", reg.r, reg.d), initTc, exp)
	}
}

// Branches with (16,17) and (16,16) as d,r
// for Muls/Mulsu/Fmul/Fmuls/Fmulsu.
func reg4Mul(t testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	initTc.am = instr.AddrMode{16, 17, instr.NoIndex}
	t.Run("r16,r17", initTc, exp)
	initTc.am = instr.AddrMode{16, 16, instr.NoIndex}
	t.Run("r16,r16", initTc, exp)
}

var addCases = [][]caseData{
	{
		{0x00, 0x01, 0x01, 0x02},
		{0x01, 0x10, 0xf1, 0x01},
		{0x02, 0x00, 0x00, 0x00},
		{0x03, 0x10, 0xf0, 0x00},
		{0x0c, 0x40, 0x40, 0x80},
		{0x14, 0x00, 0x80, 0x80},
		{0x15, 0xc0, 0xc0, 0x80},
		{0x19, 0x81, 0x81, 0x02},
		{0x1b, 0x80, 0x80, 0x00},
		{0x20, 0x08, 0x08, 0x10},
		{0x21, 0x02, 0xff, 0x01},
		{0x23, 0x01, 0xff, 0x00},
		{0x2c, 0x48, 0x48, 0x90},
		{0x34, 0x01, 0x8f, 0x90},
		{0x35, 0xc8, 0xc8, 0x90},
		{0x39, 0x88, 0x88, 0x10},
	},
	{
		{0x00, 0x00, 0x00, 0x01},
		{0x01, 0x10, 0xf0, 0x01},
		{0x0c, 0x40, 0x40, 0x81},
		{0x14, 0x00, 0x80, 0x81},
		{0x15, 0xc0, 0xc0, 0x81},
		{0x19, 0x80, 0x80, 0x01},
		{0x20, 0x08, 0x08, 0x11},
		{0x21, 0x01, 0xff, 0x01},
		{0x23, 0x00, 0xff, 0x00},
		{0x2c, 0x48, 0x48, 0x91},
		{0x34, 0x00, 0x8f, 0x90},
		{0x35, 0xc8, 0xc8, 0x91},
		{0x39, 0x88, 0x88, 0x11},
	},
}

func addIgnoreCarry(t testcase.Tree, init, exp testcase.Testable) {
	for n, c := range addCases[0] {
		ac := arithCase{t, init.(tCpu), 0x3f, c, n}
		ac.testD5R5(Add, "(ign.) Add")
	}
}

func addRespectCarry(t testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	cIdx := 0
	if initTc.FlagC {
		cIdx = 1
	}

	for n, c := range addCases[cIdx] {
		ac := arithCase{t, init.(tCpu), 0x3f, c, n}
		ac.testD5R5(Adc, "(resp.) Adc")
	}
}

func TestAddition(t *testing.T) {
	testcase.NewTree(t, "+",
		flagsOnOff, carryOnOff, regD5R5, addIgnoreCarry).Start(tCpu{})
	testcase.NewTree(t, "+",
		flagsOnOff, carryOnOff, regD5R5, addRespectCarry).Start(tCpu{})

	adiw := func(t testcase.Tree, init, exp testcase.Testable) {
		var cases = []caseData{
			{0x00, 0x0000, 0x01, 0x0001},
			{0x01, 0xffc3, 0x3e, 0x0001},
			{0x02, 0x0000, 0x00, 0x0000},
			{0x03, 0xffc2, 0x3e, 0x0000},
			{0x0c, 0x7fc2, 0x3e, 0x8000},
			{0x14, 0x8000, 0x00, 0x8000},
		}
		for n, c := range cases {
			ac := arithCase{t, init.(tCpu), 0x1f, c, n}
			ac.testDDK6(Adiw, "Adiw")
		}
	}
	testcase.NewTree(t, "+",
		flagsOnOff, carryOnOff, adiw).Start(tCpu{})
}

var subCases = [][]caseData{
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
	for n, c := range subCases[0] {
		acsub := arithCase{t, init.(tCpu), 0x3f, c, n}
		accp := acsub
		accp.res = accp.v1
		acsub.testD5R5(Sub, "(ign.) Sub")
		acsub.testD4K8(Subi, "(ign.) Subi")
		accp.testD5R5(Cp, "(ign.) Cp")
		accp.testD4K8(Cpi, "(ign.) Cpi")
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
		acsub := arithCase{t, init.(tCpu), 0x3f, c, n}
		accp := acsub
		accp.res = accp.v1
		if c.res == 0 {
			acsub.mask = 0x3d
			accp.mask = 0x3d
		}
		acsub.testD5R5(Sbc, "(resp.) Sbc")
		accp.testD5R5(Cpc, "(resp.) Cpc")
		acsub.testD4K8(Sbci, "(resp.) Sbci")
	}
}

func TestSubtraction(t *testing.T) {
	testcase.NewTree(t, "-",
		flagsOnOff, carryOnOff, regD5R5, subIgnoreCarry).Start(tCpu{})
	testcase.NewTree(t, "-",
		flagsOnOff, carryOnOff, regD5R5, subRespectCarry).Start(tCpu{})

	sbiw := func(t testcase.Tree, init, exp testcase.Testable) {
		var cases = []caseData{
			{0x00, 0x0001, 0x00, 0x0001},
			{0x02, 0x0000, 0x00, 0x0000},
			{0x14, 0x8000, 0x00, 0x8000},
			{0x15, 0x0000, 0x01, 0xffff},
			{0x18, 0x8000, 0x01, 0x7fff},
		}
		for n, c := range cases {
			ac := arithCase{t, init.(tCpu), 0x1f, c, n}
			ac.testDDK6(Sbiw, "Sbiw")
		}
	}
	testcase.NewTree(t, "-", flagsOnOff, carryOnOff, sbiw).Start(tCpu{})
}

func andAndi(t testcase.Tree, init, exp testcase.Testable) {
	var andCases = []caseData{
		{0x00, 0x01, 0x01, 0x01},
		{0x02, 0xaa, 0x55, 0x00},
		{0x14, 0x80, 0x80, 0x80},
	}
	for n, c := range andCases {
		ac := arithCase{t, init.(tCpu), 0x1e, c, n}
		ac.testD5R5(And, "And")
		ac.testD4K8(Andi, "Andi")
	}
}

func orOri(t testcase.Tree, init, exp testcase.Testable) {
	var orCases = []caseData{
		{0x00, 0x01, 0x03, 0x03},
		{0x02, 0x00, 0x00, 0x00},
		{0x14, 0x80, 0x01, 0x81},
	}
	for n, c := range orCases {
		ac := arithCase{t, init.(tCpu), 0x1e, c, n}
		ac.testD5R5(Or, "Or")
		ac.testD4K8(Ori, "Ori")
	}
}

func eorEor(t testcase.Tree, init, exp testcase.Testable) {
	var eorCases = []caseData{
		{0x00, 0x01, 0x03, 0x02},
		{0x02, 0xaa, 0xaa, 0x00},
		{0x14, 0xaa, 0x55, 0xff},
	}
	for n, c := range eorCases {
		ac := arithCase{t, init.(tCpu), 0x1e, c, n}
		ac.testD5R5(Eor, "Eor")
	}
}

func TestBoolean(t *testing.T) {
	testcase.NewTree(t, "&",
		flagsOnOff, regD5R5, andAndi).Start(tCpu{})
	testcase.NewTree(t, "|",
		flagsOnOff, regD5R5, orOri).Start(tCpu{})
	testcase.NewTree(t, "^",
		flagsOnOff, regD5R5, eorEor).Start(tCpu{})
}

func mulMul(t testcase.Tree, init, exp testcase.Testable) {
	var cases = []caseData{
		{0x00, 0xff, 0x01, 0x00ff},
		{0x00, 0x7f, 0x7f, 0x3f01},
		{0x01, 0xff, 0xff, 0xfe01},
		{0x02, 0xff, 0x00, 0x0000},
	}
	for n, c := range cases {
		ac := arithCase{t, init.(tCpu), 0x03, c, n}
		ac.testMul(Mul, "Mul")
	}
}

func mulMuls(t testcase.Tree, init, exp testcase.Testable) {
	var cases = []caseData{
		{0x00, 0xff, 0xff, 0x0001},
		{0x00, 0x7f, 0x7f, 0x3f01},
		{0x01, 0xff, 0x01, 0xffff},
		{0x02, 0xff, 0x00, 0x0000},
	}
	for n, c := range cases {
		ac := arithCase{t, init.(tCpu), 0x03, c, n}
		ac.testMul(Muls, "Muls")
	}
}

func mulMulsu(t testcase.Tree, init, exp testcase.Testable) {
	var cases = []caseData{
		{0x00, 0x01, 0xff, 0x00ff},
		{0x00, 0x7f, 0x7f, 0x3f01},
		{0x01, 0xff, 0xff, 0xff01},
		{0x02, 0xff, 0x00, 0x0000},
	}
	for n, c := range cases {
		ac := arithCase{t, init.(tCpu), 0x03, c, n}
		ac.testMul(Mulsu, "Mulsu")
	}
}

func TestMultiplication(t *testing.T) {
	testcase.NewTree(t, "*", flagsOnOff, reg5Mul, mulMul).Start(tCpu{})
	testcase.NewTree(t, "*", flagsOnOff, reg4Mul, mulMuls).Start(tCpu{})
	testcase.NewTree(t, "*", flagsOnOff, reg4Mul, mulMulsu).Start(tCpu{})
}
