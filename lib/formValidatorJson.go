package lib

import (
	"zsmvctool/persistence"

	"github.com/asaskevich/govalidator"
)

type FormHandlerJson struct {
	schema             map[int]persistence.JsonForm
	data               map[string]string
	formIsValid        bool
	formErrors         map[string]string
	formName           string            // Name of the form
	_translator        *Translator_i18n  // set the translator
	valuesFromDatabase map[string]string // Values to edit from DB
	_handleSession     *HandleSession
}

func (f *FormHandlerJson) InitFormJson(
	_schema map[int]persistence.JsonForm,
	_data map[string]string,
	formName string,
	_handleSession *HandleSession,
	translator *Translator_i18n) {

	// Set request param globally
	f.schema = _schema
	f.data = _data
	f._handleSession = _handleSession
	f.formName = formName
	f._translator = translator
	f.formErrors = map[string]string{}

}
func (f *FormHandlerJson) InitValidation() {

	f.formIsValid = true

	for _, element := range f.schema {
		element.Value = f.data[element.Key]
		f.handleValidationValue(element)

	}

}

func (f *FormHandlerJson) handleValidationValue(element persistence.JsonForm) {

	var postValue = element.Value

	// convert string values to respected formats
	var isMandatory = element.IsMandatory
	var min = element.Min
	var max = element.Max

	// Check isMandatory
	f.checkIfIsMandatory(element.Key, isMandatory, len(postValue), element.ErrorMessage)

	// Check if satisfies min
	f.checkMin(element.Key, isMandatory, len(postValue), min)

	// Check if satisfies max
	f.checkMax(element.Key, len(postValue), max)

	// Check validation type, like int or email etc.
	f.checkValidationType(element.Key, element.Type, postValue)

}

// Check Min
// Checks if the value satisfies minimum length
// For type Input pass string length, and for
// Radio/Checkboxes pass count options. If not
// mandatory and the value is empty then it is
// a valid case. If invalid populates formError
// with appropriate message and marks the form as invalid
// @ param elementName string required to set error message
// @ param validationType string potentially a string from XML
// @ param postValue string post value from http.request
// @ param return void
func (f *FormHandlerJson) checkMin(elementName string, isMandatory bool, valueLength int, min int) {

	// the field may be empty if the field is not mandatory
	if !isMandatory && valueLength == 0 {
		return
	}

	if min != 0 && valueLength < min {
		f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_underMin", elementName))
	}

}

// Check Max
// Checks if the value does not exceed max length
// For type Input pass string length, and for
// Radio/Checkboxes pass count options.
// If invalid populates formError with appropriate message
// and marks the form as invalid
// @ param elementName string required to set error message
// @ param valueLength int
// @ param max int
// @ param return void
func (f *FormHandlerJson) checkMax(elementName string, valueLength int, max int) {

	if max != 0 && valueLength > max {
		f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_exceedsMax", elementName))
	}

}

// Check Validation Type
// Uses external Lib github.com/asaskevich/govalidator
// More info : http://www.github.com/asaskevich/govalidator
// If invalid populates formError with appropriate message
// and marks the form as invalid
// @ param elementName string required to set error message
// @ param validationType string potentially a string from XML
// @ param postValue string post value from http.request
// @ param return void
func (f *FormHandlerJson) checkValidationType(elementName string, validationType string, postValue string) {

	govalidator.SetFieldsRequiredByDefault(true)

	// strip white spaces
	postValue = TrimWhiteSpaces(postValue)

	switch validationType {

	case "int":
		// fmt.Println("type int")
		if !govalidator.IsInt(postValue) && postValue != "" {
			f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_isNotInt", elementName))
		}

	case "float":
		// fmt.Println("type float")
		if !govalidator.IsFloat(postValue) && postValue != "" {
			f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_isNotFloat", elementName))
		}

	case "email":
		// fmt.Println("type email")
		if postValue != "" && (!govalidator.IsEmail(postValue) || !ValidateEmailHost(postValue)) {
			f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_isNotEmail", elementName))
		}
	case "alphabetOnly":
		if !govalidator.IsAlpha(postValue) && postValue != "" {
			f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_isNotAlphabetOnly", elementName))
		}
	case "alphabetNumericOnly":
		if !govalidator.IsAlphanumeric(postValue) && postValue != "" {
			f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_isNotAlphabetNumericOnly", elementName))
		}
		// note numeric is different from int
		// int does not allow anything beginning
		// with a zero. A zip code in germany for
		// instance can typically begin with a zero
	case "numericOnly":
		if !govalidator.IsNumeric(postValue) && postValue != "" {
			f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_isNotNumericOnly", elementName))
		}
	case "date":
		if _, err := StandardizeDateToMysqlFormat(f._handleSession.GetLanguageFromCookie(), postValue); err != nil && postValue != "" {
			f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_invalid_date", elementName))
		}
	case "commaSeperatedFloat":
		// if (!govalidator.IsInt(postValue) && !govalidator.IsFloat(FindReplace(postValue, ",", "."))) && postValue != "" {
		// 	f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_commaSeperatedFloat", elementName))
		// }
		if !IsCommaSeperatedFloat(postValue) && postValue != "" {
			f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_commaSeperatedFloat", elementName))
		}

	}

}

// Walk through all options one at a time and verify
// their value against values in range given from XML
// The verification itself takes place in f.checkForPostValueInRange()
// @ param element persistence.Formelements
// @ param postValues []string values from http.request
// @ param return void

// Checks if the field is a mandatory field
// If invalid marks the form as invalid and sets
// an appropriate formError message
// @ param elementName string required to set error message
// @ param isMandatory bool
// @ param elementErrorMessage string required to set error message
// @ param valueLength int it may be the length of input val or count checkbox options
// @ param return void
func (f *FormHandlerJson) checkIfIsMandatory(elementName string, isMandatory bool, valueLength int, elementErrorMessage string) {
	if isMandatory && valueLength == 0 {
		errorMessage := "form_error_isMandatory"
		if elementErrorMessage != "" {
			errorMessage = elementErrorMessage
		}
		f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment(errorMessage, elementName))
	}
}

// Mark the form as invalid and populate formError with
// appropriate message
// @ param elementName string needed to set error message
// @ param errorMsg string the error message
// @ param return void
func (f *FormHandlerJson) MarkFormInvalid(elementName string, errorMsg string) {
	f.formIsValid = false
	f.formErrors[elementName] = errorMsg
}

// Check the validity of form
func (f *FormHandlerJson) IsFormValid() bool {
	return f.formIsValid
}

func (f *FormHandlerJson) GetFormErrors() map[string]string {
	return f.formErrors
}

func (f *FormHandlerJson) GetData() map[string]string {
	return f.data
}
