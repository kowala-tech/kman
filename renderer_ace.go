package kman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/yosssi/ace"
)

type rendererAce struct {
	fs           afero.Fs
	templatePath string
	outputPath   string
}

type rendererAceNavigation struct {
	Title       string
	URL         string
	Active      bool
	ActiveChild bool
	Children    []rendererAceNavigation
}

func (r *rendererAceNavigation) flatten() (output []rendererAceNavigation) {
	r.flattenToList(&output)
	return
}

func (r *rendererAceNavigation) addToList(list *[]rendererAceNavigation) {

	*list = append(
		*list,
		rendererAceNavigation{
			Title:       r.Title,
			URL:         r.URL,
			Active:      r.Active,
			ActiveChild: r.ActiveChild,
		},
	)
}

func (r *rendererAceNavigation) flattenToList(list *[]rendererAceNavigation) {

	r.addToList(list)

	for _, child := range r.Children {
		child.flattenToList(list)
	}
}

func NewRendererAce(fs afero.Fs, templatePath, outputPath string) Renderer {
	return &rendererAce{
		fs:           fs,
		templatePath: templatePath,
		outputPath:   outputPath,
	}
}

func (r *rendererAce) Render(d Documentation) error {

	if err := r.executeTemplate("index", "", d, d.RootTopic.Title, d.RootTopic); err != nil {
		return err
	}

	for _, topic := range d.RootTopic.Children {
		if err := r.renderTopic("", d, topic); err != nil {
			return err
		}
	}

	if err := r.executeTemplate("glossary", "glossary", d, "Glossary", d.Glossary); err != nil {
		return err
	}

	return r.copyAssets()
}

func (r *rendererAce) navigation(d Documentation, currentPath string) (nav rendererAceNavigation) {

	nav = rendererAceNavigation{
		Title: d.RootTopic.Title,
		URL:   "/",
	}

	if currentPath == "/" {
		nav.Active = true
	} else {
		nav.ActiveChild = true
	}

	r.navigationBranch("/", currentPath, d.RootTopic.Children, &nav.Children)

	if len(d.Glossary) > 0 {

		glossary := rendererAceNavigation{
			Title: "Glossary",
			URL:   "/glossary",
		}

		if currentPath == "/glossary" {
			glossary.Active = true
		}

		nav.Children = append(nav.Children, glossary)
	}

	return
}

func (r *rendererAce) navigationBranch(root string, currentPath string, items []TopicRef, nav *[]rendererAceNavigation) {

	for _, topic := range items {

		url := filepath.Join(root, topic.Handle)

		branch := rendererAceNavigation{
			Title: topic.Title,
			URL:   url,
		}

		if currentPath == url {
			branch.Active = true
		} else if strings.HasPrefix(currentPath, url) {
			branch.ActiveChild = true
		}

		r.navigationBranch(url, currentPath, topic.Children, &branch.Children)

		*nav = append(*nav, branch)
	}
}

func (r *rendererAce) copyAssets() error {

	return afero.Walk(r.fs, r.templatePath, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir() && info.Size() > 0 && filepath.Ext(path) != ".ace" {
			dest := filepath.Join(r.outputPath, strings.TrimPrefix(path, r.templatePath+"/"))

			if err := r.fs.MkdirAll(filepath.Base(dest), os.ModePerm); err != nil {
				return err
			}

			data, err := afero.ReadFile(r.fs, path)

			if err != nil {
				return err
			}

			if err := afero.WriteFile(r.fs, dest, data, info.Mode()); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *rendererAce) renderTopic(parentPath string, doc Documentation, topic TopicRef) error {

	handle := filepath.Join(parentPath, topic.Handle)

	if err := r.executeTemplate("topic", handle, doc, topic.Title, topic); err != nil {
		return err
	}

	for _, child := range topic.Children {
		if err := r.renderTopic(handle, doc, child); err != nil {
			return err
		}
	}

	return nil
}

func (r *rendererAce) aceTemplatePath() string {
	return filepath.Join(r.templatePath, "ace")
}

func (r *rendererAce) htmlPath(path string) string {
	return filepath.Join(r.outputPath, path) + "/index.html"
}

func (r *rendererAce) asset(file string) ([]byte, error) {
	return afero.ReadFile(r.fs, file)
}

func (r *rendererAce) executeTemplate(src, dest string, d Documentation, title string, context interface{}) error {

	tpl, err := ace.Load("master", src, &ace.Options{
		Asset:         r.asset,
		DynamicReload: true,
		FuncMap:       r.templateFuncs(),
		BaseDir:       r.aceTemplatePath(),
	})

	if err != nil {
		return err
	}

	var buf bytes.Buffer

	pageURL := "/" + dest
	nav := r.navigation(d, pageURL)

	args := struct {
		Context     interface{}
		Doc         Documentation
		Navigation  rendererAceNavigation
		SearchItems []rendererAceNavigation
		Glossary    []TermRef
		Title       string
		PageURL     string
	}{
		Doc:         d,
		Context:     context,
		Navigation:  nav,
		SearchItems: nav.flatten(),
		Glossary:    d.Glossary,
		Title:       title,
		PageURL:     pageURL,
	}

	if err := tpl.Execute(&buf, args); err != nil {
		return err
	}

	if err := r.fs.MkdirAll(filepath.Base(r.htmlPath(dest)), os.ModePerm); err != nil {
		return err
	}

	return afero.WriteReader(r.fs, r.htmlPath(dest), &buf)
}

func (r *rendererAce) templateFuncs() template.FuncMap {
	return template.FuncMap{
		"json": func(inp interface{}) template.JS {

			jsn, err := json.Marshal(inp)

			if err == nil {
				return template.JS(jsn)
			}

			return template.JS(fmt.Sprintf("{\"json.Marshal failed\":\"%s\"}", err))
		},
	}
}
