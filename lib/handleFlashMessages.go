package lib

import (
	"strings"
)

const (
	keySuccess       string = "flashMsgSuccess"
	keyWarning              = "flashMsgWarning"
	keyError                = "flashMsgError"
	keyUncategorized        = "flashMsgUncategorized"
)

// Instantiate HandleFlashMessages
//	c._flashMessages = lib.HandleFlashMessages{}
// Init with HandleSession Pointer
//	c._flashMessages.Init(c._handleSession)
// Set a flash message (available types Success, Warning, Error, Uncategorized)
//	c._flashMessages.SetFlashMessage(c._flashMessages.TypeSuccess, "Ya Rab")
// Get message
//	c._flashMessages.GetFlashMessage(c._flashMessages.TypeSuccess)
type HandleFlashMessages struct {

	// Type string, store the success msg
	TypeSuccess string

	// Type string, store warning
	TypeWarning string

	// Type string, store the error
	TypeError string

	// Type string, store any special messages
	TypeUncategorized string

	// Session to save messages to cookie
	_handleSession *HandleSession
}

func (f *HandleFlashMessages) Init(handleSession *HandleSession) {

	f._handleSession = handleSession
	f.TypeSuccess = keySuccess
	f.TypeWarning = keyWarning
	f.TypeError = keyError
	f.TypeUncategorized = keyUncategorized

}

func (f *HandleFlashMessages) SetFlashMessage(flashType, m string) {
	f._handleSession.SetSessionCookie(flashType, m)
}

func (f *HandleFlashMessages) GetFlashMessage(flashType string) string {
	defer f._handleSession.UnsetSessionCookie(flashType)
	return f._handleSession.GetSessionCookie(flashType)
}

func (f *HandleFlashMessages) GetAllFlashMessages() map[string]string {

	var messagesToGet []string = []string{f.TypeSuccess, f.TypeWarning, f.TypeError, f.TypeUncategorized}
	var messagesToReturn map[string]string = map[string]string{}
	valueHolder := ""
	for _, flashKey := range messagesToGet {
		valueHolder = f.GetFlashMessage(flashKey)
		if valueHolder != "" {
			messagesToReturn[strings.TrimLeft(flashKey, "flashMsg")] = valueHolder
		}
	}

	return messagesToReturn

}
