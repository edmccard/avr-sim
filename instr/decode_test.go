package instr

import "testing"

var mnems = map[Mnemonic][]struct {
	mask, fixed     uint16
	count, minLevel int
}{
	Adc:    {{0xfc00, 0x1c00, 1024, 0}},
	Add:    {{0xfc00, 0x0c00, 1024, 0}},
	Adiw:   {{0xff00, 0x9600, 256, 1}},
	And:    {{0xfc00, 0x2000, 1024, 0}},
	Andi:   {{0xf000, 0x7000, 4096, 0}},
	Asr:    {{0xfe0f, 0x9405, 32, 0}},
	Bclr:   {{0xff8f, 0x9488, 8, 0}},
	Bld:    {{0xfe08, 0xf800, 256, 0}},
	Brbc:   {{0xfc00, 0xf400, 1024, 0}},
	Brbs:   {{0xfc00, 0xf000, 1024, 0}},
	Break:  {{0xffff, 0x9598, 1, 4}},
	Bset:   {{0xff8f, 0x9408, 8, 0}},
	Bst:    {{0xfe08, 0xfa00, 256, 0}},
	Call:   {{0xfe0e, 0x940e, 64, 2}},
	Cbi:    {{0xff00, 0x9800, 256, 0}},
	Com:    {{0xfe0f, 0x9400, 32, 0}},
	Cp:     {{0xfc00, 0x1400, 1024, 0}},
	Cpc:    {{0xfc00, 0x0400, 1024, 0}},
	Cpi:    {{0xf000, 0x3000, 4096, 0}},
	Cpse:   {{0xfc00, 0x1000, 1024, 0}},
	Dec:    {{0xfe0f, 0x940a, 32, 0}},
	Des:    {{0xff0f, 0x940b, 16, 6}},
	Eicall: {{0xffff, 0x9519, 1, 5}},
	Eijmp:  {{0xffff, 0x9419, 1, 5}},
	Elpm: {{0xffff, 0x95d8, 1, 2},
		{0xfe0f, 0x9006, 32, 2}, {0xfe0f, 0x9007, 32, 2}},
	Eor:    {{0xfc00, 0x2400, 1024, 0}},
	Fmul:   {{0xff88, 0x0308, 64, 3}},
	Fmuls:  {{0xff88, 0x0380, 64, 3}},
	Fmulsu: {{0xff88, 0x0388, 64, 3}},
	Icall:  {{0xffff, 0x9509, 1, 1}},
	Ijmp:   {{0xffff, 0x9409, 1, 1}},
	In:     {{0xf800, 0xb000, 2048, 0}},
	Inc:    {{0xfe0f, 0x9403, 32, 0}},
	Jmp:    {{0xfe0e, 0x940c, 64, 2}},
	Lac:    {{0xfe0f, 0x9206, 32, 6}},
	Las:    {{0xfe0f, 0x9205, 32, 6}},
	Lat:    {{0xfe0f, 0x9207, 32, 6}},
	Ld: {{0xfe0f, 0x8000, 32, 0},
		{0xfe0f, 0x8008, 32, 1}, {0xfe0f, 0x9001, 32, 1},
		{0xfe0f, 0x9002, 32, 1}, {0xfe0f, 0x9009, 32, 1},
		{0xfe0f, 0x900a, 32, 1}, {0xfe0f, 0x900c, 32, 1},
		{0xfe0d, 0x900d, 32, 1}, {0xfe0f, 0x900e, 32, 1}},
	Ldd: {{0xd208, 0x8000, 1984, 1}, {0xd208, 0x8008, 2048, 1}},
	Ldi: {{0xf000, 0xe000, 4096, 0}},
	Lds: {{0xfe0f, 0x9000, 32, 1}},
	Lpm: {{0xffff, 0x95c8, 1, 0},
		{0xfe0f, 0x9004, 32, 3}, {0xfe0f, 0x9005, 32, 3}},
	Lsr:   {{0xfe0f, 0x9406, 32, 0}},
	Mov:   {{0xfc00, 0x2c00, 1024, 0}},
	Movw:  {{0xff00, 0x0100, 256, 3}},
	Mul:   {{0xfc00, 0x9c00, 1024, 3}},
	Muls:  {{0xff00, 0x0200, 256, 3}},
	Mulsu: {{0xff88, 0x0300, 64, 3}},
	Neg:   {{0xfe0f, 0x9401, 32, 0}},
	Nop:   {{0xffff, 0x0000, 1, 0}},
	Or:    {{0xfc00, 0x2800, 1024, 0}},
	Ori:   {{0xf000, 0x6000, 4096, 0}},
	Out:   {{0xf800, 0xb800, 2048, 0}},
	Pop:   {{0xfe0f, 0x900f, 32, 1}},
	Push:  {{0xfe0f, 0x920f, 32, 1}},
	Rcall: {{0xf000, 0xd000, 4096, 0}},
	Ret:   {{0xffff, 0x9508, 1, 0}},
	Reti:  {{0xffff, 0x9518, 1, 0}},
	Rjmp:  {{0xf000, 0xc000, 4096, 0}},
	Ror:   {{0xfe0f, 0x9407, 32, 0}},
	Sbc:   {{0xfc00, 0x0800, 1024, 0}},
	Sbci:  {{0xf000, 0x4000, 4096, 0}},
	Sbi:   {{0xff00, 0x9a00, 256, 0}},
	Sbic:  {{0xff00, 0x9900, 256, 0}},
	Sbis:  {{0xff00, 0x9b00, 256, 0}},
	Sbiw:  {{0xff00, 0x9700, 256, 1}},
	Sbrc:  {{0xfe08, 0xfc00, 256, 0}},
	Sbrs:  {{0xfe08, 0xfe00, 256, 0}},
	Sleep: {{0xffff, 0x9588, 1, 0}},
	Spm: {{0xffff, 0x95e8, 1, 3},
		{0xffff, 0x95f8, 1, 6}},
	St: {{0xfe0f, 0x8200, 32, 0},
		{0xfe0f, 0x8208, 32, 1}, {0xfe0f, 0x9201, 32, 1},
		{0xfe0f, 0x9202, 32, 1}, {0xfe0d, 0x9209, 32, 1},
		{0xfe0f, 0x920c, 32, 1}, {0xfe0f, 0x920a, 32, 1},
		{0xfe0f, 0x920d, 32, 1}, {0xfe0f, 0x920e, 32, 1}},
	Std:  {{0xd208, 0x8200, 1984, 1}, {0xd208, 0x8208, 2048, 1}},
	Sts:  {{0xfe0f, 0x9200, 32, 1}},
	Sub:  {{0xfc00, 0x1800, 1024, 0}},
	Subi: {{0xf000, 0x5000, 4096, 0}},
	Swap: {{0xfe0f, 0x9402, 32, 0}},
	Wdr:  {{0xffff, 0x95a8, 1, 0}},
	Xch:  {{0xfe0f, 0x9204, 32, 6}},
}

