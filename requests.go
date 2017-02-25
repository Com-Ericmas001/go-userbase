package userbase

//CreateUserRequest is a request
type CreateUserRequest struct {
	Username       string
	Authentication AuthenticationInfo
	Profile        ProfileInfo
}