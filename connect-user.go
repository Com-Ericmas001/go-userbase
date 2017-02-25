package userbase

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//ValidateCredentials validates if user and password are correct
func (context DbContext) ValidateCredentials(username string, password string) ConnectUserResponse {

	id := context.IDFromUsername(username)
	if id == 0 {
		return invalidConnectUserResponse()
	}

	stmt, err := context.Db.Prepare("SELECT Password FROM UserAuthentications WHERE IdUser = ?")
	checkErr(err)
	defer stmt.Close()

	var passwordFromdb string
	err = stmt.QueryRow(id).Scan(&passwordFromdb)
	checkErr(err)

	err = bcrypt.CompareHashAndPassword([]byte(passwordFromdb), []byte(context.saltPassword(password)))
	if err != nil {
		return invalidConnectUserResponse()
	}

	return context.newTokenSuccessResponse(id)
}

//ValidateToken validates if username and token is valid
func (context DbContext) ValidateToken(username string, token string) ConnectUserResponse {

	id := context.IDFromUsername(username)
	if id == 0 {
		return invalidConnectUserResponse()
	}

	stmt, err := context.Db.Prepare("SELECT Expiration FROM UserTokens WHERE IdUser = ? AND Token = ?")
	checkErr(err)
	defer stmt.Close()

	var expiration time.Time
	err = stmt.QueryRow(id, token).Scan(&expiration)

	if err == sql.ErrNoRows {
		return invalidConnectUserResponse()
	}
	checkErr(err)

	if time.Now().After(expiration) {
		return invalidConnectUserResponse()
	}

	newExpiration := time.Now().Add(time.Minute * time.Duration(10))

	stmt2, err := context.Db.Prepare("UPDATE UserTokens SET Expiration = ? WHERE IdUser = ? AND Token = ?")
	checkErr(err)
	defer stmt2.Close()

	_, err = stmt2.Exec(newExpiration, id, token)

	return ConnectUserResponse{
		IDUser: id,
		TokenResponse: TokenSuccessResponse{
			Success: true,
			Token: Token{
				ID:         token,
				ValidUntil: newExpiration}}}
}
