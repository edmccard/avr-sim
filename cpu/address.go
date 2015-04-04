package cpu

type Addr uint32

type IndexReg uint32

const (
	None     IndexReg = 0
	Z        IndexReg = 1
	ZPostInc IndexReg = 2
	ZPreDec  IndexReg = 3
	Y        IndexReg = 9
	YPostInc IndexReg = 10
	YPreDec  IndexReg = 11
	X        IndexReg = 13
	XPostInc IndexReg = 14
	XPreDec  IndexReg = 15
)

type AddrMode struct {
	a1, a2 Addr
	ireg   IndexReg
}

type AddrDecoder func(Instruction) AddrMode

// ElpmAddrMode extracts d from _______ddddd____
// and selects the index register.
func ElpmAddrMode(inst Instruction) AddrMode {
	if inst.op1 == 0x95d8 {
		return AddrMode{}
	} else {
		am := Reg5(inst)
		if (inst.op1 & 0x1) == 0 {
			am.ireg = Z
		} else {
			am.ireg = ZPostInc
		}
		return am
	}
}

// LdStAddrMode extracts d from _______ddddd____
// and selects the index register.
func LdStAddrMode(inst Instruction) AddrMode {
	am := Reg5(inst)
	am.ireg = IndexReg((inst.op1 & 0xf)) + 1
	return am
}

// LddStdAddrMode extracts d,q from __q_qq_ddddd_qqq
// and selects the index register.
func LddStdAddrMode(inst Instruction) AddrMode {
	am := Reg5(inst)
	qa := (inst.op1 & 0x2000) >> 8
	qb := (inst.op1 & 0xc00) >> 7
	qc := (inst.op1 & 0x7)
	am.a2 = Addr(qa | qb | qc)
	if (inst.op1 & 0x8) == 0 {
		am.ireg = Z
	} else {
		am.ireg = Y
	}
	return am
}

// LpmAddrMode extracts d from _______ddddd____
// and selects the index register.
func LpmAddrMode(inst Instruction) AddrMode {
	if inst.op1 == 0x95c8 {
		return AddrMode{}
	} else {
		am := Reg5(inst)
		if (inst.op1 & 0x1) == 0 {
			am.ireg = Z
		} else {
			am.ireg = ZPostInc
		}
		return am
	}
}

// Addr16Reg extracts d,k from _______ddddd____ kkkkkkkkkkkkkkkk.
func Addr16Reg5(inst Instruction) AddrMode {
	am := Reg5(inst)
	am.a2 = Addr(inst.op2)
	return am
}

// Addr22 extracts k from _______kkkkk___k kkkkkkkkkkkkkkkk.
func Addr22(inst Instruction) AddrMode {
	ka := Addr((inst.op1 & 0x1f0)) << 13
	kb := Addr((inst.op1 & 0x1)) << 16
	return AddrMode{ka | kb | Addr(inst.op2), 0, 0}
}

// Bit3 extracts b from _________bbb____.
func Bit3(inst Instruction) AddrMode {
	b := (inst.op1 & 0x70) >> 4
	return AddrMode{Addr(b), 0, 0}
}

// Bit3Io5 extracts b,A from ________AAAAAbbb.
func Bit3Io5(inst Instruction) AddrMode {
	b := (inst.op1 & 0x7)
	A := (inst.op1 & 0xf8) >> 3
	return AddrMode{Addr(b), Addr(A), 0}
}

// Bit3Offset7 extracts b,k from ______kkkkkkkbbb.
func Bit3Offset7(inst Instruction) AddrMode {
	b := (inst.op1 & 0x7)
	k := (inst.op1 & 0x3f8) >> 3
	return AddrMode{Addr(b), Addr(k), 0}
}

// Bit3Reg5 extracts b,d from _______ddddd_bbb.
func Bit3Reg5(inst Instruction) AddrMode {
	b := (inst.op1 & 0x7)
	d := (inst.op1 & 0x1f0) >> 4
	return AddrMode{Addr(b), Addr(d), 0}
}

// Imm4 extracts K from ________KKKK____.
func Imm4(inst Instruction) AddrMode {
	K := (inst.op1 & 0xf0) >> 4
	return AddrMode{Addr(K), 0, 0}
}

// Imm6Regpair extracts K,d from ________KKddKKKK.
func Imm6Regpair(inst Instruction) AddrMode {
	Ka := (inst.op1 & 0xc0) >> 2
	Kb := (inst.op1 & 0xf)
	d := (inst.op1 & 0x30) >> 4
	return AddrMode{Addr(Ka | Kb), Addr(d*2 + 24), 0}
}

// Imm8Reg4 extracts K,d from ____KKKKddddKKKK/
func Imm8Reg4(inst Instruction) AddrMode {
	Ka := (inst.op1 & 0xf00) >> 4
	Kb := (inst.op1 & 0xf)
	d := (inst.op1 & 0xf0) >> 4
	return AddrMode{Addr(Ka | Kb), Addr(d), 0}
}

// Io6Reg5 extracts A,d from _____AAdddddAAAA.
func Io6Reg5(inst Instruction) AddrMode {
	Aa := (inst.op1 & 0x600) >> 5
	Ab := (inst.op1 & 0xf)
	r := (inst.op1 & 0x1f0) >> 4
	return AddrMode{Addr(Aa | Ab), Addr(r), 0}
}

// Offset12 extracts k from ____kkkkkkkkkkkk.
func Offset12(inst Instruction) AddrMode {
	return AddrMode{Addr(inst.op1 & 0xfff), 0, 0}
}

// Reg3Reg3 extracts r,d from _________ddd_rrr.
func Reg3Reg3(inst Instruction) AddrMode {
	r := (inst.op1 & 0x7)
	d := (inst.op1 & 0x70) >> 4
	return AddrMode{Addr(r), Addr(d), 0}
}

// Reg4Reg4 extracts r,d from ________ddddrrrr.
func Reg4Reg4(inst Instruction) AddrMode {
	r := (inst.op1 & 0xf)
	d := (inst.op1 & 0xf0) >> 4
	return AddrMode{Addr(r), Addr(d), 0}
}

// Reg5 extracts d from _______ddddd____.
func Reg5(inst Instruction) AddrMode {
	d := (inst.op1 & 0x1f0) >> 4
	return AddrMode{Addr(d), 0, 0}
}

// Reg5Reg5 extracts r,d from ______rdddddrrrr.
func Reg5Reg5(inst Instruction) AddrMode {
	ra := (inst.op1 & 0x200) >> 5
	rb := (inst.op1 & 0xf)
	d := (inst.op1 & 0x1f0) >> 4
	return AddrMode{Addr(ra | rb), Addr(d), 0}
}
