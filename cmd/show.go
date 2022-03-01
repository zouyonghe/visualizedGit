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
	//"sort"
	"time"
	"visualizedGit/lib"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [-e email|--email email]",
	Short: "Show visualized local git contributions",
	Long:  `Show visualized local git contributions.`,
	Run: func(cmd *cobra.Command, args []string) {
		if email == "" {
			color.Red("No email address specified error!")
			fmt.Println()
			fmt.Println("Using \"visualizedGit add --help\" for more information")
			return
		}
		stats(email)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	showCmd.Flags().StringVarP(&email, "email", "e", "", "The email address reference to commits")
}

const outOfRange = 99999

type column []int

// Stats calculates and prints the stats.
func stats(email string) {

	commits := processRepositories(email)
	printCommitsStats(commits)
}

// getBeginningOfDay given a time.Time calculates the start time of that day
func getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

// countDaysSinceDate counts how many days passed since the passed `date`
func countDaysSinceDate(date time.Time) int {
	days := 0
	now := getBeginningOfDay(time.Now())
	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++
		if days > lib.GetDaysInLastSixMonths() {
			return outOfRange
		}
	}
	return days
}

// fillCommits given a repository found in `path`, gets the commits and
// puts them in the `commits` map, returning it when completed
func fillCommits(email string, path string, commits map[int]int) map[int]int {
	// instantiate a git repo object from path
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}
	// get the HEAD reference
	ref, err := repo.Reference(plumbing.HEAD, false)
	//ref, err := repo.Head()
	if err != nil {
		panic(err)
	}
	// get the commits history starting from HEAD
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		//panic(err)
		return commits
	}
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := countDaysSinceDate(c.Author.When)
		if c.Author.Email != email {
			return nil
		}

		if daysAgo != outOfRange {
			commits[daysAgo]++
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return commits
}

// processRepositories given an user email, returns the
// commits made in the last 6 months
func processRepositories(email string) map[int]int {
	filePath := lib.GetDotFilePath()
	repos := lib.ParseFileLinesToSlice(filePath)
	daysInMap := lib.GetDaysInLastSixMonths()

	commits := make(map[int]int, daysInMap)
	for i := daysInMap; i > 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = fillCommits(email, path, commits)
	}
	return commits
}

// calcOffset determines and returns the amount of days missing to fill
// the last row of the stats graph
func calcOffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		offset = 0
	case time.Monday:
		offset = 1
	case time.Tuesday:
		offset = 2
	case time.Wednesday:
		offset = 3
	case time.Thursday:
		offset = 4
	case time.Friday:
		offset = 5
	case time.Saturday:
		offset = 6
	}
	return offset
}

// printCell given a cell value prints it with a different format
// based on the value amount, and on the `today` flag.
func printCell(val int, today bool) {
	escape := "\033[0;37;30m"
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m"
	case val >= 5 && val < 10:
		escape = "\033[1;30;46m"
	case val >= 10:
		escape = "\033[1;30;42m"
	}

	if today {
		escape = "\033[1;37;45m"
	}

	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	str := "  %d "
	switch {
	case val >= 10:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

// printCommitsStats prints the commits stats
func printCommitsStats(commits map[int]int) {
	cols := buildCols(commits)
	printCells(cols)
}

// buildCols generates a map with rows and columns ready to be printed to screen
func buildCols(commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}
	for i := 0; i < 7-calcOffset()-1; i++ {
		col = append(col, 0)
	}
	for i := 0; i < calcOffset()+1; i++ {
		col = append(col, commits[i])
	}
	cols[0] = col
	col = column{}
	week := 1
	for i := 0; i < len(commits)-calcOffset(); i++ {
		col = append(col, commits[i+calcOffset()+1])
		if len(col) == 7 {
			cols[week] = col
			week++
			col = column{}
		}
	}
	return cols
}

// printCells prints the cells of the graph
func printCells(cols map[int]column) {
	printMonths()
	for j := 0; j < 7; j++ {
		for i := lib.GetWeeksInLastSixMon(lib.GetDaysInLastSixMonths()); i >= 0; i-- {
			if i == lib.GetWeeksInLastSixMon(lib.GetDaysInLastSixMonths()) /*weeksInLastSixMonths*/ {
				printDayCol(j)
			}
			if col, ok := cols[i]; ok {
				//special case today
				if i == 0 && j == calcOffset() {
					printCell(col[7-j-1], true)
					continue
				} else {
					printCell(col[7-j-1], false)
					continue
				}
			}
			printCell(0, false)
		}
		fmt.Printf("\n")
	}
}

// printMonths prints the month names in the first line, determining when the month
// changed between switching weeks
func printMonths() {
	//year, month, day := time.Now().Date()
	week := getBeginningOfDay(time.Now()).Add(-(time.Duration(lib.GetDaysInLastSixMonths()) * time.Hour * 24))
	month := week.Month()
	fmt.Printf("       ")
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		week = week.Add(7 * time.Hour * 24)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

// printDayCol given the day number (0 is Sunday) prints the day name,
// alternating the rows (prints just 2,4,6)
func printDayCol(day int) {
	out := "     "
	switch day {
	case 0:
		out = " Sun "
	case 1:
		out = " Mon "
	case 2:
		out = " Tue "
	case 3:
		out = " Wed "
	case 4:
		out = " Tur "
	case 5:
		out = " Fri "
	case 6:
		out = " Sat "
	}

	fmt.Printf(out)
}
