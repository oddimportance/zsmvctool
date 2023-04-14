package persistence

// Structure to read data from XML form file
// Form elements are unmarshelled into
// Form > Elements
// @ Procedure :
// Step 1: Set Field Name in struct Form
// Step 2: Set Attribute Name in struct Option (eg.: required in Radio and Checkbox)
// Step 3: Set Persistence Param in struct FormElements
// Step 4: Set Class Params in Class FormHandler > parseXmlToPersistence()
type Form struct {
	Elements []struct {
		FieldName                     string             `xml:"fieldName"`
		FieldType                     string             `xml:"fieldType"`
		FieldOptions                  []Option           `xml:"fieldOptions>option"`
		FieldOptionsFromDatabaseQuery OptionFromDatabase `xml:"fieldOptionsFromDatabase"`
		FieldLabel                    string             `xml:"fieldLabel"`
		FieldLabelAdditional          string             `xml:"fieldLabelAdditional"`
		FieldPlaceholder              string             `xml:"fieldPlaceholder"`
		FieldCSS                      string             `xml:"fieldCSS"`
		FieldOnclick                  string             `xml:"fieldOnclick"`
		FieldDecoratorCSS             string             `xml:"fieldDecoratorCSS"`
		FieldID                       string             `xml:"fieldID"`
		FieldIcon                     string             `xml:"fieldIcon"`
		FieldDisabled                 string             `xml:"fieldDisabled"`
		FieldHintText                 string             `xml:"fieldHintText"`
		FieldValidation               ValidationParams   `xml:"validation"`
		FieldForAdminOnly             string             `xml:"fieldForAdminOnly"`
	} `xml:"element"`
}

type Option struct {
	OptionValue    string `xml:"value,attr"`
	OptionSelected string `xml:"selected,attr"`
	OptionLabel    string `xml:",chardata"`
}

type OptionFromDatabase struct {
	SqlQuery             string   `xml:"sqlQuery"`
	FieldForValue        string   `xml:"fieldForValue"`
	OptionSelected       string   `xml:"selectAtValue"`
	FieldForLable        string   `xml:"fieldForLabel"`
	RuntimeHTTPGetParams []string `xml:"runtimeParams>httpGetParam"`
	RuntimeSessionParams []string `xml:"runtimeParams>sessionParam"`
}

type ValidationParams struct {
	ValidationType    string              `xml:"validationType"`
	IsMandatory       string              `xml:"isMandatory"`
	ErrorMessage      string              `xml:"errorMessage"`
	Min               string              `xml:"min"`
	Max               string              `xml:"max"`
	AllowedExtensions []AllowedExtensions `xml:"allowedExtensions>extension"`
}

type AllowedExtensions struct {
	Extension string `xml:",chardata"`
}

// Presistence of Element structure for
// further processing, like HTML decoration
// and delegating data to Template View
type FormElements struct {
	FieldName                     string
	FieldType                     string
	FieldOptions                  []Option
	FieldOptionsFromDatabaseQuery OptionFromDatabase
	FieldLabel                    string
	FieldLabelAdditional          string
	FieldPlaceholder              string
	FieldCSS                      string
	FieldOnclick                  string
	FieldDecoratorCSS             string
	FieldID                       string
	FieldIcon                     string
	FieldDisabled                 string
	FieldHintText                 string
	FieldValidation               ValidationParams
	FieldForAdminOnly             string
}

type JsonForm struct {
	Key          string
	Value        string
	Max          int
	Min          int
	Type         string
	IsMandatory  bool
	ErrorMessage string
}
