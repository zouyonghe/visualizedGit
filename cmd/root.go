/*
Copyright © 2022 Yonghe Zou 1259085392z@gmail.com

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	email string
)

var visualizedGitVersion = "0.0.2"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "visualizedGit",
	Short: "Visualize local git contributions.",
	Long: `VisualizedGit is a CLI tool for developers to visualize their git contributions.
Developers can specify the git repository and view the visualized local git contributions.`,
	Version: visualizedGitVersion,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	/*	Run: func(cmd *cobra.Command, args []string) {
		if folder == "" && email == "" {
			//PrintVersion()
			color.Green("visualizedGit version %s", visualizedGitVersion)
			fmt.Println("Visualize local git contributions.")
			fmt.Printf("Using \"visualizedGit --help\" or \"visualizedGit -h\" for more information.\n")
			return
		}

		if folder != "" {
			scan.Scan(folder)
			return
		}
		stats.Stats(email)
	},*/
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.visualizedGit.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
