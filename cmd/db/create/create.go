package create

import (
	"database/sql"
	"fmt"
	"main/coldbrew/cmd/db/migrate"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Variables []map[string]string `yaml:"vars"`
}

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "create database",
	Long: `Initializes the database with all the current migrations. 
	Should only need to be run on installation or if the DB is corrupted.`,
	Run: func(cmd *cobra.Command, args []string) {
		var config Config

		os.MkdirAll("./InstallDir", 0755)

		os.Create("./InstallDir/data.db")

		db, err := sql.Open("sqlite3", "./InstallDir/data.db")

		if err != nil {
			logrus.Error(err.Error())
			os.Exit(1)
		}
		defer db.Close()

		err = createTable(db, "Global")

		if err != nil {
			logrus.Error(err.Error())
			os.Exit(1)
		}

		res, err := os.ReadFile("./config.yml")

		if err != nil {
			logrus.Error(err.Error())
		}

		err = yaml.Unmarshal(res, &config)

		if err != nil {
			logrus.Error(err.Error())
		}

		err = seed(db, mergeMaps(config.Variables), "Global")
		if err != nil {
			logrus.Error(err.Error())
		}

	},
}

func createTable(db *sql.DB, roast string) error {
	logrus.Info("Creating Table: ", roast)

	statement := fmt.Sprintf("CREATE TABLE `%v` (`attr` VARCHAR(64) NOT NULL UNIQUE, `data` VARCHAR(255) NOT NULL)", roast)
	_, err := db.Exec(statement)

	return err
}

func seed(db *sql.DB, vars map[string]string, roast string) error {
	insertQuery := fmt.Sprintf("INSERT INTO %v (attr, data) VALUES ($2, $3)", roast)

	for attr, data := range vars {
		_, err := db.Exec(insertQuery, attr, data)

		if err != nil {
			return err
		}
	}

	return nil
}

func mergeMaps(mapArr []map[string]string) map[string]string {
	merged := map[string]string{}

	for _, t := range mapArr {
		for k, v := range t {
			merged[k] = v
		}
	}

	return merged
}
