package scan

import (
	"fmt"
	"visualizedGit/lib"
)

func Scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repositories := lib.RecursiveScanFolder(folder)
	filePath := lib.GetDotFilePath()
	lib.AddNewSliceElementsToFile(filePath, repositories)
	fmt.Printf("\nSuccessfully added repositories\n\n")
}
