package emailvalidator

import "regexp"

var emailRegexp *regexp.Regexp

func init() {
	emailRegexp = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
}

func IsValidEmail(email string) bool {
	return emailRegexp.MatchString(email)
}
