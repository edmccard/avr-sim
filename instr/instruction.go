package instr

type Mnemonic int

const (
	Adc Mnemonic = iota
	Add
	Adiw
	And
	Andi
	Asr
	Bclr
	Bld
	Brbc
	Brbs
	Break
	Bset
	Bst
	Call
	Cbi
	Com
	Cp
	Cpc
	Cpi
	Cpse
	Dec
	Des
	Eicall
	Eijmp
	Elpm
	Eor
	Fmul
	Fmuls
	Fmulsu
	Icall
	Ijmp
	In
	Inc
	Jmp
	Lac
	Las
	Lat
	Ld
	Ldd
	Ldi
	Lds
	Lpm
	Lsr
	Mov
	Movw
	Mul
	Muls
	Mulsu
	Neg
	Nop
	Or
	Ori
	Out
	Pop
	Push
	Rcall
	Ret
	Reti
	Rjmp
	Ror
	Sbc
	Sbci
	Sbi
	Sbic
	Sbis
	Sbiw
	Sbrc
	Sbrs
	Sleep
	Spm
	St
	Std
	Sts
	Sub
	Subi
	Swap
	Wdr
	Xch
	Reserved
	Illegal
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
