package instr

type IndexMode int

const (
	NoMode  IndexMode = 0
	PostInc IndexMode = 4
	PreDec  IndexMode = 8
)

type IndexReg int

const (
	NoIndex  IndexReg = 0
	X        IndexReg = 1
	Y        IndexReg = 2
	Z        IndexReg = 3
	XPostInc IndexReg = X + IndexReg(PostInc)
	XPreDec  IndexReg = X + IndexReg(PreDec)
	YPostInc IndexReg = Y + IndexReg(PostInc)
	YPreDec  IndexReg = Y + IndexReg(PreDec)
	ZPostInc IndexReg = Z + IndexReg(PostInc)
	ZPreDec  IndexReg = Z + IndexReg(PreDec)
)

//go:generate stringer -type=IndexReg

func (r IndexReg) Base() Addr {
	return Addr(r & 0x3)
}

func (r IndexReg) Reg() Addr {
	return Addr(24 + (r&0x3)*2)
}

func (r IndexReg) Mode() IndexMode {
	return IndexMode(r & 0xc)
}

func (r IndexReg) WithMode(mode IndexMode) IndexReg {
	return (r & 0x3) + IndexReg(mode)
}

type Mode int

const (
	ModeA5B3 Mode = iota
	ModeA6D5
	ModeB3
	ModeB3K7
	ModeB3D5
	ModeD3R3
	ModeD4K7
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

type AddrDecoder func(*AddrMode, Opcode, Opcode)

// DecodeA5B3 returns AddrMode{A, b, NoIndex} extracted from
// ________AAAAAbbb.
func DecodeA5B3(am *AddrMode, op1, op2 Opcode) {
	A := (op1 & 0xf8) >> 3
	b := (op1 & 0x7)
	am.A1 = Addr(A)
	am.A2 = Addr(b)
	am.Ireg = NoIndex
}

// DecodeA6D5 returns AddrMode{A, d, NoIndex} extracted from
// _____AAdddddAAAA.
func DecodeA6D5(am *AddrMode, op1, op2 Opcode) {
	Aa := (op1 & 0x600) >> 5
	Ab := (op1 & 0xf)
	d := (op1 & 0x1f0) >> 4
	am.A1 = Addr(Aa | Ab)
	am.A2 = Addr(d)
	am.Ireg = NoIndex
}

// DecodeB3 returns AddrMode{b, 0, NoIndex} extracted from
// _________bbb____.
func DecodeB3(am *AddrMode, op1, op2 Opcode) {
	b := (op1 & 0x70) >> 4
	am.A1 = Addr(b)
	am.A2 = 0
	am.Ireg = NoIndex
}

// DecodeB3K7 returns AddrMode{b, k, NoIndex} (where -64<=k<=63)
// extracted from ______kkkkkkkbbb.
func DecodeB3K7(am *AddrMode, op1, op2 Opcode) {
	b := (op1 & 0x7)
	k := Addr((op1 & 0x3f8) >> 3)
	if (k & 0x40) != 0 {
		k -= 0x80
	}
	am.A1 = Addr(b)
	am.A2 = k
	am.Ireg = NoIndex
}

// DecodeB3D5 returns AddrMode{b, d, NoIndex} extracted from
// _______ddddd_bbb.
func DecodeB3D5(am *AddrMode, op1, op2 Opcode) {
	b := (op1 & 0x7)
	d := (op1 & 0x1f0) >> 4
	am.A1 = Addr(b)
	am.A2 = Addr(d)
	am.Ireg = NoIndex
}

// DecodeD3R3 returns AddrMode{d, r, NoIndex} (where 16<=d,r<=23)
// extracted from _________ddd_rrr.
func DecodeD3R3(am *AddrMode, op1, op2 Opcode) {
	d := ((op1 & 0x70) >> 4) | 0x10
	r := (op1 & 0x7) | 0x10
	am.A1 = Addr(d)
	am.A2 = Addr(r)
	am.Ireg = NoIndex
}

// AddModeD4K7 returns AddrMode{d, k, NoIndex} (where 16<=d<=31 )
// extracted from _____kkkddddkkkk.
func DecodeD4K7(am *AddrMode, op1, op2 Opcode) {
	d := ((op1 & 0xf0) >> 4) | 0x10
	ka := (op1 & 0x700) >> 4
	kb := (op1 & 0xf)
	am.A1 = Addr(d)
	am.A2 = Addr(ka | kb)
	am.Ireg = NoIndex
}

// AddModeD4K8 returns AddrMode{d, K, NoIndex} (where 16<=d<=31 )
// extracted from ____KKKKddddKKKK.
func DecodeD4K8(am *AddrMode, op1, op2 Opcode) {
	d := ((op1 & 0xf0) >> 4) | 0x10
	Ka := (op1 & 0xf00) >> 4
	Kb := (op1 & 0xf)
	am.A1 = Addr(d)
	am.A2 = Addr(Ka | Kb)
	am.Ireg = NoIndex
}

// DecodeD4R4 returns AddrMode{d, r, NoIndex} (where 16<=d,r<=31)
// extracted from ________ddddrrrr.
func DecodeD4R4(am *AddrMode, op1, op2 Opcode) {
	r := (op1 & 0xf) | 0x10
	d := ((op1 & 0xf0) >> 4) | 0x10
	am.A1 = Addr(d)
	am.A2 = Addr(r)
	am.Ireg = NoIndex
}

// DecodeD5 returns AddrMode{d, 0, NoIndex} extracted from
// _______ddddd____.
func DecodeD5(am *AddrMode, op1, op2 Opcode) {
	d := (op1 & 0x1f0) >> 4
	am.A1 = Addr(d)
	am.A2 = 0
	am.Ireg = NoIndex
}

// DecodeD5K16 returns AddrMode{d, k, NoIndex} extracted from
// _______ddddd____ kkkkkkkkkkkkkkkk.
func DecodeD5K16(am *AddrMode, op1, op2 Opcode) {
	d := (op1 & 0x1f0) >> 4
	k := op2
	am.A1 = Addr(d)
	am.A2 = Addr(k)
	am.Ireg = NoIndex
}

// DecodeD5R5 returns AddrMode{d, r, NoIndex} extracted from
// ______rdddddrrrr.
func DecodeD5R5(am *AddrMode, op1, op2 Opcode) {
	d := (op1 & 0x1f0) >> 4
	ra := (op1 & 0x200) >> 5
	rb := (op1 & 0xf)
	am.A1 = Addr(d)
	am.A2 = Addr(ra | rb)
	am.Ireg = NoIndex
}

// DecodeDDDDRRRR returns AddrMode{d, r, NoIndex} (where d,r are 0,2,..30)
// extracted from ________ddddrrrr.
func DecodeDDDDRRRR(am *AddrMode, op1, op2 Opcode) {
	d := (op1 & 0xf0) >> 4
	r := (op1 & 0xf)
	am.A1 = Addr(d * 2)
	am.A2 = Addr(r * 2)
	am.Ireg = NoIndex
}

// DecodeDDK6 returns AddrMode{d, K, NoIndex} (where d is one of
// 24, 26, 28, 30) extracted from ________KKddKKKK.
func DecodeDDK6(am *AddrMode, op1, op2 Opcode) {
	Ka := (op1 & 0xc0) >> 2
	Kb := (op1 & 0xf)
	d := (op1 & 0x30) >> 4
	am.A1 = Addr(d*2 + 24)
	am.A2 = Addr(Ka | Kb)
	am.Ireg = NoIndex
}

// DecodeElpm returns AddrMode{d, 0, ireg} extracted from
// _______ddddd____, or AddrMode{0, 0, Z} for the no-argument form.
func DecodeElpm(am *AddrMode, op1, op2 Opcode) {
	if op1 == 0x95d8 {
		am.A1 = 0
		am.A2 = 0
		am.Ireg = Z
		return
	}

	d := (op1 & 0x1f0) >> 4
	ireg := Z
	if (op1 & 0x1) != 0 {
		ireg = ZPostInc
	}
	am.A1 = Addr(d)
	am.A2 = 0
	am.Ireg = ireg
}

// DecodeK4 returns AddrMode{K, 0, NoIndex} extracted from
// ________KKKK____.
func DecodeK4(am *AddrMode, op1, op2 Opcode) {
	K := (op1 & 0xf0) >> 4
	am.A1 = Addr(K)
	am.A2 = 0
	am.Ireg = NoIndex
}

// DecodeK12 returns AddrMode{K, 0, NoIndex} extracted from
// ____kkkkkkkkkkkk.
func DecodeK12(am *AddrMode, op1, op2 Opcode) {
	k := Addr(op1 & 0xfff)
	if (k & 0x800) != 0 {
		k -= 0x1000
	}
	am.A1 = k
	am.A2 = 0
	am.Ireg = NoIndex
}

// DecodeK22 returns AddrMode{k, 0, NoIndex} extracted from
// _______kkkkk___k kkkkkkkkkkkkkkkk.
func DecodeK22(am *AddrMode, op1, op2 Opcode) {
	ka := Addr((op1 & 0x1f0)) << 13
	kb := Addr((op1 & 0x1)) << 16
	kc := Addr(op2)
	am.A1 = ka | kb | kc
	am.A2 = 0
	am.Ireg = NoIndex
}

// DecodeLddStd returns AddrMode{d, q, ireg} extracted from
// __q_qq_ddddd_qqq.
func DecodeLddStd(am *AddrMode, op1, op2 Opcode) {
	d := (op1 & 0x1f0) >> 4
	qa := (op1 & 0x2000) >> 8
	qb := (op1 & 0xc00) >> 7
	qc := (op1 & 0x7)
	ireg := Z
	if (op1 & 0x8) != 0 {
		ireg = Y
	}
	am.A1 = Addr(d)
	am.A2 = Addr(qa | qb | qc)
	am.Ireg = ireg
}

var ldstireg = []IndexReg{
	Z, ZPostInc, ZPreDec, NoIndex, NoIndex, NoIndex, NoIndex, NoIndex,
	Y, YPostInc, YPreDec, NoIndex, X, XPostInc, XPreDec, NoIndex,
}

// DecodeLdSt returns AddrMode{d, 0, ireg} extracted from
// _______ddddd____.
func DecodeLdSt(am *AddrMode, op1, op2 Opcode) {
	d := (op1 & 0x1f0) >> 4
	ireg := ldstireg[op1.Nibble0()]
	am.A1 = Addr(d)
	am.A2 = 0
	am.Ireg = ireg
}

// DecodeLpm returns AddrMode{d, 0, ireg} extracted from
// _______ddddd____.
func DecodeLpm(am *AddrMode, op1, op2 Opcode) {
	if op1 == 0x95c8 {
		am.A1 = 0
		am.A2 = 0
		am.Ireg = Z
		return
	}

	d := (op1 & 0x1f0) >> 4
	ireg := Z
	if (op1 & 0x1) != 0 {
		ireg = ZPostInc
	}
	am.A1 = Addr(d)
	am.A2 = 0
	am.Ireg = ireg
}

// DecodeNone returns AddrMode{} (for use by no-argument instructions).
func DecodeNone(am *AddrMode, op1, op2 Opcode) {
}

// DecodeSpm returns AddrMode{0, 0, ZPostInc}.
func DecodeSpm(am *AddrMode, op1, op2 Opcode) {
	am.A1 = 0
	am.A2 = 0
	am.Ireg = ZPostInc
}

var OpModes = [...]Mode{
	ModeNone,     // Reserved
	ModeD5R5,     // Adc
	ModeD5R5,     // AdcReduced
	ModeD5R5,     // Add
	ModeD5R5,     // AddReduced
	ModeDDK6,     // Adiw
	ModeD5R5,     // And
	ModeD5R5,     // AndReduced
	ModeD4K8,     // Andi
	ModeD5,       // Asr
	ModeD5,       // AsrReduced
	ModeB3,       // Bclr
	ModeB3D5,     // Bld
	ModeB3D5,     // BldReduced
	ModeB3K7,     // Brbc
	ModeB3K7,     // Brbs
	ModeNone,     // Break
	ModeB3,       // Bset
	ModeB3D5,     // Bst
	ModeB3D5,     // BstReduced
	ModeK22,      // Call
	ModeA5B3,     // Cbi
	ModeD5,       // Com
	ModeD5,       // ComReduced
	ModeD5R5,     // Cp
	ModeD5R5,     // CpReduced
	ModeD5R5,     // Cpc
	ModeD5R5,     // CpcReduced
	ModeD4K8,     // Cpi
	ModeD5R5,     // Cpse
	ModeD5R5,     // CpseReduced
	ModeD5,       // Dec
	ModeD5,       // DecReduced
	ModeK4,       // Des
	ModeNone,     // Eicall
	ModeNone,     // Eijmp
	ModeElpm,     // Elpm
	ModeElpm,     // ElpmEnhanced
	ModeD5R5,     // Eor
	ModeD5R5,     // EorReduced
	ModeD3R3,     // Fmul
	ModeD3R3,     // Fmuls
	ModeD3R3,     // Fmulsu
	ModeNone,     // Icall
	ModeNone,     // Ijmp
	ModeA6D5,     // In
	ModeA6D5,     // InReduced
	ModeD5,       // Inc
	ModeD5,       // IncReduced
	ModeK22,      // Jmp
	ModeD5,       // Lac
	ModeD5,       // Las
	ModeD5,       // Lat
	ModeLdSt,     // LdClassic
	ModeLdSt,     // LdClassicReduced
	ModeLdSt,     // LdMinimal
	ModeLdSt,     // LdMinimalReduced
	ModeLddStd,   // Ldd
	ModeD4K8,     // Ldi
	ModeD5K16,    // Lds
	ModeD4K7,     // Lds16
	ModeLpm,      // Lpm
	ModeLpm,      // LpmEnhanced
	ModeD5,       // Lsr
	ModeD5,       // LsrReduced
	ModeD5R5,     // Mov
	ModeD5R5,     // MovReduced
	ModeDDDDRRRR, // Movw
	ModeD5R5,     // Mul
	ModeD4R4,     // Muls
	ModeD3R3,     // Mulsu
	ModeD5,       // Neg
	ModeD5,       // NegReduced
	ModeNone,     // Nop
	ModeD5R5,     // Or
	ModeD5R5,     // OrReduced
	ModeD4K8,     // Ori
	ModeA6D5,     // Out
	ModeA6D5,     // OutReduced
	ModeD5,       // Pop
	ModeD5,       // PopReduced
	ModeD5,       // Push
	ModeD5,       // PushReduced
	ModeK12,      // Rcall
	ModeNone,     // Ret
	ModeNone,     // Reti
	ModeK12,      // Rjmp
	ModeD5,       // Ror
	ModeD5,       // RorReduced
	ModeD5R5,     // Sbc
	ModeD5R5,     // SbcReduced
	ModeD4K8,     // Sbci
	ModeA5B3,     // Sbi
	ModeA5B3,     // Sbic
	ModeA5B3,     // Sbis
	ModeDDK6,     // Sbiw
	ModeB3D5,     // Sbrc
	ModeB3D5,     // SbrcReduced
	ModeB3D5,     // Sbrs
	ModeB3D5,     // SbrsReduced
	ModeNone,     // Sleep
	ModeNone,     // Spm
	ModeSpm,      // SpmXmega
	ModeLdSt,     // StClassic
	ModeLdSt,     // StClassicReduced
	ModeLdSt,     // StMinimal
	ModeLdSt,     // StMinimalReduced
	ModeLddStd,   // Std
	ModeD5K16,    // Sts
	ModeD4K7,     // Sts16
	ModeD5R5,     // Sub
	ModeD5R5,     // SubReduced
	ModeD4K8,     // Subi
	ModeD5,       // Swap
	ModeD5,       // SwapReduced
	ModeNone,     // Wdr
	ModeD5,       // Xch
}

var decoders = []AddrDecoder{
	DecodeA5B3, DecodeA6D5, DecodeB3, DecodeB3K7, DecodeB3D5,
	DecodeD3R3, DecodeD4K7, DecodeD4K8, DecodeD4R4, DecodeD5,
	DecodeD5K16, DecodeD5R5, DecodeDDDDRRRR, DecodeDDK6,
	DecodeElpm, DecodeK4, DecodeK12, DecodeK22, DecodeLddStd,
	DecodeLdSt, DecodeLpm, DecodeNone, DecodeSpm,
}
