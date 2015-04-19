package instr

type Mnemonic int

const (
	Reserved Mnemonic = iota
	Adc
	AdcReduced
	Add
	AddReduced
	Adiw
	And
	AndReduced
	Andi
	Asr
	AsrReduced
	Bclr
	Bld
	BldReduced
	Brbc
	Brbs
	Break
	Bset
	Bst
	BstReduced
	Call
	Cbi
	Com
	ComReduced
	Cp
	CpReduced
	Cpc
	CpcReduced
	Cpi
	Cpse
	CpseReduced
	Dec
	DecReduced
	Des
	Eicall
	Eijmp
	Elpm
	ElpmEnhanced
	Eor
	EorReduced
	Fmul
	Fmuls
	Fmulsu
	Icall
	Ijmp
	In
	InReduced
	Inc
	IncReduced
	Jmp
	Lac
	Las
	Lat
	LdClassic
	LdClassicReduced
	LdMinimal
	LdMinimalReduced
	Ldd
	Ldi
	Lds
	Lds16
	Lpm
	LpmEnhanced
	Lsr
	LsrReduced
	Mov
	MovReduced
	Movw
	Mul
	Muls
	Mulsu
	Neg
	NegReduced
	Nop
	Or
	OrReduced
	Ori
	Out
	OutReduced
	Pop
	PopReduced
	Push
	PushReduced
	Rcall
	Ret
	Reti
	Rjmp
	Ror
	RorReduced
	Sbc
	SbcReduced
	Sbci
	Sbi
	Sbic
	Sbis
	Sbiw
	Sbrc
	SbrcReduced
	Sbrs
	SbrsReduced
	Sleep
	Spm
	SpmXmega
	StClassic
	StClassicReduced
	StMinimal
	StMinimalReduced
	Std
	Sts
	Sts16
	Sub
	SubReduced
	Subi
	Swap
	SwapReduced
	Wdr
	Xch
	NumMnems
)

//go:generate stringer -type=Mnemonic

type Opcode uint16

func (o Opcode) Nibble3() uint {
	return uint(o >> 12)
}

func (o Opcode) Nibble2() uint {
	return uint((o & 0xf00) >> 8)
}

func (o Opcode) Nibble1() uint {
	return uint((o & 0xf0) >> 4)
}

func (o Opcode) Nibble0() uint {
	return uint(o & 0xf)
}

type Instruction struct {
	Op1, Op2 Opcode
}
