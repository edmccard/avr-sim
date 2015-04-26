package instr

import "fmt"

// An IndexAction represents the increment/decrement action of an indexed
// address mode.
type IndexAction int

const (
	NoAction IndexAction = 0
	PostInc  IndexAction = 4
	PreDec   IndexAction = 8
)

// An IndexReg represents the index register (and increment/decrement
// action) for indexed address modes.
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

// Base returns the equivalent register "name" of an IndexReg. For
// example,
//  X.Base() == XPostInc.Base() == XPreDec.Base() == X
func (r IndexReg) Base() IndexReg {
	return r & 0x3
}

// Reg returns the (low-byte) register number of an IndexReg. For
// example,
//  Z.Reg() == 30
func (r IndexReg) Reg() int {
	return int(24 + (r&0x3)*2)
}

// Action returns the increment/decrement action (none,
// post-increment, pre-decrement) of an IndexReg.
func (r IndexReg) Action() IndexAction {
	return IndexAction(r & 0xc)
}

// WithAction creates a new IndexReg by adding an action to the base
// of an IndexReg. For example, X.WithAction(PostInc) == XPostInc
func (r IndexReg) WithAction(action IndexAction) IndexReg {
	return (r & 0x3) + IndexReg(action)
}

// An AddrMode uniquely determines the meaning of the operands for an
// instruction. (See Mnemonic for the mapping from instructions to
// address modes.)
type AddrMode int

const (
	ModeNone    AddrMode = iota
	Mode2Reg3            // Mulsu r<Dst[16,23]>, r<Src[16,23]>
	Mode2Reg4            // Muls r<Dst[16,31]>, r<Src[16,31]>
	Mode2Reg5            // Add r<Dst[0,31]>, r<Src[0,31]>
	ModeAtomic           // Lac Z, <Dst[0,31]>
	ModeBranch           // Brbs <Src[0,7>, pc+1+<Off[-64,63]>
	ModeDes              // Des <Src[0,15]>
	ModeLpmEnh           // Lpm <Dst[0,31]>, <Src[IndexReg]>
	ModeIn               // In <Dst[0,31]>, $<Src[0,64]>
	ModeIOBit            // Sbis $<Src[0,31]>, <Off[0,7]>
	ModeLd               // Ld r<Dst[0,31]>, <Src[IndexReg]>
	ModeLdd              // Ldd r<Dst[0,31]>, <Src[IndexReg]>+<Off[0,63]>
	ModeLds              // Lds r<Dst[0,31], $<Off[uint16]>
	ModeLds16            // Lds16 r<Dst[16,31]>, $<Off[64,191]>
	ModeOut              // Out $<Dst[0,64]>, r<Src[0,31]>
	ModePc               // Jmp $<Off[uint22]>
	ModePcOff            // Rcall pc+1+<Off[-2048,2047]>
	ModeReg5             // Inc r<Dst[0,31]>
	ModeRegBit           // Bld r<Dst[0,31]>, <Off[0,7>
	ModeRegImm           // Andi r<Dst[16,31]>, <Src[0,255]>
	ModeRegPair          // Movw r<Dst+1>:r<Dst>, r<Src+1>:r<Src>
	ModePairImm          // Adiw r<Dst+1>:r<Dst>, <Src[0,63]>
	ModeSBit             // Bclr <Off[0,7]>
	ModeSt               // St <Dst[IndexReg]>, r<Src[0,31]>
	ModeStd              // Std <Dst[IndexReg>+<Off[0,63]>, r<Src[0,31]>
	ModeSts              // Sts $<Off[uint16]>, r<Src[0,31]>
	ModeSts16            // Sts16 $<Off[64,191]>, r<Src[16,31]>
	ModeSpmX             // Spm <Src[IndexReg]>
)

var decoders = []operandDecoder{
	decodeNone, decode2Reg3, decode2Reg4, decode2Reg5, decodeAtomic,
	decodeBranch, decodeDes, decodeLpmEnh, decodeIn, decodeIOBit,
	decodeLd, decodeLdd, decodeLdsSts, decodeLdsSts16, decodeOut,
	decodePc, decodePcOff, decodeReg5, decodeRegBit, decodeRegImm,
	decodeRegPair, decodePairImm, decodeSBit, decodeSt, decodeStd,
	decodeLdsSts, decodeLdsSts16, decodeSpm,
}

