package kman

type Documentation struct {
	RootTopic TopicRef
	Glossary  []TermRef
}

type Item struct {
	FileName string
	Line     uint
	Title    string
	Handle   string
	Content  string
}

type itemList []Item

func (r itemList) Len() int      { return len(r) }
func (r itemList) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r itemList) Less(i, j int) bool {
	return r[i].Handle < r[j].Handle
}

type TopicRef struct {
	Item
	Children []TopicRef
}

type TermRef struct {
	Item
}
