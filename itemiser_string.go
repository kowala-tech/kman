package kman

import "strings"

const (
	topicToken  = "topic:"
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

// TODO! add Type to Items
func (s *itemiserString) Itemise(items *[]Item) error {

	lines := strings.Split(s.input, "\n")

	topic, handle, content := "", "", []string{}

	reset := func() {
		topic, handle, content = "", "", []string{}
	}

	addItem := func() {
		*items = append(*items, Item{
			FileName: s.path,
			Title:    topic,
			Handle:   handle,
			Content:  strings.Trim(strings.Join(content, "\n"), "\n"),
		})
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(strings.ToLower(line), topicToken) {

			if topic != "" && len(lines) > 0 {
				addItem()
				reset()
			}

			topic = strings.TrimSpace(line[len(topicToken):])
			handle = s.handlise(topic)

		} else if strings.HasPrefix(strings.ToLower(line), handleToken) {
			handle = strings.TrimSpace(line[len(handleToken):])
		} else if topic != "" {
			content = append(content, line)
		}
	}

	if topic != "" {
		addItem()
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
