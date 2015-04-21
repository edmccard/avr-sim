package instr

// A Mnemonic identifies an instruction, i.e. a particular combination
// of functionality and address mode.
type Mnemonic int

// Not all mnemonics are available in all instruction sets (see type
// Set).
const (
	Reserved         Mnemonic = iota
	Adc                       // Mode2Reg5
	AdcReduced                // Mode2Reg4   registers >= 16
	Add                       // Mode2Reg5
	AddReduced                // Mode2Reg5   registers >= 16
	Adiw                      // ModePairImm
	And                       // Mode2Reg5
	AndReduced                // Mode2Reg5   registers >= 16
	Andi                      // ModeRegImm
	Asr                       // ModeReg5
	AsrReduced                // ModeReg5    registers >= 16
	Bclr                      // ModeSBit
	Bld                       // ModeRegBit
	BldReduced                // ModeRegBit  registers >= 16
	Brbc                      // ModeBranch
	Brbs                      // ModeBranch
	Break                     // ModeNone
	Bset                      // ModeSBit
	Bst                       // ModeRegBit
	BstReduced                // ModeRegBit  registers >= 16
	Call                      // ModePc
	Cbi                       // ModeIOBit
	Com                       // ModeReg5
	ComReduced                // ModeReg5    registers >= 16
	Cp                        // Mode2Reg5
	CpReduced                 // Mode2Reg5   registers >= 16
	Cpc                       // Mode2Reg5
	CpcReduced                // Mode2Reg5   registers >= 16
	Cpi                       // ModeRegImm
	Cpse                      // Mode2Reg5
	CpseReduced               // Mode2Reg5   registers >= 16
	Dec                       // ModeReg5
	DecReduced                // ModeReg5    registers >= 16
	Des                       // ModeDes
	Eicall                    // ModeNone
	Eijmp                     // ModeNone
	Elpm                      // ModeNone
	ElpmEnhanced              // ModeLpmEnh
	Eor                       // Mode2Reg5
	EorReduced                // Mode2Reg5   registers >= 16
	Fmul                      // Mode2Reg3
	Fmuls                     // Mode2Reg3
	Fmulsu                    // Mode2Reg3
	Icall                     // ModeNone
	Ijmp                      // ModeNone
	In                        // ModeIn
	InReduced                 // ModeIn      registers >= 16
	Inc                       // ModeReg5
	IncReduced                // ModeReg5    registers >= 16
	Jmp                       // ModePc
	Lac                       // ModeReg5
	Las                       // ModeReg5
	Lat                       // ModeReg5
	Ld                        // ModeLd
	LdReduced                 // ModeLd      registers >=16
	LdMinimal                 // ModeLd      Z only
	LdMinimalReduced          // ModeLd      registers >= 16, Z only
	Ldd                       // ModeLdd
	Ldi                       // ModeRegImm
	Lds                       // ModeLds
	Lds16                     // ModeLds16
	Lpm                       // ModeNone
	LpmEnhanced               // ModeLpmEnh
	Lsr                       // ModeReg5
	LsrReduced                // ModeReg5    registers >= 16
	Mov                       // Mode2Reg5
	MovReduced                // Mode2Reg5   registers >= 16
	Movw                      // ModeRegPair
	Mul                       // Mode2Reg5
	Muls                      // Mode2Reg4
	Mulsu                     // Mode2Reg3
	Neg                       // ModeReg5
	NegReduced                // ModeReg5    registers >= 16
	Nop                       // ModeNone
	Or                        // Mode2Reg5
	OrReduced                 // Mode2Reg5   registers >= 16
	Ori                       // ModeRegImm
	Out                       // ModeOut
	OutReduced                // ModeOut     registers >= 16
	Pop                       // ModeReg5
	PopReduced                // ModeReg5    registers >= 16
	Push                      // ModeReg5
	PushReduced               // ModeReg5    registers >= 16
	Rcall                     // ModePcOff
	Ret                       // ModeNone
	Reti                      // ModeNone
	Rjmp                      // ModePcOff
	Ror                       // ModeReg5
	RorReduced                // ModeReg5    registers >= 16
	Sbc                       // Mode2Reg5
	SbcReduced                // Mode2Reg5   registers >= 16
	Sbci                      // ModeRegImm
	Sbi                       // ModeIOBit
	Sbic                      // ModeIOBit
	Sbis                      // ModeIOBit
	Sbiw                      // ModePairImm
	Sbrc                      // ModeRegBit
	SbrcReduced               // ModeRegBit  registers >= 16
	Sbrs                      // ModeRegBit
	SbrsReduced               // ModeRegBit  registers >= 16
	Sleep                     // ModeNone
	Spm                       // ModeNone
	SpmXmega                  // ModeSpmX
	St                        // ModeSt
	StReduced                 // ModeSt      registers >= 16
	StMinimal                 // ModeSt      Z only
	StMinimalReduced          // ModeSt      registers >= 16, Z only
	Std                       // ModeStd
	Sts                       // ModeSts
	Sts16                     // ModeSts16
	Sub                       // Mode2Reg5
	SubReduced                // Mode2Reg5   registers >= 16
	Subi                      // ModeRegImm
	Swap                      // ModeReg5
	SwapReduced               // ModeReg5    registers >= 16
	Wdr                       // ModeNone
	Xch                       // ModeReg5
	NumMnems
)

//go:generate stringer -type=Mnemonic

type Opcode uint16

func (o Opcode) nibble3() uint {
	return uint(o >> 12)
}

func (o Opcode) nibble2() uint {
	return uint((o & 0xf00) >> 8)
}

func (o Opcode) nibble1() uint {
	return uint((o & 0xf0) >> 4)
}

func (o Opcode) nibble0() uint {
	return uint(o & 0xf)
}
