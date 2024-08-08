package repo

type File struct {
	Id       int
	Filename string
	IsFolder bool
	ParentId int
}

type CipherData struct {
	Id         int
	CipherData string
	FilesId    int
}
