package objects

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

type Tag struct {
	ID          int             `json:"id,omitempty"`
	Name        string          `json:"name,omitempty"`
	Slug        string          `json:"slug,omitempty"`
	Color       constants.Color `json:"color,omitempty"`
	Description string          `json:"description,omitempty"`
}

func (t Tag) String() string {
	return fmt.Sprintf("Tag{Name: %s}", t.Name)
}

// Tag implements IDItem interface.
func (t *Tag) GetID() int {
	return t.ID
}
func (t *Tag) GetObjectType() constants.ContentType {
	return constants.ContentTypeExtrasTag
}

func (t *Tag) GetAPIPath() constants.APIPath {
	return constants.TagsAPIPath
}

// CustomFieldTypes are predefined netbox's types for CustomFields.
type CustomFieldType struct {
	Choice
}

// Predefined netbox's types for CustomFields
// https://github.com/netbox-community/netbox/blob/35be4f05ef376e28d9af4d7245ba10cc286bb62a/netbox/extras/choices.py#L10
var (
	CustomFieldTypeText     = CustomFieldType{Choice{Value: "text", Label: "Text"}}
	CustomFieldTypeLongText = CustomFieldType{Choice{Value: "longtext", Label: "Text (long)"}}
	CustomFieldTypeInteger  = CustomFieldType{Choice{Value: "integer", Label: "Integer"}}
	CustomFieldTypeDecimal  = CustomFieldType{Choice{Value: "decimal", Label: "Decimal"}}
	CustomFieldTypeBoolean  = CustomFieldType{
		Choice{Value: "boolean", Label: "Boolean (true/false)"},
	}
	CustomFieldTypeDate = CustomFieldType{Choice{Value: "date", Label: "Date"}}
)

type FilterLogic struct {
	Choice
}

var (
	FilterLogicLoose = FilterLogic{Choice{Value: "loose", Label: "Loose"}}
)

type CustomFieldUIVisible struct {
	Choice
}

var (
	CustomFieldUIVisibleAlways = CustomFieldUIVisible{Choice{Value: "always", Label: "Always"}}
	CustomFieldUIVisibleIfSet  = CustomFieldUIVisible{Choice{Value: "if-set", Label: "If set"}}
	CustomFieldUIVisibleHidden = CustomFieldUIVisible{Choice{Value: "hidden", Label: "Hidden"}}
)

type CustomFieldUIEditable struct {
	Choice
}

var (
	CustomFieldUIEditableYes    = CustomFieldUIEditable{Choice{Value: "yes", Label: "Yes"}}
	CustomFieldUIEditableNo     = CustomFieldUIEditable{Choice{Value: "no", Label: "No"}}
	CustomFieldUIEditableHidden = CustomFieldUIEditable{Choice{Value: "hidden", Label: "Hidden"}}
)

const (
	DisplayWeightDefault = 100
	SearchWeightDefault  = 1000
)

type CustomField struct {
	ID int `json:"id,omitempty"`
	// Name of the custom field (e.g. host_cpu_cores). This field is required.
	Name string `json:"name,omitempty"`
	// Label represents name of the field as displayed to users (e.g. Physical CPU cores).
	// If not provided, the name will be used instead.
	Label string `json:"label,omitempty"`
	// Type is the type of the custom field.
	// Valid choices are: text, integer, boolean, date, url, select, multiselect. This field is required.
	Type CustomFieldType `json:"type,omitempty"`
	// Type of the related object (for object/multi-object fields only) (e.g. objects.device). This field is required.
	ObjectTypes []constants.ContentType `json:"object_types,omitempty"`
	// Description is a description of the field. This field is optional.
	Description string `json:"description,omitempty"`
	// Weighting for search. Lower values are considered more important. Default (1000).
	SearchWeight int `json:"search_weight,omitempty"`
	// Filter logic. This field is required. (Default loose).
	FilterLogic FilterLogic `json:"filter_logic,omitempty"`
	// UI visible. This field is required. (Default read-write).
	CustomFieldUIVisible *CustomFieldUIVisible `json:"ui_visible,omitempty"`
	// UI editable. This field is required. (Default read-write).
	CustomFieldUIEditable *CustomFieldUIEditable `json:"ui_editable,omitempty"`
	// Display Weight. Fields with higher weights appear lower in a form (default is 100).
	DisplayWeight int `json:"weight,omitempty"`
	// Default value for the field (must be a JSON value). Encapsulate strings with double quotes (e.g. "Foo").
	Default interface{} `json:"default"`
	// If this field is required or not.
	Required bool `json:"required"`
}

func (cf CustomField) String() string {
	return fmt.Sprintf("CustomField{ID: %d, Name: %s}", cf.ID, cf.Name)
}

// CustomField implements IDItem interface.
func (cf *CustomField) GetID() int {
	return cf.ID
}
func (cf *CustomField) GetObjectType() constants.ContentType {
	return constants.ContentTypeExtrasCustomField
}
func (cf *CustomField) GetAPIPath() constants.APIPath {
	return constants.CustomFieldsAPIPath
}
