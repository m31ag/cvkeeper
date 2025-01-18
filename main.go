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
	b, err := os.ReadFile("vars.yml")
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	vars := model.NewFromYaml(b)
	p := tea.NewProgram(model.InitModel(r, vars))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Run error: %v", err)
		os.Exit(1)
	}

}
