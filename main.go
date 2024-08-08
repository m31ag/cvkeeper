package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/m31ag/cvkeeper/model"
	"os"
)

func main() {
	p := tea.NewProgram(model.InitModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
