package userbase

import "golang.org/x/crypto/bcrypt"

//CreateUser creates a user
func (context DbContext) CreateUser(request CreateUserRequest) ConnectUserResponse {

	if context.UsernameExists(request.Username) || context.EmailExists(request.Authentication.Email) {
		return invalidConnectUserResponse()
	}

	if !ValidateUsername(request.Username) || !ValidateEmail(request.Authentication.Email) || !ValidatePassword(request.Authentication.Password) || !ValidateDisplayName(request.Profile.DisplayName) {
		return invalidConnectUserResponse()
	}

	idUser := context.createUserEntity(request)
	context.createUserAthenticationEntity(request, idUser)
	context.createUserProfileEntity(request, idUser)
	context.createUserSettingEntity(request, idUser)

	return context.ValidateCredentials(request.Username, request.Authentication.Password)
}

func (context DbContext) createUserEntity(request CreateUserRequest) int64 {

	stmt, err := context.Db.Prepare("INSERT INTO Users(Name, Active) VALUES(?, 1)")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(request.Username)
	checkErr(err)

	idUser, err := res.LastInsertId()
	checkErr(err)

	return idUser
}
func (context DbContext) createUserAthenticationEntity(request CreateUserRequest, idUser int64) {

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(context.saltPassword(request.Authentication.Password)), bcrypt.DefaultCost)
	checkErr(err)

	stmt, err := context.Db.Prepare("INSERT INTO UserAuthentications(IdUser, Password, RecoveryEmail) VALUES(?, ?, ?)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(idUser, hashedPassword, request.Authentication.Email)
	checkErr(err)
}
func (context DbContext) createUserProfileEntity(request CreateUserRequest, idUser int64) {

	stmt, err := context.Db.Prepare("INSERT INTO UserProfiles(IdUser, DisplayName) VALUES(?, ?)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(idUser, request.Profile.DisplayName)
	checkErr(err)
}
func (context DbContext) createUserSettingEntity(request CreateUserRequest, idUser int64) {

	stmt, err := context.Db.Prepare("INSERT INTO UserSettings(IdUser, IdUserAccessTypeListFriends) VALUES(?, 1)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(idUser)
	checkErr(err)
}
