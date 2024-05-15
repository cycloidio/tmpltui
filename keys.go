package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type keymap = struct {
	file, quit key.Binding
}

func newkeymap() keymap {
	return keymap{
		file: key.NewBinding(
			key.WithKeys("f", "F"),
			key.WithHelp("f", "toggle file picker"),
		),
		quit: key.NewBinding(
			key.WithKeys("esc", "q"),
			key.WithHelp("q", "quit"),
		),
	}
}
