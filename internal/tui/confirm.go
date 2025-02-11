package tui

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// RunConfirmation displays a clean confirmation UI using huh
func RunConfirmation(prompt string) bool {
	var confirmed bool

	// Create a confirmation form with a group
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(prompt).
				Affirmative("Yes").
				Negative("No").
				Value(&confirmed),
		),
	)

	// Display the prompt
	err := form.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	return confirmed
}
