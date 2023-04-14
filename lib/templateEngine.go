package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/http"

	"github.com/oddimportance/zsmvctool/persistence"
)

type TemplateEngine struct {
	_getVars map[string]string

	_baseLayout string

	_handleFilePath *HandleFilePath

	ViewData map[string]interface{}

	localCSSHolder         []string
	externalCSSHolder      []string
	localJSHolderHeader    []string
	localJSHolderFooter    []string
	externalJSHolderHeader []string
	externalJSHolderFooter []string

	_parseTemplateFilesForImports ParseTemplateFilesForImports

	_w http.ResponseWriter

	_tmpl template.Template

	_translator *Translator_i18n
}

/**
 *
 *
 */
func (t *TemplateEngine) SetBaseLayout(baseLayout string) {
	t._baseLayout = baseLayout
}

func (t *TemplateEngine) CreateTemplate(_getVars map[string]string, handleFilePath *HandleFilePath, w http.ResponseWriter, translator *Translator_i18n) {

	t._handleFilePath = handleFilePath

	t.SetBaseLayout("default_base_layout")

	t.initiateViewDataInterface()

	t._w = w

	t._translator = translator

}

func (t *TemplateEngine) initiateViewDataInterface() {
	t.ViewData = map[string]interface{}{}
}

func (t *TemplateEngine) SetLocalCSSFile(cssFile string) {
	t.localCSSHolder = append(t.localCSSHolder, cssFile)
}

func (t *TemplateEngine) SetExternalCSSFile(cssUrl string) {
	t.externalCSSHolder = append(t.externalCSSHolder, cssUrl)
}

func (t *TemplateEngine) SetLocalJSToHeader(jsFile string) {
	t.localJSHolderHeader = append(t.localJSHolderHeader, jsFile)
}

func (t *TemplateEngine) SetExternalJSToHeader(jsFile string) {
	t.externalJSHolderHeader = append(t.externalJSHolderHeader, jsFile)
}

func (t *TemplateEngine) SetLocalJSToFooter(jsFile string) {
	t.localJSHolderFooter = append(t.localJSHolderFooter, jsFile)
}

func (t *TemplateEngine) SetExternalJSToFooter(jsFile string) {
	t.externalJSHolderFooter = append(t.externalJSHolderFooter, jsFile)
}

func (t *TemplateEngine) setLocalCSSFileNamesToTemplate() {
	t.ViewData["LocalCSSFileNamesToLoad"] = t.appendFileNamesToIncludes(t.localCSSHolder, t._parseTemplateFilesForImports.GetLocalCSSFiles())
}

func (t *TemplateEngine) setExternalCSSFileNamesToTemplate() {
	t.ViewData["ExternalCSSFileNamesToLoad"] = t.appendFileNamesToIncludes(t.externalCSSHolder, t._parseTemplateFilesForImports.GetExternalCSSFiles())
}

func (t *TemplateEngine) setLocalHeaderJSFileNamesToTemplate() {
	t.ViewData["LocalHeaderJSFileNamesToLoad"] = t.appendFileNamesToIncludes(t.localJSHolderHeader, t._parseTemplateFilesForImports.GetLocalJSFilesForHeader())
}

func (t *TemplateEngine) setExternalHeaderJSFileNamesToTemplate() {
	t.ViewData["ExternalHeaderJSFileNamesToLoad"] = t.appendFileNamesToIncludes(t.externalJSHolderHeader, t._parseTemplateFilesForImports.GetExternalJSFilesForHeader())
}

func (t *TemplateEngine) setLocalFooterJSFileNamesToTemplate() {
	t.ViewData["LocalFooterJSFileNamesToLoad"] = t.appendFileNamesToIncludes(t.localJSHolderFooter, t._parseTemplateFilesForImports.GetLocalJSFilesForFooter())
}

func (t *TemplateEngine) setExternalFooterJSFileNamesToTemplate() {
	t.ViewData["ExternalFooterJSFileNamesToLoad"] = t.appendFileNamesToIncludes(t.externalJSHolderFooter, t._parseTemplateFilesForImports.GetExternalJSFilesForFooter())
}

func (t *TemplateEngine) appendFileNamesToIncludes(fileNamesArray, filesToAppend []string) []string {
	for _, fileName := range filesToAppend {
		fileNamesArray = append(fileNamesArray, fileName)
	}
	return fileNamesArray
}

func (t *TemplateEngine) ExecuteParseToString(controller, action string) string {
	return t._execute(controller, action, true)
}

func (t *TemplateEngine) Execute(controller, action string) {
	_ = t._execute(controller, action, false)
}

