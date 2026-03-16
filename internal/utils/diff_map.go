package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/src-doo/netbox-ssot/internal/constants"
	"github.com/src-doo/netbox-ssot/internal/netbox/objects"
)

// Helper function to determine if a given reflect.Value contains an embedded objects.Choice.
// We assume that Choice attribute is always the first attribute of an object.
func isChoiceEmbedded(v reflect.Value) bool {
	vType := v.Type()
	return vType.Field(0).Type == reflect.TypeOf(objects.Choice{})
	// for i := 0; i < v.NumField(); i++ {
	// 	if vType.Field(i).Type == reflect.TypeOf(objects.Choice{}) {
	// 		return true
	// 	}
	// }
	// return false
}

func choiceValue(v reflect.Value) interface{} {
	return v.Field(0).FieldByName("Value").Interface()
}

// Struct used for patching objects, when attributes are structs, or slices
// we only need object with an ID of original object for patching.
type IDObject struct {
	ID int `json:"id"`
}

// hasPriorityOver returns true if newObj has priority over existingObj, false otherwise.
// newObject will always have priority over exisitngObject, unless following cases:
//  1. When sourcePriority[newObj] > sourcePriority[existingObj]
//  2. When the new IP address is an arp entry (custom field ArpEntry),
//     and the exisiting is not. If both we follow case 1. .
func hasPriorityOver(newObj, existingObj reflect.Value, source2priority map[string]int) bool {
	// Retrieve the SourceName field from CustomFields for both objects
	newObjCustomFields := newObj.FieldByName("CustomFields")
	existingObjCustomFields := existingObj.FieldByName("CustomFields")

	// Check if fields are valid and present in the sourcePriority map
	if newObjCustomFields.IsValid() && existingObjCustomFields.IsValid() {
		if newCustomFields, ok := newObjCustomFields.Interface().(map[string]interface{}); ok {
			if existingCustomFields, ok := existingObjCustomFields.Interface().(map[string]interface{}); ok {
				// 2. case
				if newCustomFields[constants.CustomFieldArpEntryName] != existingCustomFields[constants.CustomFieldArpEntryName] {
					if newCustomFields[constants.CustomFieldArpEntryName] != nil {
						return !newCustomFields[constants.CustomFieldArpEntryName].(bool) //nolint:forcetypeassert
					} else if existingCustomFields[constants.CustomFieldArpEntryName] == true {
						return true
					}
				}

				// 1. case
				if newCustomFields[constants.CustomFieldSourceName] != nil &&
					existingCustomFields[constants.CustomFieldSourceName] != nil {
					newPriority := int(^uint(0) >> 1) // max int
					if priority, newOk := source2priority[newCustomFields[constants.CustomFieldSourceName].(string)]; newOk {
						newPriority = priority
					}
					existingPriority := int(^uint(0) >> 1)
					//nolint:lll
					if priority, existingOk := source2priority[existingCustomFields[constants.CustomFieldSourceName].(string)]; existingOk {
						existingPriority = priority
					}
					// In case newPriority is lower or equal than existingPriority
					// newObj has precedence over exsitingObj
					return newPriority <= existingPriority
				}
			}
		}
	}

	return true
}

