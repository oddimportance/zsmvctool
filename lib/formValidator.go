package lib

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/oddimportance/zsmvctool/persistence"

	"github.com/asaskevich/govalidator"
)

// Get the size of uploaded file for max size limitation
type Sizer interface {
	Size() int64
}

// Form validation, handles all types of HTML fields all
// the way from Input to File
// @ param _getVars map[string]string HTTP Get vars
// @ param r *http.Request is required to read Post Values
// @ param w http.ResponseWriter is required to write data to browser
// @ param return void
func (f *FormHandler) InitValidation() {

	f.formIsValid = true

	for _, element := range f.formElements {

		if InArray(f.elementsToSkipValidation, element.FieldName) {
			continue
		}

		switch element.FieldType {
		case "input":
			f.handleValidationInput(element)
		case "password":
			f.handleValidationInput(element)
		case "file":
			f.handleValidationFile(element)
		case "select":
			f.handleValidationSelect(element)
		case "checkbox":
			f.handleValidationCheckbox(element)
		case "radio":
			f.handleValidationRadio(element)
		default:
			// fmt.Println(r.FormValue(element.FieldName))
		}
	}

	// fmt.Println(f.formErrors)

}

// Handle File Upload
// Checks if is mandatory, checks file size, check allowed file types
// @ param element persistence.FormElements
// @ param return void
func (f *FormHandler) handleValidationFile(element persistence.FormElements) {

	// parse file
	// @file : multiplepart.File
	// @handle : *multiplepart.FileHeader
	file, _, err := f.postRequest.FormFile(element.FieldName)

	// convert string values to respected formats
	isMandatory, _ := strconv.ParseBool(element.FieldValidation.IsMandatory)

	// Validate upload only if error is nil, which is,
	// if a file is uploaded
	if err == nil {

		defer file.Close()

		// convert xml string to float
		maxFileSize, _ := strconv.ParseFloat(element.FieldValidation.Max, 10)

		sizer, ok := file.(Sizer)
		fileSize := float64(0)

		if ok {
			// read file size from header
			fileSize = float64(sizer.Size())
		}

		// check file size
		f.checkFileSize(element.FieldName, fileSize, maxFileSize)
		// check file extension
		// f.checkFileAllowedExtensions(element.FieldName, handle.Header.Get("Content-Type"), element.FieldValidation.AllowedExtensions)
		f.checkFileAllowedExtensions_v1(element.FieldName, file, element.FieldValidation.AllowedExtensions)

	} else {

		// Looks like no file was uploaded. Hence, check if is Mandatory
		f.checkIfIsMandatory(element.FieldName, isMandatory, 0, element.FieldValidation.ErrorMessage)

	}

}

// Handle File Upload
// Checks if is mandatory, check min length, checks max,
// checks type like AlphabetOnly, AlphaNumeric, EMail etc.
// @ param element persistence.FormElements
// @ param return void
func (f *FormHandler) handleValidationInput(element persistence.FormElements) {

	var postValue = f.postRequest.FormValue(element.FieldName)

	// convert string values to respected formats
	var isMandatory, min, max = f.parseXMLStringToRespectedTypes(element)

	// Check isMandatory
	f.checkIfIsMandatory(element.FieldName, isMandatory, len(postValue), element.FieldValidation.ErrorMessage)

	// Check if satisfies min
	f.checkMin(element.FieldName, isMandatory, len(postValue), min)

	// Check if satisfies max
	f.checkMax(element.FieldName, len(postValue), max)

	// Check validation type, like int or email etc.
	f.checkValidationType(element.FieldName, element.FieldValidation.ValidationType, postValue)

}

