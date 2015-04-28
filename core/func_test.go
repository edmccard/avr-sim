package core

import (
	"fmt"
	it "github.com/edmccard/avr-sim/instr"
	"reflect"
	"testing"
)

func TestABCDE(t *testing.T) {
	testop := func(tag string, mnem it.Mnemonic, branches ...branch) {
		if !testing.Short() {
			branches = append([]branch{statusCases}, branches...)
		}
		casetree{t, branches}.run(tcase{tag: tag, mnem: mnem})
	}

	testop("Ld", it.Ld, loadCase, iregCases(), iregActionCases)
	testop("Ldd", it.Ldd, loadCase, iregCases(), iregDispCases)
	testop("Lds", it.Lds, loadCase, directCase)
	testop("Pop", it.Pop, loadCase, popCases)
	testop("In", it.In, loadCase, portCase)
	testop("Lpm", it.LpmEnhanced, zOnlyCase, lpmCases)
	testop("Elpm", it.ElpmEnhanced, zOnlyCase, elpmCases)
	testop("St", it.St, storeCase, iregCases(), iregActionCases)
	testop("Std", it.Std, storeCase, iregCases(), iregDispCases)
	testop("Sts", it.Sts, storeCase, directCase)
	testop("Push", it.Push, storeCase, pushCases)
	testop("Out", it.Out, storeCase, portCase)
	testop("Lac", it.Lac, zOnlyCase, lacCases)
	testop("Las", it.Las, zOnlyCase, lasCases)
	testop("Lat", it.Lat, zOnlyCase, latCases)
	testop("Xch", it.Xch, zOnlyCase, xchCases)
	testop("Jmp", it.Jmp, jmpCases)
	testop("Call", it.Call, callCases)
	testop("Rjmp", it.Rjmp, rjmpCases)
	testop("Rcall", it.Rcall, rcallCases)
	testop("Ijmp", it.Ijmp, ijmpCases)
	testop("Icall", it.Icall, icallCases)
	testop("Eijmp", it.Eijmp, eijmpCases)
	testop("Eicall", it.Eicall, eicallCases)
	testop("Ret", it.Ret, retCases)
	testop("Reti", it.Reti, retiCases)
	testop("Brbs", it.Brbs, brSetCases, goCases)
	testop("Brbs", it.Brbs, brClrCases, stayCase)
	testop("Brbc", it.Brbc, brClrCases, goCases)
	testop("Brbc", it.Brbc, brSetCases, stayCase)
}

var loadCase = branch{
	{tag: "load",
		init: cdata{mval: 0x42, dstreg: 16},
		exp:  cdata{dstval: 0x42}},
}

var storeCase = branch{
	{tag: "store",
		init: cdata{srcreg: 16, srcval: 0x42},
		exp:  cdata{mval: 0x42}},
}

var directCase = branch{
	{tag: "direct", init: cdata{addr: 0x200}},
}

var portCase = branch{
	{tag: "port", init: cdata{port: 0x10, addr: 0x30}},
}

var pushCases = branch{
	{tag: "stack",
		init: cdata{addr: 0x1000, sp: 0x1000},
		exp:  cdata{sp: 0xfff}},
	{tag: "stack wrap",
		init: cdata{addr: 0x0, sp: 0x0},
		exp:  cdata{sp: 0xffff}},
}

var popCases = branch{
	{tag: "stack",
		init: cdata{addr: 0x1000, sp: 0xfff},
		exp:  cdata{sp: 0x1000}},
	{tag: "stack wrap",
		init: cdata{addr: 0x0, sp: 0xffff},
		exp:  cdata{sp: 0x0}},
}

func iregCases() branch {
	var cases = branch{
		{tag: "X", init: cdata{ireg: it.X, action: it.NoAction}},
		{tag: "Y", init: cdata{ireg: it.Y, action: it.NoAction}},
		{tag: "Z", init: cdata{ireg: it.Z, action: it.NoAction}},
	}
	if testing.Short() {
		return branch{cases[2]}
	}
	return cases
}

