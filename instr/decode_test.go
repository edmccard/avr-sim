package instr

import "testing"

// Mnemonic Tests
// --------------
// Each value from 0-ffff is decoded and, if not reserved, tested by a
// mnemonic-specific verifier function. Distinguishing between
// reserved and non-reserved values is indirectly tested by checking
// the counts of non-reserved opcodes.

var mnems = map[Mnemonic][]struct {
	count   int
	opCheck func(Opcode) bool
}{
	Adc: {{768,
		func(o Opcode) bool {
			return o >= 0x1c00 && o < 0x1f00
		}}},
	AdcReduced: {{256,
		func(o Opcode) bool {
			return o >= 0x1f00 && o < 0x2000
		}}},
	Add: {{768,
		func(o Opcode) bool {
			return o >= 0x0c00 && o < 0x0f00
		}}},
	AddReduced: {{256,
		func(o Opcode) bool {
			return o >= 0x0f00 && o < 0x1000
		}}},
	Adiw: {{256,
		func(o Opcode) bool {
			return o >= 0x9600 && o < 0x9700
		}}},
	And: {{768,
		func(o Opcode) bool {
			return o >= 0x2000 && o < 0x2300
		}}},
	AndReduced: {{256,
		func(o Opcode) bool {
			return o >= 0x2300 && o < 0x2400
		}}},
	Andi: {{4096,
		func(o Opcode) bool {
			return o >= 0x7000 && o < 0x8000
		}}},
	Asr: {{16,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9500 && (o&0xf) == 5
		}}},
	AsrReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x9500 && o < 0x9600 && (o&0xf) == 5
		}}},
	Bclr: {{8,
		func(o Opcode) bool {
			return o >= 0x9480 && o < 0x9500 && (o&0xf) == 8
		}}},
	Bld: {{128,
		func(o Opcode) bool {
			return o >= 0xf800 && o < 0xf900 && (o&0xf) < 8
		}}},
	BldReduced: {{128,
		func(o Opcode) bool {
			return o >= 0xf900 && o < 0xfa00 && (o&0xf) < 8
		}}},
	Brbc: {{1024,
		func(o Opcode) bool {
			return o >= 0xf400 && o < 0xf800
		}}},
	Brbs: {{1024,
		func(o Opcode) bool {
			return o >= 0xf000 && o < 0xf400
		}}},
	Break: {{1, func(o Opcode) bool { return o == 0x9598 }}},
	Bset: {{8,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9480 && (o&0xf) == 8
		}}},
	Bst: {{128,
		func(o Opcode) bool {
			return o >= 0xfa00 && o < 0xfb00 && (o&0xf) < 8
		}}},
	BstReduced: {{128,
		func(o Opcode) bool {
			return o >= 0xfb00 && o < 0xfc00 && (o&0xf) < 8
		}}},
	Call: {{64,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 && (o&0xf) > 0xd
		}}},
	Cbi: {{256,
		func(o Opcode) bool {
			return o >= 0x9800 && o < 0x9900
		}}},
	Com: {{16,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9500 && (o&0xf) == 0
		}}},
	ComReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x9500 && o < 0x9600 && (o&0xf) == 0
		}}},
	Cp: {{768,
		func(o Opcode) bool {
			return o >= 0x1400 && o < 0x1700
		}}},
	CpReduced: {{256,
		func(o Opcode) bool {
			return o >= 0x1700 && o < 0x1800
		}}},
	Cpc: {{768,
		func(o Opcode) bool {
			return o >= 0x0400 && o < 0x0700
		}}},
	CpcReduced: {{256,
		func(o Opcode) bool {
			return o >= 0x0700 && o < 0x0800
		}}},
	Cpi: {{4096,
		func(o Opcode) bool {
			return o >= 0x3000 && o < 0x4000
		}}},
	Cpse: {{768,
		func(o Opcode) bool {
			return o >= 0x1000 && o < 0x1300
		}}},
	CpseReduced: {{256,
		func(o Opcode) bool {
			return o >= 0x1300 && o < 0x1400
		}}},
	Dec: {{16,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9500 && (o&0xf) == 0xa
		}}},
	DecReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x9500 && o < 0x9600 && (o&0xf) == 0xa
		}}},
	Des: {{16,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9500 && (o&0xf) == 0xb
		}}},
	Eicall: {{1, func(o Opcode) bool { return o == 0x9519 }}},
	Eijmp:  {{1, func(o Opcode) bool { return o == 0x9419 }}},
	Elpm:   {{1, func(o Opcode) bool { return o == 0x95d8 }}},
	ElpmEnhanced: {{64,
		func(o Opcode) bool {
			if o == 0x95d8 {
				return true
			}
			return (o >= 0x9000 && o < 0x9200) &&
				((o&0xf) == 6 || (o&0xf) == 7)
		}}},
	Eor: {{768,
		func(o Opcode) bool {
			return o >= 0x2400 && o < 0x2700
		}}},
	EorReduced: {{256,
		func(o Opcode) bool {
			return o >= 0x2700 && o < 0x2800
		}}},
	Fmul: {{64,
		func(o Opcode) bool {
			return o >= 0x0300 && o < 0x0380 && (o&0xf) > 7
		}}},
	Fmuls: {{64,
		func(o Opcode) bool {
			return o >= 0x380 && o < 0x400 && (o&0xf) < 8
		}}},
	Fmulsu: {{64,
		func(o Opcode) bool {
			return o >= 0x380 && o < 0x400 && (o&0xf) > 7
		}}},
	Icall: {{1, func(o Opcode) bool { return o == 0x9509 }}},
	Ijmp:  {{1, func(o Opcode) bool { return o == 0x9409 }}},
	In: {{1024,
		func(o Opcode) bool {
			return o >= 0xb000 && o < 0xb800 &&
				(o.Nibble2()&0x1) == 0
		}}},
	InReduced: {{1024,
		func(o Opcode) bool {
			return o >= 0xb000 && o < 0xb800 &&
				(o.Nibble2()&0x1) != 0
		}}},
	Inc: {{16,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9500 && (o&0xf) == 3
		}}},
	IncReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x9500 && o < 0x9600 && (o&0xf) == 3
		}}},
	Jmp: {{64,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 &&
				((o&0xf) == 0xc || (o&0xf) == 0xd)
		}}},
	Lac: {{32,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9400 && (o&0xf) == 6
		}}},
	Las: {{32,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9400 && (o&0xf) == 5
		}}},
	Lat: {{32,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9400 && (o&0xf) == 7
		}}},
	LdMinimal: {{16,
		func(o Opcode) bool {
			return o >= 0x8000 && o < 0x8100 && (o&0xf) == 0
		}}},
	LdMinimalReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x8100 && o < 0x8200 && (o&0xf) == 0
		}}},
	LdClassic: {{128,
		func(o Opcode) bool {
			if ((o & 0x1f0) >> 4) >= 16 {
				return false
			}
			on0 := o.Nibble0()
			switch {
			case o >= 0x8000 && o < 0x8200:
				return on0 == 8
			case o >= 0x9000 && o < 0x9200:
				return on0 == 1 || on0 == 2 || on0 == 9 || on0 == 0xa ||
					on0 == 0xc || on0 == 0xd || on0 == 0xe
			default:
				return false
			}
		}}},
	LdClassicReduced: {{128,
		func(o Opcode) bool {
			if ((o & 0x1f0) >> 4) < 16 {
				return false
			}
			on0 := o.Nibble0()
			switch {
			case o >= 0x8000 && o < 0x8200:
				return on0 == 8
			case o >= 0x9000 && o < 0x9200:
				return on0 == 1 || on0 == 2 || on0 == 9 || on0 == 0xa ||
					on0 == 0xc || on0 == 0xd || on0 == 0xe
			default:
				return false
			}
		}}},
	Ldd: {{4032,
		func(o Opcode) bool {
			if o < 0x8000 || o >= 0xb000 {
				return false
			}
			if ((o.Nibble2() >> 1) & 0x1) == 1 {
				return false
			}
			if o < 0x8200 && (o.Nibble0() == 0 || o.Nibble0() == 8) {
				return false
			}
			return true
		}}},
	Ldi: {{4096,
		func(o Opcode) bool {
			return o >= 0xe000 && o < 0xf000
		}}},
	Lds: {{32,
		func(o Opcode) bool {
			return o >= 0x9000 && o < 0x9200 && o.Nibble0() == 0
		}}},
	Lpm: {{1, func(o Opcode) bool { return o == 0x95c8 }}},
	LpmEnhanced: {{64,
		func(o Opcode) bool {
			return o >= 0x9000 && o < 0x9200 &&
				(o.Nibble0() == 4 || o.Nibble0() == 5)
		}}},
	Lsr: {{16,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9500 && o.Nibble0() == 6
		}}},
	LsrReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x9500 && o < 0x9600 && o.Nibble0() == 6
		}}},
	Mov: {{768,
		func(o Opcode) bool {
			return o >= 0x2c00 && o < 0x2f00
		}}},
	MovReduced: {{256,
		func(o Opcode) bool {
			return o >= 0x2f00 && o < 0x3000
		}}},
	Movw: {{256,
		func(o Opcode) bool {
			return o >= 0x0100 && o < 0x0200
		}}},
	Mul: {{1024,
		func(o Opcode) bool {
			return o >= 0x9c00 && o < 0xa000
		}}},
	Muls: {{256,
		func(o Opcode) bool {
			return o >= 0x0200 && o < 0x0300
		}}},
	Mulsu: {{64,
		func(o Opcode) bool {
			return o >= 0x0300 && o < 0x0380 && o.Nibble0() < 8
		}}},
	Neg: {{16,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9500 && o.Nibble0() == 1
		}}},
	NegReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x9500 && o < 0x9600 && o.Nibble0() == 1
		}}},
	Nop: {{1, func(o Opcode) bool { return o == 0x0000 }}},
	Or: {{768,
		func(o Opcode) bool {
			return o >= 0x2800 && o < 0x2b00
		}}},
	OrReduced: {{256,
		func(o Opcode) bool {
			return o >= 0x2b00 && o < 0x2c00
		}}},
	Ori: {{4096,
		func(o Opcode) bool {
			return o >= 0x6000 && o < 0x7000
		}}},
	Out: {{1024,
		func(o Opcode) bool {
			return o >= 0xb800 && o < 0xc000 &&
				(o.Nibble2()&0x1) == 0
		}}},
	OutReduced: {{1024,
		func(o Opcode) bool {
			return o >= 0xb800 && o < 0xc000 &&
				(o.Nibble2()&0x1) != 0
		}}},
	Pop: {{16,
		func(o Opcode) bool {
			return o >= 0x9000 && o < 0x9100 && o.Nibble0() == 0xf
		}}},
	PopReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x9100 && o < 0x9200 && o.Nibble0() == 0xf
		}}},
	Push: {{16,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9300 && o.Nibble0() == 0xf
		}}},
	PushReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x9300 && o < 0x9400 && o.Nibble0() == 0xf
		}}},
	Rcall: {{4096,
		func(o Opcode) bool {
			return o >= 0xd000 && o < 0xe000
		}}},
	Ret:  {{1, func(o Opcode) bool { return o == 0x9508 }}},
	Reti: {{1, func(o Opcode) bool { return o == 0x9518 }}},
	Rjmp: {{4096,
		func(o Opcode) bool {
			return o >= 0xc000 && o < 0xd000
		}}},
	Ror: {{16,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9500 && o.Nibble0() == 7
		}}},
	RorReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x9500 && o < 0x9600 && o.Nibble0() == 7
		}}},
	Sbc: {{768,
		func(o Opcode) bool {
			return o >= 0x0800 && o < 0x0b00
		}}},
	SbcReduced: {{256,
		func(o Opcode) bool {
			return o >= 0x0b00 && o < 0x0c00
		}}},
	Sbci: {{4096,
		func(o Opcode) bool {
			return o >= 0x4000 && o < 0x5000
		}}},
	Sbi: {{256,
		func(o Opcode) bool {
			return o >= 0x9a00 && o < 0x9b00
		}}},
	Sbic: {{256,
		func(o Opcode) bool {
			return o >= 0x9900 && o < 0x9a00
		}}},
	Sbis: {{256,
		func(o Opcode) bool {
			return o >= 0x9b00 && o < 0x9c00
		}}},
	Sbiw: {{256,
		func(o Opcode) bool {
			return o >= 0x9700 && o < 0x9800
		}}},
	Sbrc: {{128,
		func(o Opcode) bool {
			return o >= 0xfc00 && o < 0xfd00 && o.Nibble0() < 8
		}}},
	SbrcReduced: {{128,
		func(o Opcode) bool {
			return o >= 0xfd00 && o < 0xfe00 && o.Nibble0() < 8
		}}},
	Sbrs: {{128,
		func(o Opcode) bool {
			return o >= 0xfe00 && o < 0xff00 && o.Nibble0() < 8
		}}},
	SbrsReduced: {{128,
		func(o Opcode) bool {
			return o >= 0xff00 && o.Nibble0() < 8
		}}},
	Sleep:    {{1, func(o Opcode) bool { return o == 0x9588 }}},
	Spm:      {{1, func(o Opcode) bool { return o == 0x95e8 }}},
	SpmXmega: {{1, func(o Opcode) bool { return o == 0x95f8 }}},
	StMinimal: {{16,
		func(o Opcode) bool {
			return o >= 0x8200 && o < 0x8300 && o.Nibble0() == 0
		}}},
	StMinimalReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x8300 && o < 0x8400 && o.Nibble0() == 0
		}}},
	StClassic: {{128,
		func(o Opcode) bool {
			if ((o & 0x1f0) >> 4) >= 16 {
				return false
			}
			on0 := o.Nibble0()
			switch {
			case o >= 0x8200 && o < 0x8400:
				return on0 == 8
			case o >= 0x9200 && o < 0x9400:
				return on0 == 1 || on0 == 2 || on0 == 9 || on0 == 0xa ||
					on0 == 0xc || on0 == 0xd || on0 == 0xe
			default:
				return false
			}
		}}},
	StClassicReduced: {{128,
		func(o Opcode) bool {
			if ((o & 0x1f0) >> 4) < 16 {
				return false
			}
			on0 := o.Nibble0()
			switch {
			case o >= 0x8200 && o < 0x8400:
				return on0 == 8
			case o >= 0x9200 && o < 0x9400:
				return on0 == 1 || on0 == 2 || on0 == 9 || on0 == 0xa ||
					on0 == 0xc || on0 == 0xd || on0 == 0xe
			default:
				return false
			}
		}}},
	Std: {{4032,
		func(o Opcode) bool {
			if o < 0x8000 || o >= 0xb000 {
				return false
			}
			if ((o.Nibble2() >> 1) & 0x1) == 0 {
				return false
			}
			if o < 0x8400 && (o.Nibble0() == 0 || o.Nibble0() == 8) {
				return false
			}
			return true
		}}},
	Sts: {{32,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9400 && o.Nibble0() == 0
		}}},
	Sub: {{768,
		func(o Opcode) bool {
			return o >= 0x1800 && o < 0x1b00
		}}},
	SubReduced: {{256,
		func(o Opcode) bool {
			return o >= 0x1b00 && o < 0x1c00
		}}},
	Subi: {{4096,
		func(o Opcode) bool {
			return o >= 0x5000 && o < 0x6000
		}}},
	Swap: {{16,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9500 && o.Nibble0() == 2
		}}},
	SwapReduced: {{16,
		func(o Opcode) bool {
			return o >= 0x9500 && o < 0x9600 && o.Nibble0() == 2
		}}},
	Wdr: {{1, func(o Opcode) bool { return o == 0x95a8 }}},
	Xch: {{32,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9400 && o.Nibble0() == 4
		}}},
}

