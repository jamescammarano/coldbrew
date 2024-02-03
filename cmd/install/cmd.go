package install

import (
	"main/coldbrew/cmd/install/kavita"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "install",
	Short: "install new roasts",
	Long:  "install new roasts",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	opts := map[string]string{"tag": "latest", "restart": "unless_stopped", "port": "5000"}

	// 	tag_select := promptui.Select{
	// 		Label: "Select Tag",
	// 		// readFile tags for supported tags
	// 		Items: []string{"latest", "nightly"},
	// 	}

	// 	_, tag, err := tag_select.Run()

	// 	if err != nil {
	// 		log.Default()
	// 	}
	// 	//  TODO Error Handling for tag select
	// 	opts["tag"] = tag

	// 	healthchecks := promptui.Select{
	// 		Label: "Include healthchecks?",
	// 		// readFile tags for supported tags
	// 		Items: []string{"Yes", "No"},
	// 	}

	// 	_, result, err := healthchecks.Run()
	// 	// TODO err handling for healthchecks

	// 	opts["healthchecks"] = result
	// 	logrus.Info(opts)
	// 	// b, err := os.ReadFile(fmt.Sprintf(`roasts/kavita/docker-compose.yml`))
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// err := Cmd.execute()
	// if err != nil {
	// 	os.Exit(1)
	// }
}

func init() {
	Cmd.AddCommand(kavita.Cmd)
}

// install software:
// cb install [roast]
// check for install status => "run cb reinstall"
// find roast in /roasts OR just dl gitrepo/roasts/[appname] 4 config
//    what goes in a roast config?
// 		- base config
//			- dockerimage name + version
// 			- supported tags
// 			- extra files
// 		- Docker compose w/ template vars for secrets
// 		- config.go
// 			- any install opts (kavita-email w/ kavita etc)
// subs template vars
// prompts from config.go
// downloads docker image
// runs docker
