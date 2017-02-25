package userbase

import (
	"database/sql"
	"log"
	"os"

	//Justification: sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var isInit = false

//DbContext is a context
type DbContext struct {
	Db   *sql.DB
	salt string
}

//FnCreateDatabase is a function used to create the database
type FnCreateDatabase func(*sql.DB)

//FnInitDatabase is a function used to init the database after creating it
type FnInitDatabase func(*DbContext)

// Init the db
func Init(dbPath string, salt string, fnCreate FnCreateDatabase, fnInit FnInitDatabase) *DbContext {
	isFirstTime := false
	if _, err := os.Stat(dbPath); err != nil {
		isFirstTime = true
	}

	db, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	if isFirstTime {
		fnCreate(db)
	}

	context := DbContext{Db: db, salt: salt}

	if isFirstTime {
		fnInit(&context)
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
