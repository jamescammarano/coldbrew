package database

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func TestCreateTable(t *testing.T) {

	column1 := Column{Name: "attr", Type: "VARCHAR(64)", Nullable: false, Unique: true}
	column2 := Column{Name: "data", Type: "VARCHAR(64)", Nullable: true, Unique: false}
	badcolumn := Column{Name: "bad_column", Type: "oops(64)", Nullable: true, Unique: false}

	os.MkdirAll("./TestingDir", 0755)

	os.Create("./TestingDir/data.db")

	db, err := sql.Open("sqlite3", "./TestingDir/data.db")

	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	err = CreateTable(db, []Column{column1, column2}, "testing1")

	if err != nil {
		t.Error(err)
	}

	// Should return an error. Bad path
	err = CreateTable(db, []Column{badcolumn}, "Bad Column")

	if err != nil {
		logrus.Error(err)
	}

	os.RemoveAll("./TestingDir")
}
