package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/common"
)

// Helper function to determine if a given reflect.Value contains an embedded common.Choice
func isChoiceEmbedded(v reflect.Value) bool {
	vType := v.Type()
	return vType.Field(0).Type == reflect.TypeOf(common.Choice{})
	// for i := 0; i < v.NumField(); i++ {
	// 	if vType.Field(i).Type == reflect.TypeOf(common.Choice{}) {
	// 		return true
	// 	}
	// }
	// return false
}

func choiceValue(v reflect.Value) interface{} {
	return v.Field(0).FieldByName("Value").Interface()
}

// Function that converts an object to a map[string]interface{}
// which can be used to create a json body for netbox API, especially
// for POST requests.
func StructToNetboxJsonMap(obj interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	netboxJsonMap := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := v.Type().Field(i)
		jsonTag := fieldType.Tag.Get("json")
		jsonTag = strings.Split(jsonTag, ",")[0]

		if fieldType.Name == "ID" {
			continue
		}

		// Special case when object inherits from NetboxObject
		if fieldType.Name == "NetboxObject" {
			diffMap, err := StructToNetboxJsonMap(fieldValue.Interface())
			if err != nil {
				return nil, fmt.Errorf("error procesing ObjToJsonMap when procerssing NetboxObject %s", err)
			}
			for k, v := range diffMap {
				netboxJsonMap[k] = v
			}
			continue
		}

		// If field is a pointer, we need to get the element it points to
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		// If a field is empty we skip it
		if !fieldValue.IsValid() || fieldValue.IsZero() {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Slice:
			if fieldValue.Len() == 0 {
				continue
			}
			sliceItems := make([]interface{}, 0)
			for j := 0; j < fieldValue.Len(); j++ {
				attribute := fieldValue.Index(j)
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
			netboxJsonMap[jsonTag] = sliceItems
		case reflect.Struct:
			if isChoiceEmbedded(fieldValue) {
				choiceValue := fieldValue.FieldByName("Value")
				if choiceValue.IsValid() {
					netboxJsonMap[jsonTag] = choiceValue.Interface()
				}
			} else {
				id := fieldValue.FieldByName("ID")
				if id.IsValid() {
					netboxJsonMap[jsonTag] = id.Int()
				} else {
					netboxJsonMap[jsonTag] = fieldValue.Interface()
				}
			}
		default:
			netboxJsonMap[jsonTag] = fieldValue.Interface()
		}
	}
	return netboxJsonMap, nil
}

// This function takes an object pointer, and returns a json body,
// that can be used to create that object in netbox API.
// This is essential because default marshal of the object
// isn't compatible with netbox API when attributes have nested
// objects.
func NetboxJsonMarshal(obj interface{}) ([]byte, error) {
	objMap, err := StructToNetboxJsonMap(obj)
	if err != nil {
		return nil, fmt.Errorf("error converting object to json map: %s", err)
	}
	return json.Marshal(objMap)
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

		if fieldName == "ID" {
			continue
		}

		// Custom logic for all objects that inherit from NetboxObject
		if fieldName == "NetboxObject" {
			netboxObjectDiffMap, err := JsonDiffMapExceptId(newObject.Field(i).Interface(), existingObject.Field(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("error procesing JsonDiffMapExceptId when procerssing NetboxObject %s", err)
			}
			for k, v := range netboxObjectDiffMap {
				diff[k] = v
			}
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
			err := addSliceDiffToMap(newObjectField, existingObjectField, jsonTag, diff)
			if err != nil {
				return nil, fmt.Errorf("error procesing JsonDiffMapExceptId when procerssing slice %s", err)
			}

		case reflect.Struct:
			err := addStructDiffToMap(newObjectField, existingObjectField, jsonTag, diff)
			if err != nil {
				return nil, fmt.Errorf("error procesing JsonDiffMapExceptId when procerssing struct %s", err)
			}

		case reflect.Map:
			err := addMapDiffToMap(newObjectField, existingObjectField, jsonTag, diff)
			if err != nil {
				return nil, fmt.Errorf("error procesing JsonDiffMapExceptId when procerssing map %s", err)
			}

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

func addSliceDiffToMap(newSlice reflect.Value, existingSlice reflect.Value, jsonTag string, diffMap map[string]interface{}) error {

	// If first slice is nil, that means that we reset the value
	if !newSlice.IsValid() || newSlice.Len() == 0 {
		if existingSlice.IsValid() && existingSlice.Len() > 0 {
			diffMap[jsonTag] = nil // reset the value
		}
		return nil
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
					return nil
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
					return nil
				}
			}
		}
	}
	return nil
}

// Returns json form for patching the difference e.g. { "id": 1 }
func addStructDiffToMap(newObj reflect.Value, existingObj reflect.Value, jsonTag string, diffMap map[string]interface{}) error {

	// If first struct is nil, that means that we reset the attribute to nil
	if !newObj.IsValid() {
		diffMap[jsonTag] = nil
		return nil
	}

	// We use ids for comparison between structs, because for patching objects, all we need is id of attribute
	idField := newObj.FieldByName("ID")

	// If objects don't have ID field, compare them by their values
	if !idField.IsValid() {
		if !existingObj.IsValid() {
			diffMap[jsonTag] = newObj.Interface()
		} else {
			if newObj.Interface() != existingObj.Interface() {
				if isChoiceEmbedded(newObj) {
					diffMap[jsonTag] = choiceValue(newObj)
				} else {
					diffMap[jsonTag] = newObj.Interface()
				}
			}
		}
	} else {
		if !existingObj.IsValid() {
			diffMap[jsonTag] = IDObject{ID: idField.Interface().(int)}
		} else {
			// Objects have ID field, compare their ids
			if newObj.FieldByName("ID").Interface() != existingObj.FieldByName("ID").Interface() {
				id := newObj.FieldByName("ID").Interface().(int)
				diffMap[jsonTag] = IDObject{ID: id}
			}
		}
	}
	return nil
}

func addMapDiffToMap(newMap reflect.Value, existingMap reflect.Value, jsonTag string, diffMap map[string]interface{}) error {

	// If first map is nil, that means that we reset the attribute to nil
	if !newMap.IsValid() {
		diffMap[jsonTag] = nil
		return nil
	}

	// Go through all keys in new map, and check if they are in existing map
	// If they are not, add them to diff map
	mapsDiff := make(map[string]interface{})
	for _, key := range newMap.MapKeys() {
		// Keys have to be strings
		if key.Kind() != reflect.String {
			return fmt.Errorf("map keys have to be strings. Not implemented for anything else yet")
		}
		if !existingMap.MapIndex(key).IsValid() {
			mapsDiff[key.Interface().(string)] = newMap.MapIndex(key).Interface()
		} else if newMap.MapIndex(key).Interface() != existingMap.MapIndex(key).Interface() {
			mapsDiff[key.Interface().(string)] = newMap.MapIndex(key).Interface()
		}
	}

	if len(mapsDiff) > 0 {
		diffMap[jsonTag] = mapsDiff
	}
	return nil
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

// Matches input string to a regex from input map, and returns the value
// If there is no match, it returns an empty string
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
// Slugified version can only contain: lowercase letters, numbers,
// underscores or hyphens.
// e.g. "My Name" -> "my-name"
// e.g. "  @Test@ " -> "test"
func Slugify(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")

	// Remove characters except lowercase letters, numbers, underscores, hyphens
	reg, _ := regexp.Compile("[^a-z0-9_-]+")
	name = reg.ReplaceAllString(name, "")
	return name
}
