package tui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/jasonuc/gignr/internal/utils"
)

// RunConfirmation displays a clean confirmation UI using huh
func RunConfirmation(prompt string) bool {
	var confirmed bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(prompt).
				Affirmative("Yes").
				Negative("No").
				Value(&confirmed),
		),
	)

	err := form.Run()
	if err != nil {
		utils.PrintError(fmt.Sprintf("Unable to run confirmation: %v", err))
		return false
	}

	return confirmed
}
