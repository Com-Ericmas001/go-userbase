package userbase

import (
	"unicode"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
)

//ValidateDisplayName validates displayName
func ValidateDisplayName(displayName string) bool {
	if len(displayName) < 3 {
		return false
	}
	return true
}

//ValidateEmail validates email
func ValidateEmail(email string) bool {
	return govalidator.IsEmail(email)
}

//ValidatePassword validates password
func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	for i := 0; i < len(password); i++ {
		var r rune
		r, _ = utf8.DecodeRune([]byte(password)[i:1])
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) && !unicode.IsSymbol(r) {
			return false
		}
	}
	return true
}

//ValidateUsername validates username
func ValidateUsername(username string) bool {
	if len(username) < 3 {
		return false
	}
	for i := 0; i < len(username); i++ {
		var r rune
		r, _ = utf8.DecodeRune([]byte(username)[i:1])
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
