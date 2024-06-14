package reinstall

import (
	"fmt"
	"os"

	"coldbrew.go/cb/cmd/install"
	"coldbrew.go/cb/common/database"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "reinstall",
	Short: "usage: cb reinstall roast",
	Long:  "reinstall roasts",
	Run: func(cmd *cobra.Command, args []string) {
		reinstaller(args[0])
	},
}

func reinstaller(roast string) {

	// TODO BIG WARNING OF DATA LOSS
	// TODO don't hardcode this
	db := database.OpenDatabase()

	defer db.Close()

	logrus.Info(fmt.Sprintf("DROP TABLE `%v`", roast))
	res, err := db.Query(fmt.Sprintf("DROP TABLE `%v`", roast))

	if err != nil {
		logrus.Error(err)
	}
	_, err = database.ReadRows(res)

	if err != nil {
		logrus.Error(err)
	}
	// TODO Delete docker image/container

	// TODO Delete roast from Port
	sqlrows, err := db.Query("SELECT * FROM Global WHERE attr = 'InstallDir'")

	if err != nil {
		logrus.Error(err)
	}

	// TODO rename that
	variable, err := database.ReadRows(sqlrows)

	if err != nil {
		logrus.Error(err)
	}

	err = os.RemoveAll(fmt.Sprintf("%v/%v", variable["InstallDir"], roast))

	if err != nil {
		logrus.Error(err)
	}

	install.Installer(roast)

}
