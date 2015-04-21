// Package instr implements instruction decoding for 8-bit Atmel AVR opcodes.
package instr

// A Decoder finds the mnemonic, instruction length, and operands
// for opcodes in a particular instruction set.
type Decoder struct {
	isReduced bool
	set       Set
}

// NewDecoder returns a Decoder for devices other than Reduced Core
// tiny devices.
func NewDecoder(set Set) Decoder {
	return Decoder{isReduced: false, set: set}
}

// NewReducedDecoder returns a Decoder for Reduced Core tiny devices
// (ATtiny4/5/9/10; GCC avrtiny).
func NewReducedDecoder() Decoder {
	return Decoder{isReduced: true, set: setReduced}
}

// DecodeMnem returns the mnemonic and length for an opcode.
func (d Decoder) DecodeMnem(op Opcode) (Mnemonic, int) {
	mnem, ln := d.decodeAnyMnem(op)
	if d.set[mnem] {
		return mnem, ln
	} else {
		return Reserved, 1
	}
}

func (d *Decoder) decodeAnyMnem(op Opcode) (Mnemonic, int) {
	opnibble1 := op.nibble1()
	opnibble0 := op.nibble0()
	// Hand-rolled binary search.
	switch {
	case op < 0x8000:
		switch {
		case op < 0x4000:
			switch {
			case op < 0x2000:
				if op < 0x1000 {
					switch {
					case op < 0x0800:
						if op < 0x0400 {
							switch {
							case op < 0x0200:
								switch {
								case op == 0x0000:
									return Nop, 1
								case op < 0x0100:
									return Reserved, 1
								default:
									return Movw, 1
								}
							case op < 0x0300:
								return Muls, 1
							case op < 0x0380:
								if opnibble0 < 0x8 {
									return Mulsu, 1
								} else {
									return Fmul, 1
								}
							default:
								if opnibble0 < 0x8 {
									return Fmuls, 1
								} else {
									return Fmulsu, 1
								}
							}
						} else {
							return d.decodeCpc(op)
						}
					case op < 0x0c00:
						return d.decodeSbc(op)
					default:
						return d.decodeAdd(op)
					}
				} else {
					switch {
					case op < 0x1800:
						if op < 0x1400 {
							return d.decodeCpse(op)
						} else {
							return d.decodeCp(op)
						}
					case op < 0x1c00:
						return d.decodeSub(op)
					default:
						return d.decodeAdc(op)
					}
				}
			case op < 0x3000:
				switch {
				case op < 0x2800:
					if op < 0x2400 {
						return d.decodeAnd(op)
					} else {
						return d.decodeEor(op)
					}
				case op < 0x2c00:
					return d.decodeOr(op)
				default:
					return d.decodeMov(op)
				}
			default:
				return Cpi, 1
			}
		case op < 0x6000:
			if op < 0x5000 {
				return Sbci, 1
			} else {
				return Subi, 1
			}
		case op < 0x7000:
			return Ori, 1
		default:
			return Andi, 1
		}
	case op < 0xc000:
		switch {
		case op < 0xa000:
			switch {
			case op < 0x9000:
				switch {
				case op < 0x8800:
					switch {
					case op < 0x8400:
						if op < 0x8200 {
							switch {
							case opnibble0 < 0x8:
								if opnibble0 < 0x1 {
									return d.decodeLdMinimal(op)
								} else {
									return Ldd, 1
								}
							case opnibble0 < 0x9:
								return d.decodeLdClassic(op)
							default:
								return Ldd, 1
							}
						} else {
							switch {
							case opnibble0 < 0x8:
								if opnibble0 < 0x1 {
									return d.decodeStMinimal(op)
								} else {
									return Std, 1
								}
							case opnibble0 < 0x9:
								return d.decodeStClassic(op)
							default:
								return Std, 1
							}
						}
					case op < 0x8600:
						return Ldd, 1
					default:
						return Std, 1
					}
				case op < 0x8c00:
					if op < 0x8a00 {
						return Ldd, 1
					} else {
						return Std, 1
					}
				case op < 0x8e00:
					return Ldd, 1
				default:
					return Std, 1
				}
			case op < 0x9800:
				switch {
				case op < 0x9400:
					if op < 0x9200 {
						switch {
						case opnibble0 < 0x8:
							switch {
							case opnibble0 < 0x4:
								switch opnibble0 {
								case 0x0:
									return Lds, 2
								case 0x3:
									return Reserved, 1
								default:
									return d.decodeLdClassic(op)
								}
							case opnibble0 < 0x6:
								return LpmEnhanced, 1
							default:
								return ElpmEnhanced, 1
							}
						case opnibble0 < 0xc:
							if opnibble0 == 0x8 || opnibble0 == 0xb {
								return Reserved, 1
							} else {
								return d.decodeLdClassic(op)
							}
						case opnibble0 < 0xf:
							return d.decodeLdClassic(op)
						default:
							return d.decodePop(op)
						}
					} else {
						switch {
						case opnibble0 < 0x8:
							switch {
							case opnibble0 < 0x4:
								switch opnibble0 {
								case 0x0:
									return Sts, 2
								case 0x3:
									return Reserved, 1
								default:
									return d.decodeStClassic(op)
								}
							case opnibble0 < 0x6:
								if opnibble0 < 5 {
									return Xch, 1
								} else {
									return Las, 1
								}
							case opnibble0 < 0x7:
								return Lac, 1
							default:
								return Lat, 1
							}
						case opnibble0 < 0xf:
							if opnibble0 == 0x8 || opnibble0 == 0xb {
								return Reserved, 1
							} else {
								return d.decodeStClassic(op)
							}
						default:
							return d.decodePush(op)
						}
					}
				case op < 0x9600:
					switch {
					case opnibble0 < 0x8:
						switch {
						case opnibble0 < 0x4:
							switch {
							case opnibble0 < 0x2:
								if opnibble0 < 0x1 {
									return d.decodeCom(op)
								} else {
									return d.decodeNeg(op)
								}
							case opnibble0 < 0x3:
								return d.decodeSwap(op)
							default:
								return d.decodeInc(op)
							}
						case opnibble0 < 0x6:
							if opnibble0 < 0x5 {
								return Reserved, 1
							} else {
								return d.decodeAsr(op)
							}
						case opnibble0 < 0x7:
							return d.decodeLsr(op)
						default:
							return d.decodeRor(op)
						}
					case opnibble0 < 0xc:
						switch {
						case opnibble0 < 0xa:
							if opnibble0 == 0x8 {
								if op < 0x9500 {
									if opnibble1 < 0x8 {
										return Bset, 1
									} else {
										return Bclr, 1
									}
								} else {
									switch opnibble1 {
									case 0x0:
										return Ret, 1
									case 0x1:
										return Reti, 1
									case 0x8:
										return Sleep, 1
									case 0x9:
										return Break, 1
									case 0xa:
										return Wdr, 1
									case 0xc:
										return Lpm, 1
									case 0xd:
										return Elpm, 1
									case 0xe:
										return Spm, 1
									case 0xf:
										return SpmXmega, 1
									default:
										return Reserved, 1
									}
								}
							} else {
								if op < 0x9500 {
									switch opnibble1 {
									case 0x0:
										return Ijmp, 1
									case 0x1:
										return Eijmp, 1
									default:
										return Reserved, 1
									}
								} else {
									switch opnibble1 {
									case 0x0:
										return Icall, 1
									case 0x1:
										return Eicall, 1
									default:
										return Reserved, 1
									}
								}
							}
						case opnibble0 < 0xb:
							return d.decodeDec(op)
						default:
							if op < 0x9500 {
								return Des, 1
							} else {
								return Reserved, 1
							}
						}
					case opnibble0 < 0xe:
						return Jmp, 2
					default:
						return Call, 2
					}
				case op < 0x9700:
					return Adiw, 1
				default:
					return Sbiw, 1
				}
			case op < 0x9c00:
				switch {
				case op < 0x9a00:
					if op < 0x9900 {
						return Cbi, 1
					} else {
						return Sbic, 1
					}
				case op < 0x9b00:
					return Sbi, 1
				default:
					return Sbis, 1
				}
			default:
				return Mul, 1
			}
		case op < 0xb000:
			if d.isReduced {
				if op < 0xa800 {
					return Lds16, 1
				} else {
					return Sts16, 1
				}
			}
			switch {
			case op < 0xa800:
				switch {
				case op < 0xa400:
					if op < 0xa200 {
						return Ldd, 1
					} else {
						return Std, 1
					}
				case op < 0xa600:
					return Ldd, 1
				default:
					return Std, 1
				}
			case op < 0xac00:
				if op < 0xaa00 {
					return Ldd, 1
				} else {
					return Std, 1
				}
			case op < 0xae00:
				return Ldd, 1
			default:
				return Std, 1
			}
		case op < 0xb800:
			return d.decodeIn(op)
		default:
			return d.decodeOut(op)
		}
	case op < 0xe000:
		if op < 0xd000 {
			return Rjmp, 1
		} else {
			return Rcall, 1
		}
	case op < 0xf000:
		return Ldi, 1
	case op < 0xf800:
		if op < 0xf400 {
			return Brbs, 1
		} else {
			return Brbc, 1
		}
	case op < 0xfc00:
		if op < 0xfa00 {
			return d.decodeBld(op)
		} else {
			return d.decodeBst(op)
		}
	case op < 0xfe00:
		return d.decodeSbrc(op)
	default:
		return d.decodeSbrs(op)
	}
}

