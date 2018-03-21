package kman

type Itemiser interface {
	Itemise(*[]Item) error
}
