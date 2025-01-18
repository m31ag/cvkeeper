package model

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Vars struct {
	Colors Colors `yaml:"colors"`
}
type Colors struct {
	HintKey string `yaml:"hint_key"`
	History string `yaml:"history"`
}

func NewFromYaml(b []byte) Vars {
	var v Vars
	if err := yaml.Unmarshal(b, &v); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return v
}