func TestDecodeMnem(t *testing.T) {
	counts := make([]int, NumMnems)
	decoder := NewDecoder(setXmega)

	for op := 0; op < 0x10000; op++ {
		mnem, _ := decoder.DecodeMnem(Opcode(op))
		if mnem == Reserved {
			continue
		}
		opMnemByLevel, ok := mnems[mnem]
		if !ok {
			t.Errorf("unexpected %s", mnem)
		}
		found := false
		for _, opMnem := range opMnemByLevel {
			if opMnem.opCheck(Opcode(op)) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s(%x)", mnem, op)
		} else {
			counts[mnem] += 1
		}
	}

	for mn, didCount := range counts {
		shouldCount := 0
		for _, wanted := range mnems[Mnemonic(mn)] {
			shouldCount += wanted.count
		}
		if shouldCount != didCount {
			t.Errorf("%s count %d not %d",
				Mnemonic(mn), didCount, shouldCount)
		}
	}
}

// Address Mode Tests
// ------------------
// There are three ways to incorrectly decode opcode parameters:
// 1. Using a non-parameter bit as a parameter bit
// 2. Ignoring a parameter bit
// 3. When there are multiple parameters, using a bit from one to decode
//    the other.
//
// The first case is tested by synthesizing two "opcodes" for each
// address mode, differing only in the non-parameter bits, which are
// set to zero in the first and one in the second; if non-parameter
// bits are incorrectly used then the results of decoding the two
// "opcodes" will differ.
//
// The second and third cases are tested by generating every possible
// combination of parameters for each address mode.

