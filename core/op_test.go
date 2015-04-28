package core

import (
	"fmt"
	"github.com/edmccard/avr-sim/instr"
	"github.com/edmccard/testcase"
	"math/rand"
	"testing"
)

type aluData struct {
	status, v1, v2, res int
}

type aluCase struct {
	t    testcase.Tree
	init tCpu
	mask int
	aluData
	n int
}

var alucases = []struct {
	label    string
	mnem     instr.Mnemonic
	tag      string
	fmask    int
	run      aluTest
	branches []testcase.Branch
	data     []aluData
}{
	{"+-*", instr.Add, "Add", 0x3f, testD5R5,
		[]testcase.Branch{brD5R5, brCarry},
		addCClrOrIgnored},
	{"+-*", instr.Adc, "C0 Adc", 0x3f, testD5R5,
		[]testcase.Branch{brD5R5, brCarry, ifCClr},
		addCClrOrIgnored},
	{"+-*", instr.Adc, "C1 Adc", 0x3f, testD5R5,
		[]testcase.Branch{brD5R5, brCarry, ifCSet},
		addCSet},
	{"+-*", instr.Adiw, "Adiw", 0x1f, testDDK6,
		[]testcase.Branch{brCarry},
		[]aluData{
			{0x00, 0x0000, 0x01, 0x0001},
			{0x01, 0xffc3, 0x3e, 0x0001},
			{0x02, 0x0000, 0x00, 0x0000},
			{0x03, 0xffc2, 0x3e, 0x0000},
			{0x0c, 0x7fc2, 0x3e, 0x8000},
			{0x14, 0x8000, 0x00, 0x8000}}},
	{"+-*", 0, "SubIgnoreCarry", 0x3f, testSubIgnoreCarry,
		[]testcase.Branch{brD5R5, brCarry},
		subCClrOrIgnored},
	{"+-*", 0, "SubRespectCarry", 0x3f, testSubRespectCarry,
		[]testcase.Branch{brD5R5, brZero, brCarry, ifCClr},
		subCClrOrIgnored},
	{"+-*", 0, "SubRespectCarry", 0x3f, testSubRespectCarry,
		[]testcase.Branch{brD5R5, brZero, brCarry, ifCSet},
		subCSet},
	{"+-*", instr.Sbiw, "Sbiw", 0x1f, testDDK6,
		[]testcase.Branch{brCarry},
		[]aluData{
			{0x00, 0x0001, 0x00, 0x0001},
			{0x02, 0x0000, 0x00, 0x0000},
			{0x14, 0x8000, 0x00, 0x8000},
			{0x15, 0x0000, 0x01, 0xffff},
			{0x18, 0x8000, 0x01, 0x7fff}}},
	{"+-*", instr.Mul, "Mul", 0x03, testMul,
		[]testcase.Branch{brMul5},
		[]aluData{
			{0x00, 0xff, 0x01, 0x00ff},
			{0x00, 0x7f, 0x7f, 0x3f01},
			{0x01, 0xff, 0xff, 0xfe01},
			{0x02, 0xff, 0x00, 0x0000}}},
	{"+-*", instr.Muls, "Muls", 0x03, testMul,
		[]testcase.Branch{brMul34},
		[]aluData{
			{0x00, 0xff, 0xff, 0x0001},
			{0x00, 0x7f, 0x7f, 0x3f01},
			{0x01, 0xff, 0x01, 0xffff},
			{0x02, 0xff, 0x00, 0x0000}}},
	{"+-*", instr.Mulsu, "Mulsu", 0x03, testMul,
		[]testcase.Branch{brMul34},
		[]aluData{
			{0x00, 0x01, 0xff, 0x00ff},
			{0x00, 0x7f, 0x7f, 0x3f01},
			{0x01, 0xff, 0xff, 0xff01},
			{0x02, 0xff, 0x00, 0x0000}}},
	{"+-*", instr.Fmul, "Fmul", 0x03, testMul,
		[]testcase.Branch{brMul34},
		[]aluData{
			{0x00, 0xff, 0x01, 0x01fe},
			{0x00, 0x80, 0x80, 0x8000},
			{0x01, 0xd0, 0xd0, 0x5200},
			{0x01, 0xe0, 0xe0, 0x8800},
			{0x02, 0xff, 0x00, 0x0000}}},
	{"+-*", instr.Fmuls, "Fmuls", 0x03, testMul,
		[]testcase.Branch{brMul34},
		[]aluData{
			{0x00, 0x7f, 0x7f, 0x7e02},
			{0x00, 0x80, 0x80, 0x8000},
			{0x01, 0xff, 0x01, 0xfffe},
			{0x02, 0xff, 0x00, 0x0000}}},
	{"+-*", instr.Fmulsu, "Fmulsu", 0x03, testMul,
		[]testcase.Branch{brMul34},
		[]aluData{
			{0x00, 0x01, 0xff, 0x01fe},
			{0x00, 0x7f, 0xc8, 0xc670},
			{0x01, 0xff, 0xff, 0xfe02},
			{0x01, 0x9c, 0xaa, 0x7b30},
			{0x02, 0xff, 0x00, 0x0000}}},
	{"&|^", instr.And, "And", 0x1e, testD5R5,
		[]testcase.Branch{brD5R5},
		andData},
	{"&|^", instr.Andi, "Andi", 0x1e, testD4K8,
		nil,
		andData},
	{"&|^", instr.Or, "Or", 0x1e, testD5R5,
		[]testcase.Branch{brD5R5},
		orData},
	{"&|^", instr.Ori, "Ori", 0x1e, testD4K8,
		nil,
		orData},
	{"&|^", instr.Eor, "Eor", 0x1e, testD5R5,
		[]testcase.Branch{brD5R5},
		[]aluData{
			{0x00, 0x01, 0x03, 0x02},
			{0x02, 0xaa, 0xaa, 0x00},
			{0x14, 0xaa, 0x55, 0xff}}},
	{"RMW", instr.Asr, "Asr", 0x1f, testD5,
		[]testcase.Branch{brCarry},
		[]aluData{
			{0x00, 0x02, 0x02, 0x01},
			{0x02, 0x00, 0x00, 0x00},
			{0x0c, 0x80, 0x80, 0xc0},
			{0x15, 0x81, 0x81, 0xc0},
			{0x19, 0x03, 0x03, 0x01},
			{0x1b, 0x01, 0x01, 0x00}}},
	{"RMW", instr.Lsr, "Lsr", 0x1f, testD5,
		[]testcase.Branch{brCarry},
		[]aluData{
			{0x00, 0x02, 0x02, 0x01},
			{0x02, 0x00, 0x00, 0x00},
			{0x19, 0x03, 0x03, 0x01},
			{0x1b, 0x01, 0x01, 0x00}}},
	{"RMW", instr.Com, "Com", 0x1f, testD5,
		nil,
		[]aluData{
			{0x01, 0x80, 0x80, 0x7f},
			{0x03, 0xff, 0xff, 0x00},
			{0x15, 0x00, 0x00, 0xff}}},
	{"RMW", instr.Neg, "Neg", 0x3f, testD5,
		nil,
		[]aluData{
			{0x01, 0x90, 0x90, 0x70},
			{0x02, 0x00, 0x00, 0x00},
			{0x0d, 0x80, 0x80, 0x80},
			{0x15, 0x10, 0x10, 0xf0},
			{0x21, 0x81, 0x81, 0x7f},
			{0x35, 0x01, 0x01, 0xff}}},
	{"RMW", instr.Swap, "Swap", 0x00, testD5,
		nil,
		[]aluData{
			{0x00, 0xff, 0xff, 0xff},
			{0x00, 0x00, 0x00, 0x00},
			{0x00, 0x12, 0x12, 0x21}}},
	{"RMW", instr.Dec, "Dec", 0x1e, testD5,
		nil,
		[]aluData{
			{0x00, 0x02, 0x02, 0x01},
			{0x02, 0x01, 0x01, 0x00},
			{0x14, 0x00, 0x00, 0xff},
			{0x18, 0x80, 0x80, 0x7f}}},
	{"RMW", instr.Inc, "Inc", 0x1e, testD5,
		nil,
		[]aluData{
			{0x00, 0x00, 0x00, 0x01},
			{0x02, 0xff, 0xff, 0x00},
			{0x0c, 0x7f, 0x7f, 0x80},
			{0x14, 0x80, 0x80, 0x81}}},
	{"RMW", instr.Ror, "C0 Ror", 0x1f, testD5,
		[]testcase.Branch{brCarry, ifCClr},
		[]aluData{
			{0x00, 0x02, 0x02, 0x01},
			{0x02, 0x00, 0x00, 0x00},
			{0x19, 0x03, 0x03, 0x01},
			{0x1b, 0x01, 0x01, 0x00}}},
	{"RMW", instr.Ror, "C1 Ror", 0x1f, testD5,
		[]testcase.Branch{brCarry, ifCSet},
		[]aluData{
			{0x0c, 0x00, 0x00, 0x80},
			{0x15, 0x01, 0x01, 0x80}}},
	{"<->", instr.Mov, "Mov", 0x00, testD5R5,
		[]testcase.Branch{brD5R5},
		[]aluData{
			{0x00, 0x00, 0x10, 0x10},
			{0x00, 0x10, 0x10, 0x10}}},
	{"<->", instr.Ldi, "Ldi", 0x00, testD4K8,
		nil,
		[]aluData{
			{0x00, 0x00, 0xff, 0xff},
			{0x00, 0xff, 0x00, 0x00}}},
	{"<->", instr.Movw, "Movw", 0x00, testMovw,
		[]testcase.Branch{brRegPair},
		[]aluData{
			{0x00, 0x0000, 0x1234, 0x1234},
			{0x00, 0x4321, 0x4321, 0x4321}}},
}

