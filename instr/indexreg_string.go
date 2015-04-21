package instr

import "fmt"

const (
	_IndexReg_name_0 = "NoIndexXYZ"
	_IndexReg_name_1 = "X+Y+Z+"
	_IndexReg_name_2 = "-X-Y-Zc"
)

var (
	_IndexReg_index_0 = [...]uint8{0, 7, 8, 9, 10}
	_IndexReg_index_1 = [...]uint8{0, 2, 4, 6}
	_IndexReg_index_2 = [...]uint8{0, 2, 4, 6}
)

// String returns the string representation of an IndexReg. For
// example,
//  ZPostInc.String() == "Z+"
func (i IndexReg) String() string {
	switch {
	case 0 <= i && i <= 3:
		return _IndexReg_name_0[_IndexReg_index_0[i]:_IndexReg_index_0[i+1]]
	case 5 <= i && i <= 7:
		i -= 5
		return _IndexReg_name_1[_IndexReg_index_1[i]:_IndexReg_index_1[i+1]]
	case 9 <= i && i <= 11:
		i -= 9
		return _IndexReg_name_2[_IndexReg_index_2[i]:_IndexReg_index_2[i+1]]
	default:
		return fmt.Sprintf("IndexReg(%d)", i)
	}
}
