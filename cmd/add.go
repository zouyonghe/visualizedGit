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

var path string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [-p path|--path path]",
	Short: "Add a new folder to scan for Git repositories",
	Long:  `Add a new folder to scan for Git repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		if path == "" {
			color.Red("No path specified error!")
			fmt.Println()
			fmt.Println("Using \"visualizedGit add --help\" for more information")
			return
		}
		scan(path)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCmd.Flags().StringVarP(&path, "path", "p", "", "The path of repositories to be tracked")
}

func scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repositories := lib.RecursiveScanFolder(folder)
	filePath := lib.GetDotFilePath()
	lib.AddNewSliceElementsToFile(filePath, repositories)
	color.Green("\nSuccessfully added repositories\n\n")
}
