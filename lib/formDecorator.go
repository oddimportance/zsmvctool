package lib

import (
	"fmt"
	"html/template"
	"strings"

	"zsmvctool-api/persistence"
)

// By default theme path is set to templates/uielements
// If required set custom theme path
func (f *FormHandler) SetThemePath(themePath string) {
	f.themePath = themePath
}

// Initiate HTML parsing. Before the initiation
// form elements must be parsed globally from
// their respected xml files
func (f *FormHandler) InitiateDecoration() map[string]template.HTML {
	return f.parseToHTML()
}

func (f *FormHandler) SetDefaultValueForElement(element string, value string) {
	f.valuesForInidvidualElementByDefault[element] = value
}

func (f *FormHandler) GetDefaultValueForElement(element string) string {
	return f.valuesForInidvidualElementByDefault[element]
}

// Parse HTML according to element type
func (f *FormHandler) parseToHTML() map[string]template.HTML {

	contentToReturn := map[string]template.HTML{}

	for _, element := range f.formElements {

		switch element.FieldType {
		case "input":
			contentToReturn[element.FieldName] = f.handleTypeInput(element)
		case "hidden":
			contentToReturn[element.FieldName] = f.handleTypeInput(element)
		case "textarea":
			contentToReturn[element.FieldName] = f.handleTypeInput(element)
		case "radio":
			contentToReturn[element.FieldName] = f.handleTypeRadio(element)
		case "checkbox":
			contentToReturn[element.FieldName] = f.handleTypeCheckbox(element)
		case "select":
			contentToReturn[element.FieldName] = f.handleTypeSelect(element)
		case "file":
			contentToReturn[element.FieldName] = f.handleTypeFile(element)
		case "password":
			contentToReturn[element.FieldName] = f.handleTypePassword(element)
		default:
			f._handleRedirectAndPanic.TriggerPanic("00010 " + element.FieldType)
		}

	}

	return contentToReturn
}

func (f *FormHandler) pregElementInContent(element persistence.FormElements, htmlContent string) string {

	// Set Element value from Post Request
	elementValue := ""

	// Override element value if request method is POST
	if f.postRequest.Method == "POST" {
		elementValue = f.postRequest.PostFormValue(element.FieldName)
	} else if f.postRequest.Method != "POST" && f.valuesFromDatabase != nil {
		elementValue = f.cleanNullFromDb(f.valuesFromDatabase[fmt.Sprintf("%s%s", f.dbTablePrefix, strings.ToLower(element.FieldName))])
	} else if f.postRequest.Method != "POST" && f.valuesForInidvidualElementByDefault[element.FieldName] != "" {
		elementValue = f.valuesForInidvidualElementByDefault[element.FieldName]
	}

	htmlContent = f.findAndReplace(htmlContent, "NAME", element.FieldName)
	htmlContent = f.findAndReplace(htmlContent, "LABEL", f._translator.Translate(element.FieldLabel))
	htmlContent = f.findAndReplace(htmlContent, "LABELADDITIONAL", f._translator.Translate(element.FieldLabelAdditional))
	htmlContent = f.findAndReplace(htmlContent, "PLACEHOLDER", f._translator.Translate(element.FieldPlaceholder))
	htmlContent = f.findAndReplace(htmlContent, "CSS", element.FieldCSS)
	htmlContent = f.findAndReplace(htmlContent, "DECORATORCSS", element.FieldDecoratorCSS)
	htmlContent = f.findAndReplace(htmlContent, "ID", element.FieldID)
	htmlContent = f.findAndReplace(htmlContent, "ONCLICK", element.FieldOnclick)
	htmlContent = f.findAndReplace(htmlContent, "ICON", element.FieldIcon)
	htmlContent = f.findAndReplace(htmlContent, "DISABLED", element.FieldDisabled)
	htmlContent = f.findAndReplace(htmlContent, "VALUE", UnescapeHtml(elementValue))

	return htmlContent

}

func (f *FormHandler) pregRadioCheckboxOptionsInContent(option persistence.Option, element persistence.FormElements, htmlContent string) string {

	elementSelected := f.handleRadioCheckboxIsChecked(element, option)

	htmlContent = f.findAndReplace(htmlContent, "NAME", element.FieldName)
	htmlContent = f.findAndReplace(htmlContent, "LABEL", f._translator.Translate(option.OptionLabel))
	htmlContent = f.findAndReplace(htmlContent, "VALUE", option.OptionValue)
	htmlContent = f.findAndReplace(htmlContent, "SELECTED", elementSelected)

	return htmlContent

}

// Set the checked attrib according to following
// - if option is selected by default and the user
// has not changed it or if the option was selected
// by user
// @ param element persistence.Formelements
// @ param option persisitence.Option
// @ return string
func (f *FormHandler) handleRadioCheckboxIsChecked(element persistence.FormElements, option persistence.Option) string {

	// To avoid panic error look if postRequest is set
	if f.postRequest.Method == "POST" {
		// Is checked only if the option value of
		// the element is the same as the post value
		// of the element
		if f.postRequest.PostFormValue(element.FieldName) == option.OptionValue {
			return f.checkIfSelectOrChecked(element)
		} else {
			// looks like the user has either unchecked
			// the option or is not selected by default
			return ""
		}
	} else if f.valuesFromDatabase != nil {
		if f.valuesFromDatabase[fmt.Sprintf("%s%s", f.dbTablePrefix, strings.ToLower(element.FieldName))] == option.OptionValue {
			return f.checkIfSelectOrChecked(element)
		}
	}

	// the post method seems to be GET
	// therefor return the xml value,
	// which may either be checked or
	// empty by default
	return option.OptionSelected

}

