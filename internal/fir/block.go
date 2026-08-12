package fir

type BlockID uint32

type Block struct {
	ID   BlockID
	Name string

	Ins  []Instruction
	Term Terminator
}
