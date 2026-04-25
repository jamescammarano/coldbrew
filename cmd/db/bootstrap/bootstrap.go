package bootstrap

import (
	"os"

	"coldbrew.go/cb/common/database"

	"coldbrew.go/cb/common/types"
	"coldbrew.go/cb/common/utils"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var Cmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "bootstrap database",
	Long: `Initializes the database with all the current migrations. 
	Should only need to be run on installation or if the DB is corrupted.`,
	Run: func(cmd *cobra.Command, args []string) {
		var config types.Config

		// TODO Don't hardcode
		os.MkdirAll("./InstallDir", 0755)

		os.Create(database.GetDatabaseURL())

		db := database.OpenDatabase()
		defer db.Close()

		app := database.Column{Name: "app", Type: "VARCHAR(64)", Nullable: false, Unique: true}
		port := database.Column{Name: "port", Type: "INT", Nullable: true, Unique: false}

		err := database.CreateTable(db, []database.Column{app, port}, "Port")

		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		err = database.CreateAppTable(db, "Global")

		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		// TODO remove the .example.
		res, err := os.ReadFile("./cb.example.yml")

		if err != nil {
			logrus.Error(err)
		}

		err = yaml.Unmarshal(res, &config)

		if err != nil {
			logrus.Error(err)
		}

		errs := database.Insert(db, utils.GenerateVariables(config.Variables), "Global")

		for _, err := range errs {
			logrus.Error(err)
		}
	},
}
