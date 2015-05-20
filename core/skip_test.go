package core

import (
	"testing"

	it "github.com/edmccard/avr-sim/instr"
)

func TestSkip(t *testing.T) {
	testop := func(tag string, mnem it.Mnemonic, branches ...branch) {
		if !testing.Short() {
			branches = append([]branch{statusCases}, branches...)
		}
		casetree{t, branches}.run(tcase{tag: tag, mnem: mnem})
	}

	testop("Cpse", it.Cpse, cpseCases)
	testop("Sbic", it.Sbic, sbicCases)
	testop("Sbis", it.Sbis, sbisCases)
	testop("Sbrc", it.Sbrc, sbrcCases)
	testop("Sbrs", it.Sbrs, sbrsCases)
}

var cpseCases = branch{
	{tag: "!=",
		init: cdata{dstreg: 16, dstval: 0x00, srcreg: 17, srcval: 0xff,
			opcds: flash{0: 0x1301, 1: 0x0000}},
		exp: cdata{pc: 0x1}},
	{tag: "= l1",
		init: cdata{dstreg: 16, dstval: 0x01, srcreg: 17, srcval: 0x01,
			opcds: flash{0: 0x1301, 1: 0x0000}},
		exp: cdata{pc: 0x2}},
	{tag: "= l2",
		init: cdata{dstreg: 16, dstval: 0x01, srcreg: 17, srcval: 0x01,
			opcds: flash{0: 0x1301, 1: 0x940e}},
		exp: cdata{pc: 0x3}},
}

var sbrsCases = branch{
	{tag: "b0 clr",
		init: cdata{srcreg: 16, srcval: 0, bit: 0,
			opcds: flash{0: 0xff00, 1: 0x0000}},
		exp: cdata{pc: 0x1}},
	{tag: "b0 set l1",
		init: cdata{srcreg: 16, srcval: 0x1, bit: 0,
			opcds: flash{0: 0xff00, 1: 0x0000}},
		exp: cdata{pc: 0x2}},
	{tag: "b0 set l2",
		init: cdata{srcreg: 16, srcval: 0x1, bit: 0,
			opcds: flash{0: 0xff00, 1: 0x940e}},
		exp: cdata{pc: 0x3}},
}

var sbrcCases = branch{
	{tag: "b0 set",
		init: cdata{srcreg: 16, srcval: 0x1, bit: 0,
			opcds: flash{0: 0xfd00, 1: 0x0000}},
		exp: cdata{pc: 0x1}},
	{tag: "b0 clr l1",
		init: cdata{srcreg: 16, srcval: 0x0, bit: 0,
			opcds: flash{0: 0xfd00, 1: 0x0000}},
		exp: cdata{pc: 0x2}},
	{tag: "b0 clr l2",
		init: cdata{srcreg: 16, srcval: 0x0, bit: 0,
			opcds: flash{0: 0xfd00, 1: 0x940e}},
		exp: cdata{pc: 0x3}},
}

var sbicCases = branch{
	{tag: "b0 set",
		init: cdata{port: 0x10, addr: 0x30, mval: 0x01, bit: 0,
			opcds: flash{0: 0x9980, 1: 0x0000}},
		exp: cdata{pc: 0x1}},
	{tag: "b0 clr l1",
		init: cdata{port: 0x10, addr: 0x30, mval: 0x00, bit: 0,
			opcds: flash{0: 0x9980, 1: 0x0000}},
		exp: cdata{pc: 0x2}},
	{tag: "b0 clr l2",
		init: cdata{port: 0x10, addr: 0x30, mval: 0x00, bit: 0,
			opcds: flash{0: 0x9980, 1: 0x940e}},
		exp: cdata{pc: 0x3}},
}

var sbisCases = branch{
	{tag: "b0 clr",
		init: cdata{port: 0x10, addr: 0x30, mval: 0x00, bit: 0,
			opcds: flash{0: 0x9b80, 1: 0x0000}},
		exp: cdata{pc: 0x1}},
	{tag: "b0 set l1",
		init: cdata{port: 0x10, addr: 0x30, mval: 0x01, bit: 0,
			opcds: flash{0: 0x9b80, 1: 0x0000}},
		exp: cdata{pc: 0x2}},
	{tag: "b0 set l2",
		init: cdata{port: 0x10, addr: 0x30, mval: 0x01, bit: 0,
			opcds: flash{0: 0x9b80, 1: 0x940e}},
		exp: cdata{pc: 0x3}},
}