var iregActionCases = branch{
	{tag: "no ramp",
		init: cdata{action: it.NoAction, ptr: 0xff, addr: 0xff},
		exp:  cdata{ptr: 0xff}},
	{tag: "no ramp postinc",
		init: cdata{action: it.PostInc, ptr: 0xff, addr: 0xff},
		exp:  cdata{ptr: 0x100}},
	{tag: "no ramp postinc wrap",
		init: cdata{action: it.PostInc, ptr: 0xffff, addr: 0xffff},
		exp:  cdata{ptr: 0x0}},
	{tag: "no ramp predec",
		init: cdata{action: it.PreDec, ptr: 0x100, addr: 0xff},
		exp:  cdata{ptr: 0xff}},
	{tag: "no ramp predec wrap",
		init: cdata{action: it.PreDec, ptr: 0, addr: 0xffff},
		exp:  cdata{ptr: 0xffff}},
	{tag: "ramp",
		init: cdata{action: it.NoAction, ramp: 0x1, ptr: 0xff, addr: 0x100ff},
		exp:  cdata{ramp: 0x1, ptr: 0xff}},
	{tag: "ramp postinc",
		init: cdata{action: it.PostInc, ramp: 0x1, ptr: 0xff, addr: 0x100ff},
		exp:  cdata{ptr: 0x100}},
	{tag: "ramp postinc rollover",
		init: cdata{action: it.PostInc, ramp: 0x1, ptr: 0xffff, addr: 0x1ffff},
		exp:  cdata{ramp: 0x2, ptr: 0x0}},
	{tag: "ramp predec",
		init: cdata{action: it.PreDec, ramp: 0x1, ptr: 0x100, addr: 0x100ff},
		exp:  cdata{ptr: 0xff}},
	{tag: "ramp predec rollover",
		init: cdata{action: it.PreDec, ramp: 0x2, ptr: 0x0, addr: 0x1ffff},
		exp:  cdata{ramp: 0x1, ptr: 0xffff}},
}

var iregDispCases = branch{
	{tag: "disp no ramp",
		init: cdata{disp: 0x1, ptr: 0xff, addr: 0x100}},
	{tag: "disp no ramp wrap",
		init: cdata{disp: 0x2, ptr: 0xffff, addr: 0x1}},
	{tag: "disp ramp",
		init: cdata{ramp: 0x1, disp: 0x1, ptr: 0xfffe, addr: 0x1ffff}},
	{tag: "disp ramp rollover",
		init: cdata{ramp: 0x1, disp: 0x1, ptr: 0xffff, addr: 0x20000}},
}

var zOnlyCase = branch{
	{tag: "",
		init: cdata{ireg: it.Z, action: it.NoAction, dstreg: 16}},
}

var xchCases = branch{
	{tag: "no ramp",
		init: cdata{ptr: 0x100, dstval: 0xaa, mval: 0x55},
		exp:  cdata{dstval: 0x55, mval: 0xaa}},
	{tag: "ramp",
		init: cdata{ramp: 0x1, ptr: 0x100, addr: 0x10100,
			dstval: 0x55, mval: 0xaa},
		exp: cdata{dstval: 0xaa, mval: 0x55}},
}

var lacCases = branch{
	{tag: "no ramp zero",
		init: cdata{ptr: 0x200, dstval: 0x00, mval: 0x77},
		exp:  cdata{mval: 0x77, dstval: 0x77}},
	{tag: "no ramp non-zero",
		init: cdata{ptr: 0x200, dstval: 0x22, mval: 0x77},
		exp:  cdata{dstval: 0x77, mval: 0x55}},
	{tag: "ramp zero",
		init: cdata{ramp: 0x2, ptr: 0x200, addr: 0x20200,
			dstval: 0x00, mval: 0x77},
		exp: cdata{mval: 0x77, dstval: 0x77}},
}

var lasCases = branch{
	{tag: "no ramp non-zero",
		init: cdata{ptr: 0x100, dstval: 0xaa, mval: 0x55},
		exp:  cdata{dstval: 0x55, mval: 0xff}},
	{tag: "no ramp zero",
		init: cdata{ptr: 0x100, dstval: 0x00, mval: 0xaa},
		exp:  cdata{dstval: 0xaa, mval: 0xaa}},
	{tag: "ramp non-zero",
		init: cdata{ramp: 0x2, ptr: 0x100, addr: 0x20100,
			dstval: 0xaa, mval: 0x55},
		exp: cdata{dstval: 0x55, mval: 0xff}},
}

