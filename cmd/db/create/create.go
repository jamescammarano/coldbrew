package create

import (
	"cmd/cb/cmd/db/migrate"
	"database/sql"
	"fmt"
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
		create("./InstallDir", "data.db")
		migrate.Run("./InstallDir/data.db")
		seed("./InstallDir/data.db", "./config.yml")
	},
}

func create(path string, name string) {
	os.MkdirAll(path, 0755)
	dbPath := fmt.Sprintf("%v/%v", path, name)

	os.Create(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	errHandler(err, db)
	db.Close()
}

func seed(path string, config string) error {
	var globalConfig GlobalConfig
	insertQuery := "INSERT INTO `GlobalVariables` (attribute, data) VALUES ($1, $2)"

	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(config), &globalConfig)

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
