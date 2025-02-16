package main

import (
	"os"

	"github.com/jasonuc/gignr/cmd"
	"github.com/jasonuc/gignr/internal/utils"
)

var currentVersion = "v1.3.0"

func main() {
	if err := cmd.Execute(currentVersion); err != nil {
		utils.PrintError(err.Error())
		os.Exit(1)
	}
}
