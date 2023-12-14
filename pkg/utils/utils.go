package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// This function takes an object pointer, and returns a json body,
// that can be used to create that object in netbox API.
// This is essential because default marshal of the object
// isn't compatible with netbox API when attributes have nested
// objects.
func NetboxJsonMarshal(obj interface{}) ([]byte, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("object must be a pointer to a struct")
	}
	v = v.Elem()

	netboxJson := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		fieldName := fieldType.Name
		jsonTag := fieldType.Tag.Get("json")
		jsonTag = strings.Split(jsonTag, ",")[0]

		if fieldName == "ID" {
			continue
		}

		// If field is a pointer, we need to get the element it points to
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		// If it is a nil pointer, we need to set it to nil in json
		if !field.IsValid() {
			netboxJson[jsonTag] = nil
			continue
		}

		switch field.Kind() {
		case reflect.Slice:
			if field.Len() == 0 {
				continue
			}
			sliceItems := make([]interface{}, 0)
			for j := 0; j < field.Len(); j++ {
				attribute := field.Index(j)
				if attribute.Kind() == reflect.Ptr {
					attribute = attribute.Elem()
				}
				if attribute.Kind() == reflect.Struct {
					id := attribute.FieldByName("ID")
					if id.IsValid() && id.Int() != 0 {
						sliceItems = append(sliceItems, id.Int())
					} else {
						sliceItems = append(sliceItems, attribute.Interface())
					}
				} else {
					sliceItems = append(sliceItems, attribute.Interface())
				}
			}
			netboxJson[jsonTag] = sliceItems
		case reflect.Struct:
			id := field.FieldByName("ID")
			if id.IsValid() {
				netboxJson[jsonTag] = id.Int()
			} else {
				if fieldName == "Status" {
					status := field.FieldByName("Value")
					netboxJson[jsonTag] = status.String()
				} else {
					netboxJson[jsonTag] = field.Interface()
				}
			}
		default:
			netboxJson[jsonTag] = field.Interface()
		}
	}

	return json.Marshal(netboxJson)
}

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

	// There are going to be 2 kinds of comparison.
	// One where slice will contain objects, in that case we
	// will compare ids of the objects, else we will compare
	// the values of the slice
	switch newSlice.Index(0).Kind() {
	case reflect.String:
		strSet := make(map[string]bool)
		for j := 0; j < newSlice.Len(); j++ {
			strSet[newSlice.Index(j).Interface().(string)] = true
		}
		if len(strSet) != existingSlice.Len() {
			diffMap[jsonTag] = newSlice.Interface()
		} else {
			for j := 0; j < existingSlice.Len(); j++ {
				if !strSet[existingSlice.Index(j).Interface().(string)] {
					diffMap[jsonTag] = newSlice.Interface()
					return
				}
			}
		}

	default:
		idObjectsSet := make(map[int]bool, newSlice.Len())
		idObjectsArr := make([]IDObject, 0, newSlice.Len())
		var id int
		for j := 0; j < newSlice.Len(); j++ {
			if newSlice.Index(j).IsNil() {
				continue
			}
			if newSlice.Index(j).Kind() == reflect.Ptr {
				id = newSlice.Index(j).Elem().FieldByName("ID").Interface().(int)
				idObjectsSet[id] = true
				idObjectsArr = append(idObjectsArr, IDObject{ID: id})
			} else {
				id = newSlice.Index(j).FieldByName("ID").Interface().(int)
				idObjectsSet[id] = true
				idObjectsArr = append(idObjectsArr, IDObject{ID: id})
			}
		}

		if len(idObjectsSet) != existingSlice.Len() {
			diffMap[jsonTag] = idObjectsArr
		} else {
			for j := 0; j < existingSlice.Len(); j++ {
				if existingSlice.Index(j).Kind() == reflect.Ptr {
					id = existingSlice.Index(j).Elem().FieldByName("ID").Interface().(int)
				} else {
					id = existingSlice.Index(j).FieldByName("ID").Interface().(int)
				}
				if _, ok := idObjectsSet[id]; !ok {
					diffMap[jsonTag] = idObjectsArr
					return
				}
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
