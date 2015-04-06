package instr

type IndexReg int

const (
	NoIndex IndexReg = iota
	X
	XPostInc
	XPreDec
	Y
	YPostInc
	YPreDec
	Z
	ZPostInc
	ZPreDec
)

type Mode int

const (
	ModeA5B3 Mode = iota
	ModeA6D5
	ModeB3
	ModeB3K7
	ModeB3D5
	ModeD3R3
	ModeD4K8
	ModeD4R4
	ModeD5
	ModeD5K16
	ModeD5R5
	ModeDDDDRRRR
	ModeDDK6
	ModeElpm
	ModeK4
	ModeK12
	ModeK22
	ModeLddStd
	ModeLdSt
	ModeLpm
	ModeNone
	ModeSpm
)

//go:generate stringer -type=Mode

type Addr int

type AddrMode struct {
	A1, A2 Addr
	Ireg   IndexReg
}

type AddrDecoder func(Instruction) AddrMode

// DecodeA5B3 returns AddrMode{A, b, NoIndex} extracted from
// ________AAAAAbbb.
func DecodeA5B3(inst Instruction) AddrMode {
	A := (inst.Op1 & 0xf8) >> 3
	b := (inst.Op1 & 0x7)
	return AddrMode{Addr(A), Addr(b), NoIndex}
}

// DecodeA6D5 returns AddrMode{A, d, NoIndex} extracted from
// _____AAdddddAAAA.
func DecodeA6D5(inst Instruction) AddrMode {
	Aa := (inst.Op1 & 0x600) >> 5
	Ab := (inst.Op1 & 0xf)
	d := (inst.Op1 & 0x1f0) >> 4
	return AddrMode{Addr(Aa | Ab), Addr(d), NoIndex}
}

// DecodeB3 returns AddrMode{b, 0, NoIndex} extracted from
// _________bbb____.
func DecodeB3(inst Instruction) AddrMode {
	b := (inst.Op1 & 0x70) >> 4
	return AddrMode{Addr(b), 0, NoIndex}
}

// DecodeB3K7 returns AddrMode{b, k, NoIndex} (where -64<=k<=63)
// extracted from ______kkkkkkkbbb.
func DecodeB3K7(inst Instruction) AddrMode {
	b := (inst.Op1 & 0x7)
	k := Addr((inst.Op1 & 0x3f8) >> 3)
	if (k & 0x40) != 0 {
		k -= 0x80
	}
	return AddrMode{Addr(b), k, NoIndex}
}

// DecodeB3D5 returns AddrMode{b, d, NoIndex} extracted from
// _______ddddd_bbb.
func DecodeB3D5(inst Instruction) AddrMode {
	b := (inst.Op1 & 0x7)
	d := (inst.Op1 & 0x1f0) >> 4
	return AddrMode{Addr(b), Addr(d), NoIndex}
}

// DecodeD3R3 returns AddrMode{d, r, NoIndex} (where 16<=d,r<=23)
// extracted from _________ddd_rrr.
func DecodeD3R3(inst Instruction) AddrMode {
	d := ((inst.Op1 & 0x70) >> 4) | 0x10
	r := (inst.Op1 & 0x7) | 0x10
	return AddrMode{Addr(d), Addr(r), NoIndex}
}

// AddModeD4K8 returns AddrMode{d, K, NoIndex} (where 16<=d<=31 )
// extracted from ____KKKKddddKKKK.
func DecodeD4K8(inst Instruction) AddrMode {
	d := ((inst.Op1 & 0xf0) >> 4) | 0x10
	Ka := (inst.Op1 & 0xf00) >> 4
	Kb := (inst.Op1 & 0xf)
	return AddrMode{Addr(d), Addr(Ka | Kb), NoIndex}
}

// DecodeD4R4 returns AddrMode{d, r, NoIndex} (where 16<=d,r<=31)
// extracted from ________ddddrrrr.
func DecodeD4R4(inst Instruction) AddrMode {
	r := (inst.Op1 & 0xf) | 0x10
	d := ((inst.Op1 & 0xf0) >> 4) | 0x10
	return AddrMode{Addr(d), Addr(r), NoIndex}
}

// DecodeD5 returns AddrMode{d, 0, NoIndex} extracted from
// _______ddddd____.
func DecodeD5(inst Instruction) AddrMode {
	d := (inst.Op1 & 0x1f0) >> 4
	return AddrMode{Addr(d), 0, NoIndex}
}

