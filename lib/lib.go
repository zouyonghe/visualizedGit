package lib

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"time"
)

// scanGitFolders returns a list of subfolders of `folder` ending with `.git`.
// Returns the base folder of the repo, the .git folder parent.
// Recursively searches in the subfolders by passing an existing `folders` slice.
func scanGitFolders(folders []string, folder string) []string {
	// trim the last `/`
	folder = strings.TrimSuffix(folder, "/")
	//fmt.Println(folders)
	f, err := os.Open(folder)
	if err != nil {
		zap.L().Error("Error opening folder", zap.Error(err))
		os.Exit(-1)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	files, err := f.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}
			folders = scanGitFolders(folders, path)
		}
	}

	return folders
}

// RecursiveScanFolder starts the recursive search of git repositories
// living in the `folder` subtree
func RecursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

// GetDotFilePath returns the dot file of the repos list.
// Creates it and the enclosing folder if it does not exist.
func GetDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		zap.L().Error("Error fetching current User", zap.Error(err))
		os.Exit(-1)
	}

	dotFile := usr.HomeDir + "/.gogitlocalstats"
	return dotFile
}

// AddNewSliceElementsToFile given a slice of strings representing paths, stores them
// to the filesystem
func AddNewSliceElementsToFile(filePath string, newRepos []string) {
	existingRepos := ParseFileLinesToSlice(filePath)
	repos := joinSlices(newRepos, existingRepos)
	dumpStringsSliceToFile(repos, filePath)
}

// ParseFileLinesToSlice given a file path string, gets the content
// of each line and parses it to a slice of strings.
func ParseFileLinesToSlice(filePath string) []string {
	f := openFile(filePath)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var lines []string
	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			zap.L().Error("Error creating scanner", zap.Error(err))
			os.Exit(-1)
		}
	}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			zap.L().Error("Error reading file", zap.Error(err))
			os.Exit(-1)
		}
	}
	return lines
}

// openFile opens the file located at `filePath`. Creates it if not existing.
func openFile(filePath string) *os.File {
	//f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0755)
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			f, err = os.Create(filePath)
			if err != nil {
				zap.L().Error("Error creating file", zap.Error(err))
				os.Exit(-1)
			}
		} else {
			// other error
			zap.L().Error("Other error occurred", zap.Error(err))
			os.Exit(-1)
		}
	}

	return f
}

// joinSlices adds the element of the `new` slice
// into the `existing` slice, only if not already there
func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !sliceContains(existing, i) {
			existing = append(existing, i)
		}
	}
	return existing
}

// sliceContains returns true if `slice` contains `value`
func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// dumpStringsSliceToFile writes content to the file in path `filePath` (overwriting existing content)
func dumpStringsSliceToFile(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	err := ioutil.WriteFile(filePath, []byte(content), 0755)
	if err != nil {
		zap.L().Error("Error writing to file", zap.Error(err))
		os.Exit(-1)
	}
}

func GetDaysInLastSixMonths() int {
	y, month, day := time.Now().Date()
	sumDays := day
	m := int(month) - 1
	for i := 0; i < 5; i++ {
		if m > 0 {
			sumDays += daysOfMonth(y, m)
		} else {
			sumDays += daysOfMonth(y-1, 12-m)
		}
		m--
	}
	return sumDays
}

func daysOfMonth(year int, month int) int {
	var days int
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30

		} else {
			days = 31

		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}
	return days
}

func GetWeeksInLastSixMon(days int) int {
	return days/7 + 1
}
