package create

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type GlobalConfig struct {
	InstallDir  string `yaml:"InstallDir"`
	MediaDir    string `yaml:"MediaDir"`
	DownloadDir string `yaml:"DownloadDir"`
}

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "create database",
	Long:  `create database`,
	Run: func(cmd *cobra.Command, args []string) {
		os.MkdirAll("./InstallDir/", 0755)
		os.Create("./InstallDir/data.db")

		db, err := sql.Open("sqlite3", "./InstallDir/data.db")
		errHandler(err, db)

		_, err = db.Exec("CREATE TABLE `GlobalVariables` (`attribute` VARCHAR(64) NOT NULL UNIQUE, `data` VARCHAR(255) NOT NULL)")
		errHandler(err, db)

		err = seedDb(db)
		errHandler(err, db)

		db.Close()
	},
}

func seedDb(db *sql.DB) error {
	var globalConfig GlobalConfig
	insertQuery := "INSERT INTO `GlobalVariables` (attribute, data) VALUES ($1, $2)"

	b, err := os.ReadFile("./config.yml")

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &globalConfig)

	if err != nil {
		return err
	}

	_, err = db.Exec(insertQuery, "DownloadDir", globalConfig.DownloadDir)
	if err != nil {
		return err
	}

	_, err = db.Query(insertQuery, "MediaDir", globalConfig.MediaDir)
	if err != nil {
		return err
	}

	_, err = db.Query(insertQuery, "InstallDir", globalConfig.InstallDir)

	return err
}

func errHandler(err error, db *sql.DB) {
	if err != nil {
		logrus.Error(err.Error())
		db.Close()
		os.Exit(1)
	}
}
