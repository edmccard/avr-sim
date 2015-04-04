package cpu

import "testing"

func checkInst(t *testing.T, ad AddrDecoder, nm string, op Opcode, mask Opcode,
	exp AddrMode) {

	inst := Instruction{0, op, 0}
	result := ad(inst)
	if result != exp {
		t.Errorf("%s(%x) %v, not %v", nm, op, result, exp)
	}
	inst = Instruction{0, op | mask, 0}
	result = ad(inst)
	if result != exp {
		t.Errorf("%s(%x) %v, not %v", nm, op, result, exp)
	}
}

func TestElpmAddrMode(t *testing.T) {
	cases := []struct {
		suff int
		ireg IndexReg
	}{
		{0x6, Z}, {0x7, ZPostInc},
	}
	checkInst(t, ElpmAddrMode, "ElpmAddrMode", 0x95d8, 0, AddrMode{})
	for c := 0; c < 2; c++ {
		for d := 0; d < 32; d++ {
			e := AddrMode{Addr(d), 0, cases[c].ireg}
			op := (d << 4) | cases[c].suff
			checkInst(t, ElpmAddrMode, "ElpmAddrMode", Opcode(op), 0xfe00, e)
		}
	}
}

func TestLdStAddrMode(t *testing.T) {
	cases := []struct {
		suff int
		ireg IndexReg
	}{
		{0xc, X}, {0xd, XPostInc}, {0xe, XPreDec},
		{0x8, Y}, {0x9, YPostInc}, {0xa, YPreDec},
		{0x0, Z}, {0x1, ZPostInc}, {0x2, ZPreDec},
	}
	for c := 0; c < 9; c++ {
		for d := 0; d < 32; d++ {
			e := AddrMode{Addr(d), 0, cases[c].ireg}
			op := (d << 4) | cases[c].suff
			checkInst(t, LdStAddrMode, "LdStAddrMode", Opcode(op), 0xfe00, e)
		}
	}
}

func TestLddStdAddrMode(t *testing.T) {
	cases := []struct {
		suff int
		ireg IndexReg
	}{
		{0x8, Y}, {0x0, Z},
	}
	for c := 0; c < 2; c++ {
		for d := 0; d < 32; d++ {
			for q := 0; q < 64; q++ {
				e := AddrMode{Addr(d), Addr(q), cases[c].ireg}
				op := ((q & 0x20) << 8) | ((q & 0x18) << 7) | (q & 0x7)
				op |= (d << 4) | cases[c].suff
				checkInst(t, LddStdAddrMode, "LddStdAddrMode",
					Opcode(op), 0xd200, e)
			}
		}
	}
}

func TestLpmAddrMode(t *testing.T) {
	checkInst(t, LpmAddrMode, "LpmAddrMode", 0x95c8, 0, AddrMode{})
	cases := []struct {
		suff int
		ireg IndexReg
	}{
		{0x4, Z}, {0x5, ZPostInc},
	}
	for c := 0; c < 2; c++ {
		for d := 0; d < 32; d++ {
			e := AddrMode{Addr(d), 0, cases[c].ireg}
			op := (d << 4) | cases[c].suff
			checkInst(t, LpmAddrMode, "LpmAddrMode", Opcode(op), 0xfe00, e)
		}
	}
}

func TestAddr16Reg5(t *testing.T) {
	for d := 0; d < 32; d++ {
		e := AddrMode{Addr(d), 0, 0}
		op := d << 4
		checkInst(t, Addr16Reg5, "Addr16Reg5", Opcode(op), 0xfe0f, e)
	}
}

func TestAddr22(t *testing.T) {
	for k := 0; k < 64; k++ {
		e := AddrMode{Addr(k << 16), 0, 0}
		op := ((k & 0x3e) << 3) | (k & 0x1)
		checkInst(t, Addr22, "Addr22", Opcode(op), 0xfe0e, e)
	}
}

func TestBit3(t *testing.T) {
	for b := 0; b < 8; b++ {
		e := AddrMode{Addr(b), 0, 0}
		op := b << 4
		checkInst(t, Bit3, "Bit3", Opcode(op), 0xff8f, e)
	}
}

