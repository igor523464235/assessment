package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"str/bff/internal/server"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type DB struct {
	db *sql.DB
}

func New(db *sql.DB) *DB {
	return &DB{db}
}

// TODO: create migrations and move this code to first migration.
// Init creates files table to store information about uploaded files.
func (db *DB) Init(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS files (
			id SERIAL PRIMARY KEY,
			file_id UUID, 
			file_info jsonb
		);`
	_, err := db.db.ExecContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "create table")
	}

	return nil
}

// SaveFileInfo saves information about the uploaded file.
func (db *DB) SaveFileInfo(ctx context.Context, fm *server.FileMetadata) error {
	file := new(File).Marshal(fm)

	query := `INSERT INTO files (file_id, file_info) VALUES ($1, $2);`
	_, err := db.db.ExecContext(ctx, query, file.ID, file.Info)
	if err != nil {
		return errors.Wrap(err, "insert file info ")
	}

	return nil
}

// GetFileInfo retrieves information about the uploaded file.
func (db *DB) GetFileInfo(ctx context.Context, fileID uuid.UUID) (*server.FileMetadata, error) {
	query := `Select file_info FROM files where file_id = $1;`
	res := db.db.QueryRowContext(ctx, query, fileID)
	if res.Err() != nil {
		return nil, errors.Wrap(res.Err(), "get files info ")
	}

	fileInfo := &FileInfo{
		Name:     "",
		Size:     0,
		Storages: make(map[int]int),
	}

	err := res.Scan(fileInfo)
	if err != nil {
		return nil, errors.Wrap(err, "get files info scan")
	}

	file := &File{
		ID:   fileID,
		Info: fileInfo,
	}

	return file.Unmarshal(), nil
}

type FileInfo struct {
	Name     string      `json:"name"`
	Size     uint64      `json:"size"`
	Storages map[int]int `json:"storages"`
}

type File struct {
	ID   uuid.UUID `json:"id"`
	Info *FileInfo `json:"file_info"`
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct .
func (a FileInfo) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *FileInfo) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func (f *File) Marshal(fm *server.FileMetadata) *File {
	f = &File{
		ID: fm.ID,
		Info: &FileInfo{
			Name:     fm.Name,
			Size:     fm.Size,
			Storages: make(map[int]int, len(fm.Storages)),
		},
	}

	for i, storageNum := range fm.Storages {
		f.Info.Storages[int(i)] = int(storageNum)
	}

	return f
}

func (f *File) Unmarshal() *server.FileMetadata {
	fileMetadata := &server.FileMetadata{
		ID:       f.ID,
		Name:     f.Info.Name,
		Size:     f.Info.Size,
		Storages: make(map[server.FilePart]server.ConnectionID),
	}

	for i, storageNum := range f.Info.Storages {
		fileMetadata.Storages[server.FilePart(i)] = server.ConnectionID(storageNum)
	}

	return fileMetadata
}