func (d Decoder) decodeCpc(op Opcode) (Mnemonic, int) {
	if op < 0x0700 {
		return Cpc, 1
	} else {
		return CpcReduced, 1
	}
}

func (d Decoder) decodeSbc(op Opcode) (Mnemonic, int) {
	if op < 0x0b00 {
		return Sbc, 1
	} else {
		return SbcReduced, 1
	}
}

func (d Decoder) decodeAdd(op Opcode) (Mnemonic, int) {
	if op < 0x0f00 {
		return Add, 1
	} else {
		return AddReduced, 1
	}
}

func (d Decoder) decodeCpse(op Opcode) (Mnemonic, int) {
	if op < 0x1300 {
		return Cpse, 1
	} else {
		return CpseReduced, 1
	}
}

func (d Decoder) decodeCp(op Opcode) (Mnemonic, int) {
	if op < 0x1700 {
		return Cp, 1
	} else {
		return CpReduced, 1
	}
}

func (d Decoder) decodeSub(op Opcode) (Mnemonic, int) {
	if op < 0x1b00 {
		return Sub, 1
	} else {
		return SubReduced, 1
	}
}

func (d Decoder) decodeAdc(op Opcode) (Mnemonic, int) {
	if op < 0x1f00 {
		return Adc, 1
	} else {
		return AdcReduced, 1
	}
}

