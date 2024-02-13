package create

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func beforeEach() (*sql.DB, error) {
	os.MkdirAll("./test", 0755)

	os.Create("./test/data.db")

	return sql.Open("sqlite3", "./test/data.db")
}

func afterEach() {
	os.RemoveAll("./test")
}

func TestCreateTable(t *testing.T) {
	db, err := beforeEach()

	if err != nil {
		fmt.Printf("unexpected error: %s", err)
	}

	defer db.Close()

	err = CreateTable(db, "test")

	if err != nil {
		fmt.Printf("unexpected error: %s", err)
		t.Errorf("error while creating table `test`: %s", err)
	}

	err = CreateTable(db, "sqlite_test")

	if err != nil {
		fmt.Printf("expected error: %s", err)
	}

	afterEach()
}

func TestSeed(t *testing.T) {
	db, err := beforeEach()

	if err != nil {
		fmt.Printf("unexpected error: %s", err)
	}

	variables := map[string]string{"test": "done", "test2": "doneagain"}

	defer db.Close()

	CreateTable(db, "test")

	err = Seed(db, variables, "test")

	if err != nil {
		fmt.Printf("unexpected error: %s", err)
		t.Errorf("error while creating table `test`: %s", err)
	}

	afterEach()
}