var addCClrOrIgnored = []aluData{
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
}

var addCSet = []aluData{
	{0x00, 0x00, 0x00, 0x01},
	{0x01, 0x10, 0xf0, 0x01},
	{0x0c, 0x40, 0x40, 0x81},
	{0x14, 0x00, 0x80, 0x81},
	{0x15, 0xc0, 0xc0, 0x81},
	{0x19, 0x80, 0x80, 0x01},
	{0x20, 0x08, 0x08, 0x11},
	{0x21, 0x01, 0xff, 0x01},
	{0x2c, 0x48, 0x48, 0x91},
	{0x34, 0x00, 0x8f, 0x90},
	{0x35, 0xc8, 0xc8, 0x91},
	{0x39, 0x88, 0x88, 0x11},
}

var subCClrOrIgnored = []aluData{
	{0x00, 0x01, 0x00, 0x01},
	{0x01, 0x00, 0x90, 0x70},
	{0x0d, 0x00, 0x80, 0x80},
	{0x14, 0x80, 0x00, 0x80},
	{0x15, 0x00, 0x10, 0xf0},
	{0x18, 0x80, 0x10, 0x70},
	{0x20, 0x10, 0x01, 0x0f},
	{0x21, 0x00, 0x81, 0x7f},
	{0x2d, 0x10, 0x81, 0x8f},
	{0x34, 0x90, 0x01, 0x8f},
	{0x35, 0x00, 0x01, 0xff},
	{0x38, 0x80, 0x01, 0x7f},
	{0x02, 0x00, 0x00, 0x00},
}

