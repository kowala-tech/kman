package kman

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type assemblerMarkdownFilesystem struct {
	fs afero.Fs
}

func NewMarkdownAssemblerWithFilesystem(fs afero.Fs) Assembler {
	return &assemblerMarkdownFilesystem{
		fs: fs,
	}
}

func NewMarkdownAssemblerFromLocalFilesystem() Assembler {
	return NewMarkdownAssemblerWithFilesystem(afero.NewOsFs())
}

func (m *assemblerMarkdownFilesystem) Assemble() ([]Item, error) {

	docItems := []Item{}

	files := m.findMarkdownFiles()

	for _, file := range files {

		content, err := afero.ReadFile(m.fs, file)

		if err != nil {
			return docItems, err
		}

		if err := NewItemiserFromString(file, string(content)).Itemise(&docItems); err != nil {
			return docItems, err
		}
	}

	return docItems, nil
}

func (m *assemblerMarkdownFilesystem) findMarkdownFiles() (files []string) {

	afero.Walk(m.fs, ".", func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() && info.Size() > 0 && (filepath.Ext(path) == ".md" || filepath.Ext(path) == ".markdown") {
			files = append(files, path)
		}

		return nil
	})

	return
}
