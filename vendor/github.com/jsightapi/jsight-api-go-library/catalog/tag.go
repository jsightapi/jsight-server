package catalog

import (
	"encoding/json"
	"strings"
)

type Tag struct {
	ResourceMethods *TagResourceMethods
	Children        *Tags
	Name            TagName
	Title           string
	Description     string
}

var _ json.Marshaler = &Tags{}

func newEmptyTag(r ResourceMethodId) *Tag {
	title := tagTitle(r.path.String())
	return &Tag{
		ResourceMethods: &TagResourceMethods{},
		Children:        &Tags{},
		Title:           title,
		Name:            tagName(title),
	}
}

func tagTitle(path string) string {
	p := strings.Split(path, "/")
	for len(p) != 0 {
		if p[0] != "" && p[0] != "." {
			break
		}
		p = p[1:]
	}
	if len(p) == 0 {
		return "/"
	}
	return "/" + p[0]
}

func (t *Tag) appendResourceMethodId(r ResourceMethodId) {
	list, ok := t.ResourceMethods.Get(r.path)
	if !ok {
		list = newResourceMethodIdList()
		t.ResourceMethods.Set(r.path, list)
	}
	list.append(r)
}

func (t *Tag) MarshalJSON() ([]byte, error) {
	var data struct {
		ResourceMethods *TagResourceMethods `json:"resourceMethods"`
		Children        *Tags               `json:"children,omitempty"`
		Name            TagName             `json:"name"`
		Title           string              `json:"title"`
		Description     string              `json:"description,omitempty"`
	}

	data.Name = t.Name
	data.Title = t.Title
	data.Description = t.Description
	data.ResourceMethods = t.ResourceMethods
	if t.Children.Len() > 0 {
		data.Children = t.Children
	}

	return json.Marshal(data)
}
