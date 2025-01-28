package utils

import (
	"reflect"
	"strings"
)

func ExtractJSONTagsFromStructIntoString(inputStruct interface{}) string {
	jsonFields := ExtractJSONTagsFromStruct(inputStruct)
	return strings.Join(jsonFields, ",")
}

func ExtractJSONTagsFromStruct(inputStruct interface{}) []string {
	var jsonFields []string

	// Helper function to recursively extract JSON tags
	var extractFields func(reflect.Type)
	extractFields = func(t reflect.Type) {
		// If the type is a pointer, dereference it
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		// Ensure the type is a struct
		if t.Kind() != reflect.Struct {
			return
		}

		// Iterate through struct fields
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			// Check if the field is embedded
			if field.Anonymous {
				// Recursively process the embedded struct
				extractFields(field.Type)
				continue
			}

			// Get the JSON tag
			tag := field.Tag.Get("json")
			if tag != "" && tag != "-" {
				// Handle "omitempty" or other tags (split by comma)
				tagParts := strings.Split(tag, ",")
				jsonFields = append(jsonFields, tagParts[0])
			}
		}
	}

	// Start extracting fields from the input struct
	t := reflect.TypeOf(inputStruct)
	extractFields(t)

	return jsonFields
}