var subCSet = []aluData{
	{0x00, 0x02, 0x00, 0x01},
	{0x01, 0x01, 0x90, 0x70},
	{0x0d, 0x01, 0x80, 0x80},
	{0x14, 0x81, 0x00, 0x80},
	{0x15, 0x01, 0x10, 0xf0},
	{0x18, 0x81, 0x10, 0x70},
	{0x20, 0x10, 0x00, 0x0f},
	{0x21, 0x00, 0x80, 0x7f},
	{0x2d, 0x10, 0x80, 0x8f},
	{0x34, 0x90, 0x00, 0x8f},
	{0x35, 0x00, 0x00, 0xff},
	{0x38, 0x80, 0x00, 0x7f},
	{0x02, 0x01, 0x00, 0x00},
	{0x22, 0x10, 0x0f, 0x00},
	{0x23, 0x00, 0xff, 0x00},
	{0x3a, 0x80, 0x7f, 0x00},
}

var andData = []aluData{
	{0x00, 0x01, 0x01, 0x01},
	{0x02, 0xaa, 0x55, 0x00},
	{0x14, 0x80, 0x80, 0x80},
}

var orData = []aluData{
	{0x00, 0x01, 0x03, 0x03},
	{0x02, 0x00, 0x00, 0x00},
	{0x14, 0x80, 0x01, 0x81},
}