var latCases = branch{
	{tag: "no ramp non-zero",
		init: cdata{ptr: 0x10, dstval: 0x55, mval: 0x66},
		exp:  cdata{dstval: 0x66, mval: 0x33}},
	{tag: "no ramp zero",
		init: cdata{ptr: 0x20, dstval: 0x00, mval: 0x55},
		exp:  cdata{dstval: 0x55, mval: 0x55}},
	{tag: "ramp non-zero",
		init: cdata{ramp: 0x2, ptr: 0x10, addr: 0x20010,
			dstval: 0x55, mval: 0x66},
		exp: cdata{dstval: 0x66, mval: 0x33}},
}

var lpmCases = branch{
	{tag: "ramp",
		init: cdata{ramp: 0x2, ptr: 0x1001, addr: 0x800, pval: 0x1234},
		exp:  cdata{dstval: 0x12}},
	{tag: "ramp rollover",
		init: cdata{action: it.PostInc, ramp: 0x2, ptr: 0xffff,
			addr: 0x17fff, pval: 0x1234},
		exp: cdata{dstval: 0x12, ptr: 0x0}},
	{tag: "no ramp lo byte",
		init: cdata{ptr: 0x1000, addr: 0x800, pval: 0x1234},
		exp:  cdata{dstval: 0x34}},
	{tag: "no ramp hi byte",
		init: cdata{ptr: 0x1001, addr: 0x800, pval: 0x1234},
		exp:  cdata{dstval: 0x12}},
	{tag: "no ramp wrap",
		init: cdata{action: it.PostInc, ptr: 0xffff,
			addr: 0x7fff, pval: 0x1234},
		exp: cdata{ptr: 0x0, dstval: 0x12}},
}

var elpmCases = append(branch{
	{tag: "ramp",
		init: cdata{ramp: 0x2, ptr: 0x1000, addr: 0x10800, pval: 0x1234},
		exp:  cdata{dstval: 0x34}},
	{tag: "ramp rollover",
		init: cdata{action: it.PostInc, ramp: 0x2, ptr: 0xffff,
			addr: 0x17fff, pval: 0x1234},
		exp: cdata{dstval: 0x12, ramp: 0x3, ptr: 0x0}}},
	lpmCases[2:]...)

var jmpCases = branch{
	{tag: "no ramp",
		init: cdata{pc: 0x0, disp: 0x2000},
		exp:  cdata{pc: 0x2000}},
	{tag: "no ramp >16-bit",
		init: cdata{pc: 0x0, disp: 0x1ffff},
		exp:  cdata{pc: 0xffff}},
	{tag: "ramp >16-bit",
		init: cdata{ramp: 0x1, pc: 0x0, disp: 0x10000},
		exp:  cdata{pc: 0x10000}},
}

var callCases = branch{
	{tag: "no ramp >16-bit",
		init: cdata{pc: 0x1234, disp: 0x10000, sp: 0x3fff},
		exp: cdata{pc: 0x0, sp: 0x3ffd,
			savepc: stack{0x3fff: 0x34, 0x3ffe: 0x12}}},
	{tag: "no ramp",
		init: cdata{pc: 0x1234, disp: 0x5678, sp: 0x3fff},
		exp: cdata{pc: 0x5678, sp: 0x3ffd,
			savepc: stack{0x3fff: 0x34, 0x3ffe: 0x12}}},
	{tag: "ramp",
		init: cdata{ramp: 0x1, pc: 0x123456, disp: 0x20000, sp: 0x3fff},
		exp: cdata{pc: 0x20000, sp: 0x3ffc,
			savepc: stack{0x3fff: 0x56, 0x3ffe: 0x34, 0x3ffd: 0x12}}},
}

