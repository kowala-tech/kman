package kman

type documenterDefault struct {
	assemblers []Assembler
	sorter     Sorter
}

func NewDefaultDocumenter(sorter Sorter, assemblers ...Assembler) Documenter {
	return &documenterDefault{
		assemblers: assemblers,
		sorter:     sorter,
	}
}

func (d *documenterDefault) Document() (Documentation, error) {

	items := []Item{}

	for _, a := range d.assemblers {
		assembled, err := a.Assemble()

		if err != nil {
			return Documentation{}, err
		}

		items = append(items, assembled...)
	}

	return d.sorter.Sort(items), nil
}
