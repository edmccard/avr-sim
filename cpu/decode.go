package cpu

// NOTE: No ATtiny10

type Opcode uint16

func (o Opcode) Nibble3() Opcode {
	return o >> 12
}

func (o Opcode) Nibble2() Opcode {
	return (o >> 8) & 0xf
}

func (o Opcode) Nibble1() Opcode {
	return (o >> 4) & 0xf
}

func (o Opcode) Nibble0() Opcode {
	return o & 0xf
}

func Decode(opcode Opcode) InstType {
	opNibble3 := opcode.Nibble3()
	opNibble2 := opcode.Nibble2()
	opNibble1 := opcode.Nibble1()
	opNibble0 := opcode.Nibble0()
	switch opNibble3 {
	case 0x0:
		switch opNibble2 >> 2 {
		case 0:
			switch opNibble2 {
			case 0:
				if opcode == 0 {
					return Nop
				} else {
					return Reserved
				}
			case 1:
				return Movw
			case 2:
				return Muls
			case 3:
				if opNibble1 < 0x8 {
					if opNibble0 < 0x8 {
						return Mulsu
					} else {
						return Fmul
					}
				} else {
					if opNibble0 < 0x8 {
						return Fmuls
					} else {
						return Fmulsu
					}
				}
			}
		case 1:
			return Cpc
		case 2:
			return Sbc
		case 3:
			return Add
		}
	case 0x1:
		switch opNibble2 >> 2 {
		case 0:
			return Cpse
		case 1:
			return Cp
		case 2:
			return Sub
		case 3:
			return Adc
		}
	case 0x2:
		switch opNibble2 >> 2 {
		case 0:
			return And
		case 1:
			return Eor
		case 2:
			return Or
		case 3:
			return Mov
		}
	case 0x3:
		return Cpi
	case 0x4:
		return Sbci
	case 0x5:
		return Subi
	case 0x6:
		return Ori
	case 0x7:
		return Andi
	case 0x8:
		switch (opNibble2 >> 1) & 1 {
		case 0:
			if opNibble2 < 4 {
				if opNibble0 == (opNibble0 & 8) {
					return Ld
				} else {
					return Ldd
				}
			} else {
				return Ldd
			}
		case 1:
			if opNibble2 < 0x4 {
				if opNibble0 == (opNibble0 & 8) {
					return St
				} else {
					return Std
				}
			} else {
				return Std
			}
		}
	case 0x9:
		switch opNibble2 >> 1 {
		case 0:
			switch opNibble0 {
			case 0x0:
				return Lds
			case 0x3, 0x8, 0xb:
				return Reserved
			case 0x4, 0x5:
				return Lpm
			case 0x6, 0x7:
				return Elpm
			case 0xf:
				return Pop
			default:
				return Ld
			}
		case 1:
			switch opNibble0 {
			case 0x0:
				return Sts
			case 0x3, 0x8, 0xb:
				return Reserved
			case 0x4:
				return Xch
			case 0xf:
				return Push
			case 0x5:
				return Las
			case 0x6:
				return Lac
			case 0x7:
				return Lat
			default:
				return St
			}
		case 2:
			switch opNibble0 {
			case 0x0:
				return Com
			case 0x1:
				return Neg
			case 0x2:
				return Swap
			case 0x3:
				return Inc
			case 0x4:
				return Reserved
			case 0x5:
				return Asr
			case 0x6:
				return Lsr
			case 0x7:
				return Ror
			case 0x8:
				if opNibble2 == 0x4 {
					if opNibble1 < 0x8 {
						return Bset
					} else {
						return Bclr
					}
				} else {
					switch opNibble1 {
					case 0x0:
						return Ret
					case 0x1:
						return Reti
					case 0x8:
						return Sleep
					case 0x9:
						return Break
					case 0xa:
						return Wdr
					case 0xc:
						return Lpm
					case 0xd:
						return Elpm
					case 0xe:
						return Spm
					case 0xf:
						return Spm2
					default:
						return Reserved
					}
				}
			case 0x9:
				if opNibble2 == 0x4 {
					switch opNibble1 {
					case 0x0:
						return Ijmp
					case 0x1:
						return Eijmp
					default:
						return Reserved
					}
				} else {
					switch opNibble1 {
					case 0x0:
						return Icall
					case 0x1:
						return Eicall
					default:
						return Reserved
					}
				}
			case 0xa:
				return Dec
			case 0xb:
				if opNibble2 == 0x4 {
					return Des
				} else {
					return Reserved
				}
			case 0xc, 0xd:
				return Jmp
			case 0xe, 0xf:
				return Call
			}
		case 3:
			if opNibble2 == 0x6 {
				return Adiw
			} else {
				return Sbiw
			}
		case 4:
			if opNibble2 == 0x8 {
				return Cbi
			} else {
				return Sbic
			}
		case 5:
			if opNibble2 == 0xa {
				return Sbi
			} else {
				return Sbis
			}
		case 6, 7:
			return Mul
		}
	case 0xa:
		switch (opNibble2 >> 1) & 1 {
		case 0:
			return Ldd
		case 1:
			return Std
		}
	case 0xb:
		switch opNibble2 >> 3 {
		case 0:
			return In
		case 1:
			return Out
		}
	case 0xc:
		return Rjmp
	case 0xd:
		return Rcall
	case 0xe:
		return Ldi
	case 0xf:
		if opNibble2 < 0x4 {
			return Brbs
		} else if opNibble2 < 0x8 {
			return Brbc
		} else {
			if opNibble0 < 0x8 {
				switch (opNibble2 & 7) >> 1 {
				case 0:
					return Bld
				case 1:
					return Bst
				case 2:
					return Sbrc
				case 3:
					return Sbrs
				}
			} else {
				return Reserved
			}
		}
	}
	return -1
}
