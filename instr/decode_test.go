package instr

import "testing"

// Mnemonic Tests
// --------------
// Each value from 0-ffff is decoded and, if not reserved, tested by a
// mnemonic-specific verifier function. Distinguishing between
// reserved and non-reserved values is indirectly tested by checking
// the counts of non-reserved opcodes.

var mnems = map[Mnemonic][]struct {
	minLevel, count int
	opCheck         func(Opcode) bool
}{
	Adc: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0x1c00 && o < 0x2000
		}}},
	Add: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0x0c00 && o < 0x1000
		}}},
	Adiw: {{1, 256,
		func(o Opcode) bool {
			return o >= 0x9600 && o < 0x9700
		}}},
	And: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0x2000 && o < 0x2400
		}}},
	Andi: {{0, 4096,
		func(o Opcode) bool {
			return o >= 0x7000 && o < 0x8000
		}}},
	Asr: {{0, 32,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 && (o&0xf) == 5
		}}},
	Bclr: {{0, 8,
		func(o Opcode) bool {
			return o >= 0x9480 && o < 0x9500 && (o&0xf) == 8
		}}},
	Bld: {{0, 256,
		func(o Opcode) bool {
			return o >= 0xf800 && o < 0xfa00 && (o&0xf) < 8
		}}},
	Brbc: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0xf400 && o < 0xf800
		}}},
	Brbs: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0xf000 && o < 0xf400
		}}},
	Break: {{4, 1, func(o Opcode) bool { return o == 0x9598 }}},
	Bset: {{0, 8,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9480 && (o&0xf) == 8
		}}},
	Bst: {{0, 256,
		func(o Opcode) bool {
			return o >= 0xfa00 && o < 0xfc00 && (o&0xf) < 8
		}}},
	Call: {{2, 64,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 && (o&0xf) > 0xd
		}}},
	Cbi: {{0, 256,
		func(o Opcode) bool {
			return o >= 0x9800 && o < 0x9900
		}}},
	Com: {{0, 32,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 && (o&0xf) == 0
		}}},
	Cp: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0x1400 && o < 0x1800
		}}},
	Cpc: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0x0400 && o < 0x0800
		}}},
	Cpi: {{0, 4096,
		func(o Opcode) bool {
			return o >= 0x3000 && o < 0x4000
		}}},
	Cpse: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0x1000 && o < 0x1400
		}}},
	Dec: {{0, 32,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 && (o&0xf) == 0xa
		}}},
	Des: {{6, 16,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9500 && (o&0xf) == 0xb
		}}},
	Eicall: {{5, 1, func(o Opcode) bool { return o == 0x9519 }}},
	Eijmp:  {{5, 1, func(o Opcode) bool { return o == 0x9419 }}},
	Elpm: {{2, 65,
		func(o Opcode) bool {
			if o == 0x95d8 {
				return true
			}
			return (o >= 0x9000 && o < 0x9200) &&
				((o&0xf) == 6 || (o&0xf) == 7)
		}}},
	Eor: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0x2400 && o < 0x2800
		}}},
	Fmul: {{3, 64,
		func(o Opcode) bool {
			return o >= 0x0300 && o < 0x0380 && (o&0xf) > 7
		}}},
	Fmuls: {{3, 64,
		func(o Opcode) bool {
			return o >= 0x380 && o < 0x400 && (o&0xf) < 8
		}}},
	Fmulsu: {{3, 64,
		func(o Opcode) bool {
			return o >= 0x380 && o < 0x400 && (o&0xf) > 7
		}}},
	Icall: {{1, 1, func(o Opcode) bool { return o == 0x9509 }}},
	Ijmp:  {{1, 1, func(o Opcode) bool { return o == 0x9409 }}},
	In: {{0, 2048,
		func(o Opcode) bool {
			return o >= 0xb000 && o < 0xb800
		}}},
	Inc: {{0, 32,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 && (o&0xf) == 3
		}}},
	Jmp: {{2, 64,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 &&
				((o&0xf) == 0xc || (o&0xf) == 0xd)
		}}},
	Lac: {{6, 32,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9400 && (o&0xf) == 6
		}}},
	Las: {{6, 32,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9400 && (o&0xf) == 5
		}}},
	Lat: {{6, 32,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9400 && (o&0xf) == 7
		}}},
	Ld: {
		{0, 32,
			func(o Opcode) bool {
				return o >= 0x8000 && o < 0x8200 && (o&0xf) == 0
			}},
		{1, 256,
			func(o Opcode) bool {
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
	Ldd: {{1, 4032,
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
	Ldi: {{0, 4096,
		func(o Opcode) bool {
			return o >= 0xe000 && o < 0xf000
		}}},
	Lds: {{1, 32,
		func(o Opcode) bool {
			return o >= 0x9000 && o < 0x9200 && o.Nibble0() == 0
		}}},
	Lpm: {
		{0, 1,
			func(o Opcode) bool { return o == 0x95c8 }},
		{3, 64,
			func(o Opcode) bool {
				return o >= 0x9000 && o < 0x9200 &&
					(o.Nibble0() == 4 || o.Nibble0() == 5)
			}}},
	Lsr: {{0, 32,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 && o.Nibble0() == 6
		}}},
	Mov: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0x2c00 && o < 0x3000
		}}},
	Movw: {{3, 256,
		func(o Opcode) bool {
			return o >= 0x0100 && o < 0x0200
		}}},
	Mul: {{3, 1024,
		func(o Opcode) bool {
			return o >= 0x9c00 && o < 0xa000
		}}},
	Muls: {{3, 256,
		func(o Opcode) bool {
			return o >= 0x0200 && o < 0x0300
		}}},
	Mulsu: {{3, 64,
		func(o Opcode) bool {
			return o >= 0x0300 && o < 0x0380 && o.Nibble0() < 8
		}}},
	Neg: {{0, 32,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 && o.Nibble0() == 1
		}}},
	Nop: {{0, 1, func(o Opcode) bool { return o == 0x0000 }}},
	Or: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0x2800 && o < 0x2c00
		}}},
	Ori: {{0, 4096,
		func(o Opcode) bool {
			return o >= 0x6000 && o < 0x7000
		}}},
	Out: {{0, 2048,
		func(o Opcode) bool {
			return o >= 0xb800 && o < 0xc000
		}}},
	Pop: {{1, 32,
		func(o Opcode) bool {
			return o >= 0x9000 && o < 0x9200 && o.Nibble0() == 0xf
		}}},
	Push: {{1, 32,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9400 && o.Nibble0() == 0xf
		}}},
	Rcall: {{0, 4096,
		func(o Opcode) bool {
			return o >= 0xd000 && o < 0xe000
		}}},
	Ret:  {{0, 1, func(o Opcode) bool { return o == 0x9508 }}},
	Reti: {{0, 1, func(o Opcode) bool { return o == 0x9518 }}},
	Rjmp: {{0, 4096,
		func(o Opcode) bool {
			return o >= 0xc000 && o < 0xd000
		}}},
	Ror: {{0, 32,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 && o.Nibble0() == 7
		}}},
	Sbc: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0x0800 && o < 0x0c00
		}}},
	Sbci: {{0, 4096,
		func(o Opcode) bool {
			return o >= 0x4000 && o < 0x5000
		}}},
	Sbi: {{0, 256,
		func(o Opcode) bool {
			return o >= 0x9a00 && o < 0x9b00
		}}},
	Sbic: {{0, 256,
		func(o Opcode) bool {
			return o >= 0x9900 && o < 0x9a00
		}}},
	Sbis: {{0, 256,
		func(o Opcode) bool {
			return o >= 0x9b00 && o < 0x9c00
		}}},
	Sbiw: {{1, 256,
		func(o Opcode) bool {
			return o >= 0x9700 && o < 0x9800
		}}},
	Sbrc: {{0, 256,
		func(o Opcode) bool {
			return o >= 0xfc00 && o < 0xfe00 && o.Nibble0() < 8
		}}},
	Sbrs: {{0, 256,
		func(o Opcode) bool {
			return o >= 0xfe00 && o.Nibble0() < 8
		}}},
	Sleep: {{0, 1, func(o Opcode) bool { return o == 0x9588 }}},
	Spm: {
		{3, 1, func(o Opcode) bool { return o == 0x95e8 }},
		{6, 1, func(o Opcode) bool { return o == 0x95f8 }}},
	St: {
		{0, 32,
			func(o Opcode) bool {
				return o >= 0x8200 && o < 0x8400 && o.Nibble0() == 0
			}},
		{1, 256,
			func(o Opcode) bool {
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
	Std: {{1, 4032,
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
	Sts: {{1, 32,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9400 && o.Nibble0() == 0
		}}},
	Sub: {{0, 1024,
		func(o Opcode) bool {
			return o >= 0x1800 && o < 0x1c00
		}}},
	Subi: {{0, 4096,
		func(o Opcode) bool {
			return o >= 0x5000 && o < 0x6000
		}}},
	Swap: {{0, 32,
		func(o Opcode) bool {
			return o >= 0x9400 && o < 0x9600 && o.Nibble0() == 2
		}}},
	Wdr: {{0, 1, func(o Opcode) bool { return o == 0x95a8 }}},
	Xch: {{6, 32,
		func(o Opcode) bool {
			return o >= 0x9200 && o < 0x9400 && o.Nibble0() == 4
		}}},
}

func TestDecodeMnem(t *testing.T) {
	var levels = []Set{
		Minimal{}, Classic8K{}, Classic128K{},
		Enhanced8K{}, Enhanced128K{}, Enhanced4M{}, Xmega{}}
	for lvl := 0; lvl < len(levels); lvl++ {
		counts := make([]int, Reserved)

		for op := 0; op < 0x10000; op++ {
			mnem, _ := levels[lvl].DecodeMnem(Opcode(op))
			if mnem == Reserved {
				continue
			}
			opMnemByLevel, ok := mnems[mnem]
			if !ok {
				continue
				// t.Errorf("%s - unexpected %s", levels[lvl], mnem)
			}
			found := false
			for _, opMnem := range opMnemByLevel {
				if opMnem.minLevel <= lvl {
					if opMnem.opCheck(Opcode(op)) {
						found = true
						break
					}
				}
			}
			if !found {
				t.Errorf("%s - %s(%x)", levels[lvl], mnem, op)
			} else {
				counts[mnem] += 1
			}
		}

		for mn, didCount := range counts {
			shouldCount := 0
			for _, wanted := range mnems[Mnemonic(mn)] {
				if wanted.minLevel <= lvl {
					shouldCount += wanted.count
				}
			}
			if shouldCount != didCount {
				t.Errorf("%s - %s count %d not %d",
					levels[lvl], Mnemonic(mn), didCount, shouldCount)
			}
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
	{0xf800, []Mnemonic{In, Out},
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
	{0xfe08, []Mnemonic{Bld, Bst, Sbrc, Sbrs},
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
	{0xfe0f, []Mnemonic{Asr, Com, Dec, Inc, Lac, Las, Lat, Lsr, Neg, Pop, Push, Ror, Swap, Xch},
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
	{0xfc00, []Mnemonic{Adc, Add, And, Cp, Cpc, Cpse, Eor, Mov, Mul, Or, Sbc, Sub},
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
			testAD(t, Elpm, 0x95d8, 0, AddrMode{})
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
	{0xfe00, []Mnemonic{Ld, St},
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
			testAD(t, Lpm, 0x95c8, 0, AddrMode{})
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
			testAD(t, Spm, 0x95f8, 0, AddrMode{0, 0, ZPostInc})
		}},
}

func testAD(t *testing.T, mn Mnemonic, op int, mask uint16, exp AddrMode) {
	s := Minimal{}
	mode := OpModes[mn]

	inst := Instruction{Opcode(op), 0, mn}
	got := s.DecodeAddr(inst)
	if got != exp {
		t.Errorf("%s:%s (0 mask) %v not %v", mn, mode, got, exp)
	}
	if mask == 0 {
		return
	}

	inst = Instruction{Opcode(op) | Opcode(mask), 0, mn}
	got = s.DecodeAddr(inst)
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

var result int

func BenchmarkDecodeMnem(b *testing.B) {
	s := Classic128K{}
	var r Mnemonic
	for n := 0; n < b.N; n++ {
		for o := 0; o < 0x10000; o++ {
			r, _ = s.DecodeMnem(Opcode(o))
		}
	}
	result = int(r)
}

func BenchmarkDecodeAddr(b *testing.B) {
	s := Classic128K{}
	var am AddrMode
	for n := 0; n < b.N; n++ {
		for o := 0; o < 0x10000; o++ {
			op := Opcode(o)
			mnem, _ := s.DecodeMnem(op)
			am = s.DecodeAddr(Instruction{op, 0, mnem})
		}
	}
	result = int(am.A1)
}
