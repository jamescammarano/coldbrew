package migrate

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

func Run(db *sql.DB) error {
	err := createTableGlobalVariables(db)
	return err
}

func createTableGlobalVariables(db *sql.DB) error {
	logrus.Info("ATTEMPTING 20240204a: Create Table `GlobalVariables`")
	_, err := db.Exec("CREATE TABLE `GlobalVariables` (`attribute` VARCHAR(64) NOT NULL UNIQUE, `data` VARCHAR(255) NOT NULL)")

	if err != nil {
		logrus.Error(`Migration 20240204a:`, err.Error())
		// r20240204(db)
	} else {
		logrus.Info(`SUCCESSFUL Migration 20240204a `)
	}
	return err
}
