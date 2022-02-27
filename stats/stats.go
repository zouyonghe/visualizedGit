package stats

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"sort"
	"time"
	"visualizedGit/lib"
)

const outOfRange = 99999
const daysInLastSixMonths = 183
const weeksInLastSixMonths = 26

type column []int

// Stats calculates and prints the stats.
func Stats(email string) {
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
		if days > daysInLastSixMonths {
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
	// iterate the commits
	offset := calcOffset()
	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := countDaysSinceDate(c.Author.When) + offset

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
	daysInMap := daysInLastSixMonths

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
		escape = "\033[1;30;43m"
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
	keys := sortMapIntoSlice(commits)
	cols := buildCols(keys, commits)
	printCells(cols)
}

// sortMapIntoSlice returns a slice of indexes of a map, ordered
func sortMapIntoSlice(m map[int]int) []int {
	// order map
	// To store the keys in slice in sorted order
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// buildCols generates a map with rows and columns ready to be printed to screen
func buildCols(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}

	for _, k := range keys {
		week := k / 7      //26,25...1
		dayInWeek := k % 7 // 0,1,2,3,4,5,6

		if dayInWeek == 0 { //reset
			col = column{}
		}

		col = append(col, commits[k])

		if dayInWeek == 6 {
			cols[week] = col
		}
	}
	fmt.Println(cols)
	return cols
}

// printCells prints the cells of the graph
func printCells(cols map[int]column) {
	printMonths()
	for j := 6; j >= 0; j-- {
		for i := weeksInLastSixMonths + 1; i >= 0; i-- {
			if i == weeksInLastSixMonths+1 {
				printDayCol(j)
			}
			if col, ok := cols[i]; ok {
				//special case today
				if i == 0 && j == calcOffset() {
					printCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						printCell(col[j], false)
						continue
					}
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
	week := getBeginningOfDay(time.Now()).Add(-(daysInLastSixMonths * time.Hour * 24))
	month := week.Month()
	fmt.Printf("         ")
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
		out = " Sat "
	case 2:
		out = " Fri "
	case 3:
		out = " Tur "
	case 4:
		out = " Wed "
	case 5:
		out = " Tue "
	case 6:
		out = " Mon "
	}

	fmt.Printf(out)
}