var rjmpCases = branch{
	{tag: "no ramp",
		init: cdata{pc: 0x1000, disp: 0x7ff},
		exp:  cdata{pc: 0x17ff}},
	{tag: "no ramp forward wrap",
		init: cdata{pc: 0xffff, disp: 0x1},
		exp:  cdata{pc: 0x0}},
	{tag: "no ramp backwrap wrap",
		init: cdata{pc: 0x0, disp: -1},
		exp:  cdata{pc: 0xffff}},
	{tag: "ramp forward",
		init: cdata{ramp: 0x1, pc: 0x1ffff, disp: 0x1},
		exp:  cdata{pc: 0x20000}},
	{tag: "ramp backward",
		init: cdata{ramp: 0x1, pc: 0x0, disp: -1},
		exp:  cdata{pc: 0x3fffff}},
}

var rcallCases = branch{
	{tag: "no ramp",
		init: cdata{pc: 0x1234, disp: 0x200, sp: 0x1ff},
		exp: cdata{pc: 0x1434, sp: 0x1fd,
			savepc: stack{0x1ff: 0x34, 0x1fe: 0x12}}},
	{tag: "no ramp forward wrap",
		init: cdata{pc: 0xffff, disp: 0x1, sp: 0x1ff},
		exp: cdata{pc: 0x0, sp: 0x1fd,
			savepc: stack{0x1ff: 0xff, 0x1fe: 0xff}}},
	{tag: "no ramp backward wrap",
		init: cdata{pc: 0x0, disp: -1, sp: 0x1ff},
		exp: cdata{pc: 0xffff, sp: 0x1fd,
			savepc: stack{0x1ff: 0x0, 0x1fe: 0x0}}},
	{tag: "ramp forward",
		init: cdata{ramp: 0x1, pc: 0x1ffff, disp: 0x1, sp: 0x1ff},
		exp: cdata{pc: 0x20000, sp: 0x1fc,
			savepc: stack{0x1ff: 0xff, 0x1fe: 0xff, 0x1fd: 0x01}}},
	{tag: "ramp backward",
		init: cdata{ramp: 0x1, pc: 0x0, disp: -1, sp: 0x1ff},
		exp: cdata{pc: 0x3fffff, sp: 0x1fc,
			savepc: stack{0x1ff: 0x0, 0x1fe: 0x0, 0x1fd: 0x0}}},
}

var ijmpCases = branch{
	{tag: "no ramp",
		init: cdata{pc: 0x0, ptr: 0x2000},
		exp:  cdata{pc: 0x2000}},
	{tag: "ramp",
		init: cdata{ramp: 0x1, pc: 0x10000, ptr: 0x2000},
		exp:  cdata{pc: 0x2000}},
}

var icallCases = branch{
	{tag: "no ramp",
		init: cdata{pc: 0x1234, ptr: 0x2000, sp: 0x2ff},
		exp: cdata{pc: 0x2000, sp: 0x2fd,
			savepc: stack{0x2ff: 0x34, 0x2fe: 0x12}}},
	{tag: "ramp",
		init: cdata{ramp: 0x1, pc: 0x123456, ptr: 0x2000, sp: 0x2ff},
		exp: cdata{pc: 0x2000, sp: 0x2fc,
			savepc: stack{0x2ff: 0x56, 0x2fe: 0x34, 0x2fd: 0x12}}},
}

var eijmpCases = branch{
	{tag: "no ramp",
		init: cdata{pc: 0x0, ptr: 0x2000},
		exp:  cdata{pc: 0x2000}},
	{tag: "ramp",
		init: cdata{ramp: 0x1, pc: 0x0, ptr: 0x2000},
		exp:  cdata{pc: 0x12000}},
}

var eicallCases = branch{
	{tag: "no ramp",
		init: cdata{pc: 0x1234, ptr: 0x2000, sp: 0x3ff},
		exp: cdata{pc: 0x2000, sp: 0x3fd,
			savepc: stack{0x3ff: 0x34, 0x3fe: 0x12}}},
	{tag: "ramp",
		init: cdata{ramp: 0x1, pc: 0x123456, ptr: 0x2000, sp: 0x3ff},
		exp: cdata{pc: 0x12000, sp: 0x3fc,
			savepc: stack{0x3ff: 0x56, 0x3fe: 0x34, 0x3fd: 0x12}}},
}

