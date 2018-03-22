package kman

type mockAssembler struct {
	items []Item
}

func (a *mockAssembler) Assemble() ([]Item, error) {
	return a.items, nil
}
