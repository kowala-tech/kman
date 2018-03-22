package kman

import "strings"

const (
	topicToken  = "topic:"
	termToken   = "term:"
	handleToken = "handle:"
)

type itemiserString struct {
	path  string
	input string
}

func NewItemiserFromString(path, input string) Itemiser {
	return &itemiserString{
		input: input,
		path:  path,
	}
}

func (s *itemiserString) Itemise(items *[]Item) error {

	lines := strings.Split(s.input, "\n")

	title, handle, content, typ := "", "", []string{}, ItemTypeTopic

	reset := func() {
		title, handle, content, typ = "", "", []string{}, ItemTypeTopic
	}

	addItem := func(typ ItemType) {
		*items = append(*items, Item{
			Type:     typ,
			FileName: s.path,
			Title:    title,
			Handle:   handle,
			Content:  strings.Trim(strings.Join(content, "\n"), "\n"),
		})
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(strings.ToLower(line), topicToken) {

			if title != "" && len(lines) > 0 {
				addItem(typ)
				reset()
			}

			title = strings.TrimSpace(line[len(topicToken):])
			handle = s.handlise(title)
			typ = ItemTypeTopic

		} else if strings.HasPrefix(strings.ToLower(line), termToken) {

			if title != "" && len(lines) > 0 {
				addItem(typ)
				reset()
			}

			title = strings.TrimSpace(line[len(topicToken):])
			handle = s.handlise(title)
			typ = ItemTypeTerm

		} else if strings.HasPrefix(strings.ToLower(line), handleToken) {
			handle = strings.TrimSpace(line[len(handleToken):])
		} else if title != "" {
			content = append(content, line)
		}
	}

	if title != "" {
		addItem(typ)
	}

	return nil
}

func (g *itemiserString) handlise(input string) (output string) {

	input = strings.TrimSpace(input)

	var buf []rune

	for _, r := range strings.ToLower(input) {
		switch {
		case r >= 'a' && r <= 'z':
			buf = append(buf, r)

		case r >= '0' && r <= '9':
			buf = append(buf, r)

		case r == ' ':
			buf = append(buf, '_')
		}
	}

	return string(buf)
}