func (t *TemplateEngine) _execute(controller, action string, parseTemplateToString bool) string {

	t._parseTemplateFilesForImports = ParseTemplateFilesForImports{}

	var arrayOfFilesToParse = t._parseTemplateFilesForImports.InitParsing(t._handleFilePath, t._baseLayout, controller, action)

	t._tmpl = *template.New("Main").Funcs(template.FuncMap{
		"TemplateFunctionTranslate":               t.TemplateFunctionTranslate,
		"TemplateFunctionTranslateFormat":         t.TemplateFunctionTranslateFormat,
		"TemplateFunctionMakeSalutation":          t.TemplateFunctionMakeSalutation,
		"TemplateFunctionLocalizeDate":            t.TemplateFunctionLocalizeDate,
		"TemplateFunctionLocalizeDateTime":        t.TemplateFunctionLocalizeDateTime,
		"TemplateFunctionIsNullOrEmpty":           t.TemplateFunctionIsNullOrEmpty,
		"TemplateFunctionEncodeToURL":             t.TemplateFunctionEncodeToURL,
		"TemplateFunctionGetAssociatedValue":      t.TemplateFunctionGetAssociatedValue,
		"TemplateFunctionSubstr":                  t.TemplateFunctionSubstr,
		"TemplateFunctionSubStringWhitespaceSafe": t.TemplateFunctionSubStringWhitespaceSafe,
		"TemplateFunctionStringToHtml":            t.TemplateFunctionStringToHtml,
		"TemplateFunctionMakeArray":               t.TemplateFunctionMakeArray,
		"TemplateFunctionAppendToArray":           t.TemplateFunctionAppendToArray,
		"TemplateFunctionInArray":                 t.TemplateFunctionInArray,
		"TemplateFunctionIterate":                 t.TemplateFunctionIterate,
		"TemplateFunctionFindReplace":             t.TemplateFunctionFindReplace,
		"TemplateFunctionContainsString":          t.TemplateFunctionContainsString,
		"TemplateFunctionEncryptIdForUrl":         t.TemplateFunctionEncryptIdForUrl,
		"TemplateFunctionObfuscateString":         t.TemplateFunctionObfuscateString,
		"TemplateFunctionTranslateFormatString":   t.TemplateFunctionTranslateFormatString,
		"TemplateFunctionConcateStrings":          t.TemplateFunctionConcateStrings,
		"TemplateFunctionConcateHtml":             t.TemplateFunctionConcateHtml,
	})

	// setting local (theme) css files handels
	// both setters from Action and from Templates
	t.setLocalCSSFileNamesToTemplate()

	// setting externl (CDN) css files handels
	// both setters from Action and from Templates
	t.setExternalCSSFileNamesToTemplate()

	// setting local (theme) JS files for header handels
	// both setters from Action and from Templates
	t.setLocalHeaderJSFileNamesToTemplate()

	// setting extrnal (CDN) JS files for header handels
	// both setters from Action and from Templates
	t.setExternalHeaderJSFileNamesToTemplate()

	// setting local (theme) JS files for footer handels
	// both setters from Action and from Templates
	t.setLocalFooterJSFileNamesToTemplate()

	// setting external (CDN) JS files for header handels
	// both setters from Action and from Templates
	t.setExternalFooterJSFileNamesToTemplate()

	for _, file := range arrayOfFilesToParse {
		//fmt.Println(file)
		t._tmpl.ParseFiles(file)
	}

	if parseTemplateToString {
		var buf bytes.Buffer
		var err = t._tmpl.ExecuteTemplate(&buf, "base_layout", t.ViewData)
		if err != nil {
			fmt.Println(err)
		} else {
			return buf.String()
			// fmt.Println(s)
		}
	} else {

		var err = t._tmpl.ExecuteTemplate(t._w, "base_layout", t.ViewData)

		if err != nil {
			fmt.Println(err)
		}
	}

	return ""

}

func (t *TemplateEngine) getEnvConfigVars() persistence.EnvConfigVars {

	// @ persistence/EnvConfigVarsFilePath
	// load the environment variables
	envConfigFilePath := persistence.EnvConfigVarsFilePath

	var _envConfigVars = GetEnvConfigVars{}
	return _envConfigVars.Initiate(envConfigFilePath)

}

// Read the menu file for left aside
// navigation
func (t *TemplateEngine) GetAsideMenuItems(asideMenuFile string) []persistence.AsideMenu {
	// get the config vars for menu file path
	var _envConfigVars = t.getEnvConfigVars()
	// make the json path of given file
	var asideMenuFilePath = fmt.Sprintf("%s/%s.json", _envConfigVars.EnvConfigVarList["asideMenuPath"], asideMenuFile)
	//	fmt.Println(asideMenuFilePath)
	// init readfile as pointer (to discard after use to avoid memory garbage)
	var _readFile = &ReadFile{}

	var asideMenuList []persistence.AsideMenu = []persistence.AsideMenu{}
	json.Unmarshal([]byte(_readFile.Initiate(asideMenuFilePath)), &asideMenuList)

	// free memory
	_readFile = nil
	asideMenuFilePath = ""

	return asideMenuList
}