// Branches with all flags initially set/all flags initially clear.
func brAllFlags(tree testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	initTc.SregFromByte(0xff)
	tree.Run("SRff", initTc, exp)
	initTc.SregFromByte(0x00)
	tree.Run("SR00", initTc, exp)
}

// Returns a function that branches with a given flag set and cleared.
func brFlag(f Flag) testcase.Branch {
	return func(tree testcase.Tree, init, exp testcase.Testable) {
		initTc := init.(tCpu)
		initTc.SetFlag(f, true)
		tree.Run(fmt.Sprintf("%s1", f), initTc, exp)
		initTc.SetFlag(f, false)
		tree.Run(fmt.Sprintf("%s0", f), initTc, exp)
	}
}

var brCarry = brFlag(FlagC)
var brZero = brFlag(FlagZ)

// Returns a function that continues only if a flag has the given value.
func ifFlag(f Flag, cont bool) testcase.Branch {
	return func(tree testcase.Tree, init, exp testcase.Testable) {
		initTc := init.(tCpu)
		if initTc.GetFlag(f) == cont {
			tree.Run("", init, exp)
		}
	}
}

var ifCClr = ifFlag(FlagC, false)
var ifCSet = ifFlag(FlagC, true)

// Branches with (d, d+1) and (d,d) as dest/source registers
// for two-register instructions.
func brD5R5(tree testcase.Tree, init, exp testcase.Testable) {
	d := rand.Intn(32)
	r := (d + 1) & 0x1f
	initTc := init.(tCpu)
	initTc.ops = instr.Operands{Dst: d, Src: r}
	tree.Run(fmt.Sprintf("r%02d,r%02d", d, r), initTc, exp)
	initTc.ops = instr.Operands{Dst: d, Src: d}
	tree.Run(fmt.Sprintf("r%02d,r%02d", d, d), initTc, exp)
}

// Branches with various dest/source registers for Mul.
func brMul5(tree testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	regs := []struct {
		d, r int
	}{{0, 1}, {0, 0}, {1, 1}, {1, 2}, {2, 3}}
	for _, reg := range regs {
		initTc.ops = instr.Operands{Dst: reg.d, Src: reg.r}
		tree.Run(fmt.Sprintf("r%02d,r%02d", reg.r, reg.d), initTc, exp)
	}
}

// Branches with (16,17) and (16,16) as d,r
// for Muls/Mulsu/Fmul/Fmuls/Fmulsu.
func brMul34(tree testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	initTc.ops = instr.Operands{Dst: 16, Src: 17}
	tree.Run("r16,r17", initTc, exp)
	initTc.ops = instr.Operands{Dst: 16, Src: 16}
	tree.Run("r16,r16", initTc, exp)
}

// Branches with (1:0, 3:2) and (1:0, 1:0) as register pairs for Movw.
func brRegPair(tree testcase.Tree, init, exp testcase.Testable) {
	initTc := init.(tCpu)
	initTc.ops = instr.Operands{Dst: 0, Src: 2}
	tree.Run("r1:r0,r3:r2", initTc, exp)
	initTc.ops = instr.Operands{}
	tree.Run("r1:r0,r1:r0", initTc, exp)
}

type aluTest func(aluCase, instr.Mnemonic)

