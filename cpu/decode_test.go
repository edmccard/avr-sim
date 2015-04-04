package cpu

import "testing"

func TestDecode(t *testing.T) {
	insts := []struct {
		mask, fixed, wanted uint16
	}{
		{0xfc00, 0x1c00, 1024}, // Adc
		{0xfc00, 0x0c00, 1024}, // Add
		{0xff00, 0x9600, 256},  // Adiw
		{0xfc00, 0x2000, 1024}, // And
		{0xf000, 0x7000, 4096}, // Andi
		{0xfe0f, 0x9405, 32},   // Asr
		{0xff8f, 0x9488, 8},    // Bclr
		{0xfe08, 0xf800, 256},  // Bld
		{0xfc00, 0xf400, 1024}, // Brbc
		{0xfc00, 0xf000, 1024}, // Brbs
		{0xffff, 0x9598, 1},    // Break
		{0xff8f, 0x9408, 8},    // Bset
		{0xfe08, 0xfa00, 256},  // Bst
		{0xfe0e, 0x940e, 64},   // Call
		{0xff00, 0x9800, 256},  // Cbi
		{0xfe0f, 0x9400, 32},   // Com
		{0xfc00, 0x1400, 1024}, // Cp
		{0xfc00, 0x0400, 1024}, // Cpc
		{0xf000, 0x3000, 4096}, // Cpi
		{0xfc00, 0x1000, 1024}, // Cpse
		{0xfe0f, 0x940a, 32},   // Dec
		{0xff0f, 0x940b, 16},   // Des
		{0xffff, 0x9519, 1},    // Eicall
		{0xffff, 0x9419, 1},    // Eijmp
		{0x0000, 0x0000, 65},   // Elpm
		{0xfc00, 0x2400, 1024}, // Eor
		{0xff88, 0x0308, 64},   // Fmul
		{0xff88, 0x0380, 64},   // Fmuls
		{0xff88, 0x0388, 64},   // Fmulsu
		{0xffff, 0x9509, 1},    // Icall
		{0xffff, 0x9409, 1},    // Ijmp
		{0xf800, 0xb000, 2048}, // In
		{0xfe0f, 0x9403, 32},   // Inc
		{0xfe0e, 0x940c, 64},   // Jmp
		{0xfe0f, 0x9206, 32},   // Lac
		{0xfe0f, 0x9205, 32},   // Las
		{0xfe0f, 0x9207, 32},   // Lat
		{0x0000, 0x0000, 288},  // Ld
		{0x0000, 0x0000, 4032}, // Ldd
		{0xf000, 0xe000, 4096}, // Ldi
		{0xfe0f, 0x9000, 32},   // Lds
		{0x0000, 0x0000, 65},   // Lpm
		{0xfe0f, 0x9406, 32},   // Lsr
		{0xfc00, 0x2c00, 1024}, // Mov
		{0xff00, 0x0100, 256},  // Movw
		{0xfc00, 0x9c00, 1024}, // Mul
		{0xff00, 0x0200, 256},  // Muls
		{0xff88, 0x0300, 64},   // Mulsu
		{0xfe0f, 0x9401, 32},   // Neg
		{0xffff, 0x0000, 1},    // Nop
		{0xfc00, 0x2800, 1024}, // Or
		{0xf000, 0x6000, 4096}, // Ori
		{0xf800, 0xb800, 2048}, // Out
		{0xfe0f, 0x900f, 32},   // Pop
		{0xfe0f, 0x920f, 32},   // Push
		{0xf000, 0xd000, 4096}, // Rcall
		{0xffff, 0x9508, 1},    // Ret
		{0xffff, 0x9518, 1},    // Reti
		{0xf000, 0xc000, 4096}, // Rjmp
		{0xfe0f, 0x9407, 32},   // Ror
		{0xfc00, 0x0800, 1024}, // Sbc
		{0xf000, 0x4000, 4096}, // Sbci
		{0xff00, 0x9a00, 256},  // Sbi
		{0xff00, 0x9900, 256},  // Sbic
		{0xff00, 0x9b00, 256},  // Sbis
		{0xff00, 0x9700, 256},  // Sbiw
		{0xfe08, 0xfc00, 256},  // Sbrc
		{0xfe08, 0xfe00, 256},  // Sbrs
		{0xffff, 0x9588, 1},    // Sleep
		{0xffff, 0x95e8, 1},    // Spm
		{0xffff, 0x95f8, 1},    // Spm2
		{0x0000, 0x0000, 288},  // St
		{0x0000, 0x0000, 4032}, // Std
		{0xfe0f, 0x9200, 32},   // Sts
		{0xfc00, 0x1800, 1024}, // Sub
		{0xf000, 0x5000, 4096}, // Subi
		{0xfe0f, 0x9402, 32},   // Swap
		{0xffff, 0x95a8, 1},    // Wdr
		{0xfe0f, 0x9204, 32},   // Xch
		{0x0000, 0x0000, 1554}, // Reserved
	}

	counts := make([]uint16, Illegal)

	for o := 0; o < 0x10000; o++ {
		ok := false
		inst := Decode(Opcode(o))
		counts[inst] += 1
		switch inst {
		case Elpm:
			masked := uint16(o) & 0xfe0f
			ok = o == 0x95d8 || masked == 0x9006 || masked == 0x9007
		case Ld:
			switch uint16(o) & 0xfe0f {
			case 0x8000, 0x8008:
				fallthrough
			case 0x9001, 0x9002, 0x9009, 0x900a, 0x900c, 0x900d, 0x900e:
				ok = true
			}
		case Ldd:
			ok = (uint16(o)&0xd208) == 0x8000 ||
				(uint16(o)&0xd208) == 0x8008 ||
				(uint16(o)&0xf800) == 0xa000
		case Lpm:
			masked := uint16(o) & 0xfe0f
			ok = o == 0x95c8 || masked == 0x9004 || masked == 0x9005
		case St:
			switch uint16(o) & 0xfe0f {
			case 0x8200, 0x8208:
				fallthrough
			case 0x9201, 0x9202, 0x9209, 0x920c, 0x920a, 0x920d, 0x920e:
				ok = true
			}
		case Std:
			ok = (uint16(o)&0xd208) == 0x8200 ||
				(uint16(o)&0xd208) == 0x8208
		default:
			if inst < Illegal {
				mask := insts[inst].mask
				fixed := insts[inst].fixed
				ok = (uint16(o) & mask) == fixed
			}
		}
		if !ok {
			t.Errorf("%s(%x)", inst, o)
		}
	}
	for i := Adc; i < Illegal; i++ {
		wanted := insts[i].wanted
		if wanted != 0 && counts[i] != wanted {
			t.Errorf("%s count was %d, not %d", i, counts[i], wanted)
		}
	}
}