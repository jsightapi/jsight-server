package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

type Info struct {
	Title       string  `json:"title"`
	Version     string  `json:"version"`
	Description *string `json:"description,omitempty"`
}

func defaultInfo() *Info {
	return &Info{
		Title:   "",
		Version: "",
	}
}

func newInfo(i *catalog.Info) *Info {
	if i == nil {
		return defaultInfo()
	}

	return &Info{
		Title:       i.Title,
		Version:     i.Version,
		Description: i.Description,
	}
}
