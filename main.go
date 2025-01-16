package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/m31ag/cvkeeper/model"
	"github.com/m31ag/cvkeeper/repo"
	"os"
)

func main() {
	r := repo.NewRepo()
	p := tea.NewProgram(model.InitModel(r))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Run error: %v", err)
		os.Exit(1)
	}
}
