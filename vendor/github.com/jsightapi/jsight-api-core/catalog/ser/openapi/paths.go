package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

type Paths map[string]*PathItem

func defaultPaths() Paths {
	return Paths{}
}

func newPaths(c *catalog.Catalog) (Paths, Error) {
	if c.Interactions.Len() == 0 {
		return defaultPaths(), nil
	}

	p := make(Paths, c.Interactions.Len())
	if err := fillPaths(p, c); err != nil {
		return nil, err
	}
	return p, nil
}

func fillPaths(p Paths, c *catalog.Catalog) Error {
	err := c.Interactions.Each(func(k catalog.InteractionID, v catalog.Interaction) error {
		if k.Protocol() == catalog.HTTP {
			i := v.(*catalog.HTTPInteraction)

			path := i.Path().String()
			if _, exists := p[path]; !exists {
				pi, err := newPathItem(i)
				if err != nil {
					return err
				}
				p[path] = pi
			}
			op, err := httpInteractionToOperation(i, c)
			if err != nil {
				return err
			}
			p[path].assignOperation(i.HttpMethod, op)
		}
		return nil
	})

	return castErr(err)
}

func httpInteractionToOperation(i *catalog.HTTPInteraction, c *catalog.Catalog) (*Operation, Error) {
	tags := gatherTags(c, i.Tags)
	tagTitles := getTagTitles(tags)
	return newOperation(i, tagTitles)
}
