package install

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func prerun(roast string, db *sql.DB) {
	// call db?
	// TODO check for updates to cb

	statement := fmt.Sprintf("SELECT * FROM %v;", roast)
	table, err := db.Exec(statement)

	if err == nil && table != nil {
		logrus.Error(fmt.Sprintf("roast is already installed. Reinstall it with 'cb reinstall %v'", roast))
		os.Exit(1)
	}
}

// TODO LP prerun scrape all the utils functions and then dynamically call them?
// make hash of all of the utils and store in global mem or something?
func UtilCheck() {
	// grep -re '^func'
	// side effect time?
	// map[string]func or something idk
}
