package web

import (
	"regexp"
	"unicode"
)

// isEmailValid checks if the email provided is valid by regex.
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	length, number, upper, special := verifyPassword(password)

	if !length || !number || !upper || special {
		return false
	}

	return true
}

func verifyPassword(password string) (sevenOrMore, number, upper, special bool) {
	letters := 0
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			return false, false, false, false
		}
	}
	sevenOrMore = letters >= 7
	return
}