var amodes = []struct {
	mask uint16
	mn   []Mnemonic
	gen  func(*testing.T, uint16, Mnemonic)
}{
	{0xff00, []Mnemonic{Cbi, Sbi, Sbic, Sbis},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for A := 0; A < 32; A++ {
				for b := 0; b < 8; b++ {
					op := (A << 3) | b
					am := AddrMode{Addr(A), Addr(b), NoIndex}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xf800, []Mnemonic{In, InReduced, Out, OutReduced},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for A := 0; A < 64; A++ {
				for d := 0; d < 32; d++ {
					op := ((A & 0x30) << 5) | (d << 4) | (A & 0xf)
					am := AddrMode{Addr(A), Addr(d), NoIndex}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xff8f, []Mnemonic{Bclr, Bset},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for b := 0; b < 8; b++ {
				op := b << 4
				am := AddrMode{Addr(b), 0, NoIndex}
				testAD(t, mn, op, mask, am)
			}
		}},
	{0xfc00, []Mnemonic{Brbc, Brbs},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for b := 0; b < 8; b++ {
				for k := -64; k < 64; k++ {
					kc := uint(k) & 0x7f
					op := int(kc<<3) | b
					am := AddrMode{Addr(b), Addr(k), NoIndex}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xfe08, []Mnemonic{
		Bld, Bst, Sbrc, Sbrs,
		BldReduced, BstReduced, SbrcReduced, SbrsReduced,
	},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for b := 0; b < 8; b++ {
				for d := 0; d < 32; d++ {
					op := (d << 4) | b
					am := AddrMode{Addr(b), Addr(d), NoIndex}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xff88, []Mnemonic{Fmul, Fmuls, Fmulsu, Mulsu},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for d := 16; d < 24; d++ {
				for r := 16; r < 24; r++ {
					op := ((d & 0xf) << 4) | (r & 0xf)
					am := AddrMode{Addr(d), Addr(r), NoIndex}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xf000, []Mnemonic{Andi, Cpi, Ldi, Ori, Sbci, Subi},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for d := 16; d < 32; d++ {
				for K := 0; K < 256; K++ {
					op := ((K & 0xf0) << 4) | ((d & 0xf) << 4) | (K & 0xf)
					am := AddrMode{Addr(d), Addr(K), NoIndex}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xff00, []Mnemonic{Muls},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for d := 16; d < 32; d++ {
				for r := 16; r < 32; r++ {
					op := ((d & 0xf) << 4) | (r & 0xf)
					am := AddrMode{Addr(d), Addr(r), NoIndex}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xfe0f, []Mnemonic{
		Asr, Com, Dec, Inc, Lac, Las, Lat, Lsr, Neg, Pop, Push, Ror, Swap, Xch,
		AsrReduced, ComReduced, DecReduced, IncReduced, LsrReduced,
		NegReduced, PopReduced, PushReduced, RorReduced, SwapReduced,
	},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for d := 0; d < 32; d++ {
				op := d << 4
				am := AddrMode{Addr(d), 0, NoIndex}
				testAD(t, mn, op, mask, am)
			}
		}},
	{0xfe0f, []Mnemonic{Lds, Sts},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for d := 0; d < 32; d++ {
				op := d << 4
				am := AddrMode{Addr(d), 0, NoIndex}
				testAD(t, mn, op, mask, am)
			}
		}},
	{0xfc00, []Mnemonic{
		Adc, Add, And, Cp, Cpc, Cpse, Eor, Mov, Mul, Or, Sbc, Sub,
		AdcReduced, AndReduced, CpReduced, CpcReduced, CpseReduced,
		EorReduced, MovReduced, OrReduced, SbcReduced, SubReduced,
	},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for d := 0; d < 32; d++ {
				for r := 0; r < 32; r++ {
					op := ((r & 0x10) << 5) | (d << 4) | (r & 0xf)
					am := AddrMode{Addr(d), Addr(r), NoIndex}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xff00, []Mnemonic{Movw},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for r := 0; r < 32; r += 2 {
				for d := 0; d < 32; d += 2 {
					op := ((d >> 1) << 4) | (r >> 1)
					am := AddrMode{Addr(d), Addr(r), NoIndex}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xff00, []Mnemonic{Adiw, Sbiw},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for d := 24; d < 32; d += 2 {
				for K := 0; K < 64; K++ {
					op := ((K & 0x30) << 2) | (((d - 24) >> 1) << 4) | (K & 0xf)
					am := AddrMode{Addr(d), Addr(K), NoIndex}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xfe00, []Mnemonic{Elpm},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			testAD(t, Elpm, 0x95d8, 0, AddrMode{0, 0, Z})
			for d := 0; d < 32; d++ {
				op := (d << 4) | 0x6
				am := AddrMode{Addr(d), 0, Z}
				testAD(t, mn, op, mask, am)
			}
			for d := 0; d < 32; d++ {
				op := (d << 4) | 0x7
				am := AddrMode{Addr(d), 0, ZPostInc}
				testAD(t, mn, op, mask, am)
			}
		}},
	{0xff0f, []Mnemonic{Des},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for K := 0; K < 16; K++ {
				op := K << 4
				am := AddrMode{Addr(K), 0, NoIndex}
				testAD(t, mn, op, mask, am)
			}
		}},
	{0xf000, []Mnemonic{Rcall, Rjmp},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for k := -2048; k < 2048; k++ {
				op := int(uint(k) & 0xfff)
				am := AddrMode{Addr(k), 0, NoIndex}
				testAD(t, mn, op, mask, am)
			}
		}},
	{0xfe0e, []Mnemonic{Call, Jmp},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for k := 0; k < 64; k++ {
				op := ((k & 0x3e) << 3) | (k & 0x1)
				am := AddrMode{Addr(k << 16), 0, 0}
				testAD(t, mn, op, mask, am)
			}
		}},
	{0xd200, []Mnemonic{Ldd, Std},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for d := 0; d < 32; d++ {
				for q := 0; q < 64; q++ {
					op := ((q & 0x20) << 8) | ((q & 0x18) << 7) | (q & 0x7)
					op |= (d << 4) | 0x8
					am := AddrMode{Addr(d), Addr(q), Y}
					testAD(t, mn, op, mask, am)
				}
			}
			for d := 0; d < 32; d++ {
				for q := 0; q < 64; q++ {
					op := ((q & 0x20) << 8) | ((q & 0x18) << 7) | (q & 0x7)
					op |= (d << 4)
					am := AddrMode{Addr(d), Addr(q), Z}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xfe00, []Mnemonic{
		LdMinimal, LdMinimalReduced, LdClassic, LdClassicReduced,
		StMinimal, StMinimalReduced, StClassic, StClassicReduced,
	},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for suff, ireg := range ldstireg {
				for d := 0; d < 32; d++ {
					op := (d << 4) | suff
					am := AddrMode{Addr(d), 0, ireg}
					testAD(t, mn, op, mask, am)
				}
			}
		}},
	{0xfe00, []Mnemonic{Lpm},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			testAD(t, Lpm, 0x95c8, 0, AddrMode{0, 0, Z})
		}},
	{0xfe00, []Mnemonic{LpmEnhanced},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			for d := 0; d < 32; d++ {
				op := (d << 4) | 0x4
				am := AddrMode{Addr(d), 0, Z}
				testAD(t, mn, op, mask, am)
			}
			for d := 0; d < 32; d++ {
				op := (d << 4) | 0x5
				am := AddrMode{Addr(d), 0, ZPostInc}
				testAD(t, mn, op, mask, am)
			}
		}},
	{0x0000, []Mnemonic{Reserved},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			testAD(t, Break, 0x9598, 0, AddrMode{})
			testAD(t, Eicall, 0x9519, 0, AddrMode{})
			testAD(t, Eijmp, 0x9419, 0, AddrMode{})
			testAD(t, Icall, 0x9509, 0, AddrMode{})
			testAD(t, Ijmp, 0x9409, 0, AddrMode{})
			testAD(t, Nop, 0x0000, 0, AddrMode{})
			testAD(t, Ret, 0x9508, 0, AddrMode{})
			testAD(t, Reti, 0x9518, 0, AddrMode{})
			testAD(t, Sleep, 0x9588, 0, AddrMode{})
			testAD(t, Wdr, 0x95a8, 0, AddrMode{})
		}},
	{0x0000, []Mnemonic{Spm},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			testAD(t, Spm, 0x95e8, 0, AddrMode{})
		}},
	{0x0000, []Mnemonic{SpmXmega},
		func(t *testing.T, mask uint16, mn Mnemonic) {
			testAD(t, SpmXmega, 0x95f8, 0, AddrMode{0, 0, ZPostInc})
		}},
}

func testAD(t *testing.T, mn Mnemonic, op int, mask uint16, exp AddrMode) {
	d := NewDecoder(setXmega)
	mode := OpModes[mn]

	inst := Instruction{Opcode(op), 0}
	got := d.DecodeAddr(mn, inst)
	if got != exp {
		t.Errorf("%s:%s (0 mask) %v not %v", mn, mode, got, exp)
	}
	if mask == 0 {
		return
	}

	inst = Instruction{Opcode(op) | Opcode(mask), 0}
	got = d.DecodeAddr(mn, inst)
	if got != exp {
		t.Errorf("%s:%s (1 mask) %v not %v", mn, mode, got, exp)
	}
}

func TestDecodeAddr(t *testing.T) {
	for _, amode := range amodes {
		for _, mnem := range amode.mn {
			amode.gen(t, amode.mask, mnem)
		}
	}
}

var setXmega = NewSetXmega()
