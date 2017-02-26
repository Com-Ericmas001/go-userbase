package userbase

import (
	"database/sql"
	"fmt"
)

//IDFromUsername returns IdUser from a username
func (context DbContext) IDFromUsername(username string) int64 {

	stmt, err := context.Db.Prepare("SELECT IdUser FROM Users WHERE Name = ? AND Active != 0")
	checkErr(err)
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(username).Scan(&id)
	if err == sql.ErrNoRows {
		id = 0
	} else {
		checkErr(err)
	}

	return id
}

//IDFromEmail returns IdUser from a email
func (context DbContext) IDFromEmail(email string) int64 {

	stmt, err := context.Db.Prepare("SELECT IdUser FROM UserAuthentications JOIN Users USING(IdUser) WHERE RecoveryEmail = ? AND Active != 0")
	checkErr(err)
	defer stmt.Close()

	var id int64
	err = stmt.QueryRow(email).Scan(&id)
	if err == sql.ErrNoRows {
		id = 0
	} else {
		checkErr(err)
	}

	return id
}

//UsernameExists returns true if user exists
func (context DbContext) UsernameExists(username string) bool {
	return context.IDFromUsername(username) != 0
}

//EmailExists returns true if email exists
func (context DbContext) EmailExists(email string) bool {
	return context.IDFromEmail(email) != 0
}

//UserSummary returns summary of info on a user
func (context DbContext) UserSummary(username string, token string) UserSummaryResponse {

	connection := context.ValidateToken(username, token)
	if !connection.Success {
		fmt.Println("invalid", username, token)
		return invalidUserSummaryResponse()
	}

	stmt, err := context.Db.Prepare("SELECT DisplayName FROM UserProfiles WHERE IdUser = ?")
	checkErr(err)
	defer stmt.Close()

	var displayName string
	err = stmt.QueryRow(connection.IDUser).Scan(&displayName)
	checkErr(err)

	return UserSummaryResponse{
		DisplayName: displayName,
		Success:     connection.Success,
		Token:       connection.Token}
}