func TestBit3Io5(t *testing.T) {
	for A := 0; A < 32; A++ {
		for b := 0; b < 8; b++ {
			e := AddrMode{Addr(b), Addr(A), 0}
			op := (A << 3) | b
			checkInst(t, Bit3Io5, "Bit3Io5", Opcode(op), 0xff00, e)
		}
	}
}

func TestBit3Offset7(t *testing.T) {
	for k := 0; k < 128; k++ {
		for b := 0; b < 8; b++ {
			e := AddrMode{Addr(b), Addr(k), 0}
			op := (k << 3) | b
			checkInst(t, Bit3Offset7, "Bit3Offset7", Opcode(op), 0xfc00, e)
		}
	}
}

func TestBit3Reg5(t *testing.T) {
	for d := 0; d < 32; d++ {
		for b := 0; b < 8; b++ {
			e := AddrMode{Addr(b), Addr(d), 0}
			op := (d << 4) | b
			checkInst(t, Bit3Reg5, "Bit3Reg5", Opcode(op), 0xfe08, e)
		}
	}
}

func TestImm4(t *testing.T) {
	for K := 0; K < 16; K++ {
		e := AddrMode{Addr(K), 0, 0}
		op := K << 4
		checkInst(t, Imm4, "Imm4", Opcode(op), 0xff0f, e)
	}
}

func TestImm6Regpair(t *testing.T) {
	for K := 0; K < 64; K++ {
		for d := 0; d < 4; d++ {
			e := AddrMode{Addr(K), Addr(d*2 + 24), 0}
			op := ((K & 0x30) << 2) | (d << 4) | (K & 0xf)
			checkInst(t, Imm6Regpair, "Imm6Regpair", Opcode(op), 0xff00, e)
		}
	}
}

func TestImm8Reg4(t *testing.T) {
	for K := 0; K < 256; K++ {
		for d := 0; d < 16; d++ {
			e := AddrMode{Addr(K), Addr(d), 0}
			op := ((K & 0xf0) << 4) | (d << 4) | (K & 0xf)
			checkInst(t, Imm8Reg4, "Imm8Reg4", Opcode(op), 0xf000, e)
		}
	}
}

func TestIo6Reg5(t *testing.T) {
	for A := 0; A < 64; A++ {
		for d := 0; d < 32; d++ {
			e := AddrMode{Addr(A), Addr(d), 0}
			op := ((A & 0x30) << 5) | (d << 4) | (A & 0xf)
			checkInst(t, Io6Reg5, "Io6Reg5", Opcode(op), 0xf800, e)
		}
	}
}

func TestOffset12(t *testing.T) {
	for k := 0; k < 4096; k++ {
		e := AddrMode{Addr(k), 0, 0}
		checkInst(t, Offset12, "Offset12", Opcode(k), 0xf000, e)
	}
}

func TestReg3Reg3(t *testing.T) {
	for r := 0; r < 8; r++ {
		for d := 0; d < 8; d++ {
			e := AddrMode{Addr(r), Addr(d), 0}
			op := (d << 4) | r
			checkInst(t, Reg3Reg3, "Reg3Reg3", Opcode(op), 0xff88, e)
		}
	}
}

func TestReg4Reg4(t *testing.T) {
	for r := 0; r < 16; r++ {
		for d := 0; d < 16; d++ {
			e := AddrMode{Addr(r), Addr(d), 0}
			op := (d << 4) | r
			checkInst(t, Reg4Reg4, "Reg4Reg4", Opcode(op), 0xff00, e)
		}
	}
}

func TestReg5(t *testing.T) {
	for d := 0; d < 32; d++ {
		e := AddrMode{Addr(d), 0, 0}
		op := d << 4
		checkInst(t, Reg5, "Reg5", Opcode(op), 0xfe0f, e)
	}
}

func TestReg5Reg5(t *testing.T) {
	for r := 0; r < 32; r++ {
		for d := 0; d < 32; d++ {
			e := AddrMode{Addr(r), Addr(d), 0}
			op := ((r & 0x10) << 5) | (d << 4) | (r & 0xf)
			checkInst(t, Reg5Reg5, "Reg5Reg5", Opcode(op), 0xfc00, e)
		}
	}
}
