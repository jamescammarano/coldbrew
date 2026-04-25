package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type Column struct {
	Name     string
	Type     string
	Nullable bool
	Unique   bool
}

func GetDatabaseURL() string {
	// TODO DONT HARDCODE BC THINGS MOVE
	return "./InstallDir/data.db"
}

func OpenDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", GetDatabaseURL())

	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	return db
}

// TODO make this more generic maybe
func Insert(db *sql.DB, rows map[string]string, roast string) []error {
	var errs []error

	for attr, data := range rows {
		if attr == "Port" || attr == "port" {
			_, err := db.Exec("INSERT INTO Port (app, port) VALUES ($2, $3)", roast, data)

			if err != nil {
				errs = append(errs, errors.New(strings.ReplaceAll(err.Error(), "attr", attr)))
			}
		} else {
			_, err := db.Exec(fmt.Sprintf("INSERT INTO %v (attr, data) VALUES ($2, $3)", roast), attr, data)

			if err != nil {
				errs = append(errs, errors.New(strings.ReplaceAll(err.Error(), "attr", attr)))
			}
		}
	}

	return errs
}

func CreateAppTable(db *sql.DB, roast string) error {
	attr := Column{Name: "attr", Type: "VARCHAR(64)", Nullable: false, Unique: true}
	data := Column{Name: "data", Type: "VARCHAR(64)", Nullable: true, Unique: false}

	err := CreateTable(db, []Column{attr, data}, roast)

	return err
}

func CreateTable(db *sql.DB, columns []Column, title string) error {
	logrus.Info("Creating Table: ", title)
	var columnDef []string

	for _, column := range columns {
		var null string
		var unique string

		if !strings.HasPrefix(column.Type, "VARCHAR(") && column.Type != "INT" {
			return fmt.Errorf("invalid column type %v for %v", column.Type, column.Name)
		}

		if !column.Nullable {
			null = "NOT NULL"
		}

		if column.Unique {
			null = "UNIQUE"
		}

		columnDef = append(columnDef, fmt.Sprintf("`%v` %v %v %v", column.Name, column.Type, null, unique))
	}

	_, err := db.Exec(fmt.Sprintf("CREATE TABLE `%v` (%v)", title, strings.Join(columnDef, ", ")))

	return err
}

func ReadRows(rows *sql.Rows) (map[string]string, error) {
	result := map[string]string{}
	var (
		attr string
		data string
	)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&attr, &data)
		if err != nil {
			return result, err
		}
		result[attr] = data
	}
	return result, nil
}
