package utils

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintError(err string) {
	fmt.Println(color.New(color.BgRed, color.FgWhite).Sprint(" Error "), err)
}

func PrintWarning(warning string) {
	fmt.Println(color.New(color.BgYellow, color.FgBlack).Sprint(" Warning "), warning)
}

func PrintSuccess(success string) {
	fmt.Println(color.New(color.BgGreen, color.FgBlack).Sprint(" Success "), success)
}

func PrintAlert(alert string) {
	fmt.Println(color.New(color.BgMagenta, color.FgBlack).Sprint(" Alert "), alert)
}