func (d Decoder) decodeAnd(op Opcode) (Mnemonic, int) {
	if op < 0x2300 {
		return And, 1
	} else {
		return AndReduced, 1
	}
}

func (d Decoder) decodeEor(op Opcode) (Mnemonic, int) {
	if op < 0x2700 {
		return Eor, 1
	} else {
		return EorReduced, 1
	}
}

func (d Decoder) decodeOr(op Opcode) (Mnemonic, int) {
	if op < 0x2b00 {
		return Or, 1
	} else {
		return OrReduced, 1
	}
}

func (d Decoder) decodeMov(op Opcode) (Mnemonic, int) {
	if op < 0x2f00 {
		return Mov, 1
	} else {
		return MovReduced, 1
	}
}

func (d Decoder) decodeLdMinimal(op Opcode) (Mnemonic, int) {
	if op < 0x8100 {
		return LdMinimal, 1
	} else {
		return LdMinimalReduced, 1
	}
}

func (d Decoder) decodeStMinimal(op Opcode) (Mnemonic, int) {
	if op < 0x8300 {
		return StMinimal, 1
	} else {
		return StMinimalReduced, 1
	}
}

func (d Decoder) decodeLdClassic(op Opcode) (Mnemonic, int) {
	if ((op & 0x1f0) >> 4) < 16 {
		return Ld, 1
	} else {
		return LdReduced, 1
	}
}