func (f *FormHandler) checkIfSelectOrChecked(element persistence.FormElements) string {
	if element.FieldType == "select" {
		return "selected"
	}
	return "checked"
}

func (f *FormHandler) initHintAndErrorText(content string, parseType string, textToParse string) string {
	return f.findAndReplace(content, strings.ToUpper(parseType), f.parseHintAndErrorText(parseType, strings.ToUpper(parseType), textToParse))
}

func (f *FormHandler) parseHintAndErrorText(fileName string, parseKey string, parseVal string) string {
	return f.findAndReplace(string(f.readThemeFile(fileName)), parseKey, f._translator.Translate(parseVal))
}

func (f *FormHandler) findAndReplace(content string, fieldName string, fieldVal string) string {
	return strings.Replace(content, fmt.Sprintf("##%s##", strings.ToUpper(fieldName)), fieldVal, -1)
}

func (f *FormHandler) readThemeFile(fileName string) []byte {

	// Prepare the path
	var pathToThemeFile = fmt.Sprintf("%s/%s.html", f.themePath, fileName)

	// get the content from xml file
	var htmlContent = f._readFile.Initiate(pathToThemeFile)

	// look if file did not exist, if yes controller
	// will handle panic error redirection
	if len(htmlContent) == 0 {
		f._handleRedirectAndPanic.TriggerPanic("00009")
	}

	return htmlContent

}

func (f *FormHandler) initReadAndParse(element persistence.FormElements) template.HTML {

	var htmlContent string
	htmlContent = f.pregElementInContent(element, string(f.readThemeFile(element.FieldType)))
	htmlContent = f.initHintAndErrorText(htmlContent, "hintText", element.FieldHintText)
	if f.formErrors != nil {
		htmlContent = f.initHintAndErrorText(htmlContent, "errorText", f.formErrors[element.FieldName])
	}

	// // fmt.Println(htmlContent)
	return template.HTML(htmlContent)

}

func (f *FormHandler) wrapRadioCheckboxElement(element persistence.FormElements, elementOptions string) string {

	htmlContent := string(f.readThemeFile(element.FieldType))
	htmlContent = f.findAndReplace(htmlContent, "ELEMENT", elementOptions)

	htmlContent = f.initHintAndErrorText(htmlContent, "hintText", element.FieldHintText)
	if f.formErrors != nil {
		htmlContent = f.initHintAndErrorText(htmlContent, "errorText", f.formErrors[element.FieldName])
	}

	// fmt.Println(htmlContent)
	return htmlContent

}

func (f *FormHandler) parseRadioCheckboxSelectOptions(elementType string, element persistence.FormElements) string {

	var htmlTemplateContent = string(f.readThemeFile(fmt.Sprintf("%sOption", elementType)))
	htmlContent := ""

	for _, option := range element.FieldOptions {
		htmlContent = fmt.Sprintf("%s\n%s", htmlContent, f.pregRadioCheckboxOptionsInContent(option, element, htmlTemplateContent))
	}

	// Look if there were data from DB
	var optionsFromDb = f.initOptionsFromDatabase(element)

	// looks like there are options from db to consider
	if nil != optionsFromDb && len(optionsFromDb) != 0 {

		for _, optionFromDb := range optionsFromDb {
			htmlContent = fmt.Sprintf("%s\n%s", htmlContent, f.pregRadioCheckboxOptionsInContent(optionFromDb, element, htmlTemplateContent))
		}

	}

	htmlContent = f.initHintAndErrorText(htmlContent, "hintText", element.FieldHintText)
	if f.formErrors != nil {
		htmlContent = f.initHintAndErrorText(htmlContent, "errorText", f.formErrors[element.FieldName])
	}
	// fmt.Println(htmlContent)
	return htmlContent

}

func (f *FormHandler) handleTypeInput(element persistence.FormElements) template.HTML {
	return f.initReadAndParse(element)
}

func (f *FormHandler) handleTypeRadio(element persistence.FormElements) template.HTML {
	return template.HTML(f.pregElementInContent(element, f.wrapRadioCheckboxElement(element, f.parseRadioCheckboxSelectOptions("radio", element))))
}

func (f *FormHandler) handleTypeCheckbox(element persistence.FormElements) template.HTML {
	return template.HTML(f.pregElementInContent(element, f.wrapRadioCheckboxElement(element, f.parseRadioCheckboxSelectOptions("checkbox", element))))
}

func (f *FormHandler) handleTypeSelect(element persistence.FormElements) template.HTML {
	return template.HTML(f.pregElementInContent(element, f.wrapRadioCheckboxElement(element, f.parseRadioCheckboxSelectOptions("select", element))))
}

func (f *FormHandler) handleTypeFile(element persistence.FormElements) template.HTML {
	return f.initReadAndParse(element)
}

func (f *FormHandler) handleTypePassword(element persistence.FormElements) template.HTML {
	return f.initReadAndParse(element)
}

// The null value from DB is %!s(<nil>)
// retur empty string if the case
func (f *FormHandler) cleanNullFromDb(valueToClean string) string {
	// %!s(<nil>) == NULL
	if "%!s(<nil>)" == valueToClean {
		return ""
	}
	return valueToClean
}
