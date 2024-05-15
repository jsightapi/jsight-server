package openapi

import (
	"github.com/jsightapi/jsight-schema-core/notations/jschema"
	"github.com/jsightapi/jsight-schema-core/notations/regex"
	sc "github.com/jsightapi/jsight-schema-core/openapi"
)

type schemaObject interface {
	sc.SchemaObject
}

//nolint:unused
type schemaPropertyInfo interface {
	sc.PropertyInformer
}

type schemaObjectInfo interface {
	sc.ObjectInformer
}

type schemaInfo interface {
	sc.SchemaInformer
}

func getParamInfos(s *jschema.JSchema) ([]parameterInfo, Error) {
	si, err := getSchemaAsSingleObjectInfo(s)
	if err != nil {
		return nil, err
	}
	return schemaObjectInfoToParams(si), nil
}

func schemaObjectInfoToParams(si schemaObjectInfo) []parameterInfo {
	r := make([]parameterInfo, 0)
	for _, pi := range si.PropertiesInfos() {
		r = append(r, paramInfo{
			pi.Key(),
			pi.Optional(),
			pi.SchemaObject(),
			pi.Annotation(),
		})
	}
	return r
}

func getSchemaAsSingleObjectInfo(s *jschema.JSchema) (schemaObjectInfo, Error) {
	sd := dereferenceJSchema(s)
	if len(sd) > 1 {
		return nil, newErr("schema dereferences to multiple schemas (or-notation)")
	} else {
		i := sd[0]
		if i.Type() == sc.SchemaInfoTypeObject {
			return i.(sc.ObjectInformer), nil
		} else {
			return nil, newErr("schema is neither a single object, nor a reference to a single object")
		}
	}
}

func getJSchemaInfo(s *jschema.JSchema) schemaInfo {
	return sc.NewJSchemaInfo(s)
}

func getRSchemaInfo(s *regex.RSchema) schemaInfo {
	return sc.NewRSchemaInfo(s)
}

func dereferenceJSchema(s *jschema.JSchema) []sc.SchemaInformer {
	return sc.Dereference(s)
}
