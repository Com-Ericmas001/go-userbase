package userbase

//AuthenticationInfo is a model
type AuthenticationInfo struct {
	Password string
	Email    string
}

//ProfileInfo is a model
type ProfileInfo struct {
	DisplayName string
}