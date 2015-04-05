package instr

func (s Minimal) DecodeMnem(op Opcode) Mnemonic {
	switch {
	case op < 0x8000:
		return s.DecodeMnem0000x7fff(op)
	case op >= 0xb000:
		return s.DecodeMnemb000xffff(op)
	case op < 0x9000 || op >= 0xa000:
		return s.DecodeMnemLddStd(op)
	default:
		return s.DecodeMnem9000x9fff(op)
	}
}

var mn0000x7fff = []Mnemonic{
	Illegal, Cpc, Sbc, Add,
	Cpse, Cp, Sub, Adc,
	And, Eor, Or, Mov,
	Cpi, Cpi, Cpi, Cpi,
	Sbci, Sbci, Sbci, Sbci,
	Subi, Subi, Subi, Subi,
	Ori, Ori, Ori, Ori,
	Andi, Andi, Andi, Andi,
}

func (s Minimal) DecodeMnem0000x7fff(op Opcode) Mnemonic {
	mn := mn0000x7fff[(op.Nibble3()*4)+(op.Nibble2()>>2)]
	if mn != Illegal {
		return mn
	}
	if op == 0 {
		return Nop
	} else {
		return Reserved
	}
}

var mnb000xf7ff = []Mnemonic{
	In, In, In, In, In, In, In, In,
	Out, Out, Out, Out, Out, Out, Out, Out,
	Rjmp, Rjmp, Rjmp, Rjmp, Rjmp, Rjmp, Rjmp, Rjmp,
	Rjmp, Rjmp, Rjmp, Rjmp, Rjmp, Rjmp, Rjmp, Rjmp,
	Rcall, Rcall, Rcall, Rcall, Rcall, Rcall, Rcall, Rcall,
	Rcall, Rcall, Rcall, Rcall, Rcall, Rcall, Rcall, Rcall,
	Ldi, Ldi, Ldi, Ldi, Ldi, Ldi, Ldi, Ldi,
	Ldi, Ldi, Ldi, Ldi, Ldi, Ldi, Ldi, Ldi,
	Brbs, Brbs, Brbs, Brbs, Brbc, Brbc, Brbc, Brbc,
	Bld, Bld, Bst, Bst, Sbrc, Sbrc, Sbrs, Sbrs,
}

func (s Minimal) DecodeMnemb000xffff(op Opcode) Mnemonic {
	if op >= 0xf800 && op.Nibble0() >= 0x8 {
		return Reserved
	}
	return mnb000xf7ff[((op.Nibble3()-0xb)*16)+op.Nibble2()]
}

func (s Minimal) DecodeMnemLddStd(op Opcode) Mnemonic {
	if op >= 0x8400 {
		return Reserved
	} else {
		if op.Nibble0() == 0 {
			if op.Nibble2() < 2 {
				return Ld
			} else {
				return St
			}
		} else {
			return Reserved
		}
	}
}

var mn9400x95ff = []Mnemonic{
	Com, Neg, Swap, Inc, Reserved, Asr, Lsr, Ror,
	Illegal, Reserved, Dec, Reserved, Reserved, Reserved, Reserved, Reserved,
}

func (s Minimal) DecodeMnem9000x9fff(op Opcode) Mnemonic {
	switch {
	case op < 0x9400:
		return Reserved
	case op < 0x9600:
		mn := mn9400x95ff[op.Nibble0()]
		if mn != Illegal {
			return mn
		}
		if op < 0x9480 {
			return Bset
		} else if op < 0x9500 {
			return Bclr
		} else if op == 0x9508 {
			return Ret
		} else if op == 0x9518 {
			return Reti
		} else if op == 0x9588 {
			return Sleep
		} else if op == 0x95a8 {
			return Wdr
		} else if op == 0x95c8 {
			return Lpm
		} else {
			return Reserved
		}
	case op < 0x9800:
		return Reserved
	case op < 0x9900:
		return Cbi
	case op < 0x9a00:
		return Sbic
	case op < 0x9b00:
		return Sbi
	case op < 0x9c00:
		return Sbis
	default:
		return Reserved
	}
}

func (s Classic8K) DecodeMnem(op Opcode) Mnemonic {
	switch {
	case op < 0x8000:
		return s.DecodeMnem0000x7fff(op)
	case op >= 0xb000:
		return s.DecodeMnemb000xffff(op)
	case op < 0x9000 || op >= 0xa000:
		return s.DecodeMnemLddStd(op)
	default:
		return s.DecodeMnem9000x9fff(op)
	}
}

func (s Classic8K) DecodeMnemLddStd(op Opcode) Mnemonic {
	if op >= 0x8400 {
		if ((op.Nibble2() >> 1) & 0x1) == 0 {
			return Ldd
		} else {
			return Std
		}
	} else {
		if op.Nibble0() == 0 || op.Nibble0() == 8 {
			if op < 0x8200 {
				return Ld
			} else {
				return St
			}
		} else {
			if op < 0x8200 {
				return Ldd
			} else {
				return Std
			}
		}
	}
}

func (s Classic8K) DecodeMnem9000x9fff(op Opcode) Mnemonic {
	mn := s.Minimal.DecodeMnem9000x9fff(op)
	if mn != Reserved {
		return mn
	}
	switch {
	case op < 0x9200:
		switch op.Nibble0() {
		case 0x0:
			return Lds
		case 0x1, 0x2, 0x9, 0xa, 0xc, 0xd, 0xe:
			return Ld
		case 0xf:
			return Pop
		default:
			return Reserved
		}
	case op < 0x9400:
		switch op.Nibble0() {
		case 0x0:
			return Sts
		case 0x1, 0x2, 0x9, 0xa, 0xc, 0xd, 0xe:
			return St
		case 0xf:
			return Push
		default:
			return Reserved
		}
	case op < 0x9600:
		if op == 0x9409 {
			return Ijmp
		} else if op == 0x9509 {
			return Icall
		} else {
			return Reserved
		}
	case op < 0x9700:
		return Adiw
	case op < 0x9800:
		return Sbiw
	default:
		return Reserved
	}
}

