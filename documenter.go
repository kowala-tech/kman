package kman

type Documenter interface {
	Document() (Documentation, error)
}
