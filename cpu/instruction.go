package cpu

type InstType int

const (
	Adc InstType = iota
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

//go:generate stringer -type=InstType

type Instruction struct {
	itype    InstType
	op1, op2 Opcode
}
