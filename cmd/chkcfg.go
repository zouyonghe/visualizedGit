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
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"visualizedGit/lib"
)

// chkcfgCmd represents the chkcfg command
var chkcfgCmd = &cobra.Command{
	Use:   "chkcfg",
	Short: "Check the tracking repositories in configuration file",
	Long:  `Check the tracking repositories in configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		chkcfg()
	},
}

func init() {
	rootCmd.AddCommand(chkcfgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chkcfgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chkcfgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func chkcfg() {
	filePath := lib.GetDotFilePath()
	fmt.Println("Checking configuration: ", filePath)
	fmt.Println()
	content := lib.ParseFileLinesToSlice(filePath)
	if len(content) < 1 {
		color.Red("No repositories found.\n")
	} else {
		for _, val := range content {
			fmt.Println(val)
		}
	}
	fmt.Println()
	fmt.Println("Finished checking configuration.")
}
