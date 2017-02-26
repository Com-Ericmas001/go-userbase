package userbase

import "golang.org/x/crypto/bcrypt"

//ModifyProfile modifies profile
func (context DbContext) ModifyProfile(request ModifyProfileRequest) TokenSuccessResponse {

	connection := context.ValidateToken(request.Username, request.Token)
	if !connection.TokenResponse.Success {
		return invalidTokenSuccessResponse()
	}

	if len(request.Profile.DisplayName) > 0 {
		if !ValidateDisplayName(request.Profile.DisplayName) {
			return invalidTokenSuccessResponse()
		}

		stmt, err := context.Db.Prepare("UPDATE UserProfiles SET DisplayName = ? WHERE IdUser = ?")
		checkErr(err)
		defer stmt.Close()

		_, err = stmt.Exec(request.Profile.DisplayName, connection.IDUser)
		checkErr(err)
	}

	return connection.TokenResponse
}

//ModifyCredentials modifies credentials
func (context DbContext) ModifyCredentials(request ModifyCredentialsRequest) TokenSuccessResponse {

	connection := context.ValidateToken(request.Username, request.Token)
	if !connection.TokenResponse.Success {
		return invalidTokenSuccessResponse()
	}

	if len(request.Authentication.Password) > 0 {
		if !ValidatePassword(request.Authentication.Password) {
			return invalidTokenSuccessResponse()
		}

		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(context.saltPassword(request.Authentication.Password)), bcrypt.DefaultCost)
		checkErr(err)

		stmt, err := context.Db.Prepare("UPDATE UserAuthentications SET Password = ? WHERE IdUser = ?")
		checkErr(err)
		defer stmt.Close()

		_, err = stmt.Exec(hashedPassword, connection.IDUser)
		checkErr(err)
	}

	if len(request.Authentication.Email) > 0 {
		if !ValidateEmail(request.Authentication.Email) {
			return invalidTokenSuccessResponse()
		}

		stmt, err := context.Db.Prepare("UPDATE UserAuthentications SET RecoveryEmail = ? WHERE IdUser = ?")
		checkErr(err)
		defer stmt.Close()

		_, err = stmt.Exec(request.Authentication.Email, connection.IDUser)
		checkErr(err)
	}

	return connection.TokenResponse
}