package userbase

import (
	"database/sql"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

	//Justification: sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var isInit = false

//DbContext is a context
type DbContext struct {
	Db   *sql.DB
	salt string
}

//AuthenticationInfo is a model
type AuthenticationInfo struct {
	Password string
	Email    string
}

//ProfileInfo is a model
type ProfileInfo struct {
	DisplayName string
}

//CreateUserRequest is a request
type CreateUserRequest struct {
	Username       string
	Authentication AuthenticationInfo
	Profile        ProfileInfo
}

//FnCreateDatabase is a function used to create the database
type FnCreateDatabase func(*sql.DB)

// Init the db
func Init(dbPath string, salt string, fnCreate FnCreateDatabase) *DbContext {
	isFirstTime := false
	if _, err := os.Stat(dbPath); err != nil {
		isFirstTime = true
	}

	db, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	context := DbContext{Db: db, salt: salt}

	if isFirstTime {
		fnCreate(db)
		context.CreateUser(CreateUserRequest{Username: "root", Authentication: AuthenticationInfo{Password: "abcd1234", Email: "root@ericmas001.com"}, Profile: ProfileInfo{DisplayName: "ADMIN"}})
	}

	return &context
}

// Close the db
func (context DbContext) Close() {
	context.Db.Close()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (context DbContext) saltPassword(password string) string {
	return context.salt + password
}

//IDFromUsername returns IdUser from a username
func (context DbContext) IDFromUsername(username string) int {

	stmt, err := context.Db.Prepare("SELECT IdUser FROM Users WHERE Name = ?")
	checkErr(err)
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(username).Scan(&id)
	if err == sql.ErrNoRows {
		id = 0
	} else {
		checkErr(err)
	}

	return id
}

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

// CreateDatabase creates the db
func CreateDatabase(currDb *sql.DB) {
	sqlStmt := `
	create table UserAccessTypes (IdUserAccessType INTEGER PRIMARY KEY, Name TEXT, Value INT);
	create table UserRelationTypes (IdUserRelationType INTEGER PRIMARY KEY, Name TEXT);
	create table Users (IdUser INTEGER PRIMARY KEY, Name TEXT, Active BOOLEAN);
	create table UserAuthentications (IdUser INTEGER UNIQUE REFERENCES Users(IdUser), Password TEXT, RecoveryEmail TEXT);
	create table UserProfiles (IdUser INTEGER UNIQUE REFERENCES Users(IdUser), DisplayName TEXT);
	create table UserRecoveryTokens (IdUserRecoveryToken INTEGER PRIMARY KEY, IdUser INTEGER REFERENCES Users(IdUser), Token TEXT, Expiration DATETIME);
	create table UserRelations (IdUserRelation INTEGER PRIMARY KEY, IdUserOwner INTEGER REFERENCES Users(IdUser), IdUserLinked INTEGER REFERENCES Users(IdUser), IdUserRelationType INTEGER REFERENCES UserRelationTypes(IdUserRelationType));
	create table UserSettings (IdUser INTEGER UNIQUE REFERENCES Users(IdUser), IdUserAccessTypeListFriends INTEGER REFERENCES UserAccessTypes(IdUserAccessType));
    INSERT INTO UserAccessTypes(IdUserAccessType, Name, Value) VALUES(1, 'Everybody', 10);
    INSERT INTO UserAccessTypes(IdUserAccessType, Name, Value) VALUES(2, 'EverybodyNotBlocked', 20);
    INSERT INTO UserAccessTypes(IdUserAccessType, Name, Value) VALUES(3, 'Friends', 30);
    INSERT INTO UserAccessTypes(IdUserAccessType, Name, Value) VALUES(4, 'JustMe', 40);
    INSERT INTO UserRelationTypes(IdUserRelationType, Name) VALUES(1, 'Friend');
    INSERT INTO UserRelationTypes(IdUserRelationType, Name) VALUES(2, 'Blocked');
	`
	_, err := currDb.Exec(sqlStmt)
	checkErr(err)
}
