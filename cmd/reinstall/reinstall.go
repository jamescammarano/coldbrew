package install

import (
	"database/sql"

	"coldbrew.go/cb/cmd/install"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "reinstall",
	Short: "usage: cb reinstall roast",
	Long:  "reinstall new roasts",
	Run: func(cmd *cobra.Command, args []string) {
		reinstaller(args[0])
	},
}

func reinstaller(roast string) {

	// TODO BIG WARNING OF DATA LOSS
	// TODO don't hardcode this
	db, err := sql.Open("sqlite3", "./InstallDir/data.db")

	if err != nil {
		logrus.Error(err)
	}

	defer db.Close()

	_, err = db.Query("DROP TABLE $1 FROM `Global`", roast)

	if err != nil {
		logrus.Error(err)
	}

	// TODO Expand this to remove files too

	install.Installer(roast)

}
