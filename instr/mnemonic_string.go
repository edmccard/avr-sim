// generated by stringer -type=Mnemonic; DO NOT EDIT

package instr

import "fmt"

const _Mnemonic_name = "AdcAddAdiwAndAndiAsrBclrBldBrbcBrbsBreakBsetBstCallCbiComCpCpcCpiCpseDecDesEicallEijmpElpmEorFmulFmulsFmulsuIcallIjmpInIncJmpLacLasLatLdLddLdiLdsLpmLsrMovMovwMulMulsMulsuNegNopOrOriOutPopPushRcallRetRetiRjmpRorSbcSbciSbiSbicSbisSbiwSbrcSbrsSleepSpmStStdStsSubSubiSwapWdrXchReservedIllegal"

var _Mnemonic_index = [...]uint16{0, 3, 6, 10, 13, 17, 20, 24, 27, 31, 35, 40, 44, 47, 51, 54, 57, 59, 62, 65, 69, 72, 75, 81, 86, 90, 93, 97, 102, 108, 113, 117, 119, 122, 125, 128, 131, 134, 136, 139, 142, 145, 148, 151, 154, 158, 161, 165, 170, 173, 176, 178, 181, 184, 187, 191, 196, 199, 203, 207, 210, 213, 217, 220, 224, 228, 232, 236, 240, 245, 248, 250, 253, 256, 259, 263, 267, 270, 273, 281, 288}

func (i Mnemonic) String() string {
	if i < 0 || i+1 >= Mnemonic(len(_Mnemonic_index)) {
		return fmt.Sprintf("Mnemonic(%d)", i)
	}
	return _Mnemonic_name[_Mnemonic_index[i]:_Mnemonic_index[i+1]]
}
