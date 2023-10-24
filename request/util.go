package request

import "strings"

// IsJSONType method is to check JSON content type or not
func IsJSONType(ct string) bool {
	return strings.Contains(ct, "json")
}

// IsXMLType method is to check XML content type or not
func IsXMLType(ct string) bool {
	return strings.Contains(ct, "xml")
}
