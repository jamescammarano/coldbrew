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

	err = createTable(db, "test")

	if err != nil {
		fmt.Printf("unexpected error: %s", err)
		t.Errorf("error while creating table `test`: %s", err)
	}

	err = createTable(db, "sqlite_test")

	if err != nil {
		fmt.Printf("expected error: %s", err)
	}

	afterEach()
}

func TestSeed(t *testing.T) {
	db, err := beforeEach()

	variables := map[string]string{"test": "done", "test2": "doneagain"}

	if err != nil {
		fmt.Printf("unexpected error: %s", err)
	}

	defer db.Close()

	createTable(db, "test")

	err = seed(db, variables, "test")

	if err != nil {
		fmt.Printf("unexpected error: %s", err)
		t.Errorf("error while creating table `test`: %s", err)
	}

	afterEach()
}
