package core

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
)

func removeDirectiveFromSlice(slice []*directive.Directive, i int) []*directive.Directive {
	return append(slice[:i], slice[i+1:]...)
}
