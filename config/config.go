package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/m31ag/cvkeeper/model"
)

type Config struct {
	DBPath string
	Vars   model.Vars
}

func Load() (*Config, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, fmt.Errorf("failed to setup db path: %w", err)
	}

	varsBytes, err := loadVars()
	if err != nil {
		return nil, fmt.Errorf("failed to load vars: %w", err)
	}

	vars := model.NewFromYaml(varsBytes)

	return &Config{
		DBPath: dbPath,
		Vars:   vars,
	}, nil
}

// getDBPath
// get db path from os enviroment CVKEEPER_DB
// if env is empty, creates .db in ~/.config/cvkeeper directory
func getDBPath() (string, error) {
	if dbPath := os.Getenv("CVKEEPER_DB"); dbPath != "" {
		return dbPath, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dbPath := filepath.Join(home, ".config", "cvkeeper", "cvkeeper.db")

	//create dir if not exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return "", err
	}

	return dbPath, nil
}

// loadVars
// gets variables for app from os env CVKEEPER_VARS
// if env is empty, tries to get from ~/.config/cvkeeper/vars.yml
// if config is empty, tries to fetch template of vars.yml from github gist and save into ~/.config/cvkeeper/ directory
func loadVars() ([]byte, error) {
	if varsPath := os.Getenv("CVKEEPER_VARS"); varsPath != "" {
		return os.ReadFile(varsPath)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cachedVarsPath := filepath.Join(home, ".config", "cvkeeper", "vars.yml")

	if data, err := os.ReadFile(cachedVarsPath); err == nil {
		return data, nil
	}

	data, err := fetchVars()
	if err != nil {
		return nil, err
	}

	_ = os.WriteFile(cachedVarsPath, data, 0644)

	return data, nil
}

func fetchVars() ([]byte, error) {
	url := "https://gist.githubusercontent.com/m31ag/5aa5cb9344b17abc40b9c9849044f0ff/raw/940bd750339ea2a84ffe7142d73d904806f7def9/vars.yml"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vars: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch vars: status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
