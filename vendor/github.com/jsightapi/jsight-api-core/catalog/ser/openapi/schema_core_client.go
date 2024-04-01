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

func getParamInfo(s *jschema.JSchema) []parameterInfo {
	r := make([]parameterInfo, 0)
	schemaInfos := dereferenceJSchema(s)
	if len(schemaInfos) > 1 {
		panic("or-references conversion not supported for parameter directives")
	} else {
		si := schemaInfos[0]
		switch si.Type() {
		case sc.SchemaInfoTypeObject: // TODO: get rid of sc. import?
			properties := si.(sc.ObjectInformer).PropertiesInfos()
			for _, pi := range properties {
				r = append(r, paramInfo{
					pi.Key(),
					pi.Optional(),
					pi.SchemaObject(),
					pi.Annotation(),
				})
			}
		default:
			panic("parameters directive's schema is not an object")
		}
	}

	return r
}

func getSchemaObjectInfo(s *jschema.JSchema) schemaObjectInfo {
	sd := dereferenceJSchema(s)
	if len(sd) > 1 {
		panic("or-references not supported")
	} else {
		ei := sd[0]
		if ei.Type() == sc.SchemaInfoTypeObject {
			return ei.(sc.ObjectInformer)
		} else {
			panic("schema is not an object")
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