//go:generate stringer -type=AddrMode

// Operands stores the operands for an instruction. The Mode member
// determines the meaning of Dst, Src, and Off.
type Operands struct {
	Dst, Src, Off int
	Mode          AddrMode
}

type nreg int

func (r nreg) String() string {
	switch r {
	case 26:
		return "XL"
	case 27:
		return "XH"
	case 28:
		return "YL"
	case 29:
		return "YH"
	case 30:
		return "ZL"
	case 31:
		return "ZH"
	default:
		return fmt.Sprintf("r%d", int(r))
	}
}

// String returns a string representation of an Operands in
// "disassembly" format. For example,
//  ops.Mode == Mode2Reg5 && ops.String() == "r5, r6"
func (o Operands) String() string {
	switch o.Mode {
	case ModeNone:
		return "<implied>"
	case ModeIOBit:
		return fmt.Sprintf("$%02x, %d", o.Dst, o.Off)
	case ModeSBit:
		return fmt.Sprintf("%d", o.Off)
	case ModeBranch:
		return fmt.Sprintf("%d, PC%+d", o.Src, o.Off+1)
	case Mode2Reg3, Mode2Reg4, Mode2Reg5:
		return fmt.Sprintf("%s, %s", nreg(o.Dst), nreg(o.Src))
	case ModeRegImm:
		return fmt.Sprintf("%s, $%02x", nreg(o.Dst), o.Src)
	case ModeReg5:
		return fmt.Sprintf("%s", nreg(o.Dst))
	case ModeIn:
		return fmt.Sprintf("%s, $%02x", nreg(o.Dst), o.Src)
	case ModeOut:
		return fmt.Sprintf("$%02x, %s", o.Dst, nreg(o.Src))
	case ModeRegBit:
		return fmt.Sprintf("%s, %d", nreg(o.Dst), o.Off)
	case ModeLds:
		return fmt.Sprintf("%s, $%04x", nreg(o.Dst), o.Off)
	case ModeRegPair:
		return fmt.Sprintf("%s:%s, %s:%s", nreg(o.Dst+1), nreg(o.Dst),
			nreg(o.Src+1), nreg(o.Src))
	case ModePairImm:
		return fmt.Sprintf("%s:%s, %d", nreg(o.Dst+1), nreg(o.Dst), o.Src)
	case ModeLpmEnh, ModeLd:
		return fmt.Sprintf("%s, %s", nreg(o.Dst), IndexReg(o.Src))
	case ModeSt:
		return fmt.Sprintf("%s, %s", IndexReg(o.Dst), nreg(o.Src))
	case ModeDes:
		return fmt.Sprintf("%02x", o.Src)
	case ModeLds16:
		return fmt.Sprintf("%s, $%02x", nreg(o.Dst), o.Off)
	case ModeSts16:
		return fmt.Sprintf("$%02x, %s", o.Off, nreg(o.Src))
	case ModePcOff:
		return fmt.Sprintf("PC%+d", o.Off)
	case ModeSts:
		return fmt.Sprintf("$%04x, %s", o.Off, nreg(o.Src))
	case ModePc:
		return fmt.Sprintf("%04x", o.Off)
	case ModeLdd:
		return fmt.Sprintf("%s, %s+%d", nreg(o.Dst), IndexReg(o.Src), o.Off)
	case ModeStd:
		return fmt.Sprintf("%s+%d, %s", IndexReg(o.Dst), o.Off, nreg(o.Src))
	case ModeSpmX:
		return "Z+"
	}
	return fmt.Sprintf("Mode(%d)", int(o.Mode))
}

type operandDecoder func(*Operands, Opcode, Opcode)

// decodeIOBit returns Operands{Dst:A, Src:A, Off:b} extracted from
// ________AAAAAbbb.
func decodeIOBit(o *Operands, op1, op2 Opcode) {
	A := (int(op1) & 0xf8) >> 3
	b := (int(op1) & 0x7)
	o.Dst = A
	o.Src = A
	o.Off = b
}

// decodeIn returns Operands{Src:A, Dst:d, NoIndex} extracted from
// _____AAdddddAAAA.
func decodeIn(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0x1f0) >> 4
	Aa := (int(op1) & 0x600) >> 5
	Ab := (int(op1) & 0xf)
	o.Dst = d
	o.Src = Aa | Ab
}

