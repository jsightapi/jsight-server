package openapi

import "github.com/jsightapi/jsight-api-core/catalog"

type mediaType string

// JSIght 0.3 supports: "json", "plainString", "binary"
const (
	MediaTypeRangeAny    mediaType = "*/*"
	MediaTypeJson        mediaType = "application/json"
	MediaTypeTextPlain   mediaType = "text/plain"
	MediaTypeOctetStream mediaType = "application/octet-stream" // TODO: discuss
)

// In JSight we currently can only provide one schema for a payload, so
// any payload regardles of media-type should be matched against this schema.
// In OAS terms it is described with a range wildcard "*\/*"
func formatToMediaType(f catalog.SerializeFormat) mediaType {
	switch f {
	case catalog.SerializeFormatJSON:
		return MediaTypeJson
	case catalog.SerializeFormatPlainString:
		return MediaTypeTextPlain
	case catalog.SerializeFormatBinary:
		return MediaTypeRangeAny
	default:
		return MediaTypeRangeAny
	}
}
