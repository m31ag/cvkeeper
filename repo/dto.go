package repo

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