// decodeOut returns Operands{Dst:A, Src:d, NoIndex} extracted from
// _____AAdddddAAAA.
func decodeOut(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0x1f0) >> 4
	Aa := (int(op1) & 0x600) >> 5
	Ab := (int(op1) & 0xf)
	o.Src = d
	o.Dst = Aa | Ab
}

// decodeSBit returns Operands{Src:b, Dst:b} extracted from
// _________bbb____.
func decodeSBit(o *Operands, op1, op2 Opcode) {
	b := (int(op1) & 0x70) >> 4
	o.Src = b
	o.Dst = b
}

// decodeBranch returns Operands{Src: b, Off: k} (where -64<=k<=63)
// extracted from ______kkkkkkkbbb.
func decodeBranch(o *Operands, op1, op2 Opcode) {
	b := (int(op1) & 0x7)
	k := (int(op1) & 0x3f8) >> 3
	if (k & 0x40) != 0 {
		k -= 0x80
	}
	o.Src = b
	o.Off = k
}

// decodeRegBit returns Operands{Off:b, Src:d, Dst:d} extracted from
// _______ddddd_bbb.
func decodeRegBit(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0x1f0) >> 4
	b := (int(op1) & 0x7)
	o.Src = d
	o.Dst = d
	o.Off = b
}

// decode2Reg3 returns Operands{Dst:d, Src:r} (where 16<=d,r<=23)
// extracted from _________ddd_rrr.
func decode2Reg3(o *Operands, op1, op2 Opcode) {
	d := ((int(op1) & 0x70) >> 4) | 0x10
	r := (int(op1) & 0x7) | 0x10
	o.Dst = d
	o.Src = r
}

// decodeLdsSts16 returns Operands{Dst:d, Src:d, Off:k} (where
// 16<=d<=31) extracted from _____kkkddddkkkk.
func decodeLdsSts16(o *Operands, op1, op2 Opcode) {
	b8 := (int(op1) & 0x100) >> 2
	ka := (^(b8 << 1) & 0x80) | b8
	kb := (int(op1) & 0x600) >> 5
	kc := (int(op1) & 0xf)
	d := ((int(op1) & 0xf0) >> 4) | 0x10
	o.Off = ka | kb | kc
	o.Src = d
	o.Dst = d
}

// decodeRegImm returns Operands{Dst:d, Src:K} (where 16<=d<=31 )
// extracted from ____KKKKddddKKKK.
func decodeRegImm(o *Operands, op1, op2 Opcode) {
	d := ((int(op1) & 0xf0) >> 4) | 0x10
	Ka := (int(op1) & 0xf00) >> 4
	Kb := (int(op1) & 0xf)
	o.Dst = d
	o.Src = Ka | Kb
}

// decode2Reg4 returns Operands{Dst:d, Src:r} (where 16<=d,r<=31)
// extracted from ________ddddrrrr.
func decode2Reg4(o *Operands, op1, op2 Opcode) {
	r := (int(op1) & 0xf) | 0x10
	d := ((int(op1) & 0xf0) >> 4) | 0x10
	o.Dst = d
	o.Src = r
}

// decodeReg5 returns Operands{Dst:d, Src:d} extracted from
// _______ddddd____.
func decodeReg5(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0x1f0) >> 4
	o.Dst = d
	o.Src = d
}

// decodeAtomic returns Operands{Dst:d, Src:Z} extracted from
// _______ddddd____.
func decodeAtomic(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0x1f0) >> 4
	o.Dst = d
	o.Src = int(Z)
}

// decodeLdsSts returns Operands{Dst:d, Src:d, Off:k} extracted from
// _______ddddd____ kkkkkkkkkkkkkkkk.
func decodeLdsSts(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0x1f0) >> 4
	k := int(op2)
	o.Dst = d
	o.Src = d
	o.Off = k
}

// decode2Reg5 returns Operands{Dst:d, Src:r} extracted from
// ______rdddddrrrr.
func decode2Reg5(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0x1f0) >> 4
	ra := (int(op1) & 0x200) >> 5
	rb := (int(op1) & 0xf)
	o.Dst = d
	o.Src = ra | rb
}

// decodeRegPair returns Operands{Dst:d, Src:r} (where d,r are
// 0,2,..30) extracted from ________ddddrrrr.
func decodeRegPair(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0xf0) >> 4
	r := (int(op1) & 0xf)
	o.Dst = d * 2
	o.Src = r * 2
}