var retCases = branch{
	{tag: "no ramp",
		init: cdata{pc: 0x1000, sp: 0x3fd,
			savepc: stack{0x3ff: 0x34, 0x3fe: 0x12}},
		exp: cdata{pc: 0x1234, sp: 0x3ff}},
	{tag: "ramp",
		init: cdata{ramp: 0x1, pc: 0x1000, sp: 0x3fc,
			savepc: stack{0x3ff: 0x56, 0x3fe: 0x34, 0x3fd: 0x12}},
		exp: cdata{pc: 0x123456, sp: 0x3ff}},
}

var retiCases = branch{
	{tag: "no ramp",
		init: cdata{pc: 0x1000, sp: 0x3fd,
			savepc: stack{0x3ff: 0x34, 0x3fe: 0x12}},
		exp: cdata{pc: 0x1234, sp: 0x3ff, set: flags{FlagI}}},
	{tag: "ramp",
		init: cdata{ramp: 0x1, pc: 0x1000, sp: 0x3fc,
			savepc: stack{0x3ff: 0x56, 0x3fe: 0x34, 0x3fd: 0x12}},
		exp: cdata{pc: 0x123456, sp: 0x3ff, set: flags{FlagI}}},
}

var brSetCases = branch{
	{tag: "C set", init: cdata{bit: 0, set: flags{FlagC}}},
	{tag: "Z set", init: cdata{bit: 1, set: flags{FlagZ}}},
	{tag: "N set", init: cdata{bit: 2, set: flags{FlagN}}},
	{tag: "V set", init: cdata{bit: 3, set: flags{FlagV}}},
	{tag: "S set", init: cdata{bit: 4, set: flags{FlagS}}},
	{tag: "H set", init: cdata{bit: 5, set: flags{FlagH}}},
	{tag: "T set", init: cdata{bit: 6, set: flags{FlagT}}},
	{tag: "I set", init: cdata{bit: 7, set: flags{FlagI}}},
}

var brClrCases = branch{
	{tag: "C clr", init: cdata{bit: 0, clr: flags{FlagC}}},
	{tag: "Z clr", init: cdata{bit: 1, clr: flags{FlagZ}}},
	{tag: "N clr", init: cdata{bit: 2, clr: flags{FlagN}}},
	{tag: "V clr", init: cdata{bit: 3, clr: flags{FlagV}}},
	{tag: "S clr", init: cdata{bit: 4, clr: flags{FlagS}}},
	{tag: "H clr", init: cdata{bit: 5, clr: flags{FlagH}}},
	{tag: "T clr", init: cdata{bit: 6, clr: flags{FlagT}}},
	{tag: "I clr", init: cdata{bit: 7, clr: flags{FlagI}}},
}

var goCases = branch{
	{tag: "no ramp forward",
		init: cdata{pc: 0x1000, disp: 32},
		exp:  cdata{pc: 0x1020}},
	{tag: "no ramp backward",
		init: cdata{pc: 0x1000, disp: -32},
		exp:  cdata{pc: 0xfe0}},
	{tag: "no ramp forward wrap",
		init: cdata{pc: 0xffff, disp: 1},
		exp:  cdata{pc: 0x0}},
	{tag: "no ramp backward wrap",
		init: cdata{pc: 0x0, disp: -1},
		exp:  cdata{pc: 0xffff}},
	{tag: "ramp forward",
		init: cdata{ramp: 0x1, pc: 0xffff, disp: 1},
		exp:  cdata{pc: 0x10000}},
	{tag: "ramp forward wrap",
		init: cdata{ramp: 0x1, pc: 0x3fffff, disp: 1},
		exp:  cdata{pc: 0x0}},
	{tag: "ramp backward wrap",
		init: cdata{ramp: 0x1, pc: 0x0, disp: -1},
		exp:  cdata{pc: 0x3fffff}},
}

var stayCase = branch{
	{tag: "stay", init: cdata{pc: 0x1000, disp: 0x20}, exp: cdata{pc: 0x1000}},
}

