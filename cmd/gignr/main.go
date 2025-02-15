package main

import (
	"os"

	"github.com/jasonuc/gignr/internal/utils"
)

var currentVersion = "v1.1.0"

func main() {
	if err := execute(); err != nil {
		utils.PrintError(err.Error())
		os.Exit(1)
	}
}
