package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

// // Function that checks if any field from two objects are different
// // If so it returns true and a slice of fields that are different
// func ObjDiff(obj1, obj2 interface{}) (bool, []string) {

// 	var fields []string

// 	val1 := reflect.ValueOf(obj1).Elem()
// 	val2 := reflect.ValueOf(obj2).Elem()

// 	for i := 0; i < val1.NumField(); i++ {
// 		if val1.Field(i).Interface() != val2.Field(i).Interface() {
// 			fields = append(fields, val1.Type().Field(i).Name)
// 		}
// 	}

// 	return len(fields) > 0, fields
// }

// // Function that checks if any field except ID field, from two objects are different
// // If so it returns true and a slice of fields that are different
// func ObjDiffExceptID(obj1, obj2 interface{}) (bool, []string) {

// 	var fields []string

// 	val1 := reflect.ValueOf(obj1).Elem()
// 	val2 := reflect.ValueOf(obj2).Elem()

// 	for i := 0; i < val1.NumField(); i++ {
// 		if val1.Field(i).Interface() != val2.Field(i).Interface() {
// 			if val1.Type().Field(i).Name != "ID" {
// 				fields = append(fields, val1.Type().Field(i).Name)
// 			}
// 		}
// 	}

//		return len(fields) > 0, fields
//	}
//
// Struct used for patching objects, when attributes are structs, or slices
// we only need object with an id for patching.
type IDObject struct {
	ID int `json:"id"`
}

// JsonDiffMapExceptID checks if any field except ID field, from two objects, are different.
// It returns a map of fields (represented by their JSON tag names) that are different with their values from newObj.
func JsonDiffMapExceptId(newObj, existingObj interface{}) (map[string]interface{}, error) {
	diff := make(map[string]interface{})

	newObject := reflect.ValueOf(newObj)
	existingObject := reflect.ValueOf(existingObj)

	// Ensure that both objects are of the same kind (e.g. struct)
	if newObject.Kind() != existingObject.Kind() {
		return nil, fmt.Errorf("arguments are not of the same type")
	}

	// Check if the values are pointers and get the element they point to
	if newObject.Kind() == reflect.Ptr {
		newObject = newObject.Elem()
		existingObject = existingObject.Elem()
	}

	// Ensure that we are dealing with structs
	if newObject.Kind() != reflect.Struct {
		return nil, fmt.Errorf("arguments are not structs")
	}

	for i := 0; i < newObject.NumField(); i++ {
		fieldName := newObject.Type().Field(i).Name
		jsonTag := newObject.Type().Field(i).Tag.Get("json")

		// We skip the ID field, because newly created objects, won't have this field, which is netbox specific
		if fieldName == "ID" {
			continue
		}

		// Get json tag, so we can get the json field name
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = fieldName
		} else {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		// Ensure that both fields are of the same kind
		if newObject.Field(i).Kind() != existingObject.Field(i).Kind() {
			return nil, fmt.Errorf("field %s is not of the same type in both objects", jsonTag)
		}

		// Check if elements are pointers, in that case get the elements they are pointing to
		newObjectField := newObject.Field(i)
		existingObjectField := existingObject.Field(i)
		if newObjectField.Kind() == reflect.Ptr {
			newObjectField = newObjectField.Elem()
			existingObjectField = existingObjectField.Elem()
		}

		switch newObjectField.Kind() {
		case reflect.Slice:
			addDiffSliceToMap(newObjectField, existingObjectField, jsonTag, diff)

		case reflect.Struct:
			addDiffStructToMap(newObjectField, existingObjectField, jsonTag, diff)

		default:
			if !newObjectField.IsValid() {
				if existingObjectField.IsValid() {
					diff[jsonTag] = nil
				}
				continue
			}
			if newObjectField.Interface() != existingObjectField.Interface() {
				diff[jsonTag] = newObjectField.Interface()
			}
		}
	}

	return diff, nil
}

