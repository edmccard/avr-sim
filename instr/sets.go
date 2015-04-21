package instr

// A Set records which opcodes are supported by a device or class of
// devices.
type Set []bool

// NewSetMinimal corresponds to GCC -mmcu=avr1
func NewSetMinimal() Set {
	// ATtiny11/12/15/28
	set := make([]bool, NumMnems)
	copy(set, setMinimal)
	return set
}

// NewSetClassic8k corresponds to GCC -mmcu=avr2
func NewSetClassic8k() Set {
	// ATtiny22/26
	set := make([]bool, NumMnems)
	copy(set, setMinimal)
	set[Adiw] = true
	set[Sbiw] = true
	set[Ijmp] = true
	set[Icall] = true
	set[Ld] = true
	set[LdReduced] = true
	set[Ldd] = true
	set[Lds] = true
	set[St] = true
	set[StReduced] = true
	set[Std] = true
	set[Sts] = true
	set[Pop] = true
	set[PopReduced] = true
	set[Push] = true
	set[PushReduced] = true
	return set
}

// NewSetClassic128k corresponds to GCC -mmcu=avr3/avr31
func NewSetClassic128k() Set {
	// ATmega103/603
	set := make([]bool, NumMnems)
	copy(set, NewSetClassic8k())
	set[Jmp] = true
	set[Call] = true
	// TODO: elpm only active when rampz present?
	set[Elpm] = true
	return set
}

// NewSetEnhanced8k corresponds to GCC -mmcu=avr4
func NewSetEnhanced8k() Set {
	// ATmega8/83/85/8515
	set := make([]bool, NumMnems)
	copy(set, NewSetClassic128k())
	set[Mul] = true
	set[Muls] = true
	set[Mulsu] = true
	set[Fmul] = true
	set[Fmuls] = true
	set[Fmulsu] = true
	set[Movw] = true
	set[LpmEnhanced] = true
	set[Spm] = true
	set[Elpm] = false
	return set
}

// NewSetEnhanced128k corresponds to GCC -mmcu=avr5/avr51
func NewSetEnhanced128k() Set {
	set := make([]bool, NumMnems)
	copy(set, NewSetEnhanced8k())
	// TODO: elpm only active when rampz present?
	set[Elpm] = true
	set[ElpmEnhanced] = true
	set[Break] = true
	return set
}

// NewSetEnhanced4m corresponds to GCC -mmcu=avr6
func NewSetEnhanced4m() Set {
	set := make([]bool, NumMnems)
	copy(set, NewSetEnhanced128k())
	// TODO: elpm only active when rampz present?
	set[Elpm] = true
	set[ElpmEnhanced] = true
	set[Eijmp] = true
	set[Eicall] = true
	return set
}

// NewSetXmega corresponds to GCC -mmcu=avrxmega2/4/5/6/7
func NewSetXmega() Set {
	set := make([]bool, NumMnems)
	copy(set, NewSetEnhanced4m())
	// TODO: elpm only active when rampz present?
	set[Elpm] = true
	set[ElpmEnhanced] = true
	set[SpmXmega] = true
	set[Des] = true
	set[Lac] = true
	set[Las] = true
	set[Lat] = true
	set[Xch] = true
	return set
}

