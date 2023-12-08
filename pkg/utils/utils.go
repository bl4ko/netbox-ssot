package utils

import (
	"fmt"
	"reflect"
	"regexp"
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

// 	return len(fields) > 0, fields
// }

// JsonDiffMapExceptID checks if any field except ID field, from two objects, are different.
// It returns a map of fields (represented by their JSON tag names) that are different with their values from newObj.
func JsonDiffMapExceptId(newObj, existingObj interface{}) (map[string]interface{}, error) {
	diff := make(map[string]interface{})

	val1 := reflect.ValueOf(newObj)
	val2 := reflect.ValueOf(existingObj)

	// Check if the values are pointers and get the element they point to
	if val1.Kind() == reflect.Ptr {
		val1 = val1.Elem()
	}
	if val2.Kind() == reflect.Ptr {
		val2 = val2.Elem()
	}

	// Ensure that we are dealing with structs
	if val1.Kind() != reflect.Struct || val2.Kind() != reflect.Struct {
		return nil, fmt.Errorf("arguments are not structs")
	}

	for i := 0; i < val1.NumField(); i++ {
		field := val1.Type().Field(i)
		fieldName := field.Name
		jsonTag := field.Tag.Get("json")

		if fieldName == "ID" {
			continue
		}

		// Handle fields without a json tag or with "-" as json tag
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = fieldName
		} else {
			// Extract the name part before the comma, if any
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		// If Obj's Attribute is a slice we comapre id's in that slice
		if val1.Field(i).Kind() == reflect.Slice {
			if val1.Field(i).Len() != val2.Field(i).Len() {
				diff[jsonTag] = val1.Field(i).Interface()
			} else {
				firstIdSet := map[int]bool{}
				secondIdSet := map[int]bool{}
				for j := 0; j < val1.Field(i).Len(); j++ {
					// Elements can be pointers, so we need to get the element they point to
					firstElem := val1.Field(i).Index(j)
					secondElem := val2.Field(i).Index(j)
					if firstElem.Kind() == reflect.Ptr {
						firstElem = firstElem.Elem()
					}
					if secondElem.Kind() == reflect.Ptr {
						secondElem = secondElem.Elem()
					}
					firstIdSet[firstElem.FieldByName("ID").Interface().(int)] = true
					secondIdSet[secondElem.FieldByName("ID").Interface().(int)] = true
				}
				if !reflect.DeepEqual(firstIdSet, secondIdSet) {
					diff[jsonTag] = val1.Field(i).Interface()
				}
			}
			// Else we compare the values of the fields
		} else if val1.Field(i).Interface() != val2.Field(i).Interface() {
			diff[jsonTag] = val1.Field(i).Interface()
		}
	}

	return diff, nil
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
