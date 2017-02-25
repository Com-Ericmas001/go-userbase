package userbase

import "database/sql"

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
