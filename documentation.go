package kman

import (
	"html/template"

	"github.com/russross/blackfriday"
)

type Documentation struct {
	RootTopic TopicRef
	Glossary  []TermRef
}

//go:generate stringer -type=ItemType
type ItemType int

const (
	ItemTypeTopic ItemType = iota
	ItemTypeTerm
)

type Item struct {
	Type     ItemType
	FileName string
	Line     uint
	Title    string
	Handle   string
	Content  string
}

func (i Item) HTML() template.HTML {
	return template.HTML(blackfriday.Run([]byte(i.Content)))
}

type itemListHandleSorter []Item

func (r itemListHandleSorter) Len() int      { return len(r) }
func (r itemListHandleSorter) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r itemListHandleSorter) Less(i, j int) bool {
	return r[i].Handle < r[j].Handle
}

type itemListTitleSorter []Item

func (r itemListTitleSorter) Len() int      { return len(r) }
func (r itemListTitleSorter) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r itemListTitleSorter) Less(i, j int) bool {
	return r[i].Title < r[j].Title
}

type TopicRef struct {
	Item
	Children []TopicRef
}

type TermRef struct {
	Item
}
