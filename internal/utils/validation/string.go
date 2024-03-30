package validation

import "regexp"

var (
	userNamePattern = regexp.MustCompile(`^[A-Za-z ]+$`)
	emailPattern    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func IsValidUserName(value string) bool {
	return len(value) >= 2 && userNamePattern.MatchString(value)
}

func IsValidEmail(value string) bool {
	return emailPattern.MatchString(value)
}
