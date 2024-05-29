package install

import (
	"bufio"
	"database/sql"
	"fmt"
	"main/coldbrew/cmd/db/create"
	"main/coldbrew/utils"
	"os"
	"os/exec"

	"github.com/cbroglie/mustache"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spaceweasel/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type File struct {
	Name string `yaml:"name"`
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

type Directory struct {
	Name string `yaml:"name"`
	Dest string `yaml:"dest"`
}

type Addons []struct {
	Roast string `yaml:"roast"`
}

type Config struct {
	Tags        []string            `yaml:"tags"`
	Variables   []map[string]string `yaml:"vars"`
	Directories []Directory         `yaml:"dirs,omitempty"`
	Files       []File              `yaml:"files,omitempty"`
	Restart     string              `yaml:"restart,omitempty"`
	Addons      Addons              `yaml:"addons,omitempty"`
}

var Cmd = &cobra.Command{
	Use:   "install",
	Short: "usage: cb install roast",
	Long:  "install new roasts",
	Run: func(cmd *cobra.Command, args []string) {
		installer(args[0])
	},
}

func installer(roast string) {
	var config = Config{}

	b, err := os.ReadFile(fmt.Sprintf("./roasts/%s/config.yml", roast))

	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}

	err = yaml.Unmarshal([]byte(b), &config)

	if err != nil {
		logrus.Error(err.Error())
	}

	db, err := sql.Open("sqlite3", "./InstallDir/data.db")

	if err != nil {
		logrus.Error(err.Error())
	}

	defer db.Close()

	err = create.CreateTable(db, roast)

	if err != nil {
		logrus.Error(err.Error())
	}

	variables := generateVariables(utils.MergeMaps(config.Variables))

	if variables["Restart"] == "" {
		variables["Restart"] = "unless-stopped"
	}

	err = create.Seed(db, variables, roast)

	if err != nil {
		logrus.Error(err.Error())
	}

	rows, err := db.Query("SELECT attr, data FROM `Global`")

	if err != nil {
		logrus.Error(err.Error())
	}

	res, err := readRows(rows)

	if err != nil {
		logrus.Error(err.Error())
	}

	for attr, data := range res {
		variables[attr] = data

	}

	tag, err := setTag(config.Tags)

	for err != nil {
		logrus.Error(err.Error())
		tag, err = setTag(config.Tags)
	}

	variables["tag"] = tag

	createDirs(config.Directories, variables)

	fileManager(config.Files, variables, roast)

	_, err = os.Stat(fmt.Sprintf("%s/%s/config.sh", variables["InstallDir"], roast))

	if err == nil {
		configScript := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s/%s/config.sh", variables["InstallDir"], roast))

		stdout, err := configScript.StderrPipe()

		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		err = configScript.Start()

		scanner := bufio.NewScanner(stdout)

		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}

		configScript.Wait()

		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	}
}

func setTag(tags []string) (string, error) {

	tagSelect := promptui.Select{
		Label: "Select tag",
		Items: tags,
	}

	_, tag, err := tagSelect.Run()

	return tag, err
}

func readRows(rows *sql.Rows) (map[string]string, error) {
	result := map[string]string{}
	var (
		attr string
		data string
	)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&attr, &data)
		if err != nil {
			return result, err
		}
		result[attr] = data
	}
	return result, nil
}

func createDirs(dirs []Directory, variables map[string]string) []error {
	errors := []error{}
	for _, dir := range dirs {
		sub, err := mustache.Render(dir.Dest, variables)

		if err != nil {
			errors = append(errors, err)
			logrus.Error(err.Error())
		}

		dest := fmt.Sprintf("%s/%s", sub, dir.Name)

		err = os.MkdirAll(dest, 0777)

		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func fileManager(files []File, variables map[string]string, roast string) {
	_, err := os.Stat(fmt.Sprintf("%s/%s/config.sh", variables["InstallDir"], roast))

	if err != nil {
		configScript := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s/%s/config.sh", variables["InstallDir"], roast))

		err = configScript.Run()

		if err != nil {
			logrus.Error(err.Error())
		}

		str, err := mustache.RenderFile(fmt.Sprintf("./roasts/%v/config.sh", roast), variables)

		if err != nil {
			logrus.Error(err.Error())
		}

		err = os.MkdirAll(fmt.Sprintf("%s/%s", variables["InstallDir"], roast), 0777)

		if err != nil {
			logrus.Error(err.Error())
		}

		err = os.WriteFile(fmt.Sprintf("%s/%s/config.sh", variables["InstallDir"], roast), []byte(str), 0777)

		if err != nil {
			logrus.Error(err.Error())
		}
	}

	for _, file := range files {

		src, err := mustache.Render(file.Src, variables)
		if err != nil {
			logrus.Error(err.Error())
		}

		dest, err := mustache.Render(file.Dest, variables)
		if err != nil {
			logrus.Error("mustache fileDest: ", err.Error())
		}

		err = os.MkdirAll(dest, 0777)

		if err != nil {
			logrus.Error(err.Error())
		}

		str, err := mustache.RenderFile(fmt.Sprintf(`%v/%v`, src, file.Name), variables)

		if err != nil {
			logrus.Error(err.Error())
		}

		err = os.WriteFile(fmt.Sprintf(`%v/%v`, dest, file.Name), []byte(str), 0777)

		if err != nil {
			logrus.Error(err.Error())
		}
	}
}

func generateVariables(variables map[string]string) map[string]string {
	generated := map[string]string{}

	for attr, val := range variables {
		if val == "__func(PromptInput)" {
			generated[attr] = utils.PromptInput(attr)
		} else {
			generated[attr] = val
		}
	}
	return generated
}
