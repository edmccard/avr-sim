// generated by stringer -type=AddrMode; DO NOT EDIT

package instr

import "fmt"

const _AddrMode_name = "ModeNoneMode2Reg3Mode2Reg4Mode2Reg5ModeAtomicModeBranchModeDesModeLpmEnhModeInModeIOBitModeLdModeLddModeLdsModeLds16ModeOutModePcModePcOffModeReg5ModeRegBitModeRegImmModeRegPairModePairImmModeSBitModeStModeStdModeStsModeSts16ModeSpmX"

var _AddrMode_index = [...]uint8{0, 8, 17, 26, 35, 45, 55, 62, 72, 78, 87, 93, 100, 107, 116, 123, 129, 138, 146, 156, 166, 177, 188, 196, 202, 209, 216, 225, 233}

func (i AddrMode) String() string {
	if i < 0 || i+1 >= AddrMode(len(_AddrMode_index)) {
		return fmt.Sprintf("AddrMode(%d)", i)
	}
	return _AddrMode_name[_AddrMode_index[i]:_AddrMode_index[i+1]]
}
