package create

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"coldbrew.go/cb/common"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Variables map[string]common.Variable `yaml:"vars"`
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

		// TODO don't hardcode this
		db, err := sql.Open("sqlite3", "./InstallDir/data.db")

		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		defer db.Close()

		err = CreatePortsTable(db)

		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		err = CreateTable(db, "Global")

		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		res, err := os.ReadFile("./cb.example.yml")

		if err != nil {
			logrus.Error(err)
		}

		err = yaml.Unmarshal(res, &config)

		if err != nil {
			logrus.Error(err)
		}

		errs := Seed(db, common.GenerateVariables(config.Variables), "Global")

		for _, err := range errs {
			logrus.Error(err)
		}
	},
}

func CreateTable(db *sql.DB, roast string) error {
	logrus.Info("Creating Table: ", roast)

	statement := fmt.Sprintf("CREATE TABLE `%v` (`attr` VARCHAR(64) NOT NULL UNIQUE, `data` VARCHAR(255) NOT NULL)", roast)
	_, err := db.Exec(statement)

	return err
}

func CreatePortsTable(db *sql.DB) error {
	logrus.Info("Creating Table: Port")

	_, err := db.Exec("CREATE TABLE `Port` (`app` VARCHAR(64) NOT NULL UNIQUE, `port` INT NOT NULL UNIQUE)")

	return err
}

func Seed(db *sql.DB, vars map[string]string, roast string) []error {
	insertQuery := fmt.Sprintf("INSERT INTO %v (attr, data) VALUES ($2, $3)", roast)
	var errs []error
	var err error

	for attr, data := range vars {
		if attr == "Port" || attr == "port" {
			_, err = db.Exec("INSERT INTO Port (app, port) VALUES ($2, $3)", roast, data)
		} else {
			_, err = db.Exec(insertQuery, attr, data)
		}
		if err != nil {
			errs = append(errs, errors.New(strings.ReplaceAll(err.Error(), "attr", attr)))
		}
	}

	return errs
}