func testD5R5(ac aluCase, mnem instr.Mnemonic) {
	initCpu := ac.init
	d := initCpu.ops.Dst
	r := initCpu.ops.Src
	if d == r && ac.v1 != ac.v2 {
		return
	}
	initCpu.reg[r] = ac.v2
	initCpu.reg[d] = ac.v1
	expCpu := ac.init
	expCpu.reg[r] = ac.v2
	expCpu.reg[d] = ac.res
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	opFuncs[mnem](&initCpu.Cpu, &initCpu.ops, nil)
	ac.t.Run(fmt.Sprintf("%s [%d]", mnem, ac.n), initCpu, expCpu)
}

func testD5(ac aluCase, mnem instr.Mnemonic) {
	initCpu := ac.init
	d := initCpu.ops.Dst
	initCpu.reg[d] = ac.v1
	expCpu := ac.init
	expCpu.reg[d] = ac.res
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	opFuncs[mnem](&initCpu.Cpu, &initCpu.ops, nil)
	ac.t.Run(fmt.Sprintf("%s [%d]", mnem, ac.n), initCpu, expCpu)
}

func testD4K8(ac aluCase, mnem instr.Mnemonic) {
	expCpu := ac.init
	expCpu.reg[16] = ac.res
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	initCpu := ac.init
	initCpu.reg[16] = ac.v1
	initCpu.ops = instr.Operands{Dst: 16, Src: ac.v2}
	opFuncs[mnem](&initCpu.Cpu, &initCpu.ops, nil)
	ac.t.Run(fmt.Sprintf("%s [%d]", mnem, ac.n), initCpu, expCpu)
}

func testDDK6(ac aluCase, mnem instr.Mnemonic) {
	expCpu := ac.init
	expCpu.reg[24] = ac.res & 0xff
	expCpu.reg[25] = ac.res >> 8
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	initCpu := ac.init
	initCpu.reg[24] = ac.v1 & 0xff
	initCpu.reg[25] = ac.v1 >> 8
	initCpu.ops = instr.Operands{Dst: 24, Src: ac.v2}
	opFuncs[mnem](&initCpu.Cpu, &initCpu.ops, nil)
	ac.t.Run(fmt.Sprintf("%s [%d]", mnem, ac.n), initCpu, expCpu)
}

func testMul(ac aluCase, mnem instr.Mnemonic) {
	initCpu := ac.init
	d := initCpu.ops.Dst
	r := initCpu.ops.Src
	if d == r && ac.v1 != ac.v2 {
		return
	}
	initCpu.reg[r] = ac.v2
	initCpu.reg[d] = ac.v1
	expCpu := ac.init
	expCpu.reg[r] = ac.v2
	expCpu.reg[d] = ac.v1
	expCpu.reg[0] = ac.res & 0xff
	expCpu.reg[1] = ac.res >> 8
	expCpu.setStatus(byte(ac.status), byte(ac.mask))
	opFuncs[mnem](&initCpu.Cpu, &initCpu.ops, nil)
	ac.t.Run(fmt.Sprintf("%s [%d]", mnem, ac.n), initCpu, expCpu)
}

func testMovw(ac aluCase, mnem instr.Mnemonic) {
	initCpu := ac.init
	d := initCpu.ops.Dst
	r := initCpu.ops.Src
	if d == r && ac.v1 != ac.v2 {
		return
	}
	initCpu.reg[r] = ac.v2 & 0xff
	initCpu.reg[r+1] = ac.v2 >> 8
	initCpu.reg[d] = ac.v1 & 0xff
	initCpu.reg[d+1] = ac.v1 >> 8
	expCpu := ac.init
	expCpu.reg[r] = ac.v2 & 0xff
	expCpu.reg[r+1] = ac.v2 >> 8
	expCpu.reg[d] = ac.res & 0xff
	expCpu.reg[d+1] = ac.res >> 8
	movw(&initCpu.Cpu, &initCpu.ops, nil)
	ac.t.Run(fmt.Sprintf("Movw [%d]", ac.n), initCpu, expCpu)
}

func testSubIgnoreCarry(ac aluCase, mnem instr.Mnemonic) {
	testD5R5(ac, instr.Sub)
	testD4K8(ac, instr.Subi)
	ac.res = ac.v1
	testD4K8(ac, instr.Cpi)
	testD5R5(ac, instr.Cp)
}