var statusCases = branch{
	{tag: "Srff", init: cdata{status: 0xff}},
	{tag: "Sr00", init: cdata{status: 0x00}},
}

type key int

const (
	action key = iota
	addr
	bit
	clr
	disp
	dstreg
	dstval
	ireg
	iregop
	mval
	pc
	port
	ptr
	pval
	ramp
	savepc
	set
	sp
	srcreg
	srcval
	status
)

type cdata map[key]interface{}

type stack map[int]int

type flash map[int]int

type flags []Flag

func (this cdata) merge(that cdata) cdata {
	merged := make(map[key]interface{})
	for k, v := range this {
		merged[k] = v
	}
	for k, v := range that {
		if prev, ok := merged[k]; ok && (k == set || k == clr) {
			merged[k] = append(prev.(flags), v.(flags)...)
		} else {
			merged[k] = v
		}
	}
	return merged
}

func (data cdata) musthave(k key) interface{} {
	if val, ok := data[k]; ok {
		return val
	}
	panic(fmt.Sprintf("missing case data %s", k))
}

type branch []tcase

type branches []branch

type tcase struct {
	tag  string
	init cdata
	exp  cdata
	mnem it.Mnemonic
}

func (this tcase) merge(that tcase) tcase {
	return tcase{
		tag:  this.tag + " " + that.tag,
		init: this.init.merge(that.init),
		exp:  this.exp.merge(that.exp),
		mnem: this.mnem,
	}
}

func (tc tcase) run(t *testing.T) {
	init := newsystem()
	init.apply(tc.init)
	// TODO: is clone needed when applying merged data?
	exp := init.clone()
	// exp := newsystem()
	exp.apply(tc.init.merge(tc.exp))
	opFuncs[tc.mnem](&init.cpu, &init.cpu.ops, &init.mem)
	if !init.equals(&exp) {
		t.Error(tc.tag)
		fmt.Println("INIT:", init)
		fmt.Println("EXP: ", exp)
	}
}

type tmem struct {
	data map[Addr]byte
	prog map[Addr]uint16
}

func newtmem() tmem {
	return tmem{data: make(map[Addr]byte), prog: make(map[Addr]uint16)}
}

func (m tmem) clone() tmem {
	dup := newtmem()
	for k, v := range m.data {
		dup.data[k] = v
	}
	for k, v := range m.prog {
		dup.prog[k] = v
	}
	return dup
}

func (this tmem) equals(that tmem) bool {
	return reflect.DeepEqual(this.data, that.data) &&
		reflect.DeepEqual(this.prog, that.prog)
}

func (m *tmem) ReadData(addr Addr) byte {
	if val, ok := m.data[addr]; ok {
		return val
	}
	return 0x9e
}

func (m *tmem) WriteData(addr Addr, val byte) {
	m.data[addr] = val
}

func (m *tmem) ReadProgram(addr Addr) uint16 {
	if val, ok := m.prog[addr]; ok {
		return val
	}
	return 0
}

func (m *tmem) LoadProgram(addr Addr) byte {
	if val, ok := m.prog[addr>>1]; ok {
		return byte(val >> ((uint(addr) & 0x1) * 8))
	}
	return 0
}

type system struct {
	cpu Cpu
	mem tmem
}

func newsystem() system {
	return system{cpu: Cpu{}, mem: newtmem()}
}

func (s *system) apply(data cdata) {
	// do "global" things like status here first
	if statval, ok := data[status]; ok {
		s.cpu.SregFromByte(byte(statval.(int)))
	}
	if sflags, ok := data[set]; ok {
		for _, f := range sflags.(flags) {
			s.cpu.flags[f] = true
		}
	}
	if cflags, ok := data[clr]; ok {
		for _, f := range cflags.(flags) {
			s.cpu.flags[f] = false
		}
	}
	if spval, ok := data[sp]; ok {
		s.cpu.sp = spval.(int)
	}
	if dsp, ok := data[disp]; ok {
		s.cpu.ops.Off = dsp.(int)
	}

	if val, ok := data[pc]; ok {
		s.applyoffset(val.(int), data)
		return
	}
	if val, ok := data[ireg]; ok {
		s.applyindirect(val.(it.IndexReg), data)
		return
	}
	if val, ok := data[addr]; ok {
		s.applydirect(Addr(val.(int)), data)
		return
	}
}

