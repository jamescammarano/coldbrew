package install

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"coldbrew.go/cb/common/database"
	"coldbrew.go/cb/common/types"
	"coldbrew.go/cb/common/utils"

	"github.com/cbroglie/mustache"
	_ "github.com/mattn/go-sqlite3"
	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
	"github.com/spaceweasel/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var Cmd = &cobra.Command{
	Use:   "install",
	Short: "usage: cb install roast",
	Long:  "install new roasts",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO make a flag "--skip-prerun" for testing?
		Installer(args[0])
	},
}

func Installer(roast string) {
	var config = types.Config{}

	db := database.OpenDatabase()

	defer db.Close()
	prerun(roast, db)

	b, err := os.ReadFile(fmt.Sprintf("./roasts/%s/config.yml", roast))

	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	err = yaml.Unmarshal([]byte(b), &config)

	if err != nil {
		logrus.Error(err)
	}

	err = database.CreateAppTable(db, roast)

	if err != nil {
		logrus.Error(err)
	}

	// TODO should use the variables loaded from the db and only fill in the missing vars
	// (this should never actually happen though)
	variables := utils.GenerateVariables(config.Variables)

	if variables["Restart"] == "" {
		variables["Restart"] = "unless-stopped"
	}

	rows, err := db.Query("SELECT attr, data FROM `Global`")

	if err != nil {
		logrus.Error(err)
	}

	res, err := database.ReadRows(rows)

	if err != nil {
		logrus.Error(err)
	}

	// This prefers the global var to anything the roast requires for its own needs where names clash
	for attr, data := range res {
		variables[attr] = data
	}

	tag, err := setTag(config.Tags)

	for err != nil {
		logrus.Error(err)
		tag, err = setTag(config.Tags)
	}

	variables["Tag"] = tag

	errs := database.Insert(db, variables, roast)

	for _, err := range errs {
		logrus.Error(err)
	}

	// TODO Subvars in files before mkdir, I would just mv it to line 144 but I don't want to overwrite the
	// subbed files. no "onFileExists" opt on that copy package
	errs = mkdirs(config.Directories, variables)

	for _, err := range errs {
		logrus.Error(err)
	}
	// TODO test 0600 perms
	postInstaller := types.File{Name: "post-install.sh", Dest: fmt.Sprintf("%s/%s", variables["InstallDir"], roast), Src: fmt.Sprintf("./roasts/%v/", roast), Perms: 0777}
	_, fileExistErr := os.Stat(postInstaller.Src)

	if fileExistErr == nil {
		config.Files = append(config.Files, postInstaller)
	}

	for _, file := range config.Files {
		fileFormatter(file, variables)
	}

	if fileExistErr == nil {
		runScript(fmt.Sprintf("%s/%s", postInstaller.Dest, postInstaller.Name))
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

func mkdirs(dirs []types.Directory, variables map[string]string) []error {
	errors := []error{}
	for _, dir := range dirs {
		// TODO default file permission
		// TODO I don't need to sub as I copy because fileManager will sub as it copies
		// WILL create redundancy in copying files multiple times so not ideal

		// replaces the variables in the filepath so it knows what to make where
		path, err := mustache.Render(dir.Dest, variables)

		if err != nil {
			errors = append(errors, err)
			logrus.Error(err)
		}

		// makes the path string
		dest := fmt.Sprintf("%s/%s", path, dir.Name)

		// I need to add an "if src attr exist: copy"
		if dir.Src != "" {
			err = copy.Copy(dir.Src, path)

			if err != nil {
				errors = append(errors, err)
			}

		} else {
			err = os.MkdirAll(dest, fs.FileMode(dir.Perms))

			if err != nil {
				errors = append(errors, err)
			}
		}
	}
	return errors
}

func fileFormatter(file types.File, variables map[string]string) {
	// make the source filepath (var subs)
	src, err := mustache.Render(file.Src, variables)

	if err != nil {
		logrus.Error(err)
	}

	// make the dest filepath (var subs)
	dest, err := mustache.Render(file.Dest, variables)
	if err != nil {
		logrus.Error("While subbing the file destination path: ", err)
	}

	// TODO not 777
	err = os.MkdirAll(dest, 0777)

	if err != nil {
		logrus.Error(err)
	}

	str, err := mustache.RenderFile(fmt.Sprintf(`%v/%v`, src, file.Name), variables)

	if err != nil {
		logrus.Error(err)
	}

	// TODO openable by bin and by docker w/ perms 600. currently mkdirall makes a directory that isn't writable
	err = os.WriteFile(fmt.Sprintf(`%v/%v`, dest, file.Name), []byte(str), fs.FileMode(file.Perms))

	if err != nil {
		logrus.Error(err)
	}
}

func runScript(path string) (string, error) {
	args := strings.Split("-c "+path, " ")
	r, w, _ := os.Pipe()
	os.Stdout = w
	postinstaller := exec.Command("/bin/bash", args...)
	postinstaller.Stdout = os.Stdout
	postinstaller.Stderr = os.Stderr

	err := postinstaller.Run()

	outC := make(chan string)
	// TODO MAKE THIS MY OWN
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	out := <-outC

	return out, err
}