var setReduced = []bool{
	true,  // Reserved
	false, // Adc
	true,  // AdcReduced
	false, // Add
	true,  // AddReduced
	false, // Adiw
	false, // And
	true,  // AndReduced
	true,  // Andi
	false, // Asr
	true,  // AsrReduced
	true,  // Bclr
	false, // Bld
	true,  // BldReduced
	true,  // Brbc
	true,  // Brbs
	true,  // Break
	true,  // Bset
	false, // Bst
	true,  // BstReduced
	false, // Call
	true,  // Cbi
	false, // Com
	true,  // ComReduced
	false, // Cp
	true,  // CpReduced
	false, // Cpc
	true,  // CpcReduced
	true,  // Cpi
	false, // Cpse
	true,  // CpseReduced
	false, // Dec
	true,  // DecReduced
	false, // Des
	false, // Eicall
	false, // Eijmp
	false, // Elpm
	false, // ElpmEnhanced
	false, // Eor
	true,  // EorReduced
	false, // Fmul
	false, // Fmuls
	false, // Fmulsu
	true,  // Icall
	true,  // Ijmp
	false, // In
	true,  // InReduced
	false, // Inc
	true,  // IncReduced
	false, // Jmp
	false, // Lac
	false, // Las
	false, // Lat
	false, // Ld
	true,  // LdReduced
	false, // LdMinimal
	true,  // LdMinimalReduced
	false, // Ldd
	true,  // Ldi
	false, // Lds
	true,  // Lds16
	false, // Lpm
	false, // LpmEnhanced
	false, // Lsr
	true,  // LsrReduced
	false, // Mov
	true,  // MovReduced
	false, // Movw
	false, // Mul
	false, // Muls
	false, // Mulsu
	false, // Neg
	true,  // NegReduced
	true,  // Nop
	false, // Or
	true,  // OrReduced
	true,  // Ori
	false, // Out
	true,  // OutReduced
	false, // Pop
	true,  // PopReduced
	false, // Push
	true,  // PushReduced
	true,  // Rcall
	true,  // Ret
	true,  // Reti
	true,  // Rjmp
	false, // Ror
	true,  // RorReduced
	false, // Sbc
	true,  // SbcReduced
	true,  // Sbci
	true,  // Sbi
	true,  // Sbic
	true,  // Sbis
	false, // Sbiw
	false, // Sbrc
	true,  // SbrcReduced
	false, // Sbrs
	true,  // SbrsReduced
	true,  // Sleep
	false, // Spm
	false, // SpmXmega
	false, // St
	true,  // StReduced
	false, // StMinimal
	true,  // StMinimalReduced
	false, // Std
	false, // Sts
	true,  // Sts16
	false, // Sub
	true,  // SubReduced
	true,  // Subi
	false, // Swap
	true,  // SwapReduced
	true,  // Wdr
	false, // Xch
}

var setMinimal = []bool{
	true,  // Reserved
	true,  // Adc
	true,  // AdcReduced
	true,  // Add
	true,  // AddReduced
	false, // Adiw
	true,  // And
	true,  // AndReduced
	true,  // Andi
	true,  // Asr
	true,  // AsrReduced
	true,  // Bclr
	true,  // Bld
	true,  // BldReduced
	true,  // Brbc
	true,  // Brbs
	false, // Break
	true,  // Bset
	true,  // Bst
	true,  // BstReduced
	false, // Call
	true,  // Cbi
	true,  // Com
	true,  // ComReduced
	true,  // Cp
	true,  // CpReduced
	true,  // Cpc
	true,  // CpcReduced
	true,  // Cpi
	true,  // Cpse
	true,  // CpseReduced
	true,  // Dec
	true,  // DecReduced
	false, // Des
	false, // Eicall
	false, // Eijmp
	false, // Elpm
	false, // ElpmEnhanced
	true,  // Eor
	true,  // EorReduced
	false, // Fmul
	false, // Fmuls
	false, // Fmulsu
	false, // Icall
	false, // Ijmp
	true,  // In
	true,  // InReduced
	true,  // Inc
	true,  // IncReduced
	false, // Jmp
	false, // Lac
	false, // Las
	false, // Lat
	false, // Ld
	false, // LdReduced
	true,  // LdMinimal
	true,  // LdMinimalReduced
	false, // Ldd
	true,  // Ldi
	false, // Lds
	false, // Lds16
	true,  // Lpm
	false, // LpmEnhanced
	true,  // Lsr
	true,  // LsrReduced
	true,  // Mov
	true,  // MovReduced
	false, // Movw
	false, // Mul
	false, // Muls
	false, // Mulsu
	true,  // Neg
	true,  // NegReduced
	true,  // Nop
	true,  // Or
	true,  // OrReduced
	true,  // Ori
	true,  // Out
	true,  // OutReduced
	false, // Pop
	false, // PopReduced
	false, // Push
	false, // PushReduced
	true,  // Rcall
	true,  // Ret
	true,  // Reti
	true,  // Rjmp
	true,  // Ror
	true,  // RorReduced
	true,  // Sbc
	true,  // SbcReduced
	true,  // Sbci
	true,  // Sbi
	true,  // Sbic
	true,  // Sbis
	false, // Sbiw
	true,  // Sbrc
	true,  // SbrcReduced
	true,  // Sbrs
	true,  // SbrsReduced
	true,  // Sleep
	false, // Spm
	false, // SpmXmega
	false, // St
	false, // StReduced
	true,  // StMinimal
	true,  // StMinimalReduced
	false, // Std
	false, // Sts
	false, // Sts16
	true,  // Sub
	true,  // SubReduced
	true,  // Subi
	true,  // Swap
	true,  // SwapReduced
	true,  // Wdr
	false, // Xch
}
