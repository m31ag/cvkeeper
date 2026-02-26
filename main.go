package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m31ag/cvkeeper/config"
	"github.com/m31ag/cvkeeper/model"
	"github.com/m31ag/cvkeeper/repo"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	r := repo.NewRepo(cfg.DBPath)

	p := tea.NewProgram(model.InitModel(r, cfg.Vars))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Run error: %v", err)
		os.Exit(1)
	}

}
