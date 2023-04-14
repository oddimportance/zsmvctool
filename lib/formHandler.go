package lib

import (
	"fmt"
	//"html/template"
	"encoding/xml"
	"net/http"
	"strings"

	"zsmvctool-api/persistence"
)

type FormHandler struct {
	formName                            string // Name of the form
	_handleSession                      *HandleSession
	_handleFilePath                     *HandleFilePath                     // File path to handle readfile
	_handleRedirectAndPanic             *HandleRedirectAndPanic             // panic error
	_translator                         *Translator_i18n                    // set the translator
	formElements                        map[string]persistence.FormElements // Form elements
	themePath                           string                              // setable theme path
	_readFile                           ReadFile                            // read file
	formErrors                          map[string]string                   // Post vars postVars map[string]string Form Errors
	postRequest                         *http.Request                       // http request
	muxGetVars                          map[string]string                   // Mux Get Vars
	formIsValid                         bool                                // form validation
	valuesForInidvidualElementByDefault map[string]string                   // Values setted for inidividual element
	valuesFromDatabase                  map[string]string                   // Values to edit from DB
	elementsToSkipValidation            []string                            // Store names of elements, which need no validation
	_db                                 DbAdapter                           // required to seek options from database
	// Db table prefix to substr from
	// field ex.: Value for First_name from
	// DB will be zsusrrsnt_ + element tolower
	dbTablePrefix string
}

func (f *FormHandler) InitForm(
	_getVars map[string]string,
	r *http.Request,
	w http.ResponseWriter,
	formName string,
	_handleSession *HandleSession,
	handleFilePath *HandleFilePath,
	_handleRedirectAndPanic *HandleRedirectAndPanic,
	translator *Translator_i18n) map[string]persistence.FormElements {

	// Set request param globally
	f.setHttpRequest(r)

	f.formName = formName
	f._handleSession = _handleSession
	f._handleRedirectAndPanic = _handleRedirectAndPanic
	f._handleFilePath = handleFilePath
	f._translator = translator
	f.formErrors = map[string]string{}
	f.valuesForInidvidualElementByDefault = map[string]string{}

	f.formElements = f.parseXmlToPersistence(f.readDataFromXml())

	// Set Default theme path
	f.setDefaultThemePath()

	// Set global DB Adapter
	f._db = DbAdapter{}

	return f.formElements

}

func (f *FormHandler) readDataFromXml() []byte {

	// Prepare the path
	var pathToFormFile = fmt.Sprintf("%s/%s.xml", f._handleFilePath.GetFormPath(), f.formName)

	// Initate the class gloablly for futher prossessing
	f._readFile = ReadFile{}

	// get the content from xml file
	var xmlContent = f._readFile.Initiate(pathToFormFile)

	// look if file did not exist, if yes controller
	// will handle panic error redirection
	if len(xmlContent) == 0 {
		f._handleRedirectAndPanic.TriggerPanic("00009")
	}

	return xmlContent

}

func (f *FormHandler) parseXmlToPersistence(xmlRaw []byte) map[string]persistence.FormElements {

	// initiate unmarshal
	var form persistence.Form

	// unmarshal
	xml.Unmarshal([]byte(xmlRaw), &form)

	// persistence var
	formPersistence := map[string]persistence.FormElements{}

	// walk through, rearrange elements accordingly
	for _, val := range form.Elements {
		formPersistence[val.FieldName] = persistence.FormElements{
			FieldName:                     val.FieldName,
			FieldCSS:                      val.FieldCSS,
			FieldDecoratorCSS:             val.FieldDecoratorCSS,
			FieldID:                       val.FieldID,
			FieldLabel:                    val.FieldLabel,
			FieldLabelAdditional:          val.FieldLabelAdditional,
			FieldIcon:                     val.FieldIcon,
			FieldOnclick:                  val.FieldOnclick,
			FieldOptions:                  val.FieldOptions,
			FieldOptionsFromDatabaseQuery: val.FieldOptionsFromDatabaseQuery,
			FieldPlaceholder:              val.FieldPlaceholder,
			FieldType:                     val.FieldType,
			FieldDisabled:                 val.FieldDisabled,
			FieldHintText:                 val.FieldHintText,
			FieldValidation:               val.FieldValidation,
			FieldForAdminOnly:             val.FieldForAdminOnly,
		}
	}

	return formPersistence

}

func (f *FormHandler) GetFormElements() map[string]persistence.FormElements {
	return f.formElements
}

func (f *FormHandler) setDefaultThemePath() {
	f.themePath = fmt.Sprintf("%s/%s", f._handleFilePath.GetTemplatePath(), "uielements")
}

func (f *FormHandler) GetFormErrors() map[string]string {
	return f.formErrors
}

func (f *FormHandler) SetDbTablePrefixForSubStr(dbTablePrefix string) {
	f.dbTablePrefix = dbTablePrefix
}

func (f *FormHandler) SetValuesForFormElementsFormDatabase(dbValues map[string]string) {
	f.valuesFromDatabase = f.unescapeValuesToHtmlSafe(dbValues)
}

func (f *FormHandler) unescapeValuesToHtmlSafe(dbValues map[string]string) map[string]string {
	var valuesToReturn map[string]string = map[string]string{}
	for key, val := range dbValues {
		if IsNullOrEmpty(val) {
			valuesToReturn[key] = ""
		} else {
			valuesToReturn[key] = EscapeHtml(val)
		}
	}
	return valuesToReturn
}

// Globalise http request to enable reading post values globally
// @ param r http.Request
// @ param return void
func (f *FormHandler) setHttpRequest(r *http.Request) {
	f.postRequest = r
}

// Globalise Mux Vars to enable reading post values globally
// @ param muxGetVars map[string]string
// @ param return void
func (f *FormHandler) setMuxGetVars(muxGetVars map[string]string) {
	f.muxGetVars = muxGetVars
}

func (f *FormHandler) MakeDBFieldsFromFormElements(dbPrefix string, fieldsToSkip []string) ([]string, []string) {

	var fieldsToReturn = []string{}
	var valuesToReturn = []string{}

	for _, formElement := range f.formElements {
		if fieldsToSkip == nil || !InArray(fieldsToSkip, formElement.FieldName) {
			fieldsToReturn = append(fieldsToReturn, fmt.Sprintf("%s%s", dbPrefix, StringToLower(formElement.FieldName)))
			valuesToReturn = append(valuesToReturn, strings.TrimSpace(f.postRequest.FormValue(formElement.FieldName)))
		}
	}
	return fieldsToReturn, valuesToReturn
}
