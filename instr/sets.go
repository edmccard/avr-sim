package instr

type Set interface {
	DecodeMnem(Opcode) Mnemonic
	DecodeAddr(Instruction) AddrMode
	String() string
}

type Minimal struct{}

type Classic8K struct {
	Minimal
}

type Classic128K struct {
	Classic8K
}

type Enhanced8K struct {
	Classic128K
}

type Enhanced128K struct {
	Enhanced8K
}

type Enhanced4M struct {
	Enhanced128K
}

type Xmega struct {
	Enhanced4M
}

func (s Minimal) String() string {
	return "Minimal"
}

func (s Classic8K) String() string {
	return "Classic8K"
}

func (s Classic128K) String() string {
	return "Classic128K"
}

func (s Enhanced8K) String() string {
	return "Enhanced8K"
}

func (s Enhanced128K) String() string {
	return "Enhanced128K"
}

func (s Enhanced4M) String() string {
	return "Enhanced4M"
}

func (s Xmega) String() string {
	return "Xmega"
}
