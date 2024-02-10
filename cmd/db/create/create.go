package create

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type GlobalConfig struct {
	InstallDir  string `yaml:"InstallDir"`
	MediaDir    string `yaml:"MediaDir"`
	DownloadDir string `yaml:"DownloadDir"`
}

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "create database",
	Long: `Initializes the database with all the current migrations. 
	Should only need to be run on installation or if the DB is corrupted.`,
	Run: func(cmd *cobra.Command, args []string) {
		os.MkdirAll("./InstallDir", 0755)

		os.Create("./InstallDir/data.db")

		db, err := sql.Open("sqlite3", "./InstallDir/data.db")

		if err != nil {
			logrus.Error(err.Error())
			db.Close()
			os.Exit(1)
		}

		err = createTable(db, "Global")

		if err != nil {
			logrus.Error(err.Error())
			db.Close()
			os.Exit(1)
		}

		err = seed(db, "./config.yml", "Global")
		if err != nil {
			logrus.Error(err.Error())
		}
		db.Close()

	},
}

func createTable(db *sql.DB, roast string) error {
	logrus.Info("Creating Table: ", roast)

	statement := fmt.Sprintf("CREATE TABLE `%v` (`attr` VARCHAR(64) NOT NULL UNIQUE, `data` VARCHAR(255) NOT NULL)", roast)
	_, err := db.Exec(statement)

	return err
}

func seed(db *sql.DB, config string, roast string) error {
	var variables map[string]string
	insertQuery := fmt.Sprintf("INSERT INTO %v (attr, data) VALUES ($2, $3)", roast)

	res, err := os.ReadFile(config)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(res, &variables)

	if err != nil {
		return err
	}

	for attr, data := range variables {
		_, err = db.Exec(insertQuery, attr, data)

		if err != nil {
			return err
		}
	}

	return err
}
