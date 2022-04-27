package catalog

import (
	"encoding/json"
)

type ResourceMethodIdList struct {
	list []ResourceMethodId
}

func newResourceMethodIdList() *ResourceMethodIdList {
	return &ResourceMethodIdList{
		make([]ResourceMethodId, 0, 3),
	}
}

func (l *ResourceMethodIdList) append(r ResourceMethodId) {
	l.list = append(l.list, r)
}

func (l ResourceMethodIdList) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.list)
}
