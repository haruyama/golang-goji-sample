package models

import (
	"database/sql"
	"log"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/coopernurse/gorp"
	"github.com/golang/glog"
)

type User struct {
	Id       int64 `db:"user_id"`
	Email    string
	Username string
	Password []byte
}

func (user *User) HashPassword(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		glog.Fatalf("Couldn't hash password: %v", err)
		panic(err)
	}
	user.Password = hash
}

func GetUserByEmail(dbMap *gorp.DbMap, email string) (user *User) {
	err := dbMap.SelectOne(&user, "SELECT * FROM users where email = ?", email)

	if err != nil {
		glog.Warningf("Can't get user by email: %v", err)
	}
	return
}

func InsertUser(dbMap *gorp.DbMap, user *User) error {
	return dbMap.Insert(&user)
}

func GetDbMap() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	//TODO: Get user, password and database from config.
	db, err := sql.Open("mysql", "root:password@/goji")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8MB4"}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
