package create

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "create database",
	Long:  `create database`,
	Run: func(cmd *cobra.Command, args []string) {
		os.MkdirAll("./InstallDir/", 0755)
		os.Create("./InstallDir/data.db")

		db, err := sql.Open("sqlite3", "./InstallDir/data.db")

		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}
		_, err = db.Exec("CREATE TABLE `global` (`attribute` VARCHAR(64) NULL, `value` VARCHAR(255) NOT NULL, `created_on` DATETIME NULL, `app` VARCHAR(1) NULL)")

		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}

		db.Close()
	},
}
