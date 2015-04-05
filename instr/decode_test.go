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
