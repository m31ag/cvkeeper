package model

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

type Vars struct {
	Colors Colors `yaml:"colors"`
	Sizes  Sizes  `yaml:"sizes"`
	Icons  Icons  `yaml:"icons"`
}
type Colors struct {
	HintKey             lipgloss.Color `yaml:"hint_key"`
	History             lipgloss.Color `yaml:"history"`
	Selected            lipgloss.Color `yaml:"selected"`
	DefaultTextColor    lipgloss.Color `yaml:"default_text_color"`
	BoxBorderColor      lipgloss.Color `yaml:"box_border_color"`
	ErrorTextColor      lipgloss.Color `yaml:"error_text_color"`
	SuccessTextColor    lipgloss.Color `yaml:"success_text_color"`
	TitleBackground     lipgloss.Color `yaml:"title_background"`
	HorizontalSeparator lipgloss.Color `yaml:"horizontal_separator"`
	ViewFileTitle       lipgloss.Color `yaml:"view_file_title"`
}
type Sizes struct {
}
type Icons struct {
	Folder     string `yaml:"folder"`
	CipherData string `yaml:"cipher_data"`
	File       string `yaml:"file"`
	MasterKey  string `yaml:"register_key"`
}

func NewFromYaml(b []byte) Vars {
	var v Vars
	if err := yaml.Unmarshal(b, &v); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return v
}
