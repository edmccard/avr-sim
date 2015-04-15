package cpu

import "fmt"

const _Flag_name = "CZNVSHTI"

func (i Flag) String() string {
	if i < 0 || i > 7 {
		return fmt.Sprintf("Flag(%d)", i)
	}
	return fmt.Sprintf("%c", _Flag_name[i])
}
