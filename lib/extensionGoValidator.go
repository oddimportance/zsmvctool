package lib

import (
	"net"
	"regexp"
	"strings"
)

var (
	userRegexp = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	hostRegexp = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
	// As per RFC 5332 secion 3.2.3: https://tools.ietf.org/html/rfc5322#section-3.2.3
	// Dots are not allowed in the beginning, end or in occurances of more than 1 in the email address
	userDotRegexp = regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")
)

// Validate checks format of a given email and resolves its host name.
func ValidateEmailHost(email string) bool {
	if len(email) < 6 || len(email) > 254 {
		return false
	}

	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return false
	}

	user := email[:at]
	host := email[at+1:]

	if len(user) > 64 {
		return false
	}
	if userDotRegexp.MatchString(user) || !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		return false
	}

	switch host {
	case "localhost", "example.com":
		return false
	}

	if _, err := net.LookupMX(host); err != nil {
		if _, err := net.LookupIP(host); err != nil {
			// Only fail if both MX and A records are missing - any of the
			// two is enough for an email to be deliverable
			return false
		}
	}

	return true
}
