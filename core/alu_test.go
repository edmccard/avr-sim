package core

import (
	"fmt"
	"testing"

	it "github.com/edmccard/avr-sim/instr"
)

func TestAlu(t *testing.T) {
	testop := func(tag string, mnem it.Mnemonic, branches ...branch) {
		if !testing.Short() {
			branches = append([]branch{statusCases}, branches...)
		}
		casetree{t, branches}.run(tcase{tag: tag, mnem: mnem})
	}

	testop("Adc", it.Adc, fclr(FlagC), twoRegs(addC0Cases))
	testop("Adc", it.Adc, fclr(FlagC), oneReg(addC0Cases))
	testop("Adc", it.Adc, fset(FlagC), twoRegs(addC1Cases))
	testop("Adc", it.Adc, fset(FlagC), oneReg(addC1Cases))
	testop("Add", it.Add, twoRegs(addC0Cases))
	testop("Add", it.Add, oneReg(addC0Cases))
	testop("And", it.And, twoRegs(andCases))
	testop("And", it.And, oneReg(andCases))
	testop("Andi", it.Andi, immCase, andCases)
	testop("Adiw", it.Adiw, dstPairs, adiwCases)
	testop("Asr", it.Asr, asrCases)
	testop("Com", it.Com, comCases)
	testop("Cp", it.Cp, twoRegs(cpCase(subC0Cases)))
	testop("Cp", it.Cp, oneReg(cpCase(subC0Cases)))
	testop("Cpc", it.Cpc, fclr(FlagC), twoRegs(cpCase(sbcMunge(subC0Cases))))
	testop("Cpc", it.Cpc, fclr(FlagC), oneReg(cpCase(sbcMunge(subC0Cases))))
	testop("Cpc", it.Cpc, fset(FlagC), twoRegs(cpCase(sbcMunge(subC1Cases))))
	testop("Cpc", it.Cpc, fset(FlagC), oneReg(cpCase(sbcMunge(subC1Cases))))
	testop("Cpi", it.Cpi, immCase, cpCase(subC0Cases))
	testop("Dec", it.Dec, decCases)
	testop("Eor", it.Eor, twoRegs(eorCases))
	testop("Eor", it.Eor, oneReg(eorCases))
	testop("Fmul", it.Fmul, twoRegs(fmulCases))
	testop("Fmul", it.Fmul, oneReg(fmulCases))
	testop("Fmuls", it.Fmuls, twoRegs(fmulsCases))
	testop("Fmuls", it.Fmuls, oneReg(fmulsCases))
	testop("Fmulsu", it.Fmulsu, twoRegs(fmulsuCases))
	testop("Fmulsu", it.Fmulsu, oneReg(fmulsuCases))
	testop("Inc", it.Inc, incCases)
	testop("Ldi", it.Ldi, immCase, ldiCases)
	testop("Lsr", it.Lsr, lsrCases)
	testop("Mov", it.Mov, twoRegs(movCases))
	testop("Mov", it.Mov, oneReg(movCases))
	testop("Movw", it.Movw, movwCases)
	testop("Mul", it.Mul, mulRegs(mulCases))
	testop("Muls", it.Muls, mulRegs(mulsCases))
	testop("Mulsu", it.Mulsu, mulRegs(mulsuCases))
	testop("Neg", it.Neg, negCases)
	testop("Or", it.Or, twoRegs(orCases))
	testop("Or", it.Or, oneReg(orCases))
	testop("Ori", it.Ori, immCase, orCases)
	testop("Ror", it.Ror, fclr(FlagC), rorC0Cases)
	testop("Ror", it.Ror, fset(FlagC), rorC1Cases)
	testop("Sbiw", it.Sbiw, dstPairs, sbiwCases)
	testop("Sbc", it.Sbc, fclr(FlagC), twoRegs(sbcMunge(subC0Cases)))
	testop("Sbc", it.Sbc, fclr(FlagC), oneReg(sbcMunge(subC0Cases)))
	testop("Sbc", it.Sbc, fset(FlagC), twoRegs(sbcMunge(subC1Cases)))
	testop("Sbc", it.Sbc, fset(FlagC), oneReg(sbcMunge(subC1Cases)))
	testop("Sbci", it.Sbci, fclr(FlagC), immCase, sbcMunge(subC0Cases))
	testop("Sbci", it.Sbci, fset(FlagC), immCase, sbcMunge(subC1Cases))
	testop("Sub", it.Sub, twoRegs(subC0Cases))
	testop("Sub", it.Sub, oneReg(subC0Cases))
	testop("Subi", it.Subi, immCase, subC0Cases)
	testop("Swap", it.Swap, swapCases)
}

