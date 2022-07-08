package catalog

import (
	"encoding/json"
	"strings"
)

type Tag struct {
	InteractionGroups map[Protocol]TagInteractionGroup
	Children          *Tags
	Name              TagName
	Title             string
	Description       string
}

var _ json.Marshaler = &Tags{}

func newEmptyTag(r InteractionId) *Tag {
	title := tagTitle(r.Path().String())
	return &Tag{
		InteractionGroups: make(map[Protocol]TagInteractionGroup),
		Children:          &Tags{},
		Title:             title,
		Name:              tagName(title),
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

func (t *Tag) appendInteractionId(k InteractionId) {
	list, ok := t.InteractionGroups[k.Protocol()]
	if !ok {
		list = newTagInteractionGroup(k.Protocol())
		t.InteractionGroups[k.Protocol()] = list
	}
	list.append(k)
}

func (t *Tag) MarshalJSON() ([]byte, error) {
	var data struct {
		Children          *Tags                 `json:"children,omitempty"`
		Name              TagName               `json:"name"`
		Title             string                `json:"title"`
		Description       string                `json:"description,omitempty"`
		InteractionGroups []TagInteractionGroup `json:"interactionGroups"`
	}

	data.Name = t.Name
	data.Title = t.Title
	data.Description = t.Description
	data.InteractionGroups = make([]TagInteractionGroup, 0, len(t.InteractionGroups))

	// it is important to keep the sorting
	if v, ok := t.InteractionGroups[HTTP]; ok {
		data.InteractionGroups = append(data.InteractionGroups, v)
	}
	if v, ok := t.InteractionGroups[JsonRpc]; ok {
		data.InteractionGroups = append(data.InteractionGroups, v)
	}

	if t.Children.Len() > 0 {
		data.Children = t.Children
	}

	return json.Marshal(data)
}
