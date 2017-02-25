package userbase

//CreateUserRequest is a request
type CreateUserRequest struct {
	Username       string
	Authentication AuthenticationInfo
	Profile        ProfileInfo
}

//ModifyCredentialsRequest is a request
type ModifyCredentialsRequest struct {
	Username       string
	Token          string
	Authentication AuthenticationInfo
}

//ModifyProfileRequest is a request
type ModifyProfileRequest struct {
	Username string
	Token    string
	Profile  ProfileInfo
}
