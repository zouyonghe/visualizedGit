package removeConfig

import (
	"fmt"
	"log"
	"os"
	"visualizedGit/lib"
)

func Rmcfg() {
	filePath := lib.GetDotFilePath()
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	fmt.Printf("\nSuccessfully removed configuration file\n")
}
