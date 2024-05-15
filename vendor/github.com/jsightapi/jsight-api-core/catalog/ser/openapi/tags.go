package openapi

import "github.com/jsightapi/jsight-api-core/catalog"

func gatherTags(c *catalog.Catalog, names []catalog.TagName) []*catalog.Tag {
	tags := make([]*catalog.Tag, 0, len(names))
	for _, n := range names {
		if t, ok := c.Tags.Get(n); ok {
			tags = append(tags, t)
		}
	}
	return tags
}

func getTagTitles(tags []*catalog.Tag) []string {
	titles := make([]string, 0, len(tags))
	for _, t := range tags {
		titles = append(titles, t.Title)
	}
	return titles
}
