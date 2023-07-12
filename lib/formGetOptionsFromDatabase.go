package lib

import (
	"fmt"

	"github.com/oddimportance/zsmvctool/persistence"
)

func (f *FormHandler) initOptionsFromDatabase(element persistence.FormElements) []persistence.Option {

	if element.FieldOptionsFromDatabaseQuery.SqlQuery == "" {
		return nil
	}

	// Initiate an empty DB Instance
	f._db.Connect("", "", f._handleRedirectAndPanic)

	// Set the query string
	f._db.queryString = f.replaceSessionParams(element, f.replaceHTTPGetParams(element, element.FieldOptionsFromDatabaseQuery.SqlQuery))

	// execute, filter and return result
	return f.readValuesFromGivenOptionKey(element, f._db.ExecSelect())

}

func (f *FormHandler) readValuesFromGivenOptionKey(element persistence.FormElements, dbRes map[int]map[string]string) []persistence.Option {
	var optionsToReturn = []persistence.Option{}
	row := map[string]string{}
	optionSelected := ""

	// to respect mysql order by use
	// iteration and not for loop
	for i := 0; i < len(dbRes); i++ {
		row = dbRes[i]
		// log.Println(element.FieldOptionsFromDatabaseQuery.OptionSelected, row[element.FieldOptionsFromDatabaseQuery.FieldForValue])

		// if element.FieldOptionsFromDatabaseQuery.OptionSelected == row[element.FieldOptionsFromDatabaseQuery.FieldForValue] {
		// 	// if InArray(element.FieldOptionsFromDatabaseQuery.OptionSelected, row[element.FieldOptionsFromDatabaseQuery.FieldForValue]) {
		// 	optionSelected = f.checkIfSelectOrChecked(element)
		// }
		optionSelected = f.isCheckboxOrSelectOptionSelected(element, element.FieldOptionsFromDatabaseQuery.OptionSelected, row[element.FieldOptionsFromDatabaseQuery.FieldForValue])

		optionsToReturn = append(optionsToReturn,
			persistence.Option{
				OptionLabel:    row[element.FieldOptionsFromDatabaseQuery.FieldForLable],
				OptionValue:    row[element.FieldOptionsFromDatabaseQuery.FieldForValue],
				OptionSelected: optionSelected,
			})

		optionSelected = ""
	}
	return optionsToReturn
}

func (f *FormHandler) replaceHTTPGetParams(element persistence.FormElements, query string) string {

	queryToReturn := query

	for _, param := range element.FieldOptionsFromDatabaseQuery.RuntimeHTTPGetParams {
		if "" != f.postRequest.FormValue(param) {
			queryToReturn = StringReplace(queryToReturn, fmt.Sprintf("$httpGetParam_%s", param), f.postRequest.FormValue(param), -1)
		}
	}

	return queryToReturn
}

func (f *FormHandler) replaceSessionParams(element persistence.FormElements, query string) string {

	queryToReturn := query

	for _, param := range element.FieldOptionsFromDatabaseQuery.RuntimeSessionParams {
		if "" != f._handleSession.GetSessionCookie(param) {
			queryToReturn = StringReplace(queryToReturn, fmt.Sprintf("$sessionParam_%s", param), f._handleSession.GetSessionCookie(param), -1)
		}
	}

	return queryToReturn
}

func (f *FormHandler) getParamValueFormHttpGetOrMuxVars() {

}
