package kman

import (
	"fmt"
	"testing"

	"github.com/kowala-tech/snaptest"
)

func Test_ADefaultSorterShouldBeAbleToSortRawItemsIntoTopicTrees(t *testing.T) {

	for cycle, test := range []struct {
		description string

		input []Item
	}{
		{
			description: "No items",
		},
		{
			description: "One item",
			input: []Item{
				Item{Handle: "Anything"},
			},
		},
		{
			description: "two items: shorter should be root",
			input: []Item{
				Item{Handle: "should_not_be_root"},
				Item{Handle: "should_be_root"},
			},
		},
		{
			description: "Root takes precence",
			input: []Item{
				Item{Handle: "a"},
				Item{Handle: "b"},
				Item{Handle: "root"},
			},
		},
		{
			description: "basic tree sort",
			input: []Item{
				Item{Handle: "a"},
				Item{Handle: "abc"},
				Item{Handle: "ab"},
				Item{Handle: "abcd"},
				Item{Handle: "ac"},
				Item{Handle: "acb"},
			},
		},
		{
			description: "basic tree sort with underscores",
			input: []Item{
				Item{Handle: "a"},
				Item{Handle: "a_b_c"},
				Item{Handle: "a_b"},
				Item{Handle: "a_b_c_d"},
				Item{Handle: "a_c"},
				Item{Handle: "a_c_b"},
			},
		},
		{
			description: "basic tree sort with weird root",
			input: []Item{
				Item{Handle: "a"},
				Item{Handle: "a_b_c"},
				Item{Handle: "_"},
				Item{Handle: "a_b_c_d"},
				Item{Handle: "a_c"},
				Item{Handle: "a_c_b"},
			},
		},
	} {
		t.Run(fmt.Sprintf("Cycle %d: %s", cycle, test.description), func(t *testing.T) {

			sorter := &sorter{}

			snaptest.Snapshot(t, sorter.sortItemsToTopicTree(test.input))
		})
	}
}

func Test_ADefaultSorterShouldBeAbleToOrganiseTermsIntoAGlossary(t *testing.T) {

	for cycle, test := range []struct {
		description string

		input []Item
	}{
		{
			description: "No items",
		},
		{
			description: "One item",
			input: []Item{
				Item{Title: "A", Handle: "Anything"},
			},
		},
		{
			description: "Three items",
			input: []Item{
				Item{Handle: "a", Title: "c"},
				Item{Handle: "c", Title: "a"},
				Item{Handle: "b", Title: "b"},
			},
		},
	} {
		t.Run(fmt.Sprintf("Cycle %d: %s", cycle, test.description), func(t *testing.T) {

			sorter := &sorter{}

			snaptest.Snapshot(t, sorter.sortItemsToGlossary(test.input))
		})
	}
}

func Test_ADefaultSortedCanSortMixedItemsIntoDocumentation(t *testing.T) {

	sorter := &sorter{}

	snaptest.Snapshot(t, sorter.Sort([]Item{
		Item{Type: ItemTypeTopic, Handle: "A", Title: "A"},
		Item{Type: ItemTypeTopic, Handle: "A_B", Title: "B"},
		Item{Type: ItemTypeTopic, Handle: "A_B_C", Title: "C"},
		Item{Type: ItemTypeTerm, Handle: "D", Title: "D"},
	}))
}