func fset(f Flag) branch {
	return branch{
		{tag: fmt.Sprintf("%s set", f),
			init: cdata{status: flags{f: true}}}}
}

func fclr(f Flag) branch {
	return branch{
		{tag: fmt.Sprintf("%s clr", f),
			init: cdata{status: flags{f: false}}}}
}

func sreg(fs string) flags {
	sr := flags{}
	for _, v := range fs {
		switch v {
		case 'c', 'C':
			sr[FlagC] = (v == 'C')
		case 'z', 'Z':
			sr[FlagZ] = (v == 'Z')
		case 'n', 'N':
			sr[FlagN] = (v == 'N')
		case 'v', 'V':
			sr[FlagV] = (v == 'V')
		case 's', 'S':
			sr[FlagS] = (v == 'S')
		case 'h', 'H':
			sr[FlagH] = (v == 'H')
		case 't', 'T':
			sr[FlagT] = (v == 'T')
		case 'i', 'I':
			sr[FlagI] = (v == 'I')
		}
	}
	return sr
}

func arithCase(tag string, v1, v2, res int, fs string) tcase {
	return tcase{tag: tag,
		init: cdata{dstval: v1, srcval: v2},
		exp:  cdata{dstval: res, status: sreg(fs)}}
}

func rmwCase(tag string, v1, res int, fs string) tcase {
	return tcase{tag: tag,
		init: cdata{dstreg: 16, dstval: v1},
		exp:  cdata{dstval: res, status: sreg(fs)}}
}

func mulRegs(b branch) (cases branch) {
	for _, c := range b {
		if c.init[dstval] == c.init[srcval] {
			cases = append(cases, c.merge(
				tcase{tag: "d0 r0", init: cdata{dstreg: 0, srcreg: 0}}))
			cases = append(cases, c.merge(
				tcase{tag: "d1 r1", init: cdata{dstreg: 1, srcreg: 1}}))
		}
		cases = append(cases, c.merge(
			tcase{tag: "d0 r1", init: cdata{dstreg: 0, srcreg: 1}}))
		cases = append(cases, c.merge(
			tcase{tag: "d1 r0", init: cdata{dstreg: 1, srcreg: 0}}))
		cases = append(cases, c.merge(
			tcase{tag: "d1 r2", init: cdata{dstreg: 1, srcreg: 2}}))
		cases = append(cases, c.merge(
			tcase{tag: "d2 r3", init: cdata{dstreg: 2, srcreg: 3}}))
	}
	return
}

func mulCase(tag string, v1, v2, res int, fs string) tcase {
	return tcase{tag: tag,
		init: cdata{dstval: v1, srcval: v2},
		exp:  cdata{status: sreg(fs), mulval: res}}
}

var mulCases = branch{
	mulCase("s00", 0xff, 0x01, 0x00ff, "cz"),
	mulCase("s00", 0x7f, 0x7f, 0x3f01, "cz"),
	mulCase("s01", 0xff, 0xff, 0xfe01, "Cz"),
	mulCase("s02", 0xff, 0x00, 0x0000, "cZ"),
}

