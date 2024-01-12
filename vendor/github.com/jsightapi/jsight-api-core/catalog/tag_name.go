package catalog

import (
	"net/url"
	"strings"
)

type TagName string

func tagName(title string) TagName {
	if title == "/" {
		return "@_"
	}
	title = strings.Replace(title, "/", "@", 1)
	title = strings.ReplaceAll(title, "_", "__")
	title = url.PathEscape(title)
	title = strings.ReplaceAll(title, "%", "_")
	return TagName(title)
}

func (t TagName) MarshalText() ([]byte, error) {
	return []byte(t), nil
}