func addDiffSliceToMap(newSlice reflect.Value, existingSlice reflect.Value, jsonTag string, diffMap map[string]interface{}) {

	// If first slice is nil, that means that we reset the value
	if !newSlice.IsValid() {
		if existingSlice.IsValid() {
			diffMap[jsonTag] = nil // reset the value
		}
		return
	}

	idObjects := make([]IDObject, 0, newSlice.Len())
	var id int
	for j := 0; j < newSlice.Len(); j++ {
		if newSlice.Index(j).Kind() == reflect.Ptr {
			id = newSlice.Index(j).Elem().FieldByName("ID").Interface().(int)
			idObjects = append(idObjects, IDObject{ID: id})
		} else {
			id = newSlice.Index(j).FieldByName("ID").Interface().(int)
			idObjects = append(idObjects, IDObject{ID: id})
		}
	}
	// We always store the IDs in ascending order, because netbox api
	// returns them in ascending order
	slices.SortFunc(idObjects, func(i IDObject, j IDObject) int {
		return i.ID - j.ID
	})

	if newSlice.Len() != existingSlice.Len() {
		diffMap[jsonTag] = idObjects
	} else {
		for j := 0; j < existingSlice.Len(); j++ {
			if existingSlice.Index(j).Kind() == reflect.Ptr {
				id = existingSlice.Index(j).Elem().FieldByName("ID").Interface().(int)
			} else {
				id = existingSlice.Index(j).FieldByName("ID").Interface().(int)
			}
			if id != idObjects[j].ID {
				diffMap[jsonTag] = idObjects
				return
			}
		}
	}
}

// Returns json form for patching the difference e.g. { "id": 1 }
func addDiffStructToMap(newObj reflect.Value, existingObj reflect.Value, jsonTag string, diffMap map[string]interface{}) {

	// If first struct is nil, that means that we reset the attribute to nil
	if !newObj.IsValid() {
		diffMap[jsonTag] = nil
		return
	}

	// We use ids for comparison between structs, because for patching objects, all we need is id of attribute
	idField := newObj.FieldByName("ID")

	if !idField.IsValid() {
		// Both objects doesn't have ID field, compare them directly
		if newObj.Interface() != existingObj.Interface() {
			diffMap[jsonTag] = newObj.Interface()
		}
	} else {
		// Objects have ID field, compare their ids
		if newObj.FieldByName("ID").Interface() != existingObj.FieldByName("ID").Interface() {
			id := newObj.FieldByName("ID").Interface().(int)
			diffMap[jsonTag] = IDObject{ID: id}
		}
	}
}

// Validates array of regex relations
// Regex relation is a string of format "regex = value"
func ValidateRegexRelations(regexRelations []string) error {
	for _, regexRelation := range regexRelations {
		relation := strings.Split(regexRelation, "=")
		if len(relation) != 2 {
			return fmt.Errorf("invalid regex relation: %s. Should be of format: regex = value", regexRelation)
		}
		regexStr := strings.TrimSpace(relation[0])
		_, err := regexp.Compile(regexStr)
		if err != nil {
			return fmt.Errorf("invalid regex: %s, in relation: %s", regexStr, regexRelation)
		}
	}
	return nil
}

// Converts array of strings, that are of form "regex = value", to a map
// where key is regex and value is value
func ConvertStringsToRegexPairs(input []string) map[string]string {
	output := make(map[string]string, len(input))
	for _, s := range input {
		pair := strings.Split(s, "=")
		output[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
	}
	return output
}

func MatchStringToValue(input string, patterns map[string]string) (string, error) {
	for regex, value := range patterns {
		matched, err := regexp.MatchString(regex, input)
		if err != nil {
			return "", err // Handle regex compilation error
		}
		if matched {
			return value, nil
		}
	}
	return "", nil // Return an empty string or an error if no match is found
}

// Converts string name to its slugified version.
// e.g. "My Name" -> "my-name"
// e.g. "   Test  " -> "test"
func Slugify(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	return name
}
