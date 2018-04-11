package kman

type Renderer interface {
	Render(Documentation) error
}
