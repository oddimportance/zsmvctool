package lib

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"strings"

	"zsmvctool-api/persistence"
)

type Translator_i18n struct {

	// Seted language
	language string
	// Environment vars
	_envConfigVars GetEnvConfigVars
	// File path to handle readfile
	_handleFilePath *HandleFilePath
	// Handle Session for language cookie
	_handleSession *HandleSession
	// panic error
	_handleRedirectAndPanic *HandleRedirectAndPanic
	// read file
	_readFile ReadFile
	// string holding
	translatedTags map[string]string
	// available Languages
	availableLanguages []string
}

func (t *Translator_i18n) InitTranslator(handleFilePath *HandleFilePath, handleSession *HandleSession, _handleRedirectAndPanic *HandleRedirectAndPanic) {

	t._handleSession = handleSession
	t._handleRedirectAndPanic = _handleRedirectAndPanic
	t._handleFilePath = handleFilePath

	t.handleLanguage()

	t.populateTranslatedTags((t.parseXmlToPersistence(t.readDataFromXml())))

}

func (t *Translator_i18n) readDataFromXml() []byte {

	// Prepare the path
	var pathToFormFile = fmt.Sprintf("%s/%s.xml", t._handleFilePath.GetLanguagePath(), t.language)

	// Initate the class gloablly for futher prossessing
	t._readFile = ReadFile{}

	// get the content from xml file
	var xmlContent = t._readFile.Initiate(pathToFormFile)

	// look if file did not exist, if yes controller
	// will handle panic error redirection
	if len(xmlContent) == 0 {
		t._handleRedirectAndPanic.TriggerPanic("00009")
	}

	return xmlContent

}

func (t *Translator_i18n) parseXmlToPersistence(xmlRaw []byte) persistence.LanguageTags {

	// initiate unmarshal
	var languageTags persistence.LanguageTags

	// unmarshal
	xml.Unmarshal([]byte(xmlRaw), &languageTags)

	return languageTags

}

func (t *Translator_i18n) populateTranslatedTags(languageTags persistence.LanguageTags) {

	t.translatedTags = map[string]string{}

	for _, val := range languageTags.Tags {
		t.translatedTags[val.Id] = val.Message
	}

}

// Translates the give string
func (t *Translator_i18n) Translate(key string) string {
	if t.translatedTags[key] == "" {
		return key
	} else {
		return t.translatedTags[key]
	}
}

// Translates the give string
func (t *Translator_i18n) TranslateSafeWithAttachment(key string, attachment string) string {
	var keyToFetch string = fmt.Sprintf("%s_%s", key, strings.ToLower(attachment))
	if t.translatedTags[keyToFetch] != "" {
		return t.translatedTags[keyToFetch]
	} else {
		return t.Translate(key)
	}
}

func (t *Translator_i18n) StringToHTML(stringToConvert string) template.HTML {
	return template.HTML(stringToConvert)
}

func (t *Translator_i18n) handleLanguage() {

	t.setEnvConfigVars()

	var languageInCookie string = t._handleSession.GetLanguageFromCookie()

	if t.verifyLangInCookie(languageInCookie) {
		t.language = languageInCookie
	} else {
		t.language = t._envConfigVars.GetConfVar("languageDefault")
	}

}

func (t *Translator_i18n) setEnvConfigVars() {

	t._envConfigVars = GetEnvConfigVars{}
	t._envConfigVars.Initiate(persistence.EnvConfigVarsFilePath)

}

func (t *Translator_i18n) verifyLangInCookie(languageInCookie string) bool {

	//	var availableLanguages = strings.Split(t._envConfigVars.GetConfVar("languagesAvailable"), ",")
	//	return InArray(availableLanguages, languageInCookie)
	return VerifyLanguage(languageInCookie, t._envConfigVars.GetConfVar("languagesAvailable"))

}