var mulsCases = branch{
	mulCase("s00", 0xff, 0xff, 0x0001, "cz"),
	mulCase("s00", 0x7f, 0x7f, 0x3f01, "cz"),
	mulCase("s01", 0xff, 0x01, 0xffff, "Cz"),
	mulCase("s02", 0xff, 0x00, 0x0000, "cZ"),
}

var mulsuCases = branch{
	mulCase("s00", 0x01, 0xff, 0x00ff, "cz"),
	mulCase("s00", 0x7f, 0x7f, 0x3f01, "cz"),
	mulCase("s01", 0xff, 0xff, 0xff01, "Cz"),
	mulCase("s02", 0xff, 0x00, 0x0000, "cZ"),
}

var fmulCases = branch{
	mulCase("s00", 0xff, 0x01, 0x01fe, "cz"),
	mulCase("s00", 0x80, 0x80, 0x8000, "cz"),
	mulCase("s01", 0xd0, 0xd0, 0x5200, "Cz"),
	mulCase("s01", 0xe0, 0xe0, 0x8800, "Cz"),
	mulCase("s02", 0xff, 0x00, 0x0000, "cZ"),
}

var fmulsCases = branch{
	mulCase("s00", 0x7f, 0x7f, 0x7e02, "cz"),
	mulCase("s00", 0x80, 0x80, 0x8000, "cz"),
	mulCase("s01", 0xff, 0x01, 0xfffe, "Cz"),
	mulCase("s02", 0xff, 0x00, 0x0000, "cZ"),
}

var fmulsuCases = branch{
	mulCase("s00", 0x01, 0xff, 0x01fe, "cz"),
	mulCase("s00", 0x7f, 0xc8, 0xc670, "cz"),
	mulCase("s01", 0xff, 0xff, 0xfe02, "Cz"),
	mulCase("s01", 0x9c, 0xaa, 0x7b30, "Cz"),
	mulCase("s02", 0xff, 0x00, 0x0000, "cZ"),
}

var asrCases = branch{
	rmwCase("s00", 0x02, 0x01, "cznvs"),
	rmwCase("s02", 0x00, 0x00, "cZnvs"),
	rmwCase("s0c", 0x80, 0xc0, "czNVs"),
	rmwCase("s15", 0x81, 0xc0, "CzNvS"),
	rmwCase("s19", 0x03, 0x01, "CznVS"),
	rmwCase("s1b", 0x01, 0x00, "CZnVS"),
}

var lsrCases = branch{
	rmwCase("s00", 0x02, 0x01, "cznvs"),
	rmwCase("s02", 0x00, 0x00, "cZnvs"),
	rmwCase("s19", 0x03, 0x01, "CznVS"),
	rmwCase("s1b", 0x01, 0x00, "CZnVS"),
}

var comCases = branch{
	rmwCase("s01", 0x80, 0x7f, "Cznvs"),
	rmwCase("s03", 0xff, 0x00, "CZnvs"),
	rmwCase("s15", 0x00, 0xff, "CzNvS"),
}

var negCases = branch{
	rmwCase("s01", 0x90, 0x70, "Cznvsh"),
	rmwCase("s02", 0x00, 0x00, "cZnvsh"),
	rmwCase("s0d", 0x80, 0x80, "CzNVsh"),
	rmwCase("s15", 0x10, 0xf0, "CzNvSh"),
	rmwCase("s21", 0x81, 0x7f, "CznvsH"),
	rmwCase("s35", 0x01, 0xff, "CzNvSH"),
}

var swapCases = branch{
	rmwCase("", 0xff, 0xff, ""),
	rmwCase("", 0x00, 0x00, ""),
	rmwCase("", 0x12, 0x21, ""),
}

var decCases = branch{
	rmwCase("s00", 0x02, 0x01, "znvs"),
	rmwCase("s02", 0x01, 0x00, "Znvs"),
	rmwCase("s14", 0x00, 0xff, "zNvS"),
	rmwCase("s18", 0x80, 0x7f, "znVS"),
}