// decodePairImm returns Operands{Dst:d, Src:K} (where d is one of 24,
// 26, 28, 30) extracted from ________KKddKKKK.
func decodePairImm(o *Operands, op1, op2 Opcode) {
	Ka := (int(op1) & 0xc0) >> 2
	Kb := (int(op1) & 0xf)
	d := (int(op1) & 0x30) >> 4
	o.Dst = d*2 + 24
	o.Src = Ka | Kb
}

// decodeDes returns Operands{Src:K} extracted from
// ________KKKK____.
func decodeDes(o *Operands, op1, op2 Opcode) {
	K := (int(op1) & 0xf0) >> 4
	o.Src = K
}

// decodePcOff returns Operands{Off:K} extracted from
// ____kkkkkkkkkkkk.
func decodePcOff(o *Operands, op1, op2 Opcode) {
	k := int(op1) & 0xfff
	if (k & 0x800) != 0 {
		k -= 0x1000
	}
	o.Off = k
}

// decodePc returns Operands{Off:k} extracted from
// _______kkkkk___k kkkkkkkkkkkkkkkk.
func decodePc(o *Operands, op1, op2 Opcode) {
	ka := (int(op1) & 0x1f0) << 13
	kb := (int(op1) & 0x1) << 16
	kc := int(op2)
	o.Off = ka | kb | kc
}

// decodeLdd returns Operands{Dst:d, Off:q, Src:ireg} extracted from
// __q_qq_ddddd_qqq.
func decodeLdd(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0x1f0) >> 4
	qa := (int(op1) & 0x2000) >> 8
	qb := (int(op1) & 0xc00) >> 7
	qc := (int(op1) & 0x7)
	ireg := Z
	if (op1 & 0x8) != 0 {
		ireg = Y
	}
	o.Dst = d
	o.Off = qa | qb | qc
	o.Src = int(ireg)
}

// decodeStd returns Operands{Off:q, Src:d, Dst:ireg} extracted from
// __q_qq_ddddd_qqq.
func decodeStd(o *Operands, op1, op2 Opcode) {
	qa := (int(op1) & 0x2000) >> 8
	qb := (int(op1) & 0xc00) >> 7
	qc := (int(op1) & 0x7)
	d := (int(op1) & 0x1f0) >> 4
	ireg := Z
	if (op1 & 0x8) != 0 {
		ireg = Y
	}
	o.Off = qa | qb | qc
	o.Src = d
	o.Dst = int(ireg)
}

var ldstireg = []IndexReg{
	Z, ZPostInc, ZPreDec, NoIndex, NoIndex, NoIndex, NoIndex, NoIndex,
	Y, YPostInc, YPreDec, NoIndex, X, XPostInc, XPreDec, NoIndex,
}

// decodeLd returns Operands{Dst:d, Src:ireg} extracted from
// _______ddddd____.
func decodeLd(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0x1f0) >> 4
	ireg := ldstireg[op1.nibble0()]
	o.Dst = d
	o.Src = int(ireg)
}

// decodeSt returns Operands{Src:d, Dst:ireg} extracted from
// _______ddddd____.
func decodeSt(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0x1f0) >> 4
	ireg := ldstireg[op1.nibble0()]
	o.Src = d
	o.Dst = int(ireg)
}

// decodeLpmEnhreturns Operands{Dst:d, Src:ireg} extracted from
// _______ddddd____.
func decodeLpmEnh(o *Operands, op1, op2 Opcode) {
	d := (int(op1) & 0x1f0) >> 4
	ireg := Z
	if (op1 & 0x1) != 0 {
		ireg = ZPostInc
	}
	o.Dst = d
	o.Src = int(ireg)
}

// decodeNone returns Operands{} (for use by no-argument instructions).
func decodeNone(o *Operands, op1, op2 Opcode) {
}

// decodeSpm returns Operands{Src:ZPostInc}.
func decodeSpm(o *Operands, op1, op2 Opcode) {
	o.Src = int(ZPostInc)
}

