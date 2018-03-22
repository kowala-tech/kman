package kman

import (
	"sort"
	"strings"
)

type Sorter interface {
	Sort([]Item) Documentation
}

type sorter struct{}

func NewDefaultSorter() Sorter {
	return &sorter{}
}

func (s *sorter) Sort(input []Item) Documentation {

	doc := Documentation{}

	topicItems, termItems := []Item{}, []Item{}

	for _, i := range input {

		switch i.Type {

		case ItemTypeTopic:
			topicItems = append(topicItems, i)

		case ItemTypeTerm:
			termItems = append(termItems, i)
		}
	}

	doc.RootTopic = s.sortItemsToTopicTree(topicItems)
	doc.Glossary = s.sortItemsToGlossary(termItems)

	return doc
}

func (s *sorter) sortItemsToTopicTree(items []Item) (root TopicRef) {

	if len(items) == 0 {
		return
	}

	shortestIndex, shortestHandle := 0, 999
	foundIndex := -1

	// Step one, pick the root item
	for i := 0; i < len(items); i++ {

		if short := len(items[i].Handle); short < shortestHandle {
			shortestIndex = i
			shortestHandle = short
		}

		switch items[i].Handle {
		case "_", "root", "index":
			foundIndex = i
		}
	}

	if foundIndex == -1 {
		foundIndex = shortestIndex
	}

	root.Item = items[foundIndex]
	root.Handle = ""

	// Step two, run a sort on all the remaining items
	items = append(items[:foundIndex], items[foundIndex+1:]...)
	s.treeSort(&root, &items)

	// Finally, strip all parent prefixes
	for i := 0; i < len(root.Children); i++ {
		s.stripTopicHandlePrefixes(root, &root.Children[i])
	}

	return
}

/*
Given a list of Items with handles, sort them into a tree structure by their
handles, with child nodes stemming from prefixes.

For example, given the list of handles:

a
a_b_c
a_b_c_d
a_c

the function creates the tree:

a
|-a_b
| |-a_b_c
|   |-a_b_c_d
|-a_c

The handle names are preserved.
*/
func (s *sorter) treeSort(root *TopicRef, items *[]Item) {

	group, nongroup := []Item{}, []Item{}

	for _, item := range *items {

		if strings.HasPrefix(item.Handle, root.Handle) {
			group = append(group, item)
		} else {
			nongroup = append(nongroup, item)
		}
	}

	sort.Sort(itemListHandleSorter(group))

	for len(group) > 0 {

		child := TopicRef{Item: group[0]}

		group = group[1:]
		(*root).Children = append((*root).Children, child)
		s.treeSort(&root.Children[len(root.Children)-1], &group)
	}

	*items = nongroup
}

func (s *sorter) stripTopicHandlePrefixes(parent TopicRef, child *TopicRef) {

	for i := 0; i < len(child.Children); i++ {
		s.stripTopicHandlePrefixes(*child, &child.Children[i])
	}

	child.Handle = strings.Trim(strings.TrimPrefix(child.Handle, parent.Handle), "_")
}

func (s *sorter) sortItemsToGlossary(input []Item) (output []TermRef) {

	sort.Sort(itemListTitleSorter(input))

	for _, i := range input {
		output = append(output, TermRef{Item: i})
	}

	return
}