func (s *system) applyoffset(offset int, data cdata) {
	if eindval, ok := data[ramp]; ok {
		// eixxx
		s.setramp(Eind, eindval.(int))
	}
	if val, ok := data[ptr]; ok {
		// ixxx/eixxx
		s.setindex(it.Z.Reg(), val.(int))
	}
	if stk, ok := data[savepc]; ok {
		// calls/returns
		for a, v := range stk.(stack) {
			s.mem.WriteData(Addr(a), byte(v))
		}
	}
	if bitnum, ok := data[bit]; ok {
		// branches
		s.cpu.ops.Src = bitnum.(int)
	}
	s.cpu.pc = offset
}

func (s *system) applydirect(maddr Addr, data cdata) {
	ioport, hasioport := data[port]
	if !hasioport {
		s.cpu.ops.Off = int(maddr)
	}
	if memval, ok := data[mval]; ok {
		s.mem.WriteData(maddr, byte(memval.(int)))
	}
	if reg, ok := data[srcreg]; ok {
		s.cpu.reg[reg.(int)] = data.musthave(srcval).(int)
		s.cpu.ops.Src = reg.(int)
		s.cpu.ops.Dst = reg.(int)
		if hasioport {
			s.cpu.ops.Dst = ioport.(int)
		}
	} else {
		reg := data.musthave(dstreg).(int)
		s.cpu.ops.Src = reg
		s.cpu.ops.Dst = reg
		if hasioport {
			s.cpu.ops.Src = ioport.(int)
		}
		if dval, ok := data[dstval]; ok {
			s.cpu.reg[reg] = dval.(int)
		}
	}
}

func (s *system) applyindirect(base it.IndexReg, data cdata) {
	action := data.musthave(action).(it.IndexAction)
	indexreg := base.WithAction(action)
	if reg, ok := data[srcreg]; ok {
		s.cpu.reg[reg.(int)] = data.musthave(srcval).(int)
		s.cpu.ops.Src = reg.(int)
		s.cpu.ops.Dst = int(indexreg)
	} else {
		reg := data.musthave(dstreg).(int)
		s.cpu.ops.Src = int(indexreg)
		s.cpu.ops.Dst = reg
		if dval, ok := data[dstval]; ok {
			s.cpu.reg[reg] = dval.(int)
		}
	}
	iptr := data.musthave(ptr).(int)
	s.setindex(indexreg.Reg(), iptr)
	if rmp, ok := data[ramp]; ok {
		s.setramp(Ramp(base), rmp.(int))
	}
	if memval, ok := data[mval]; ok {
		if maddr, ok := data[addr]; ok {
			s.mem.WriteData(Addr(maddr.(int)), byte(memval.(int)))
		} else {
			s.mem.WriteData(Addr(iptr), byte(memval.(int)))
		}
	} else if progval, ok := data[pval]; ok {
		maddr := data.musthave(addr).(int)
		s.mem.prog[Addr(maddr)] = uint16(progval.(int))
	}
}

func (s *system) setramp(base Ramp, val int) {
	s.cpu.setRmask(base, 0x3f)
	s.cpu.SetRamp(base, byte(val))
}

func (s *system) setindex(reg, val int) {
	s.cpu.reg[reg] = val & 0xff
	s.cpu.reg[reg+1] = val >> 8
}

func (s *system) clone() system {
	return system{cpu: s.cpu, mem: s.mem.clone()}
}

func (this *system) equals(that *system) bool {
	return this.cpu == that.cpu && this.mem.equals(that.mem)
}

type casetree struct {
	t        *testing.T
	branches branches
}

func (tree casetree) run(builder tcase) {
	if len(tree.branches) == 0 {
		builder.run(tree.t)
	} else {
		next := casetree{t: tree.t, branches: tree.branches[1:]}
		for _, tc := range tree.branches[0] {
			next.run(builder.merge(tc))
		}
	}
}

var skipdecoder = it.NewDecoder(it.NewSetXmega())
