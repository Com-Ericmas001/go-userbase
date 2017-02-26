package userbase

import "github.com/asaskevich/govalidator"

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

	return govalidator.IsPrintableASCII(password)
}

//ValidateUsername validates username
func ValidateUsername(username string) bool {
	if len(username) < 3 {
		return false
	}

	return govalidator.IsAlphanumeric(username)
}
