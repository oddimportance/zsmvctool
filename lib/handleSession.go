package lib

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

type HandleSession struct {
	w     http.ResponseWriter
	r     *http.Request
	store sessions.CookieStore
	key   []byte

	cookieTypeStaticStorage  string
	cookieTypeSessionStorage string

	sessionSessionStorage *sessions.Session
	sessionStaticStorage  *sessions.Session

	keySessionID string
}

func (h *HandleSession) Init(w http.ResponseWriter, r *http.Request) {

	h.keySessionID = "sessionID"

	h.initiateRequestAndWriterVars(w, r)

	h.initiateSesionStore()

	h.setCookieTypeValues()

	// set session cookie (expires on browser close)
	h.sessionSessionStorage = h.setSession(0, h.cookieTypeSessionStorage)

	// set static cookie for 14 days
	h.sessionStaticStorage = h.setSession(14, h.cookieTypeStaticStorage)

	//fmt.Println()

}

func (h *HandleSession) initiateSesionStore() {
	h.key = []byte("NTMG0Xc6Vref6H7h1Q62NN4KcJB9rKF2ZMqlg4HBzuIch8r7YrRJMEFTvQOqvv6NiONXa4l8o2AdbxOkgIQhVV5UGZFb9jdvnEdQdtSlNYHXe5GKHSYI3TDqshxFV4DJ")
	h.store = *sessions.NewCookieStore(h.key)
}

func (h *HandleSession) IsUserAuthenticated() bool {
	return false
}

func (h *HandleSession) initiateRequestAndWriterVars(w http.ResponseWriter, r *http.Request) {

	h.w = w
	h.r = r

}

func (h *HandleSession) setCookie(session *sessions.Session, key string, value string) {
	if session != nil {
		session.Values[key] = value
		session.Save(h.r, h.w)
	}
}

func (h *HandleSession) SetSessionCookie(key string, value string) {
	h.setCookie(h.sessionSessionStorage, key, value)
}

func (h *HandleSession) SetStaticCookie(key string, value string) {
	h.setCookie(h.sessionStaticStorage, key, value)
}

func (h *HandleSession) GetSessionCookie(key string) string {
	return h.getCookie(h.sessionSessionStorage, key)
}

func (h *HandleSession) GetStaticCookie(key string) string {
	return h.getCookie(h.sessionStaticStorage, key)
}

func (h *HandleSession) UnsetSessionCookie(key string) {
	h.unsetCookie(h.sessionSessionStorage, key)
}

func (h *HandleSession) UnsetStaticCookie(key string) {
	h.unsetCookie(h.sessionStaticStorage, key)
}

func (h *HandleSession) unsetCookie(session *sessions.Session, key string) {
	delete(session.Values, key)
	session.Save(h.r, h.w)
}

func (h *HandleSession) getCookie(session *sessions.Session, key string) string {
	// log.Println("accessing session key:", key)
	if session == nil {
		log.Printf("session is empty, key:-%s- could not be accessed\n", key)
		log.Println("trying to re-initialize sessiion")
		h.initiateSesionStore()
		return ""
	}
	var rawValue, ok = session.Values[key]
	if !ok {
		return ""
	} else {
		return rawValue.(string)
	}
}

func (h *HandleSession) GetLanguageFromCookie() string {
	return h.getCookie(h.sessionStaticStorage, "_defLang")
}

func (h *HandleSession) SetLanguageToCookie(language string) {
	h.setCookie(h.sessionStaticStorage, "_defLang", language)
}

func (h *HandleSession) setCookieTypeValues() {

	// _sestrg : session storage
	h.cookieTypeSessionStorage = "_sestrg"

	// _stadat : static storage
	h.cookieTypeStaticStorage = "_stastrg"

}

func (h *HandleSession) setSession(sessionAgeInDays int, sessionName string) *sessions.Session {
	session, err := h.store.Get(h.r, sessionName)

	if err != nil {
		fmt.Println("Error creating session store :", err)
		return nil
	}

	h.generateSessionId(session)

	// set max age in seconds
	session.Options.MaxAge = h.convertDaysToSeconds(sessionAgeInDays)
	return session
}

func (h *HandleSession) generateSessionId(session *sessions.Session) {
	if session.Values[h.keySessionID] == nil || session.Values[h.keySessionID] == "" {
		h.setCookie(session, h.keySessionID, GenerateRandomAlphaString(32))
	}
}

func (h *HandleSession) UnsetSessionStorageId() {
	h.UnsetSessionCookie(h.keySessionID)
}

func (h *HandleSession) SessionIDStaticStorage() string {
	return h.GetStaticCookie(h.keySessionID)
}

func (h *HandleSession) SessionIDSessionStorage() string {
	return h.GetSessionCookie(h.keySessionID)
}

func (h *HandleSession) convertDaysToSeconds(days int) int {
	return (((60 * 60) * 24) * days)
}