// DecodeD5K16 returns AddrMode{d, k, NoIndex} extracted from
// _______ddddd____ kkkkkkkkkkkkkkkk.
func DecodeD5K16(inst Instruction) AddrMode {
	d := (inst.Op1 & 0x1f0) >> 4
	k := inst.Op2
	return AddrMode{Addr(d), Addr(k), NoIndex}
}

// DecodeD5R5 returns AddrMode{d, r, NoIndex} extracted from
// ______rdddddrrrr.
func DecodeD5R5(inst Instruction) AddrMode {
	d := (inst.Op1 & 0x1f0) >> 4
	ra := (inst.Op1 & 0x200) >> 5
	rb := (inst.Op1 & 0xf)
	return AddrMode{Addr(d), Addr(ra | rb), NoIndex}
}

// DecodeDDDDRRRR returns AddrMode{d, r, NoIndex} (where d,r are 0,2,..30)
// extracted from ________ddddrrrr.
func DecodeDDDDRRRR(inst Instruction) AddrMode {
	d := (inst.Op1 & 0xf0) >> 4
	r := (inst.Op1 & 0xf)
	return AddrMode{Addr(d * 2), Addr(r * 2), NoIndex}
}

// DecodeDDK6 returns AddrMode{d, K, NoIndex} (where d is one of
// 24, 26, 28, 30) extracted from ________KKddKKKK.
func DecodeDDK6(inst Instruction) AddrMode {
	Ka := (inst.Op1 & 0xc0) >> 2
	Kb := (inst.Op1 & 0xf)
	d := (inst.Op1 & 0x30) >> 4
	return AddrMode{Addr(d*2 + 24), Addr(Ka | Kb), NoIndex}
}

// DecodeElpm returns AddrMode{d, 0, ireg} extracted from
// _______ddddd____, or AddrMode{} for the no-argument form.
func DecodeElpm(inst Instruction) AddrMode {
	if inst.Op1 == 0x95d8 {
		return AddrMode{}
	}

	d := (inst.Op1 & 0x1f0) >> 4
	ireg := Z
	if (inst.Op1 & 0x1) != 0 {
		ireg = ZPostInc
	}
	return AddrMode{Addr(d), 0, ireg}
}

// DecodeK4 returns AddrMode{K, 0, NoIndex} extracted from
// ________KKKK____.
func DecodeK4(inst Instruction) AddrMode {
	K := (inst.Op1 & 0xf0) >> 4
	return AddrMode{Addr(K), 0, NoIndex}
}

// DecodeK12 returns AddrMode{K, 0, NoIndex} extracted from
// ____kkkkkkkkkkkk.
func DecodeK12(inst Instruction) AddrMode {
	k := Addr(inst.Op1 & 0xfff)
	if (k & 0x800) != 0 {
		k -= 0x1000
	}
	return AddrMode{k, 0, NoIndex}
}

// DecodeK22 returns AddrMode{k, 0, NoIndex} extracted from
// _______kkkkk___k kkkkkkkkkkkkkkkk.
func DecodeK22(inst Instruction) AddrMode {
	ka := Addr((inst.Op1 & 0x1f0)) << 13
	kb := Addr((inst.Op1 & 0x1)) << 16
	kc := Addr(inst.Op2)
	return AddrMode{ka | kb | kc, 0, NoIndex}
}

// DecodeLddStd returns AddrMode{d, q, ireg} extracted from
// __q_qq_ddddd_qqq.
func DecodeLddStd(inst Instruction) AddrMode {
	d := (inst.Op1 & 0x1f0) >> 4
	qa := (inst.Op1 & 0x2000) >> 8
	qb := (inst.Op1 & 0xc00) >> 7
	qc := (inst.Op1 & 0x7)
	ireg := Z
	if (inst.Op1 & 0x8) != 0 {
		ireg = Y
	}
	return AddrMode{Addr(d), Addr(qa | qb | qc), ireg}
}

var ldstireg = []IndexReg{
	Z, ZPostInc, ZPreDec, NoIndex, NoIndex, NoIndex, NoIndex, NoIndex,
	Y, YPostInc, YPreDec, NoIndex, X, XPostInc, XPreDec, NoIndex,
}

// DecodeLdSt returns AddrMode{d, 0, ireg} extracted from
// _______ddddd____.
func DecodeLdSt(inst Instruction) AddrMode {
	d := (inst.Op1 & 0x1f0) >> 4
	ireg := ldstireg[inst.Op1.Nibble0()]
	return AddrMode{Addr(d), 0, ireg}
}

