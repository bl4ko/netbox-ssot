package utils

import (
	"encoding/json"
	"reflect"
	"strings"
)

// This function takes an object pointer, and returns a json body,
// that can be used to create that object in netbox API.
// This is essential because default marshal of the object
// isn't compatible with netbox API when attributes have nested
// objects.
func NetboxJSONMarshal(obj interface{}) ([]byte, error) {
	objMap := StructToNetboxJSONMap(obj)
	json, err := json.Marshal(objMap)
	return json, err
}

// Function that converts an object to a map[string]interface{}
// which can be used to create a json body for netbox API, especially
// for POST requests.
func StructToNetboxJSONMap(obj interface{}) map[string]interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	netboxJSONMap := make(map[string]interface{})
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
			diffMap := StructToNetboxJSONMap(fieldValue.Interface())
			for k, v := range diffMap {
				netboxJSONMap[k] = v
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
			netboxJSONMap[jsonTag] = sliceItems
		case reflect.Struct:
			if isChoiceEmbedded(fieldValue) {
				choiceValue := fieldValue.FieldByName("Value")
				if choiceValue.IsValid() {
					netboxJSONMap[jsonTag] = choiceValue.Interface()
				}
			} else {
				id := fieldValue.FieldByName("ID")
				if id.IsValid() {
					netboxJSONMap[jsonTag] = id.Int()
				} else {
					netboxJSONMap[jsonTag] = fieldValue.Interface()
				}
			}
		default:
			netboxJSONMap[jsonTag] = fieldValue.Interface()
		}
	}
	return netboxJSONMap
}
