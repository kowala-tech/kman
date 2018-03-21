package kman

type Documenter interface {
	Document() (Documentation, error)
}

// FIXME! Move; default HTML
type Renderer interface {
	Render(...Documentation) error
}
