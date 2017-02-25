package userbase

import "time"

//AuthenticationInfo is a model
type AuthenticationInfo struct {
	Password string
	Email    string
}

//ProfileInfo is a model
type ProfileInfo struct {
	DisplayName string
}

//Token is a model
type Token struct {
	ID         string
	ValidUntil time.Time
}