func (d Decoder) decodeStClassic(op Opcode) (Mnemonic, int) {
	if ((op & 0x1f0) >> 4) < 16 {
		return St, 1
	} else {
		return StReduced, 1
	}
}

func (d Decoder) decodePop(op Opcode) (Mnemonic, int) {
	if op < 0x9100 {
		return Pop, 1
	} else {
		return PopReduced, 1
	}
}

func (d Decoder) decodePush(op Opcode) (Mnemonic, int) {
	if op < 0x9300 {
		return Push, 1
	} else {
		return PushReduced, 1
	}
}

func (d Decoder) decodeCom(op Opcode) (Mnemonic, int) {
	if op < 0x9500 {
		return Com, 1
	} else {
		return ComReduced, 1
	}
}

func (d Decoder) decodeNeg(op Opcode) (Mnemonic, int) {
	if op < 0x9500 {
		return Neg, 1
	} else {
		return NegReduced, 1
	}
}

func (d Decoder) decodeSwap(op Opcode) (Mnemonic, int) {
	if op < 0x9500 {
		return Swap, 1
	} else {
		return SwapReduced, 1
	}
}

func (d Decoder) decodeInc(op Opcode) (Mnemonic, int) {
	if op < 0x9500 {
		return Inc, 1
	} else {
		return IncReduced, 1
	}
}

func (d Decoder) decodeAsr(op Opcode) (Mnemonic, int) {
	if op < 0x9500 {
		return Asr, 1
	} else {
		return AsrReduced, 1
	}
}

func (d Decoder) decodeLsr(op Opcode) (Mnemonic, int) {
	if op < 0x9500 {
		return Lsr, 1
	} else {
		return LsrReduced, 1
	}
}

func (d Decoder) decodeRor(op Opcode) (Mnemonic, int) {
	if op < 0x9500 {
		return Ror, 1
	} else {
		return RorReduced, 1
	}
}

func (d Decoder) decodeDec(op Opcode) (Mnemonic, int) {
	if op < 0x9500 {
		return Dec, 1
	} else {
		return DecReduced, 1
	}
}

func (d Decoder) decodeIn(op Opcode) (Mnemonic, int) {
	if (op.nibble2() & 0x1) == 0 {
		return In, 1
	} else {
		return InReduced, 1
	}
}

func (d Decoder) decodeOut(op Opcode) (Mnemonic, int) {
	if (op.nibble2() & 0x1) == 0 {
		return Out, 1
	} else {
		return OutReduced, 1
	}
}

func (d Decoder) decodeBld(op Opcode) (Mnemonic, int) {
	if op.nibble0() > 0x7 {
		return Reserved, 1
	}
	if op < 0xf900 {
		return Bld, 1
	} else {
		return BldReduced, 1
	}
}

func (d Decoder) decodeBst(op Opcode) (Mnemonic, int) {
	if op.nibble0() > 0x7 {
		return Reserved, 1
	}
	if op < 0xfb00 {
		return Bst, 1
	} else {
		return BstReduced, 1
	}
}

func (d Decoder) decodeSbrc(op Opcode) (Mnemonic, int) {
	if op.nibble0() > 0x7 {
		return Reserved, 1
	}
	if op < 0xfd00 {
		return Sbrc, 1
	} else {
		return SbrcReduced, 1
	}
}

func (d Decoder) decodeSbrs(op Opcode) (Mnemonic, int) {
	if op.nibble0() > 0x7 {
		return Reserved, 1
	}
	if op < 0xff00 {
		return Sbrs, 1
	} else {
		return SbrsReduced, 1
	}
}

// DecodeOperands fills in the operands for an instruction. For
// single-opcode instructions, op2 is ignored.
func (d *Decoder) DecodeOperands(ops *Operands, mn Mnemonic, op1, op2 Opcode) {
	mode := opModes[mn]
	decoders[mode](ops, op1, op2)
	ops.Mode = mode
}
