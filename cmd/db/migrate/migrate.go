package migrate

import (
	"database/sql"
	"os"

	"github.com/sirupsen/logrus"
)

func Run(path string) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		logrus.Error(err.Error())
		db.Close()
		os.Exit(1)
	}

	m20240204a(db)

	if err != nil {
		logrus.Error(`Migration 20240204a:`, err.Error())
		// r20240204(db)

	} else {
		logrus.Info(`SUCCESSFUL Migration 20240204a `)

	}

	db.Close()
}

func m20240204a(db *sql.DB) (sql.Result, error) {
	logrus.Info("ATTEMPTING 20240204a: Create Table `GlobalVariables`")
	return db.Exec("CREATE TABLE `GlobalVariables` (`attribute` VARCHAR(64) NOT NULL UNIQUE, `data` VARCHAR(255) NOT NULL)")
}

// Rollback function
// func r20240204(db *sql.DB) (sql.Result, error) {
// 	logrus.Error("ROLLINGBACK 20240204a: Create Table `GlobalVariables`")

// 	return db.Exec("CREATE TABLE `GlobalVariables` (`attribute` VARCHAR(64) NOT NULL UNIQUE, `data` VARCHAR(255) NOT NULL)")
// }