// Translates the give string
// Note that due to CDATA in XML it returns
// type HTML and not a String
func (t *TemplateEngine) TemplateFunctionTranslate(stringToTranslate string) template.HTML {
	return template.HTML(t._translator.Translate(stringToTranslate))
}

func (t *TemplateEngine) TemplateFunctionTranslateFormat(stringToTranslate string, args ...interface{}) template.HTML {
	return template.HTML(fmt.Sprintf(t._translator.Translate(stringToTranslate), args...))
}

func (t *TemplateEngine) TemplateFunctionTranslateFormatString(stringToTranslate string, args ...interface{}) string {
	return t._translator.Translate(fmt.Sprintf(stringToTranslate, args...))
}

// Translates the given salutation
func (t *TemplateEngine) TemplateFunctionMakeSalutation(salutation string) string {
	return t._translator.Translate(fmt.Sprintf("salutation_%s", salutation))
}

// concates the given params and return string
func (t *TemplateEngine) TemplateFunctionConcateStrings(string1 string, args ...interface{}) string {
	return fmt.Sprintf(string1, args...)
}

// concates the given params and return HTML
func (t *TemplateEngine) TemplateFunctionConcateHtml(string1 string, args ...interface{}) template.HTML {
	return template.HTML(fmt.Sprintf(string1, args...))
}

// Translates the given date
func (t *TemplateEngine) TemplateFunctionLocalizeDate(date string) string {
	return LocalizeDate(t._translator.Translate("date"), date)
}

// Translates the given date time
func (t *TemplateEngine) TemplateFunctionLocalizeDateTime(dateTime string) string {
	return LocalizeDateTime(t._translator.Translate("dateTime"), dateTime)
}

func (t *TemplateEngine) TemplateFunctionIsNullOrEmpty(stringToVerify string) bool {
	return IsNullOrEmpty(stringToVerify)
}

func (t *TemplateEngine) TemplateFunctionEncodeToURL(s string) template.URL {
	return template.URL(s)
}

func (t *TemplateEngine) TemplateFunctionGetAssociatedValue(array map[string]string, associativeId string) string {
	return array[associativeId]
}

func (t *TemplateEngine) TemplateFunctionSubstr(s string, from, to int) string {
	if to > len(s) {
		to = len(s)
	}
	return s[from:to]
}

func (t *TemplateEngine) TemplateFunctionSubStringWhitespaceSafe(s string, subFrom, subTo int, appendix string) string {
	// have the string substringed
	strinToReturn := t.TemplateFunctionSubstr(s, subFrom, subTo)
	// set the length placeholder with max string length
	nextWhitespaceAt := len(strinToReturn)
	// append trailing ... if text was shortend
	safeAppendix := ""

	// if the len of str is less than substring to
	// the retrun the string unchanged
	if len(strinToReturn) != len(s) {
		// walk through the string from end to begin
		// and look for closest whitespace
		for i := len(strinToReturn) - 1; i >= 0; i-- {
			// if whitespace found then break
			if string(strinToReturn[i]) == " " {
				// set the string breaker
				nextWhitespaceAt = i
				// fmt.Printf("break now at %d\n", nextWhitespaceAt)
				// set the trailing appendix
				safeAppendix = appendix
				break
			}
		}
	}
	return fmt.Sprintf("%s%s", strinToReturn[0:nextWhitespaceAt], safeAppendix)
}

func (t *TemplateEngine) TemplateFunctionStringToHtml(s string) template.HTML {
	//return template.HTML(s)
	return template.HTML(html.UnescapeString(s))
}

func (t *TemplateEngine) TemplateFunctionMakeArray() []string {
	return []string{}
}

func (t *TemplateEngine) TemplateFunctionAppendToArray(arr []string, val string) []string {
	return append(arr, val)
}

func (t *TemplateEngine) TemplateFunctionInArray(arr []string, val string) bool {
	return InArray(arr, val)
}

func (t *TemplateEngine) TemplateFunctionIterate(i int) int {
	return i + 1
}

func (t *TemplateEngine) TemplateFunctionFindReplace(s, stringToFind, stringToReplace string) string {
	return FindReplace(s, stringToFind, stringToReplace)
}

func (t *TemplateEngine) TemplateFunctionContainsString(s, match string) bool {
	return ContainsString(match, s)
}

func (t *TemplateEngine) TemplateFunctionEncryptIdForUrl(idToEncrypt string) string {
	return EncryptIdForUrl(idToEncrypt)
}

func (t *TemplateEngine) TemplateFunctionObfuscateString(stringToObfuscate string, obfuscationLength int) string {
	return ObfuscateString(stringToObfuscate, obfuscationLength)
}