// JSONDiffMapExceptID compares two objects and returns a map of fields
// (represented by their JSON tag names) that are different with their
// values from newObj.
// If resetFields is set to true, the function will also include fields
// that are empty in newObj but might have a value in existingObj.
// Also we check for priority, if newObject has priority over existingObject
// we use the fields from newObject, otherwise we use the fields from exisingObject.
func JSONDiffMapExceptID(
	newObj, existingObj interface{},
	resetFields bool,
	source2priority map[string]int,
) (map[string]interface{}, error) {
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

	// Check for priority
	hasPriority := hasPriorityOver(newObject, existingObject, source2priority)

	for i := 0; i < newObject.NumField(); i++ {
		fieldName := newObject.Type().Field(i).Name
		jsonTag := newObject.Type().Field(i).Tag.Get("json")

		if fieldName == "ID" {
			continue
		}

		// Custom logic for all objects that inherit from NetboxObject
		if fieldName == "NetboxObject" {
			netboxObjectDiffMap, err := JSONDiffMapExceptID(
				newObject.Field(i).Interface(),
				existingObject.Field(i).Interface(),
				resetFields,
				source2priority,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"error processing JsonDiffMapExceptID when processing NetboxObject %s",
					err,
				)
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

		// If reset is set to false and the newObjectField is empty,
		// we don't do anything (like omitempty), see tests
		if !resetFields && (!newObjectField.IsValid() || newObjectField.IsZero()) {
			continue
		}

		switch newObjectField.Kind() {
		// Reset the field (when it is set to nil),
		// this only happens if flag resetFields is set to true.
		case reflect.Invalid:
			if existingObjectField.IsValid() {
				diff[jsonTag] = nil
			}

		case reflect.Slice:
			err := addSliceDiff(newObjectField, existingObjectField, jsonTag, hasPriority, diff)
			if err != nil {
				return nil, fmt.Errorf(
					"error processing JsonDiffMapExceptID when processing slice %s",
					err,
				)
			}

		case reflect.Struct:
			err := addStructDiff(newObjectField, existingObjectField, jsonTag, hasPriority, diff)
			if err != nil {
				return nil, fmt.Errorf(
					"error processing JsonDiffMapExceptID when processing struct %s",
					err,
				)
			}

		case reflect.Map:
			err := addMapDiff(newObjectField, existingObjectField, jsonTag, hasPriority, diff)
			if err != nil {
				return nil, fmt.Errorf(
					"error processing JsonDiffMapExceptID when processing map %s",
					err,
				)
			}

		default:
			addPrimaryDiff(newObjectField, existingObjectField, jsonTag, hasPriority, diff)
		}
	}

	return diff, nil
}

// mergeTagSlices merges two slices of *objects.Tag.
// It keeps all tags from existingSlice that are NOT managed by netbox-ssot
// (i.e. tags whose name does not start with "Source:"),
// and replaces/adds all tags from newSlice.
// Returns the merged slice as []IDObject, whether it changed, and an error.
func mergeTagSlices(
	newSlice reflect.Value,
	existingSlice reflect.Value,
) ([]IDObject, bool, error) {
	// Collect new tags into a map by ID for deduplication
	newTagsByID := make(map[int]*objects.Tag)
	for i := 0; i < newSlice.Len(); i++ {
		elem := newSlice.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		tag, ok := newSlice.Index(i).Interface().(*objects.Tag)
		if !ok {
			// Try non-pointer
			tagVal, ok2 := elem.Interface().(objects.Tag)
			if !ok2 {
				return nil, false, fmt.Errorf("tag slice contains non-Tag element")
			}
			newTagsByID[tagVal.ID] = &tagVal
		} else {
			newTagsByID[tag.ID] = tag
		}
	}

	// Collect IDs of new tags for quick lookup
	newTagIDs := make(map[int]bool, len(newTagsByID))
	for id := range newTagsByID {
		newTagIDs[id] = true
	}

	// Build merged list: start with new tags
	mergedIDs := make(map[int]bool)
	result := make([]IDObject, 0)
	for id := range newTagsByID {
		mergedIDs[id] = true
		result = append(result, IDObject{ID: id})
	}

	// Add existing tags that are NOT managed by netbox-ssot source tagging
	// A tag is considered "managed by source" if its name starts with "Source: "
	if existingSlice.IsValid() {
		for i := 0; i < existingSlice.Len(); i++ {
			elem := existingSlice.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			var tag objects.Tag
			switch v := elem.Interface().(type) {
			case *objects.Tag:
				tag = *v
			case objects.Tag:
				tag = v
			default:
				return nil, false, fmt.Errorf("existing tag slice contains non-Tag element")
			}
			isManagedBySource := strings.HasPrefix(tag.Name, "Source: ") ||
				tag.Name == constants.SsotTagName ||
				tag.Name == constants.OrphanTagName ||
				tag.Name == constants.IgnoreDeviceTypeTagName
			if !isManagedBySource && !mergedIDs[tag.ID] {
				mergedIDs[tag.ID] = true
				result = append(result, IDObject{ID: tag.ID})
			}
		}
	}

	existingIDs := make(map[int]bool)
	if existingSlice.IsValid() {
		for i := 0; i < existingSlice.Len(); i++ {
			elem := existingSlice.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			idField := elem.FieldByName("ID")
			if idField.IsValid() {
				existingIDs[int(idField.Int())] = true
			}
		}
	}
	changed := !reflect.DeepEqual(mergedIDs, existingIDs)

	return result, changed, nil
}

// Function that takes two objects (of type slice) and returns a map
// that can be easily used with json.Marshal
// To achieve this the map is of the following format:
// map[jsonTag] = [id1, id2, id3] // If the slice contains objects with ID field
// map[jsonTag] = [value1, value2, value3] // If the slice contains strings.
//
// This function only works for slices with comparable elements (e.g. ints, strings) and
// slices that contain objects that are structs but have ID attribute.
func addSliceDiff(
	newSlice reflect.Value,
	existingSlice reflect.Value,
	jsonTag string,
	hasPriority bool,
	diffMap map[string]interface{},
) error {
	if !hasPriority {
		return nil
	}

	if !newSlice.IsValid() || newSlice.Len() == 0 {
		if existingSlice.IsValid() && existingSlice.Len() > 0 {
			diffMap[jsonTag] = []interface{}{}
		}
		return nil
	}

	// Merge tags
	if jsonTag == "tags" {
		mergedTags, changed, err := mergeTagSlices(newSlice, existingSlice)
		if err != nil {
			return fmt.Errorf("merge tag slices: %s", err)
		}
		if changed {
			diffMap[jsonTag] = mergedTags
		}
		return nil
	}

	newSlice, err := convertSliceToComparableSlice(newSlice)
	if err != nil {
		return fmt.Errorf("error converting slice to comparable slice: %s", err)
	}
	if !existingSlice.IsValid() || existingSlice.Len() == 0 {
		diffMap[jsonTag] = newSlice.Interface()
		return nil
	}
	existingSlice, err = convertSliceToComparableSlice(existingSlice)
	if err != nil {
		return fmt.Errorf("error converting slice to comparable slice: %s", err)
	}
	newSet := sliceToSet(newSlice)
	existingSet := sliceToSet(existingSlice)
	if !reflect.DeepEqual(newSet, existingSet) {
		diffMap[jsonTag] = newSlice.Interface()
	}
	return nil
}

// Converts slice of structs with IDs to slice containing only ids. If slice
// contains comparable elements don't do anything.
func convertSliceToComparableSlice(slice reflect.Value) (reflect.Value, error) {
	// We determine the types of elements of the slice by checking the first element.
	firstElement := slice.Index(0)
	if firstElement.Kind() == reflect.Pointer {
		firstElement = firstElement.Elem()
	}
	if firstElement.Kind() == reflect.Struct {
		if !firstElement.FieldByName("ID").IsValid() {
			return reflect.ValueOf(
					nil,
				), fmt.Errorf(
					"slice contains struct that don't contain id field",
				)
		}
		idSlice := make([]int, 0)
		for i := 0; i < slice.Len(); i++ {
			element := slice.Index(i)
			if element.Kind() == reflect.Ptr {
				element = element.Elem()
			}
			idField := element.FieldByName("ID").Interface()
			switch value := idField.(type) {
			case int:
				idSlice = append(idSlice, value)
			default:
				return reflect.ValueOf(nil), fmt.Errorf("id is not int")
			}
		}
		return reflect.ValueOf(idSlice), nil
	}

	return slice, nil
}

// Converts slice to a set.
func sliceToSet(slice reflect.Value) map[interface{}]bool {
	set := make(map[interface{}]bool)
	for i := 0; i < slice.Len(); i++ {
		element := slice.Index(i)
		if element.Kind() == reflect.Ptr {
			element = element.Elem()
		}
		set[element.Interface()] = true
	}
	return set
}

// Returns json form for patching the difference e.g. { "ID": 1 }.
func addStructDiff(
	newObj reflect.Value,
	existingObj reflect.Value,
	jsonTag string,
	hasPriority bool,
	diffMap map[string]interface{},
) error {
	// If first struct is nil, that means that we reset the attribute to nil
	if !newObj.IsValid() {
		diffMap[jsonTag] = nil
		return nil
	}

	// We check if struct is a objects.Choice (special netbox struct)
	if isChoiceEmbedded(newObj) {
		if !existingObj.IsValid() || newObj.Interface() != existingObj.Interface() {
			diffMap[jsonTag] = choiceValue(newObj)
		}
		return nil
	}

	// We use ids for comparison between structs, because for patching objects, all we need is id of attribute
	idField := newObj.FieldByName("ID")

	// If objects don't have ID field, compare them by their values
	if !idField.IsValid() {
		if !existingObj.IsValid() {
			diffMap[jsonTag] = newObj.Interface()
		} else if newObj.Interface() != existingObj.Interface() {
			if hasPriority {
				diffMap[jsonTag] = newObj.Interface()
			}
		}
	} else {
		if !existingObj.IsValid() {
			idValue, ok := idField.Interface().(int)
			if !ok {
				return fmt.Errorf("id field is not an int")
			}
			diffMap[jsonTag] = IDObject{ID: idValue}
		} else if newObj.FieldByName("ID").Interface() != existingObj.FieldByName("ID").Interface() {
			// Objects have ID field, compare their ids
			idValue, ok := idField.Interface().(int)
			if !ok {
				return fmt.Errorf("id field is not an int")
			}
			diffMap[jsonTag] = IDObject{ID: idValue}
		}
	}
	return nil
}

func addMapDiff(
	newMap reflect.Value,
	existingMap reflect.Value,
	jsonTag string,
	hasPriority bool,
	diffMap map[string]interface{},
) error {
	// If the new map is not set, we don't change anything
	if !newMap.IsValid() {
		return nil
	}

	// Go through all keys in new map, and check if they are in existing map
	// If they are not, add them to diff map
	mapsDiff := make(map[string]interface{})
	for _, key := range newMap.MapKeys() {
		// Keys have to be strings
		if keyValue, ok := key.Interface().(string); ok {
			if !existingMap.MapIndex(key).IsValid() {
				mapsDiff[keyValue] = newMap.MapIndex(key).Interface()
			} else if !reflect.DeepEqual(newMap.MapIndex(key).Interface(), existingMap.MapIndex(key).Interface()) {
				if hasPriority {
					mapsDiff[keyValue] = newMap.MapIndex(key).Interface()
				}
			}
		} else {
			return fmt.Errorf("map keys have to be strings. Not implemented for anything else yet")
		}
	}

	if len(mapsDiff) > 0 {
		for _, key := range existingMap.MapKeys() {
			if keyValue, ok := key.Interface().(string); ok {
				if !newMap.MapIndex(key).IsValid() {
					if !existingMap.MapIndex(key).IsNil() {
						mapsDiff[keyValue] = existingMap.MapIndex(key).Interface()
					}
				}
			}
		}
		diffMap[jsonTag] = mapsDiff
	}
	return nil
}

func addPrimaryDiff(
	newField reflect.Value,
	existingField reflect.Value,
	jsonTag string,
	hasPriority bool,
	diffMap map[string]interface{},
) {
	switch {
	case newField.IsZero():
		if !existingField.IsZero() {
			diffMap[jsonTag] = reflect.Zero(newField.Type()).Interface()
		}
	case existingField.IsZero():
		diffMap[jsonTag] = newField.Interface()
	case newField.Interface() != existingField.Interface():
		if hasPriority {
			diffMap[jsonTag] = newField.Interface()
		}
	}
}

func ExtractFieldsFromDiffMap(
	diffMap map[string]interface{},
	field []string,
) map[string]interface{} {
	extractedFields := make(map[string]interface{})
	if len(diffMap) == 0 {
		return extractedFields
	}
	for _, f := range field {
		if value, ok := diffMap[f]; ok {
			extractedFields[f] = value
		}
	}
	return extractedFields
}
