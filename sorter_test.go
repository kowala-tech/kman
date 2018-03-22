package kman

import (
	"fmt"
	"testing"

	"github.com/kowala-tech/snaptest"
)

// TODO! test sort topics and terms
func Test_TheDefaultSorterShouldBeAbleToSortRawItemsIntoTopicTrees(t *testing.T) {
	for cycle, test := range []struct {
		description string

		input  []Item
		output TopicRef
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

			generator := &sorter{}

			snaptest.Snapshot(t, generator.sortItemsToTopicTree(test.input))
		})
	}
}
