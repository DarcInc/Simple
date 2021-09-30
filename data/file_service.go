package data

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/jackc/pgx/v4"
)

type FileInfo struct {
	ID       int64
	FullPath string
	FileHash string
	Filename string
	Size     int64
}

var (
	ErrFileInfoNotFound = errors.New("file info not found")
)

const (
	selectAllPaging = `SELECT id, full_path, file_hash, filename, size 
		FROM all_files
		ORDER BY id
		OFFSET $1 LIMIT $2`
	selectFileInfoById = `SELECT id, full_path, file_hash, filename, size 
		FROM all_files
		WHERE id = $1`
	selectFileInfoByName = `SELECT id, full_path, file_hash, filename, size 
		FROM all_files
		WHERE name = $1`
	selectFileInfoByHash = `SELECT id, full_path, file_hash, filename, size
		FROM all_files
		WHERE hash = $1`
)

// FileServer is the interface to the data for the files.
type FileServer interface {
	OpenFile(filePath string) (io.ReadCloser, error)
	All(ctx context.Context, start, pageSize int) ([]FileInfo, error)
	FindById(ctx context.Context, id int64) (FileInfo, error)
	FindByHash(ctx context.Context, hash string) ([]FileInfo, error)
	FindByName(ctx context.Context, name string) ([]FileInfo, error)
}

type fileSystemFileService struct {
	db DBCaller
}

func NewFileService(db DBCaller) FileServer {
	return fileSystemFileService{
		db: db,
	}
}

// OpenFile is used to open a stream to the file data.
func (fs fileSystemFileService) OpenFile(filePath string) (io.ReadCloser, error) {
	return os.Open(filePath)
}

func (fs fileSystemFileService) processRows(rows pgx.Rows) ([]FileInfo, error) {
	var result []FileInfo
	for rows.Next() {
		var info FileInfo
		if err := rows.Scan(&info.ID, &info.FullPath, &info.FileHash, &info.Filename, &info.Size); err != nil {
			return nil, err
		}

		result = append(result, info)
	}
	return result, nil
}

func (fs fileSystemFileService) All(ctx context.Context, start, pageSize int) ([]FileInfo, error) {
	rows, err := fs.db.Query(ctx, selectAllPaging, start, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := fs.processRows(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (fs fileSystemFileService) FindById(ctx context.Context, id int64) (FileInfo, error) {
	row := fs.db.QueryRow(ctx, selectFileInfoById, id)
	result := FileInfo{}

	err := row.Scan(&result.ID, &result.FullPath, &result.FileHash, &result.Filename, &result.Size)
	return result, err
}

func (fs fileSystemFileService) FindByHash(ctx context.Context, hash string) ([]FileInfo, error) {
	rows, err := fs.db.Query(ctx, selectFileInfoByHash, hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := fs.processRows(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (fs fileSystemFileService) FindByName(ctx context.Context, name string) ([]FileInfo, error) {
	rows, err := fs.db.Query(ctx, selectFileInfoByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := fs.processRows(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}
