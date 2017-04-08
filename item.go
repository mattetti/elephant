package elephant

import "fmt"

type Item struct {
	Title        string
	Link         string
	Desc         string
	Tags         []string
	Source       string
	InternalLink string
}

func (it *Item) String() string {
	if it == nil {
		return ""
	}
	return fmt.Sprintf("Title: %s\nLink: %s\nDesc: %s\nSource: %s\nTags: %q\n", it.Title, it.Link, it.Desc, it.Source, it.Tags)
}
