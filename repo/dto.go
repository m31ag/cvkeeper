package repo

import (
	"fmt"
	"strings"
)

type File struct {
	Id       int
	ParentId int
	Filename string
	IsFolder bool
}

type CipherData struct {
	Id         int
	CipherData string
	FilesId    int
}
type deleteDto struct {
	Id       int
	IsFolder bool
}
type Content struct {
	Filename    string
	FileContent string
}

func (c *Content) View() string {
	return fmt.Sprintf("name: %s\n%s\n%s", c.Filename, strings.Repeat("_", 20), c.FileContent)
}
