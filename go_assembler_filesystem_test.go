package kman

import (
	"fmt"
	"os"
	"testing"

	"github.com/kowala-tech/snaptest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func newMockGoFileSystem(t *testing.T, files map[string]string) afero.Fs {
	fs := afero.NewMemMapFs()

	for path, content := range files {
		require.Nil(t, afero.WriteFile(fs, path, []byte(content), os.ModePerm))
	}

	return fs
}

func Test_AValidGoFileSystemAssemblerShouldbeAbleToFindGoFiles(t *testing.T) {

	for cycle, test := range []struct {
		description string

		input  afero.Fs
		output []string
	}{
		{
			description: "Empty filesystem",
			input:       newMockGoFileSystem(t, map[string]string{}),
		},
		{
			description: "Non-empty filesystem",
			input: newMockGoFileSystem(t, map[string]string{
				"valid.go":            `package test`,
				"empty.go":            ``,
				"not-go.txt":          `package test`,
				"directory.go/non-go": `123`,
				"directory.go/go.go":  `package test`,
			}),
			output: []string{
				"directory.go/go.go",
				"valid.go",
			},
		},
	} {
		t.Run(fmt.Sprintf("Cycle %d: %s", cycle, test.description), func(t *testing.T) {

			generator := &goAssemblerFileSystem{
				fs: test.input,
			}

			files := generator.findGoFiles()
			require.Equal(t, test.output, files)
		})
	}

}

func Test_AValidGoFileSystemAssemblerShouldbeAbleToParseGivenFiles(t *testing.T) {

	type input struct {
		fs    afero.Fs
		files []string
	}

	type output struct {
		err bool
	}

	for cycle, test := range []struct {
		description string

		input  input
		output output
	}{
		{
			description: "Happy path: empty filesystem",
			input: input{
				fs:    newMockGoFileSystem(t, map[string]string{}),
				files: []string{},
			},
			output: output{},
		},
		{
			description: "Happy path: non-empty filesystem",
			input: input{
				fs: newMockGoFileSystem(t, map[string]string{
					"path/to/first.go": `package test`,
				}),
				files: []string{"path/to/first.go"},
			},
			output: output{},
		},
		{
			description: "Unhappy path: Inaccessible file",
			input: input{
				fs:    newMockGoFileSystem(t, map[string]string{}),
				files: []string{"path/to/first.go"},
			},
			output: output{
				err: true,
			},
		},
		{
			description: "Unhappy path: Syntax error",
			input: input{
				fs: newMockGoFileSystem(t, map[string]string{
					"path/to/first.go": `package1 test`,
				}),
				files: []string{"path/to/first.go"},
			},
			output: output{
				err: true,
			},
		},
	} {
		t.Run(fmt.Sprintf("Cycle %d: %s", cycle, test.description), func(t *testing.T) {

			generator := &goAssemblerFileSystem{
				fs: test.input.fs,
			}

			astFiles, err := generator.parseFiles(test.input.files)

			if !test.output.err {
				require.Nil(t, err)
			} else {
				snaptest.Snapshot(t, err)
			}

			snaptest.Snapshot(t, astFiles)
		})
	}

}

func Test_AValidGoFileSystemAssemblerShouldFindTopicsAndTerms(t *testing.T) {

	type input struct {
		fs afero.Fs
	}

	type output struct {
		err bool
	}

	for cycle, test := range []struct {
		description string

		input  input
		output output
	}{
		{
			description: "Happy path: Empty filesystem",
			input: input{
				fs: newMockGoFileSystem(t, map[string]string{}),
			},
			output: output{},
		},
		{
			description: "Happy path: non-empty filesystem",
			input: input{
				fs: newMockGoFileSystem(t, map[string]string{
					"first.go": `package first
/*
This is the root
*/
var _ = kman.Topic("Root")
`,
					"path/second.go": `

/*
Topic: godoc level
Hello
*/
package second

/*
One thing
*/
var topic = kman.Topic("Title")

/*
Another thing
*/
var topic_subtopic = kman.Topic("Title")

/*
Topic: topic 3
Handle: my-handle

This is my content
*/
type a interface{}

/*
Topic: topic 4

Handle should be implied.

Line 2

Topic: topic 5
Line 1
Line 2
*/

/*
Topic is malformed; shouled be ignored.
*/

type c interface{}
func d(){}

`,
				}),
			},
			output: output{},
		},
		{
			description: "Unhappy path: parser error",
			input: input{
				fs: newMockGoFileSystem(t, map[string]string{
					"first.go": `package first
/*
This is the root
*/
var 323_ = kman.Topic("Root")
`,
				}),
			},
			output: output{err: true},
		},
	} {
		t.Run(fmt.Sprintf("Cycle %d: %s", cycle, test.description), func(t *testing.T) {

			generator := &goAssemblerFileSystem{
				fs: test.input.fs,
			}

			doc, err := generator.Assemble()

			if !test.output.err {
				require.Nil(t, err)
			} else {
				snaptest.Snapshot(t, err)
			}

			snaptest.Snapshot(t, doc)
		})
	}
}