// Handle Select
// Checks if is mandatory, check min length, checks max,
// checks type like AlphabetOnly, AlphaNumeric, EMail etc.
// Checks if the post values are other than the options
// in XML Configurator (extra security)
// @ param element persistence.FormElements
// @ param return void
func (f *FormHandler) handleValidationSelect(element persistence.FormElements) {

	var postValue = f.postRequest.Form[element.FieldName]

	var isMandatory, min, max = f.parseXMLStringToRespectedTypes(element)

	// fix for a strange bug
	postLength := len(postValue)
	if len(postValue) != 0 && postValue[0] == "" {
		postLength = 0
	}

	// Check isMandatory
	f.checkIfIsMandatory(element.FieldName, isMandatory, postLength, element.FieldValidation.ErrorMessage)

	// Check if satisfies min
	f.checkMin(element.FieldName, isMandatory, postLength, min)

	// Check if satisfies max
	f.checkMax(element.FieldName, postLength, max)

	// Walk through options for validation
	// Read every single option and validate against type
	f.walkThroughOptionsForValidation(element.FieldName, postValue, element.FieldValidation.ValidationType)

	// check if the value is out of range
	// Walk through all options to verify if the post value
	// is manupilated
	//f.checkForMultiplePostValueInRange(element.FieldName, postValue, element.FieldOptions)
	f.checkForMultiplePostValueInRange(element, postValue)

}

// Handle Checkboxes
// Checks if is mandatory, check min length, checks max,
// checks type like AlphabetOnly, AlphaNumeric, EMail etc.
// Checks if the post values are other than the options
// in XML Configurator (extra security)
// @ param element persistence.FormElements
// @ param return void
func (f *FormHandler) handleValidationCheckbox(element persistence.FormElements) {

	var postValue = f.postRequest.Form[element.FieldName]

	var isMandatory, min, max = f.parseXMLStringToRespectedTypes(element)

	// Check isMandatory
	f.checkIfIsMandatory(element.FieldName, isMandatory, len(postValue), element.FieldValidation.ErrorMessage)

	// Check if satisfies min
	f.checkMin(element.FieldName, isMandatory, len(postValue), min)

	// Check if satisfies max
	f.checkMax(element.FieldName, len(postValue), max)

	// Walk through options for validation
	// Read every single option and validate against type
	f.walkThroughOptionsForValidation(element.FieldName, postValue, element.FieldValidation.ValidationType)

	// check if the value is out of range
	// Walk through all options to verify if the post value
	// is manupilated
	//f.checkForMultiplePostValueInRange(element.FieldName, postValue, element.FieldOptions)
	f.checkForMultiplePostValueInRange(element, postValue)

}

