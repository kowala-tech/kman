package kman

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

const (
	topicRef = "kman.Topic"
	termRef  = "kman.Term"
)

type goAssemblerFileSystem struct {
	fs afero.Fs
}

func NewGoAssemblerWithFilesystem(fs afero.Fs) Assembler {
	return &goAssemblerFileSystem{
		fs: fs,
	}
}

func NewGoAssemblerFromLocalFilesystem() Assembler {
	return NewGoAssemblerWithFilesystem(afero.NewOsFs())
}

func (g *goAssemblerFileSystem) Assemble() ([]Item, error) {

	docItems := []Item{}

	astFiles, err := g.parseFiles(g.findGoFiles())

	if err != nil {
		return docItems, err
	}

	for path, f := range astFiles {
		for _, d := range f.Comments {
			if err := g.findCommentReference(path, d, &docItems); err != nil {
				return docItems, nil
			}
		}

		for _, d := range f.Decls {
			g.findReference(path, topicRef, d.(ast.Node), &docItems, ItemTypeTopic)
			g.findReference(path, termRef, d.(ast.Node), &docItems, ItemTypeTerm)
		}
	}

	return docItems, nil
}

func (g *goAssemblerFileSystem) findGoFiles() (files []string) {

	afero.Walk(g.fs, ".", func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() && info.Size() > 0 && filepath.Ext(path) == ".go" {
			files = append(files, path)
		}

		return nil
	})

	return
}

func (g *goAssemblerFileSystem) parseFiles(paths []string) (map[string]*ast.File, error) {

	fileSet := token.NewFileSet()
	astFiles := make(map[string]*ast.File)

	for _, path := range paths {

		contents, err := afero.ReadFile(g.fs, path)

		if err != nil {
			return astFiles, err
		}

		parsed, err := parser.ParseFile(fileSet, path, contents, parser.ParseComments)

		if err != nil {
			return astFiles, err
		}

		astFiles[path] = parsed
	}

	return astFiles, nil
}

func (g *goAssemblerFileSystem) findReference(path, symbol string, n ast.Node, items *[]Item, itemType ItemType) {

	switch x := n.(type) {
	case *ast.GenDecl:
		switch x.Tok {
		case token.VAR:
			g.findVarReference(path, symbol, x, items, itemType)
		}
	}
}

func (g *goAssemblerFileSystem) findCommentReference(path string, x *ast.CommentGroup, items *[]Item) error {

	if x == nil {
		return nil
	}

	for _, com := range x.List {

		if err := NewItemiserFromString(path, strings.TrimSuffix(com.Text, "*/")).Itemise(items); err != nil {
			return err
		}
	}

	return nil
}

func (g *goAssemblerFileSystem) findVarReference(path, symbol string, x *ast.GenDecl, items *[]Item, itemType ItemType) {

	if len(x.Specs) == 0 {
		return
	}

	value, ok := x.Specs[0].(*ast.ValueSpec)
	if !ok {
		return
	}

	if len(value.Names) == 0 {
		return
	}

	if len(value.Values) == 0 {
		return
	}

	callexp, ok := value.Values[0].(*ast.CallExpr)
	if !ok {
		return
	}

	selectorexp, ok := callexp.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	selector := fmt.Sprintf("%s.%s", selectorexp.X, selectorexp.Sel)
	if selector != symbol {
		return
	}

	if len(callexp.Args) != 1 {
		return
	}

	refname, ok := callexp.Args[0].(*ast.BasicLit)
	if !ok {
		return
	}

	if refname.Kind != token.STRING {
		return
	}

	name := fmt.Sprintf("%s", value.Names[0])
	comment := strings.Trim(x.Doc.Text(), "\n")
	ref := strings.Trim(refname.Value, "\"")

	if ref != "" && name != "" && comment != "" {
		*items = append(*items, Item{
			Type:     itemType,
			FileName: path,
			Title:    ref,
			Handle:   name,
			Content:  comment,
		})
	}
}
