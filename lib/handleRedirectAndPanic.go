package lib

import (
	"fmt"
	"net/http"
)

type HandleRedirectAndPanic struct {
	HasPanicError               bool
	PanicCode                   string
	IgnoreRedirect              bool
	_responseWriter             http.ResponseWriter
	_request                    *http.Request
	HttpStatusPermanentRedirect int
	HttpStatusTemporaryRedirect int
	HttpStatusFound             int
	HttpStatusSeeOther          int
}

func (h *HandleRedirectAndPanic) GetHttpRequest() *http.Request {
	return h._request
}

func (h *HandleRedirectAndPanic) RedirectToPanicUrl() {
	// To avoid following http error
	// error: http: multiple response.WriteHeader calls
	// look if the header is already set
	// If the header is already set, then do not call redirect
	if h.IgnoreRedirect {
		return
	}
	h.IgnoreRedirect = true
	http.Redirect(h._responseWriter, h._request, fmt.Sprintf("%s/%s", "/error/error-500", h.PanicCode), h.HttpStatusSeeOther)
	return
}

func (h *HandleRedirectAndPanic) SetHttpParams(w http.ResponseWriter, r *http.Request) {
	h._responseWriter = w
	h._request = r

	// set the http status
	h.HttpStatusTemporaryRedirect = http.StatusTemporaryRedirect
	h.HttpStatusPermanentRedirect = http.StatusPermanentRedirect
	h.HttpStatusFound = http.StatusFound
	h.HttpStatusSeeOther = http.StatusSeeOther
}

// sets HasPanic to true, sets panic code
// PanicCode is basically of type Int,
// but to support a leading 0, its type
// is set to string
// @ panicCode string
func (h *HandleRedirectAndPanic) TriggerPanic(panicCode string) {
	if !h.HasPanicError {
		h.HasPanicError = true
		h.PanicCode = panicCode
	}
}

func (h *HandleRedirectAndPanic) RedirectToUrl(url string, httpStatus int) {
	// To avoid multiple header error (http: multiple response.WriteHeader calls)
	// Set the ignore panic redirect flag to true
	h.IgnoreRedirect = true
	http.Redirect(h._responseWriter, h._request, url, httpStatus)
}

// Redirect to given controller/action
// @param controller string
// @param action string
// @param httpStatus HandleRedirectAndPanic.HttpStatus
func (h *HandleRedirectAndPanic) RedirectToControllerAction(controller, action string, httpStatus int) {
	h.IgnoreRedirect = true
	http.Redirect(h._responseWriter, h._request, fmt.Sprintf("/%s/%s", controller, action), httpStatus)
}