var opModes = [...]AddrMode{
	ModeNone,    // Reserved
	Mode2Reg5,   // Adc
	Mode2Reg5,   // AdcReduced
	Mode2Reg5,   // Add
	Mode2Reg5,   // AddReduced
	ModePairImm, // Adiw
	Mode2Reg5,   // And
	Mode2Reg5,   // AndReduced
	ModeRegImm,  // Andi
	ModeReg5,    // Asr
	ModeReg5,    // AsrReduced
	ModeSBit,    // Bclr
	ModeRegBit,  // Bld
	ModeRegBit,  // BldReduced
	ModeBranch,  // Brbc
	ModeBranch,  // Brbs
	ModeNone,    // Break
	ModeSBit,    // Bset
	ModeRegBit,  // Bst
	ModeRegBit,  // BstReduced
	ModePc,      // Call
	ModeIOBit,   // Cbi
	ModeReg5,    // Com
	ModeReg5,    // ComReduced
	Mode2Reg5,   // Cp
	Mode2Reg5,   // CpReduced
	Mode2Reg5,   // Cpc
	Mode2Reg5,   // CpcReduced
	ModeRegImm,  // Cpi
	Mode2Reg5,   // Cpse
	Mode2Reg5,   // CpseReduced
	ModeReg5,    // Dec
	ModeReg5,    // DecReduced
	ModeDes,     // Des
	ModeNone,    // Eicall
	ModeNone,    // Eijmp
	ModeNone,    // Elpm
	ModeLpmEnh,  // ElpmEnhanced
	Mode2Reg5,   // Eor
	Mode2Reg5,   // EorReduced
	Mode2Reg3,   // Fmul
	Mode2Reg3,   // Fmuls
	Mode2Reg3,   // Fmulsu
	ModeNone,    // Icall
	ModeNone,    // Ijmp
	ModeIn,      // In
	ModeIn,      // InReduced
	ModeReg5,    // Inc
	ModeReg5,    // IncReduced
	ModePc,      // Jmp
	ModeReg5,    // Lac
	ModeReg5,    // Las
	ModeReg5,    // Lat
	ModeLd,      // Ld
	ModeLd,      // LdReduced
	ModeLd,      // LdMinimal
	ModeLd,      // LdMinimalReduced
	ModeLdd,     // Ldd
	ModeRegImm,  // Ldi
	ModeLds,     // Lds
	ModeLds16,   // Lds16
	ModeNone,    // Lpm
	ModeLpmEnh,  // LpmEnhanced
	ModeReg5,    // Lsr
	ModeReg5,    // LsrReduced
	Mode2Reg5,   // Mov
	Mode2Reg5,   // MovReduced
	ModeRegPair, // Movw
	Mode2Reg5,   // Mul
	Mode2Reg4,   // Muls
	Mode2Reg3,   // Mulsu
	ModeReg5,    // Neg
	ModeReg5,    // NegReduced
	ModeNone,    // Nop
	Mode2Reg5,   // Or
	Mode2Reg5,   // OrReduced
	ModeRegImm,  // Ori
	ModeOut,     // Out
	ModeOut,     // OutReduced
	ModeReg5,    // Pop
	ModeReg5,    // PopReduced
	ModeReg5,    // Push
	ModeReg5,    // PushReduced
	ModePcOff,   // Rcall
	ModeNone,    // Ret
	ModeNone,    // Reti
	ModePcOff,   // Rjmp
	ModeReg5,    // Ror
	ModeReg5,    // RorReduced
	Mode2Reg5,   // Sbc
	Mode2Reg5,   // SbcReduced
	ModeRegImm,  // Sbci
	ModeIOBit,   // Sbi
	ModeIOBit,   // Sbic
	ModeIOBit,   // Sbis
	ModePairImm, // Sbiw
	ModeRegBit,  // Sbrc
	ModeRegBit,  // SbrcReduced
	ModeRegBit,  // Sbrs
	ModeRegBit,  // SbrsReduced
	ModeNone,    // Sleep
	ModeNone,    // Spm
	ModeSpmX,    // SpmXmega
	ModeSt,      // St
	ModeSt,      // StReduced
	ModeSt,      // StMinimal
	ModeSt,      // StMinimalReduced
	ModeStd,     // Std
	ModeSts,     // Sts
	ModeSts16,   // Sts16
	Mode2Reg5,   // Sub
	Mode2Reg5,   // SubReduced
	ModeRegImm,  // Subi
	ModeReg5,    // Swap
	ModeReg5,    // SwapReduced
	ModeNone,    // Wdr
	ModeReg5,    // Xch
}