// DecodeLpm returns AddrMode{d, 0, ireg} extracted from
// from _______ddddd____, or AddrMode{} for the no-argument form.
func DecodeLpm(inst Instruction) AddrMode {
	if inst.Op1 == 0x95c8 {
		return AddrMode{}
	}

	d := (inst.Op1 & 0x1f0) >> 4
	ireg := Z
	if (inst.Op1 & 0x1) != 0 {
		ireg = ZPostInc
	}
	return AddrMode{Addr(d), 0, ireg}
}

// DecodeNone returns AddrMode{} (for use by no-argument instructions).
func DecodeNone(inst Instruction) AddrMode {
	return AddrMode{}
}

// DecodeSpm returns either AddrMode{} or AddrMode{0, 0, ZPostInc}
// depending on the form of the Spm opcode
func DecodeSpm(inst Instruction) AddrMode {
	if inst.Op1 == 0x95e8 {
		return AddrMode{}
	} else {
		return AddrMode{0, 0, ZPostInc}
	}
}

var OpModes = []Mode{
	ModeD5R5,     // Adc
	ModeD5R5,     // Add
	ModeDDK6,     // Adiw
	ModeD5R5,     // And
	ModeD4K8,     // Andi
	ModeD5,       // Asr
	ModeB3,       // Bclr
	ModeB3D5,     // Bld
	ModeB3K7,     // Brbc
	ModeB3K7,     // Brbs
	ModeNone,     // Break
	ModeB3,       // Bset
	ModeB3D5,     // Bst
	ModeK22,      // Call
	ModeA5B3,     // Cbi
	ModeD5,       // Com
	ModeD5R5,     // Cp
	ModeD5R5,     // Cpc
	ModeD4K8,     // Cpi
	ModeD5R5,     // Cpse
	ModeD5,       // Dec
	ModeK4,       // Des
	ModeNone,     // Eicall
	ModeNone,     // Eijmp
	ModeElpm,     // Elpm
	ModeD5R5,     // Eor
	ModeD3R3,     // Fmul
	ModeD3R3,     // Fmuls
	ModeD3R3,     // Fmulsu
	ModeNone,     // Icall
	ModeNone,     // Ijmp
	ModeA6D5,     // In
	ModeD5,       // Inc
	ModeK22,      // Jmp
	ModeD5,       // Lac
	ModeD5,       // Las
	ModeD5,       // Lat
	ModeLdSt,     // Ld
	ModeLddStd,   // Ldd
	ModeD4K8,     // Ldi
	ModeD5K16,    // Lds
	ModeLpm,      // Lpm
	ModeD5,       // Lsr
	ModeD5R5,     // Mov
	ModeDDDDRRRR, // Movw
	ModeD5R5,     // Mul
	ModeD4R4,     // Muls
	ModeD3R3,     // Mulsu
	ModeD5,       // Neg
	ModeNone,     // Nop
	ModeD5R5,     // Or
	ModeD4K8,     // Ori
	ModeA6D5,     // Out
	ModeD5,       // Pop
	ModeD5,       // Push
	ModeK12,      // Rcall
	ModeNone,     // Ret
	ModeNone,     // Reti
	ModeK12,      // Rjmp
	ModeD5,       // Ror
	ModeD5R5,     // Sbc
	ModeD4K8,     // Sbci
	ModeA5B3,     // Sbi
	ModeA5B3,     // Sbic
	ModeA5B3,     // Sbis
	ModeDDK6,     // Sbiw
	ModeB3D5,     // Sbrc
	ModeB3D5,     // Sbrs
	ModeNone,     // Sleep
	ModeSpm,      // Spm
	ModeLdSt,     // St
	ModeLddStd,   // Std
	ModeD5K16,    // Sts
	ModeD5R5,     // Sub
	ModeD4K8,     // Subi
	ModeD5,       // Swap
	ModeNone,     // Wdr
	ModeD5,       // Xch
	ModeNone,     // Reserved
}

var decoders = []AddrDecoder{
	DecodeA5B3, DecodeA6D5, DecodeB3, DecodeB3K7, DecodeB3D5,
	DecodeD3R3, DecodeD4K8, DecodeD4R4, DecodeD5, DecodeD5K16,
	DecodeD5R5, DecodeDDDDRRRR, DecodeDDK6, DecodeElpm, DecodeK4,
	DecodeK12, DecodeK22, DecodeLddStd, DecodeLdSt, DecodeLpm,
	DecodeNone, DecodeSpm,
}

func (s Minimal) DecodeAddr(inst Instruction) AddrMode {
	return decoders[OpModes[inst.Mnem]](inst)
}