var incCases = branch{
	rmwCase("s00", 0x00, 0x01, "znvs"),
	rmwCase("s02", 0xff, 0x00, "Znvs"),
	rmwCase("s0c", 0x7f, 0x80, "zNVs"),
	rmwCase("s14", 0x80, 0x81, "zNvS"),
}

var rorC0Cases = branch{
	rmwCase("s00", 0x02, 0x01, "cznvs"),
	rmwCase("s02", 0x00, 0x00, "cZnvs"),
	rmwCase("s19", 0x03, 0x01, "CznVS"),
	rmwCase("s1b", 0x01, 0x00, "CZnVS"),
}

var rorC1Cases = branch{
	rmwCase("s0c", 0x00, 0x80, "czNVs"),
	rmwCase("s15", 0x01, 0x80, "CzNvS"),
}

var andCases = branch{
	arithCase("s00", 0x01, 0x01, 0x01, "znvs"),
	arithCase("s02", 0xaa, 0x55, 0x00, "Znvs"),
	arithCase("s14", 0x80, 0x80, 0x80, "zNvS"),
}

var orCases = branch{
	arithCase("s00", 0x01, 0x03, 0x03, "znvs"),
	arithCase("s02", 0x00, 0x00, 0x00, "Znvs"),
	arithCase("s14", 0x80, 0x01, 0x81, "zNvS"),
}

var eorCases = branch{
	arithCase("s00", 0x01, 0x03, 0x02, "znvs"),
	arithCase("s02", 0xaa, 0xaa, 0x00, "Znvs"),
	arithCase("s14", 0xaa, 0x55, 0xff, "zNvS"),
}

var addC0Cases = branch{
	arithCase("s00", 0x01, 0x01, 0x02, "cznvsh"),
	arithCase("s01", 0x10, 0xf1, 0x01, "Czvnsh"),
	arithCase("s02", 0x00, 0x00, 0x00, "cZnvsh"),
	arithCase("s03", 0x10, 0xf0, 0x00, "CZnvsh"),
	arithCase("s0c", 0x40, 0x40, 0x80, "czNVsh"),
	arithCase("s14", 0x00, 0x80, 0x80, "czNvSh"),
	arithCase("s15", 0xc0, 0xc0, 0x80, "CzNvSh"),
	arithCase("s19", 0x81, 0x81, 0x02, "CznVSh"),
	arithCase("s1b", 0x80, 0x80, 0x00, "CZnVSh"),
	arithCase("s20", 0x08, 0x08, 0x10, "cznvsH"),
	arithCase("s21", 0x02, 0xff, 0x01, "CznvsH"),
	arithCase("s23", 0x01, 0xff, 0x00, "CZnvsH"),
	arithCase("s2c", 0x48, 0x48, 0x90, "czNVsH"),
	arithCase("s34", 0x01, 0x8f, 0x90, "czNvSH"),
	arithCase("s35", 0xc8, 0xc8, 0x90, "CzNvSH"),
	arithCase("s39", 0x88, 0x88, 0x10, "CznVSH"),
}

var addC1Cases = branch{
	arithCase("s00", 0x00, 0x00, 0x01, "czvnsh"),
	arithCase("s01", 0x10, 0xf0, 0x01, "Czvnsh"),
	arithCase("s0c", 0x40, 0x40, 0x81, "czVNsh"),
	arithCase("s14", 0x00, 0x80, 0x81, "czNvSh"),
	arithCase("s15", 0xc0, 0xc0, 0x81, "CzNvSh"),
	arithCase("s19", 0x80, 0x80, 0x01, "CznVSh"),
	arithCase("s20", 0x08, 0x08, 0x11, "czvnsH"),
	arithCase("s21", 0x01, 0xff, 0x01, "CzvnsH"),
	arithCase("s2c", 0x48, 0x48, 0x91, "czNVsH"),
	arithCase("s34", 0x00, 0x8f, 0x90, "czNvSH"),
	arithCase("s35", 0xc8, 0xc8, 0x91, "CzNvSH"),
	arithCase("s39", 0x88, 0x88, 0x11, "CznVSH"),
}

