package repo

import "database/sql"
import _ "github.com/mattn/go-sqlite3"

type Repository interface {
	Save(filename, data string, parentId int) error
	GetFilesByParentId(id int) []File
	GetRoot() []File
}
type repository struct {
	db *sql.DB
}

func (r repository) Save(filename, data string, parentId int) error {
	var id int
	if err := r.db.QueryRow("insert into files (filename, is_folder, parent_id) values ($1,$2,$3) returning id", filename, false, parentId).Scan(&id); err != nil {
		return err
	}
	_, err := r.db.Exec("insert into cipher_data (cipher_data, files_id) values ($1,$2)", data, id)
	return err
}
func (r repository) GetRoot() []File {
	row, _ := r.db.Query("select id, filename, is_folder from files where parent_id=-1")
	defer row.Close()
	var files []File

	for row.Next() {
		f := File{}
		_ = row.Scan(&f.Id, &f.Filename, &f.IsFolder)
		files = append(files, f)

	}
	return files
}
func (r repository) GetFilesByParentId(id int) []File {
	row, _ := r.db.Query("select id, filename, is_folder from files where parent_id = $2", id, id)
	defer row.Close()
	var files []File

	for row.Next() {
		f := File{}
		_ = row.Scan(&f.Id, &f.Filename, &f.IsFolder)
		files = append(files, f)

	}
	return files
}
func NewRepo() Repository {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	initTables(db)
	initRoot(db)
	return &repository{db: db}
}
func initTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS files (
		id integer PRIMARY KEY autoincrement,
		filename varchar not null,
		is_folder bool not null,
		parent_id int not null
);`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cipher_data (
		id integer PRIMARY KEY autoincrement,
		cipher_data varchar not null,
		files_id int not null,
		foreign key (files_id) references files (id)
	
);`)
	if err != nil {
		panic(err)
	}

}
func initRoot(db *sql.DB) {
	_, err := db.Exec(`
	INSERT INTO files (id, filename, is_folder, parent_id)
	SELECT -1, 'root', true, 0
	WHERE NOT EXISTS (
		SELECT 1 FROM files WHERE filename = 'root' AND parent_id = 0
	);
`)
	if err != nil {
		panic(err)
	}
}
