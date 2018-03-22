package kman

import (
	"fmt"
	"testing"

	"github.com/kowala-tech/snaptest"
	"github.com/stretchr/testify/require"
)

func Test_AValidItemiserStringCanHandliseStrings(t *testing.T) {

	for _, test := range [][2]string{
		{"A", "a"},
		{"1", "1"},
		{"My Long String", "my_long_string"},
		{"Someone's string", "someones_string"},
		{" Someone's string", "someones_string"},
	} {
		itemiser := &itemiserString{}
		require.Equal(t, test[1], itemiser.handlise(test[0]))
	}
}

func Test_AValidItemiserStringShouldFindTopicsAndTerms(t *testing.T) {

	for cycle, test := range []struct {
		description string

		input string
		err   bool
	}{
		{
			description: "No topics",
			input:       "",
			err:         false,
		},
		{
			description: "One topic; default handle",
			input: `
	Topic: test 1
	Line 1
		Line 2
`,
			err: false,
		},
		{
			description: "Two topics; default handle",
			input: `
	Topic: test 1
	Line 1
		Line 2

	Topic: test 1
	Line 1
		Line 2
`,
			err: false,
		},
		{
			description: "One topic; specific handle",
			input: `
	Topic: test 1
	Handle: my_handle
	Line 1
		Line 2
`,
			err: false,
		},
		{
			description: "Two topics; one specific handle",
			input: `
	Topic: test 1
	Line 1
		Line 2

	Topic: test 1
	Line 1
	Handle: my_other_handle
		Line 2
`,
			err: false,
		},
		{
			description: "One term; default handle",
			input: `
	Term: test 1
	Line 1
		Line 2
`,
			err: false,
		},
		{
			description: "One term; specific handle",
			input: `
	Term: test 1
	Handle: some_other_title
	Line 1
`,
			err: false,
		},
		{
			description: "Two terms; one specific handle",
			input: `
	Term: test 1
	Handle: some_other_title
	Line 1

	Term: test 2
	Line 2
`,
			err: false,
		},
		{
			description: "Mixed",
			input: `

	Topic: test A
	Line A

	Topic: test B
	Handle: some_title
	Line B

	Term: test 1
	Handle: some_other_title
	Line 1

	Term: test 2
	Line 2
`,
			err: false,
		},
	} {
		t.Run(fmt.Sprintf("Cycle %d: %s", cycle, test.description), func(t *testing.T) {

			itemiser := NewItemiserFromString("some-path.ext", test.input)
			items := []Item{}

			err := itemiser.Itemise(&items)

			if !test.err {
				require.Nil(t, err)
			} else {
				snaptest.Snapshot(t, err)
			}

			snaptest.Snapshot(t, items)
		})
	}
}
