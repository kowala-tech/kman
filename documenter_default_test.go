package kman

import (
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/stretchr/testify/require"
)

func Test_ADefaultDocumenterShouldBeAbleToDocument(t *testing.T) {

	a0, a1 := &mockAssembler{[]Item{
		Item{
			Type:   ItemTypeTopic,
			Handle: "A",
			Title:  "A",
		},
	}},
		&mockAssembler{[]Item{
			Item{
				Type:   ItemTypeTerm,
				Handle: "B",
				Title:  "B",
			},
		}}

	docer := NewDefaultDocumenter(NewDefaultSorter(), a0, a1)
	doc, err := docer.Document()
	require.Nil(t, err)
	snaptest.Snapshot(t, doc)
}
