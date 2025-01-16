package repo

import (
	"fmt"
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

func (c *Content) Viewed() string {
	return fmt.Sprintf("%s: %s", c.Filename, c.FileContent)
}
