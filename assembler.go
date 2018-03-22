package kman

type Assembler interface {
	Assemble() ([]Item, error)
}