func (s Classic128K) DecodeMnem(op Opcode) Mnemonic {
	switch {
	case op < 0x8000:
		return s.DecodeMnem0000x7fff(op)
	case op >= 0xb000:
		return s.DecodeMnemb000xffff(op)
	case op < 0x9000 || op >= 0xa000:
		return s.DecodeMnemLddStd(op)
	default:
		return s.DecodeMnem9000x9fff(op)
	}
}

func (s Classic128K) DecodeMnem9000x9fff(op Opcode) Mnemonic {
	mn := s.Classic8K.DecodeMnem9000x9fff(op)
	if mn != Reserved {
		return mn
	}
	if op >= 0x9400 && op < 0x9600 {
		switch op.Nibble0() {
		case 0xc, 0xd:
			return Jmp
		case 0xe, 0xf:
			return Call
		default:
			if op == 0x95d8 {
				return Elpm
			} else {
				return Reserved
			}
		}
	} else if op < 0x9200 && (op.Nibble0() == 0x6 || op.Nibble0() == 0x7) {
		return Elpm
	} else {
		return Reserved
	}
}

func (s Enhanced8K) DecodeMnem(op Opcode) Mnemonic {
	switch {
	case op < 0x8000:
		return s.DecodeMnem0000x7fff(op)
	case op >= 0xb000:
		return s.DecodeMnemb000xffff(op)
	case op < 0x9000 || op >= 0xa000:
		return s.DecodeMnemLddStd(op)
	default:
		return s.DecodeMnem9000x9fff(op)
	}
}

func (s Enhanced8K) DecodeMnem0000x7fff(op Opcode) Mnemonic {
	mn := s.Minimal.DecodeMnem0000x7fff(op)
	if mn != Reserved {
		return mn
	}
	switch op.Nibble2() {
	case 0:
		return Reserved
	case 1:
		return Movw
	case 2:
		return Muls
	case 3:
		if op.Nibble1() < 8 {
			if op.Nibble0() < 8 {
				return Mulsu
			} else {
				return Fmul
			}
		} else {
			if op.Nibble0() < 8 {
				return Fmuls
			} else {
				return Fmulsu
			}
		}
	default:
		return Reserved
	}
}

func (s Enhanced8K) DecodeMnem9000x9fff(op Opcode) Mnemonic {
	mn := s.Classic128K.DecodeMnem9000x9fff(op)
	if mn != Reserved {
		return mn
	}
	switch {
	case op >= 0x9c00:
		return Mul
	case op == 0x95e8:
		return Spm
	case op < 0x9200 && (op.Nibble0() == 4 || op.Nibble0() == 5):
		return Lpm
	default:
		return Reserved
	}
}

func (s Enhanced128K) DecodeMnem(op Opcode) Mnemonic {
	switch {
	case op < 0x8000:
		return s.DecodeMnem0000x7fff(op)
	case op >= 0xb000:
		return s.DecodeMnemb000xffff(op)
	case op < 0x9000 || op >= 0xa000:
		return s.DecodeMnemLddStd(op)
	default:
		return s.DecodeMnem9000x9fff(op)
	}
}

func (s Enhanced128K) DecodeMnem9000x9fff(op Opcode) Mnemonic {
	mn := s.Enhanced8K.DecodeMnem9000x9fff(op)
	if mn != Reserved {
		return mn
	}
	if op == 0x9598 {
		return Break
	} else {
		return Reserved
	}
}

func (s Enhanced4M) DecodeMnem(op Opcode) Mnemonic {
	switch {
	case op < 0x8000:
		return s.DecodeMnem0000x7fff(op)
	case op >= 0xb000:
		return s.DecodeMnemb000xffff(op)
	case op < 0x9000 || op >= 0xa000:
		return s.DecodeMnemLddStd(op)
	default:
		return s.DecodeMnem9000x9fff(op)
	}
}

func (s Enhanced4M) DecodeMnem9000x9fff(op Opcode) Mnemonic {
	mn := s.Enhanced128K.DecodeMnem9000x9fff(op)
	if mn != Reserved {
		return mn
	}
	switch op {
	case 0x9419:
		return Eijmp
	case 0x9519:
		return Eicall
	default:
		return Reserved
	}
}

func (s Xmega) DecodeMnem(op Opcode) Mnemonic {
	switch {
	case op < 0x8000:
		return s.DecodeMnem0000x7fff(op)
	case op >= 0xb000:
		return s.DecodeMnemb000xffff(op)
	case op < 0x9000 || op >= 0xa000:
		return s.DecodeMnemLddStd(op)
	default:
		return s.DecodeMnem9000x9fff(op)
	}
}

func (s Xmega) DecodeMnem9000x9fff(op Opcode) Mnemonic {
	mn := s.Enhanced4M.DecodeMnem9000x9fff(op)
	switch {
	case mn != Reserved:
		return mn
	case op >= 0x9200 && op < 0x9400:
		switch op.Nibble0() {
		case 4:
			return Xch
		case 5:
			return Las
		case 6:
			return Lac
		case 7:
			return Lat
		default:
			return Reserved
		}
	case op >= 0x9400 && op < 0x9500 && op.Nibble0() == 0xb:
		return Des
	case op == 0x95f8:
		return Spm
	default:
		return Reserved
	}
}
