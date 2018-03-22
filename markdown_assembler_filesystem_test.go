package kman

import (
	"fmt"
	"os"
	"testing"

	"github.com/kowala-tech/snaptest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func newMockMarkdownFileSystem(t *testing.T, files map[string]string) afero.Fs {
	fs := afero.NewMemMapFs()

	for path, content := range files {
		require.Nil(t, afero.WriteFile(fs, path, []byte(content), os.ModePerm))
	}

	return fs
}

func Test_AValidMarkdownFileSystemAssemblerShouldbeAbleToFindGoFiles(t *testing.T) {

	for cycle, test := range []struct {
		description string

		input  afero.Fs
		output []string
	}{
		{
			description: "Empty filesystem",
			input:       newMockMarkdownFileSystem(t, map[string]string{}),
		},
		{
			description: "Non-empty filesystem",
			input: newMockGoFileSystem(t, map[string]string{
				"valid.md":            `hello`,
				"valid.markdown":      `hello`,
				"empty.md":            ``,
				"not-md.txt":          `hello`,
				"directory.md/non-md": `123`,
				"directory.md/md.md":  `hello`,
			}),
			output: []string{
				"directory.md/md.md",
				"valid.markdown",
				"valid.md",
			},
		},
	} {
		t.Run(fmt.Sprintf("Cycle %d: %s", cycle, test.description), func(t *testing.T) {

			assembler := &markdownAssemblerFileSystem{
				fs: test.input,
			}

			files := assembler.findMarkdownFiles()
			require.Equal(t, test.output, files)
		})
	}

}

func Test_AValidMarkdownFileSystemAssemblerShouldFindTopicsAndTerms(t *testing.T) {

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
				fs: newMockMarkdownFileSystem(t, map[string]string{}),
			},
		},
		{
			description: "Happy path: non-empty filesystem",
			input: input{
				fs: newMockMarkdownFileSystem(t, map[string]string{
					"first.md": `
Topic: A
Line 1

Term: B
Line 2
`,
				}),
			},
		},
	} {
		t.Run(fmt.Sprintf("Cycle %d: %s", cycle, test.description), func(t *testing.T) {

			assembler := &markdownAssemblerFileSystem{
				fs: test.input.fs,
			}

			doc, err := assembler.Assemble()

			if !test.output.err {
				require.Nil(t, err)
			} else {
				snaptest.Snapshot(t, err)
			}

			snaptest.Snapshot(t, doc)
		})
	}
}
