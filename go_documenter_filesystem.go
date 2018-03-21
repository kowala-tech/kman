package kman

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/afero"
)

const (
	topicRef = "kman.Topic"
)

type goDocumenterFileSystem struct {
	fs afero.Fs
}

func NewGoDocumenterWithFileSystem(fs afero.Fs) Documenter {
	return &goDocumenterFileSystem{
		fs: fs,
	}
}

func NewGoDocumenterFromLocalFilesystem() Documenter {
	return NewGoDocumenterWithFileSystem(afero.NewOsFs())
}

func (g *goDocumenterFileSystem) Document() (Documentation, error) {

	doc := Documentation{}

	astFiles, err := g.parseFiles(g.findGoFiles())

	if err != nil {
		return doc, err
	}

	topicItems := []Item{}

	for path, f := range astFiles {
		for _, d := range f.Comments {
			if err := g.findCommentReference(path, d, &topicItems); err != nil {
				return doc, nil
			}
		}

		for _, d := range f.Decls {
			g.findReference(path, topicRef, d.(ast.Node), &topicItems)
		}
	}

	doc.RootTopic = g.sortItemsToTopicTree(topicItems)

	return doc, nil
}

func (g *goDocumenterFileSystem) findGoFiles() (files []string) {

	afero.Walk(g.fs, ".", func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() && info.Size() > 0 && filepath.Ext(path) == ".go" {
			files = append(files, path)
		}

		return nil
	})

	return
}

func (g *goDocumenterFileSystem) sortItemsToTopicTree(items []Item) (root TopicRef) {

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
	g.treeSort(&root, &items)

	// Finally, strip all parent prefixes
	for i := 0; i < len(root.Children); i++ {
		g.stripTopicHandlePrefixes(root, &root.Children[i])
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
func (g *goDocumenterFileSystem) treeSort(root *TopicRef, items *[]Item) {

	group, nongroup := []Item{}, []Item{}

	for _, item := range *items {

		if strings.HasPrefix(item.Handle, root.Handle) {
			group = append(group, item)
		} else {
			nongroup = append(nongroup, item)
		}
	}

	sort.Sort(itemList(group))

	for len(group) > 0 {

		child := TopicRef{Item: group[0]}

		group = group[1:]
		(*root).Children = append((*root).Children, child)
		g.treeSort(&root.Children[len(root.Children)-1], &group)
	}

	*items = nongroup
}

func (g *goDocumenterFileSystem) stripTopicHandlePrefixes(parent TopicRef, child *TopicRef) {
	for i := 0; i < len(child.Children); i++ {
		g.stripTopicHandlePrefixes(*child, &child.Children[i])
	}

	child.Handle = strings.Trim(strings.TrimPrefix(child.Handle, parent.Handle), "_")
}

func (g *goDocumenterFileSystem) parseFiles(paths []string) (map[string]*ast.File, error) {

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

func (g *goDocumenterFileSystem) findReference(path, symbol string, n ast.Node, items *[]Item) {

	switch x := n.(type) {
	case *ast.GenDecl:
		switch x.Tok {
		case token.VAR:
			g.findVarReference(path, symbol, x, items)
		}
	}
}

func (g *goDocumenterFileSystem) findCommentReference(path string, x *ast.CommentGroup, items *[]Item) error {

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

func (g *goDocumenterFileSystem) findVarReference(path, symbol string, x *ast.GenDecl, items *[]Item) {

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
			FileName: path,
			Title:    ref,
			Handle:   name,
			Content:  comment,
		})
	}
}
