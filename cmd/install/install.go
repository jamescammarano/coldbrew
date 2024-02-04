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
