package userbase

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//FnSendEmail is a function used to send e-mail
type FnSendEmail func(Token, string)

//SendRecoveryToken sends a recoveryToken
func (context *DbContext) SendRecoveryToken(username string, fnSendEmail FnSendEmail) bool {

	id := context.IDFromUsername(username)
	if id == 0 {
		return false
	}

	stmt, err := context.Db.Prepare("SELECT RecoveryEmail FROM UserAuthentications WHERE IdUser = ?")
	checkErr(err)
	defer stmt.Close()

	var email string
	err = stmt.QueryRow(id).Scan(&email)
	checkErr(err)

	token := context.newRecoveryTokenSuccessResponse(id)
	fnSendEmail(token.Token, email)

	return true
}

//ResetPassword resets a password
func (context *DbContext) ResetPassword(username string, recoveryToken string, newPassword string) ConnectUserResponse {

	id := context.IDFromUsername(username)
	if id == 0 {
		return invalidConnectUserResponse()
	}

	if !ValidatePassword(newPassword) {
		return invalidConnectUserResponse()
	}

	stmt, err := context.Db.Prepare("SELECT Expiration FROM UserRecoveryTokens WHERE IdUser = ? AND Token = ?")
	checkErr(err)
	defer stmt.Close()

	var expiration time.Time
	err = stmt.QueryRow(id, recoveryToken).Scan(&expiration)

	if err == sql.ErrNoRows {
		return invalidConnectUserResponse()
	}
	checkErr(err)

	if time.Now().After(expiration) {
		return invalidConnectUserResponse()
	}

	newExpiration := time.Now()

	stmt2, err := context.Db.Prepare("UPDATE UserRecoveryTokens SET Expiration = ? WHERE IdUser = ? AND Token = ?")
	checkErr(err)
	defer stmt2.Close()

	_, err = stmt2.Exec(newExpiration, id, recoveryToken)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(context.saltPassword(newPassword)), bcrypt.DefaultCost)
	checkErr(err)

	stmt3, err := context.Db.Prepare("UPDATE UserAuthentications SET Password = ? WHERE IdUser = ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt3.Exec(hashedPassword, id)
	checkErr(err)

	return context.ValidateCredentials(username, newPassword)
}
