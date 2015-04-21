package instr

import "fmt"

func Example() {
	var op1, op2 Opcode = 0x1234, 0
	ops := Operands{}

	decoder := NewDecoder(NewSetMinimal())
	mnemonic, length := decoder.DecodeMnem(op1)
	decoder.DecodeOperands(&ops, mnemonic, op1, op2)

	fmt.Println(mnemonic, length)
	fmt.Println(ops.Mode, ops)
	// Output:
	// Cpse 1
	// Mode2Reg5 r3, r20
}
