package trapi2raml

import (
	"github.com/RangelReale/trapi"
)

func TRType(datatype trapi.DataType) string {
	switch datatype {
	case trapi.DATATYPE_NONE:
		return "any"
	case trapi.DATATYPE_STRING:
		return "string"
	case trapi.DATATYPE_NUMBER:
		return "number"
	case trapi.DATATYPE_INTEGER:
		return "integer"
	case trapi.DATATYPE_BOOLEAN:
		return "boolean"
	case trapi.DATATYPE_DATE:
		return "date-only"
	case trapi.DATATYPE_TIME:
		return "time-only"
	case trapi.DATATYPE_DATETIME:
		return "datetime"
	case trapi.DATATYPE_OBJECT:
		return "object"
	case trapi.DATATYPE_ARRAY:
		return "array"
	case trapi.DATATYPE_BINARY:
		return "any"
	case trapi.DATATYPE_CUSTOM:
		return "any"
	}
	return "any"
}
