package main

import (
	"database/sql"
	"fmt"
	"go-userbase"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello")

	context := userbase.Init("D:\\test.db", "DummyUserbaseSalt", CreateDatabase, InitDatabase)
	defer context.Close()
	fmt.Println("ericmas001:", context.IDFromUsername("ericmas001"), context.UsernameExists("ericmas001"))
	fmt.Println("root:", context.IDFromUsername("root"), context.UsernameExists("root"))
	fmt.Println("ericmas001@hotmail.com:", context.IDFromEmail("ericmas001@hotmail.com"), context.EmailExists("ericmas001@hotmail.com"))
	fmt.Println("root@ericmas001.com:", context.IDFromEmail("root@ericmas001.com"), context.EmailExists("root@ericmas001.com"))
	ok := context.ValidateCredentials("root", "abcd1234")
	fmt.Println("Connect root ok:", ok)
	fmt.Println("Connect root wrong:", context.ValidateCredentials("root", "abcd12345"))

	dumpUserTokens(context)

	fmt.Println("Validate root wrong:", context.ValidateToken("root", "fe8e5991-58e1-48d8-ad6b-9e836d1695c8"))
	fmt.Println("Validate root ok:", context.ValidateToken("root", ok.TokenResponse.Token.ID))

	dumpUserTokens(context)

	dumpUsers(context)
	fmt.Println("ModifyCredentials:", context.ModifyCredentials(userbase.ModifyCredentialsRequest{Username: "root", Token: ok.TokenResponse.Token.ID, Authentication: userbase.AuthenticationInfo{Email: "user@ericmas001.com"}}))
	fmt.Println("ModifyProfile:", context.ModifyProfile(userbase.ModifyProfileRequest{Username: "root", Token: ok.TokenResponse.Token.ID, Profile: userbase.ProfileInfo{DisplayName: "BOB"}}))
	dumpUsers(context)
	fmt.Println("ModifyCredentials:", context.ModifyCredentials(userbase.ModifyCredentialsRequest{Username: "root", Token: ok.TokenResponse.Token.ID, Authentication: userbase.AuthenticationInfo{Password: "qwerty12345"}}))
	dumpUsers(context)
	fmt.Println("ModifyCredentials:", context.ModifyCredentials(userbase.ModifyCredentialsRequest{Username: "root", Token: ok.TokenResponse.Token.ID, Authentication: userbase.AuthenticationInfo{Email: "root@ericmas001.com", Password: "abcd1234"}}))
	fmt.Println("ModifyProfile:", context.ModifyProfile(userbase.ModifyProfileRequest{Username: "root", Token: ok.TokenResponse.Token.ID, Profile: userbase.ProfileInfo{DisplayName: "ADMIN"}}))
	dumpUsers(context)

	fmt.Println("Connect root ok:", context.ValidateCredentials("root", "abcd1234"))
	fmt.Println("Disconnect old root:", context.Disconnect("root", ok.TokenResponse.Token.ID))
	dumpUserTokens(context)
	//dumpAnotherTable(context)
}

func dumpAnotherTable(context *userbase.DbContext) {
	rows, err := context.Db.Query("select IdAnotherTable, Name from AnotherTable")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
func dumpUserTokens(context *userbase.DbContext) {
	rows, err := context.Db.Query("select IdUser, Token, Expiration from UserTokens")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var token string
		var expiration time.Time
		err = rows.Scan(&id, &token, &expiration)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, token, expiration)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
func dumpUsers(context *userbase.DbContext) {
	rows, err := context.Db.Query("SELECT IdUser, Name, Password, RecoveryEmail, DisplayName, IdUserAccessTypeListFriends FROM Users JOIN UserAuthentications USING(IdUser) JOIN UserProfiles USING(IdUser) JOIN UserSettings USING(IdUser)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var password string
		var email string
		var display string
		var idAccess int
		err = rows.Scan(&id, &name, &password, &email, &display, &idAccess)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, password, email, display, idAccess)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

// InitDatabase inits the db
func InitDatabase(context *userbase.DbContext) {

	context.CreateUser(userbase.CreateUserRequest{
		Username: "root",
		Authentication: userbase.AuthenticationInfo{
			Password: "abcd1234",
			Email:    "root@ericmas001.com"},
		Profile: userbase.ProfileInfo{
			DisplayName: "ADMIN"}})

	dumpUsers(context)

	stmt, err := context.Db.Prepare("insert INTO UserTokens (IdUser, Token, Expiration) VALUES (1, 'fe8e5991-58e1-48d8-ad6b-9e836d1695c8', ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now().Add(time.Minute * time.Duration(-42)))
	if err != nil {
		log.Fatal(err)
	}
}

// CreateDatabase creates the db
func CreateDatabase(currDb *sql.DB) {

	//Get basic Userbase database
	userbase.CreateDatabase(currDb)

	//Create more stuff
	sqlStmt := `
	create table AnotherTable (IdAnotherTable INTEGER PRIMARY KEY, Name TEXT);
	insert INTO AnotherTable (IdAnotherTable, Name) VALUES (42,"Answer to Life, the universe, and everything");
	insert INTO AnotherTable (IdAnotherTable, Name) VALUES (84,"Double");
	insert INTO AnotherTable (IdAnotherTable, Name) VALUES (21,"Half");
	insert INTO AnotherTable (Name) VALUES ("Next");
	`
	_, err := currDb.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}
