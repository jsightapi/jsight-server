package openapi

type schemaInfoList struct {
	lst []SchemaInformer
}

func newSchemaInfoList() *schemaInfoList {
	return &schemaInfoList{
		lst: make([]SchemaInformer, 0, 5),
	}
}

func (l *schemaInfoList) append(r SchemaInformer) {
	l.lst = append(l.lst, r)
}

func (l schemaInfoList) list() []SchemaInformer {
	return l.lst
}