func TestDecodeMnem(t *testing.T) {
	var levels = []Set{
		Minimal{}, Classic8K{}, Classic128K{},
		Enhanced8K{}, Enhanced128K{}, Enhanced4M{}, Xmega{}}
	for lvl := 0; lvl < len(levels); lvl++ {
		counts := make([]int, Reserved)

		for o := 0; o < 0x10000; o++ {
			mnem := levels[lvl].DecodeMnem(Opcode(o))
			if mnem == Reserved {
				continue
			}
			cases, mok := mnems[mnem]
			if !mok {
				// t.Errorf("%s - unexpected %s", levels[lvl], mnem)
			}
			ook := false
			for _, wanted := range cases {
				if wanted.minLevel <= lvl {
					if (uint16(o) & wanted.mask) == wanted.fixed {
						ook = true
					}
				}
			}
			if !ook {
				t.Errorf("%s - %s(%x)", levels[lvl], mnem, o)
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

var result Mnemonic

func BenchmarkDecodeMnem(b *testing.B) {
	s := Classic128K{}
	var r Mnemonic
	for n := 0; n < b.N; n++ {
		for o := 0; o < 0x10000; o++ {
			r = s.DecodeMnem(Opcode(o))
		}
	}
	result = r
}

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
		t.Errorf("%s:%s %v not %v", mn, mode, got, exp)
	}
	if mask == 0 {
		return
	}

	inst = Instruction{Opcode(op) | Opcode(mask), 0, mn}
	got = s.DecodeAddr(inst)
	if got != exp {
		t.Errorf("%s:%s %v not %v", mn, mode, got, exp)
	}
}

func TestDecodeAddr(t *testing.T) {
	for _, amode := range amodes {
		for _, mnem := range amode.mn {
			amode.gen(t, amode.mask, mnem)
		}
	}
}
