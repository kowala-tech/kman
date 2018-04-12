package kman

import (
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func newValidTemplateFilesystem(t *testing.T) afero.Fs {
	return newMockFilesystem(t,
		map[string]string{
			"template/ace/master.ace": `
= doctype html
html lang=en
  head
    meta charset=utf-8
    title Example
  body
    h1 This is a base template
    = yield main
`,
			"template/ace/index.ace": `
= content main
  h2 Index
  .topic {{.Context.HTML}}
`,
			"template/ace/topic.ace": `
= content main
  h2 Topic
`,
			"template/ace/glossary.ace": `
= content main
  h2 Glossary
`,
			"template/robots.txt":      "A",
			"template/images/logo.svg": "B",
		},
	)
}

func newValidDocumentation(t *testing.T) Documentation {
	return Documentation{
		RootTopic: TopicRef{
			Item: Item{
				Type:     0,
				FileName: "doc/topics.md",
				Line:     0,
				Title:    "k-man: intuitive documentation parser and presenter",
				Handle:   "",
				Content:  "This is an example topic which forms the root",
			},
			Children: []TopicRef{
				TopicRef{
					Item: Item{
						Type:     0,
						FileName: "doc/topics.md",
						Line:     0,
						Title:    "Usage",
						Handle:   "usage",
						Content:  "This is a topic with an explicit handle",
					},
					Children: []TopicRef{
						TopicRef{
							Item: Item{
								Type:     0,
								FileName: "doc/topics.md",
								Line:     0,
								Title:    "Usage: advanced",
								Handle:   "advanced",
								Content:  "This lives under 'usage'",
							},
							Children: nil,
						},
					},
				},
			},
		},
		Glossary: []TermRef{
			TermRef{
				Item: Item{
					Type:     1,
					FileName: "doc/terms.md",
					Line:     0,
					Title:    "Another example",
					Handle:   "another_example",
					Content:  "Another markdown-parsed example",
				},
			},
			TermRef{
				Item: Item{
					Type:     1,
					FileName: "doc/terms.md",
					Line:     0,
					Title:    "Example",
					Handle:   "example",
					Content:  "An example term, parsed from markdown",
				},
			},
		},
	}
}

func Test_AnAceRendererCanAssembleTheNav(t *testing.T) {

	fs := newValidTemplateFilesystem(t)
	renderer := NewRendererAce(fs, "template", "public")

	snaptest.Snapshot(t, renderer.(*rendererAce).navigation(newValidDocumentation(t), "/usage/advanced"))
}

func Test_ARendererAceCanRenderAWebsite(t *testing.T) {

	fs := newValidTemplateFilesystem(t)
	renderer := NewRendererAce(fs, "template", "public")

	require.Nil(t, renderer.Render(newValidDocumentation(t)))

	snapshotFilesystem(t, fs)
}
