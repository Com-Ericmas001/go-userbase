package userbase

import "time"

//PurgeRecoveryTokens purges recovery tokens
func (context *DbContext) PurgeRecoveryTokens() {

	stmt, err := context.Db.Prepare("DELETE FROM UserRecoveryTokens WHERE Expiration < ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(time.Now())
	checkErr(err)

}

//PurgeConnectionTokens purges connection tokens
func (context *DbContext) PurgeConnectionTokens() {

	stmt, err := context.Db.Prepare("DELETE FROM UserTokens WHERE Expiration < ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(time.Now())
	checkErr(err)
}

//PurgeUsers purges deactivated users
func (context *DbContext) PurgeUsers() {
	rows, err := context.Db.Query("select IdUser from Users WHERE Active = 0")
	checkErr(err)

	var toKill []int64

	defer rows.Close()
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		checkErr(err)
		toKill = append(toKill, id)

	}
	err = rows.Err()
	checkErr(err)
	rows.Close()

	for _, id := range toKill {
		context.purgeProfile(id)
		context.purgeAuthentication(id)
		context.purgeSetting(id)
		context.purgeTokens(id)
		context.purgeRecoveryTokens(id)

		context.purgeUser(id)
	}
}

func (context *DbContext) purgeProfile(id int64) {

	stmt, err := context.Db.Prepare("DELETE FROM UserProfiles WHERE IdUser = ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(id)
	checkErr(err)
}

func (context *DbContext) purgeAuthentication(id int64) {

	stmt, err := context.Db.Prepare("DELETE FROM UserAuthentications WHERE IdUser = ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(id)
	checkErr(err)
}

func (context *DbContext) purgeSetting(id int64) {

	stmt, err := context.Db.Prepare("DELETE FROM UserSettings WHERE IdUser = ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(id)
	checkErr(err)
}

func (context *DbContext) purgeTokens(id int64) {

	stmt, err := context.Db.Prepare("DELETE FROM UserTokens WHERE IdUser = ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(id)
	checkErr(err)
}

func (context *DbContext) purgeRecoveryTokens(id int64) {

	stmt, err := context.Db.Prepare("DELETE FROM UserRecoveryTokens WHERE IdUser = ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(id)
	checkErr(err)
}

func (context *DbContext) purgeUser(id int64) {

	stmt, err := context.Db.Prepare("DELETE FROM Users WHERE IdUser = ?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(id)
	checkErr(err)
}
