package system

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func InitMysql() (db *sql.DB, err error) {

	database := "issues"
	username := "root"
	password := "password"

	db, err = sql.Open("mysql", username + ":" + password + "@/" + database)

	if err != nil {
		panic(err.Error())
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db, err

}
