package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

type Paths map[string]*PathItem

func defaultPaths() Paths {
	return Paths{}
}

func newPaths(c *catalog.Catalog) Paths {
	if c.Interactions.Len() == 0 {
		return defaultPaths()
	}

	p := make(Paths, c.Interactions.Len())
	fillPaths(p, c.Interactions)
	return p
}

func fillPaths(p Paths, ii *catalog.Interactions) {
	_ = ii.Each(func(k catalog.InteractionID, v catalog.Interaction) error {
		if k.Protocol() == catalog.HTTP {
			i := v.(*catalog.HTTPInteraction)
			addOperation(p, i)
		}
		return nil
	})
}

func addOperation(p Paths, i *catalog.HTTPInteraction) {
	path := i.Path().String()
	if _, exists := p[path]; exists {
		p[path].assignOperation(i.HttpMethod, newOperation(i))
	} else {
		p[path] = newPathItem(i)
	}
}