func testSubRespectCarry(ac aluCase, mnem instr.Mnemonic) {
	if ac.res == 0 && !ac.init.GetFlag(FlagZ) {
		ac.status -= 2
	}
	testD5R5(ac, instr.Sbc)
	testD4K8(ac, instr.Sbci)
	ac.res = ac.v1
	testD5R5(ac, instr.Cpc)
}

func TestALU(t *testing.T) {
	for _, opcase := range alucases {
		run := func(tree testcase.Tree, init, exp testcase.Testable) {
			for n, c := range opcase.data {
				ac := aluCase{tree, init.(tCpu), opcase.fmask, c, n}
				opcase.run(ac, opcase.mnem)
			}
		}
		branches := []testcase.Branch{brAllFlags}
		branches = append(branches, opcase.branches...)
		branches = append(branches, run)
		root := testcase.Tree{opcase.label, t, branches}
		root.Start(tCpu{})
	}
}

// Branches with 0-7 as bit number in Src/Dst.
func brBit(tree testcase.Tree, init, exp testcase.Testable) {
	for bit := 0; bit < 7; bit++ {
		initTc := init.(tCpu)
		initTc.ops = instr.Operands{Src: bit, Dst: bit}
		tree.Run(fmt.Sprintf("b%d", bit), initTc, exp)
	}
}

// Branches with 0-7 as bit number in Off.
func brBitOff(tree testcase.Tree, init, exp testcase.Testable) {
	for bit := 0; bit < 7; bit++ {
		initTc := init.(tCpu)
		initTc.ops = instr.Operands{Off: bit}
		tree.Run(fmt.Sprintf("b%d", bit), initTc, exp)
	}
}

func TestFlag(t *testing.T) {
	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		initCpu := init.(tCpu)
		bit := Flag(initCpu.ops.Src)
		expCpuS := initCpu
		expCpuS.SetFlag(bit, true)
		bset(&initCpu.Cpu, &initCpu.ops, nil)
		tree.Run("Bset", initCpu, expCpuS)
		initCpu = init.(tCpu)
		expCpuC := initCpu
		expCpuC.SetFlag(bit, false)
		bclr(&initCpu.Cpu, &initCpu.ops, nil)
		tree.Run("Blcr", initCpu, expCpuC)
	}
	testcase.NewTree(t, "FLAG", brAllFlags, brBit, run).Start(tCpu{})
}

func TestBst(t *testing.T) {
	// Branches with reg[0] = 1 << bit, ~(1 << bit).
	brSetClr := func(tree testcase.Tree, init, exp testcase.Testable) {
		initCpu := init.(tCpu)
		bit := uint(initCpu.ops.Off)
		initCpu.reg[0] = 1 << bit
		tree.Run("bOn", initCpu, exp)
		initCpu.reg[0] = ^(1 << bit)
		tree.Run("bOff", initCpu, exp)
	}

	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		initCpu := init.(tCpu)
		mask := 1 << uint(initCpu.ops.Off)
		expCpu := initCpu
		expCpu.SetFlag(FlagT, (initCpu.reg[0]&mask) != 0)
		bst(&initCpu.Cpu, &initCpu.ops, nil)
		tree.Run("Bst", initCpu, expCpu)
	}
	testcase.NewTree(t, "XFR",
		brAllFlags, brBitOff, brSetClr, run).Start(tCpu{})
}

func TestBld(t *testing.T) {
	brXfer := brFlag(FlagT)
	run := func(tree testcase.Tree, init, exp testcase.Testable) {
		initCpu := init.(tCpu)
		bit := uint(initCpu.ops.Off)
		expCpu := initCpu
		if initCpu.GetFlag(FlagT) {
			initCpu.reg[0] = 0x00
			expCpu.reg[0] = 1 << bit
		} else {
			initCpu.reg[0] = 0xff
			expCpu.reg[0] = ^(1 << bit) & 0xff
		}
		bld(&initCpu.Cpu, &initCpu.ops, nil)
		tree.Run("Bld", initCpu, expCpu)
	}
	testcase.NewTree(t, "XFR", brAllFlags, brBitOff, brXfer, run).Start(tCpu{})
}