// Handle Radio
// Checks if is mandatory, check min length, checks max,
// checks type like AlphabetOnly, AlphaNumeric, EMail etc.
// Checks if the post values are other than the options
// in XML Configurator (extra security)
// @ param element persistence.FormElements
// @ param return void
func (f *FormHandler) handleValidationRadio(element persistence.FormElements) {

	var postValue = f.postRequest.Form[element.FieldName]

	var isMandatory, min, max = f.parseXMLStringToRespectedTypes(element)

	// Check isMandatory
	f.checkIfIsMandatory(element.FieldName, isMandatory, len(postValue), element.FieldValidation.ErrorMessage)

	// Check if satisfies min
	f.checkMin(element.FieldName, isMandatory, len(postValue), min)

	// Check if satisfies max
	f.checkMax(element.FieldName, len(postValue), max)

	// Walk through options for validation
	// Read every single option and validate against type
	f.walkThroughOptionsForValidation(element.FieldName, postValue, element.FieldValidation.ValidationType)

	// check if the value is out of range
	// Walk through all options to verify if the post value
	// is manupilated
	//f.checkForMultiplePostValueInRange(element.FieldName, postValue, element.FieldOptions)
	f.checkForMultiplePostValueInRange(element, postValue)

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
func (f *FormHandler) checkMin(elementName string, isMandatory bool, valueLength int, min int) {

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
func (f *FormHandler) checkMax(elementName string, valueLength int, max int) {

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
func (f *FormHandler) checkValidationType(elementName string, validationType string, postValue string) {

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
func (f *FormHandler) checkForMultiplePostValueInRange(element persistence.FormElements, postValues []string) {
	// @ param elementName string required to set error message
	var elementName = element.FieldName

	// @ param valueToBeInRange []persistence.Option options unmarshelled from XML
	// concat options from xml and from database
	var valuesToBeInRange []persistence.Option = f.appendOptionsFromDatabase(element.FieldOptions, f.initOptionsFromDatabase(element))

	for _, postValue := range postValues {

		f.checkForPostValueInRange(elementName, postValue, valuesToBeInRange)

	}

}

func (f *FormHandler) appendOptionsFromDatabase(fieldOption, optionsFromDatabase []persistence.Option) []persistence.Option {
	valuesToReturn := fieldOption
	for _, optionFromDb := range optionsFromDatabase {
		valuesToReturn = append(valuesToReturn, optionFromDb)
	}

	return valuesToReturn
}

// Walk through all XML Options one at a time and verify
// their value against the given value in post. If invalid
// populates formError with appropirate message
// @ param elementName string required to set error message
// @ param postValues string values from http.request
// @ param valueToBeInRange []persistence.Option options unmarshelled from XML
// @ param return void
func (f *FormHandler) checkForPostValueInRange(elementName string, postValue string, valueToBeInRange []persistence.Option) {
	// fmt.Println("teste ", postValue)
	for _, val := range valueToBeInRange {
		if val.OptionValue == postValue {
			return
		}
	}
	f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_wrongPostValue", elementName))

}

// Checks if the uploaded file is of valid type
// If invalid marks the form as invalid and sets
// an appropriate formError message
// @ param elementName string required to set error message
// @ param fileExtension string typically mime-content-type
// @ param allowedExtension []persistence.AllowedExtensions
// @ param return void
func (f *FormHandler) checkFileAllowedExtensions(elementName string, fileExtension string, allowedExtensions []persistence.AllowedExtensions) {

	// Walk through array of allowed extension
	// if valid extension was found, break the
	// loop. If not register an error
	for _, allowedExtension := range allowedExtensions {
		if allowedExtension.Extension == fileExtension {
			return
		}
	}

	f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_invalidFileExtension", elementName))

}

// use detect content type method to validte file extension
// it is safer than the header Content-Type method as
// header type can be manupilated
func (f *FormHandler) checkFileAllowedExtensions_v1(elementName string, file multipart.File, allowedExtensions []persistence.AllowedExtensions) {
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_invalidFileExtension", elementName))
		return
	}

	var fileExtension = http.DetectContentType(buff)

	// Walk through array of allowed extension
	// if valid extension was found, break the
	// loop. If not register an error
	for _, allowedExtension := range allowedExtensions {
		if allowedExtension.Extension == fileExtension {
			return
		}
	}

	f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment("form_error_invalidFileExtension", elementName))

}

// Checks if file size does not exceed the max upload size
// convert Bytes to MB by dividing filesize into 1000000
// and then check if filesize exceeds max limit
// If invalid marks the form as invalid and sets
// an appropriate formError message
// @ param elementName string required to set error message
// @ param fileSize float64 size of uploaded file in bytes
// @ param maxFileSize float64 size in MB
// @ param return void
func (f *FormHandler) checkFileSize(elementName string, fileSize float64, maxFileSize float64) {

	// there might have been an error while
	// calculating the file size. Hence, file
	// size might be 0. If so, maŕk form as
	// error
	if fileSize == 0 {
		f.MarkFormInvalid(elementName, "form_error_fileSizeNotReadable")
		return
	}

	// Validate file size
	// convert KB to MB by dividing filesize
	// into 1000000 and then check if filesize
	// exceeds max limit
	var fileSizeInMb float64 = (fileSize / 1000000)

	// If maxfilesize was empty it returns 0.00 therefore ignore it
	// or look if the maxfilesize exceeds
	if maxFileSize != 0 && fileSizeInMb > maxFileSize {
		f.MarkFormInvalid(elementName, fmt.Sprintf(f._translator.TranslateSafeWithAttachment("form_error_exceedsMaxFileSize", elementName), fileSizeInMb, maxFileSize))
	}
}

// Checks if the field is a mandatory field
// If invalid marks the form as invalid and sets
// an appropriate formError message
// @ param elementName string required to set error message
// @ param isMandatory bool
// @ param elementErrorMessage string required to set error message
// @ param valueLength int it may be the length of input val or count checkbox options
// @ param return void
func (f *FormHandler) checkIfIsMandatory(elementName string, isMandatory bool, valueLength int, elementErrorMessage string) {
	if isMandatory && valueLength == 0 {
		errorMessage := "form_error_isMandatory"
		if elementErrorMessage != "" {
			errorMessage = elementErrorMessage
		}
		f.MarkFormInvalid(elementName, f._translator.TranslateSafeWithAttachment(errorMessage, elementName))
	}
}

// @ param elementName string required to set error message
// @ param postValues string values from http.request
// @ param validationType string potentially a string from XML
// @ param return void
func (f *FormHandler) walkThroughOptionsForValidation(elementName string, postValue []string, validationType string) {
	for _, val := range postValue {
		f.checkValidationType(elementName, validationType, val)
	}
}

// Mark the form as invalid and populate formError with
// appropriate message
// @ param elementName string needed to set error message
// @ param errorMsg string the error message
// @ param return void
func (f *FormHandler) MarkFormInvalid(elementName string, errorMsg string) {
	f.formIsValid = false
	f.formErrors[elementName] = errorMsg
}

// Use String Conversion to convert string values
// from XML to their reespected types
// @ param element persistence.FormElements
// @ param return bool, int, int
func (f *FormHandler) parseXMLStringToRespectedTypes(element persistence.FormElements) (bool, int, int) {

	// convert string values to respected formats

	// bool
	isMandatory, _ := strconv.ParseBool(element.FieldValidation.IsMandatory)
	// int
	min, _ := strconv.Atoi(element.FieldValidation.Min)
	// int
	max, _ := strconv.Atoi(element.FieldValidation.Max)

	// fmt.Printf("Element : %s, %v, %d, %d\n", element.FieldName, isMandatory, min, max)

	return isMandatory, min, max

}

// Check the validity of form
func (f *FormHandler) IsFormValid() bool {
	return f.formIsValid
}

// Set form element to skip validation
func (f *FormHandler) SetFormElementsToSkipValidition(formElements []string) {
	f.elementsToSkipValidation = formElements
}

func (f *FormHandler) ValidateMobileNumberGermanFormat(mobileNumber string) (string, string, bool) {

	if len(mobileNumber) < 10 {
		return "", "", false
	}
	// Germnan carrier provider codes
	var germanCarrierCodes = []string{"0151", "0152", "0157", "0159", "0160", "0162", "0163", "0170", "0171", "0172", "0173", "0174", "0175", "0176", "0177", "0178", "0179"}

	var germanCountryCode = "0049"
	var carrierPrefix string
	var first4Digits = StringSubString(mobileNumber, 0, 4)
	var first3Digits = StringSubString(first4Digits, 0, 3)
	// log.Println(first3Digits, first4Digits)

	if InArray(germanCarrierCodes, fmt.Sprintf("0%s", first3Digits)) { //chekc if the number is 176 format
		return germanCountryCode, mobileNumber, true
	} else if InArray(germanCarrierCodes, first4Digits) { //chekc if the number is 0176 format
		return germanCountryCode, StringSubString(mobileNumber, 1, len(mobileNumber)), true
	} else if first3Digits == "+49" { // check if the number has +49
		carrierPrefix = StringSubString(mobileNumber, 4, 6)
		if InArray(germanCarrierCodes, fmt.Sprintf("01%s", carrierPrefix)) {
			return germanCountryCode, StringSubString(mobileNumber, 3, len(mobileNumber)), true
		}
	} else if first4Digits == germanCountryCode { // check if the number has 0049
		carrierPrefix = StringSubString(mobileNumber, 4, 7)
		if InArray(germanCarrierCodes, fmt.Sprintf("0%s", carrierPrefix)) {
			return germanCountryCode, StringSubString(mobileNumber, 4, len(mobileNumber)), true
		}
	}

	return "", "", false

}