var subC0Cases = branch{
	arithCase("s00", 0x01, 0x00, 0x01, "czvnsh"),
	arithCase("s01", 0x00, 0x90, 0x70, "Czvnsh"),
	arithCase("s02", 0x00, 0x00, 0x00, "cZvnsh"),
	arithCase("s0d", 0x00, 0x80, 0x80, "CzNVsh"),
	arithCase("s14", 0x80, 0x00, 0x80, "czNvSh"),
	arithCase("s15", 0x00, 0x10, 0xf0, "CzNvSh"),
	arithCase("s18", 0x80, 0x10, 0x70, "cznVSh"),
	arithCase("s20", 0x10, 0x01, 0x0f, "czvnsH"),
	arithCase("s21", 0x00, 0x81, 0x7f, "CzvnsH"),
	arithCase("s2d", 0x10, 0x81, 0x8f, "CzNVsH"),
	arithCase("s34", 0x90, 0x01, 0x8f, "czNvSH"),
	arithCase("s35", 0x00, 0x01, 0xff, "CzNvSH"),
	arithCase("s38", 0x80, 0x01, 0x7f, "cznVSH"),
}

var subC1Cases = branch{
	arithCase("s00", 0x02, 0x00, 0x01, "czvnsh"),
	arithCase("s01", 0x01, 0x90, 0x70, "Czvnsh"),
	arithCase("s02", 0x01, 0x00, 0x00, "cZvnsh"),
	arithCase("s0d", 0x01, 0x80, 0x80, "CzNVsh"),
	arithCase("s14", 0x81, 0x00, 0x80, "czNvSh"),
	arithCase("s15", 0x01, 0x10, 0xf0, "CzNvSh"),
	arithCase("s18", 0x81, 0x10, 0x70, "cznVSh"),
	arithCase("s20", 0x10, 0x00, 0x0f, "czvnsH"),
	arithCase("s21", 0x00, 0x80, 0x7f, "CzvnsH"),
	arithCase("s22", 0x10, 0x0f, 0x00, "cZvnsH"),
	arithCase("s23", 0x00, 0xff, 0x00, "CZvnsH"),
	arithCase("s2d", 0x10, 0x80, 0x8f, "CzNVsH"),
	arithCase("s34", 0x90, 0x00, 0x8f, "czNvSH"),
	arithCase("s35", 0x00, 0x00, 0xff, "CzNvSH"),
	arithCase("s38", 0x80, 0x00, 0x7f, "cznVSH"),
	arithCase("s3a", 0x80, 0x7f, 0x00, "cZnVSH"),
}

func sbcMunge(b branch) (cases branch) {
	for _, c := range b {
		if c.exp[dstval] == 0 {
			cases = append(cases, c.merge(tcase{
				init: cdata{status: flags{FlagZ: true}},
				exp:  cdata{status: flags{FlagZ: true}}}))
			cases = append(cases, c.merge(tcase{
				init: cdata{status: flags{FlagZ: false}},
				exp:  cdata{status: flags{FlagZ: false}}}))
		} else {
			cases = append(cases, c)
		}
	}
	return cases
}

func cpCase(b branch) (cases branch) {
	for _, c := range b {
		cases = append(cases, c.merge(tcase{
			exp: cdata{dstval: c.init[dstval]}}))
	}
	return
}

func twoRegs(b branch) (cases branch) {
	for _, c := range b {
		cases = append(cases,
			c.merge(tcase{tag: "d16 r17",
				init: cdata{dstreg: 16, srcreg: 17}}))
	}
	return
}

