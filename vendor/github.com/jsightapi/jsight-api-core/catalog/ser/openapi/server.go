package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
)

type Server struct {
	Url         string `json:"url"`
	Description string `json:"description,omitempty"` // def "", md support
}

func defaultServers() []Server {
	return nil
}

func newServers(ss *catalog.Servers) []Server {
	if ss.Len() == 0 {
		return defaultServers()
	}

	r := make([]Server, 0, ss.Len())
	_ = ss.Each(func(k string, v *catalog.Server) error {
		r = append(r, newServer(v))
		return nil
	})

	return r
}

func newServer(s *catalog.Server) Server {
	return Server{
		Url:         s.BaseUrl,
		Description: s.Annotation,
	}
}
