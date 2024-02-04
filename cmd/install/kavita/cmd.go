package kavita

import (
	"fmt"
	"os"

	"github.com/cbroglie/mustache"
	"github.com/sirupsen/logrus"
	"github.com/spaceweasel/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Handlebars struct {
	DownloadDir string `yaml:"download_dir,omitempty"`
	MediaDir    string `yaml:"media_dir,omitempty"`
	InstallDir  string `yaml:"install_dir,omitempty"`
	Domain      string `yaml:"domain,omitempty"`
	KavitaToken string `yaml:"kavita_token,omitempty"`
	Port        int    `yaml:"kavita_port,omitempty"`
	Tag         string
	Restart     string
}

type File struct {
	Name string `yaml:"name"`
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

type Addons []struct {
	Name          string   `yaml:"name"`
	ContainerName string   `yaml:"container_name"`
	Tags          []string `yaml:"tags"`
	Image         string   `yaml:"image"`
	Volumes       []string `yaml:"volumes"`
}

type KavitaConfig struct {
	Tags    []string `yaml:"tags"`
	Vars    []string `yaml:"vars,omitempty"`
	Files   []File   `yaml:"files,omitempty"`
	Restart string   `yaml:"restart,omitempty"`
	Addons  Addons   `yaml:"addons,omitempty"`
	Port    int      `yaml:"kavita_port"`
	Token   int      `yaml:"kavita_token,omitempty"`
}

var Cmd = &cobra.Command{
	Use:   "kavita",
	Short: "install kavita",
	Long:  "install kavita",
	Run: func(cmd *cobra.Command, args []string) {
		var config = KavitaConfig{}

		b, err := os.ReadFile(`./cmd/install/kavita/config.yml`)

		if err != nil {
			logrus.Error(err.Error())
			os.Exit(1)
		}

		err = yaml.Unmarshal([]byte(b), &config)

		if err != nil {
			logrus.Error(err.Error())
			os.Exit(1)
		}

		substitutions, err := loadSubstitutions()

		if err != nil {
			logrus.Error(err.Error())
			os.Exit(1)
		}

		substitutions, err = setTag(config, substitutions)

		if err != nil {
			logrus.Error(err.Error())
			os.Exit(1)
		}

		fileManager(config.Files, substitutions)

		selectAddons(config.Addons)
	},
}

func loadSubstitutions() (Handlebars, error) {
	var substitions Handlebars

	b, err := os.ReadFile("cb.yml")

	if err != nil {
		return substitions, err
	}

	err = yaml.Unmarshal([]byte(b), &substitions)

	return substitions, err
}

func setTag(t KavitaConfig, variables Handlebars) (Handlebars, error) {

	tagSelect := promptui.Select{
		Label: "Select tag",
		Items: t.Tags,
	}

	_, tag, err := tagSelect.Run()

	if err != nil {
		logrus.Error(err.Error())
		return Handlebars{}, err
	}

	variables.Tag = tag
	variables.Port = 5000

	return variables, nil
}

func fileManager(files []File, variables Handlebars) {
	for _, file := range files {

		src, err := mustache.Render(file.Src, variables)
		if err != nil {
			logrus.Error(err.Error())
		}

		dest, err := mustache.Render(file.Dest, variables)
		if err != nil {
			logrus.Error(err.Error())
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

func selectAddons(t Addons) (Addons, error) {
	var addons Addons

	addonsTemplates := &promptui.MultiSelectTemplates{
		Unselected: " {{ .Name }}",
		Selected:   " {{ .Name | green }}",
	}

	addonsSelect := promptui.MultiSelect{
		Label:     "Select addons to include",
		Items:     t,
		Templates: addonsTemplates,
	}

	addonsIndex, err := addonsSelect.Run()

	if err != nil {
		return nil, err
	}

	for _, num := range addonsIndex {
		addons = append(addons, t[num])
	}

	return addons, nil
}
