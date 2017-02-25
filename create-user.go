package userbase

import "golang.org/x/crypto/bcrypt"

//CreateUser creates a user
func (context DbContext) CreateUser(request CreateUserRequest) {

	//User
	stmt, err := context.Db.Prepare("INSERT INTO Users(Name, Active) VALUES(?, 1)")
	checkErr(err)

	res, err := stmt.Exec(request.Username)
	checkErr(err)

	idUser, err := res.LastInsertId()
	checkErr(err)

	//Auth
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(context.saltPassword(request.Authentication.Password)), bcrypt.DefaultCost)
	checkErr(err)

	stmt, err = context.Db.Prepare("INSERT INTO UserAuthentications(IdUser, Password, RecoveryEmail) VALUES(?, ?, ?)")
	checkErr(err)

	res, err = stmt.Exec(idUser, hashedPassword, request.Authentication.Email)
	checkErr(err)

	//Profile
	stmt, err = context.Db.Prepare("INSERT INTO UserProfiles(IdUser, DisplayName) VALUES(?, ?)")
	checkErr(err)

	res, err = stmt.Exec(idUser, request.Profile.DisplayName)
	checkErr(err)

	//Settings
	stmt, err = context.Db.Prepare("INSERT INTO UserSettings(IdUser, IdUserAccessTypeListFriends) VALUES(?, 1)")
	checkErr(err)

	res, err = stmt.Exec(idUser)
	checkErr(err)
}
