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
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os"
	"visualizedGit/lib"
)

// rmcfgCmd represents the rmcfg command
var rmcfgCmd = &cobra.Command{
	Use:   "rmcfg",
	Short: "Remove the existing configuration file",
	Long:  `Remove the existing configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		rmcfg()
		return
	},
}

func init() {
	rootCmd.AddCommand(rmcfgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmcfgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmcfgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func rmcfg() {
	filePath := lib.GetDotFilePath()
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	color.Green("Successfully removed configuration file\n")
}