func oneReg(b branch) (cases branch) {
	for _, c := range b {
		if c.init[dstval] == c.init[srcval] {
			cases = append(cases,
				c.merge(tcase{tag: "d16 r16",
					init: cdata{dstreg: 16, srcreg: 16}}))
		}
	}
	return
}

var immCase = branch{
	{tag: "imm", init: cdata{dstreg: 16, srcreg: imm}},
}

var dstPairs = branch{
	{tag: "d16,d17", init: cdata{dstreg: pair{17, 16}}},
}

var adiwCases = branch{
	{tag: "s00",
		init: cdata{srcreg: imm, srcval: 0x01, dstval: pair{0, 0}},
		exp:  cdata{dstval: pair{0, 1}, status: sreg("cznvs")}},
	{tag: "s01",
		init: cdata{srcreg: imm, srcval: 0x3e, dstval: pair{0xff, 0xc3}},
		exp:  cdata{dstval: pair{0, 1}, status: sreg("Cznvs")}},
	{tag: "s02",
		init: cdata{srcreg: imm, srcval: 0x00, dstval: pair{0, 0}},
		exp:  cdata{dstval: pair{0, 0}, status: sreg("cZnvs")}},
	{tag: "s03",
		init: cdata{srcreg: imm, srcval: 0x3e, dstval: pair{0xff, 0xc2}},
		exp:  cdata{dstval: pair{0, 0}, status: sreg("CZnvs")}},
	{tag: "s0c",
		init: cdata{srcreg: imm, srcval: 0x3e, dstval: pair{0x7f, 0xc2}},
		exp:  cdata{dstval: pair{0x80, 0}, status: sreg("czNVs")}},
	{tag: "s14",
		init: cdata{srcreg: imm, srcval: 0x00, dstval: pair{0x80, 0}},
		exp:  cdata{dstval: pair{0x80, 0}, status: sreg("czNvS")}},
}

var sbiwCases = branch{
	{tag: "s00",
		init: cdata{srcreg: imm, srcval: 0x00, dstval: pair{0, 1}},
		exp:  cdata{dstval: pair{0, 1}, status: sreg("cznvs")}},
	{tag: "s02",
		init: cdata{srcreg: imm, srcval: 0x00, dstval: pair{0, 0}},
		exp:  cdata{dstval: pair{0, 0}, status: sreg("cZnvs")}},
	{tag: "s14",
		init: cdata{srcreg: imm, srcval: 0x00, dstval: pair{0x80, 0}},
		exp:  cdata{dstval: pair{0x80, 0}, status: sreg("czNvS")}},
	{tag: "s15",
		init: cdata{srcreg: imm, srcval: 0x01, dstval: pair{0, 0}},
		exp:  cdata{dstval: pair{0xff, 0xff}, status: sreg("CzNvS")}},
	{tag: "s18",
		init: cdata{srcreg: imm, srcval: 0x01, dstval: pair{0x80, 0}},
		exp:  cdata{dstval: pair{0x7f, 0xff}, status: sreg("cznVS")}},
}

var movwCases = branch{
	{tag: "d17:16 r19:18",
		init: cdata{dstreg: pair{17, 16}, dstval: pair{0, 0},
			srcreg: pair{19, 18}, srcval: pair{0x12, 0x34}},
		exp: cdata{dstval: pair{0x12, 0x34}}},
	{tag: "d17:16 r17:16",
		init: cdata{dstreg: pair{17, 16}, dstval: pair{0x43, 0x21},
			srcreg: pair{17, 16}, srcval: pair{0x43, 0x21}},
		exp: cdata{dstval: pair{0x43, 0x21}}},
}

var movCases = branch{
	arithCase("", 0x00, 0x10, 0x10, ""),
	arithCase("", 0x10, 0x10, 0x10, ""),
}

var ldiCases = branch{
	arithCase("", 0x00, 0xff, 0xff, ""),
	arithCase("", 0xff, 0x00, 0x00, ""),
}
