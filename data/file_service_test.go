package data

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
)

func TestNewFileService(t *testing.T) {
	mockDb := &TestDBCaller{}
	fileService := NewFileService(mockDb)
	if fsfs, ok := fileService.(fileSystemFileService); !ok {
		t.Errorf("Did not return the correct underlying type")
		if fsfs.db != mockDb {
			t.Errorf("Expected nil for conn")
		}
	}
}

func fileServiceTestSetup() (*TestDBCaller, FileServer, context.Context) {
	mockDB, ctx := createTestDBCaller()
	return mockDB, NewFileService(mockDB), ctx
}

func TestFileSystemFileService_All(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size").
		WillReturnRows(pgxmock.NewRows([]string{"id", "file_path", "file_hash", "filename", "size"}).
			AddRow(int64(1), "/foo/bar1.jpg", "ABCD1234", "bar1.jpg", int64(1234567890)).
			AddRow(int64(2), "/foo/bar2.jpg", "ABCD1234", "bar2.jpg", int64(1234567890)).
			AddRow(int64(3), "/foo/bar3.jpg", "ABCD1234", "bar3.jpg", int64(1234567890)))

	infos, err := fs.All(ctx, 1, 3)
	if err != nil {
		t.Fatalf("Unable to page files: %v", err)
	}

	if len(infos) != 3 {
		t.Errorf("Expected 3 infos but got %d", len(infos))
	}

	if infos[0].ID != 1 || infos[1].ID != 2 || infos[2].ID != 3 {
		t.Errorf("Expected ids 1, 2, 3 but got %d, %d, %d", infos[0].ID, infos[1].ID, infos[2].ID)
	}
}

func TestFileSystemFileService_AllScanError(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size").
		WillReturnError(errors.New("random database error"))

	infos, err := fs.All(ctx, 1, 3)
	if err == nil {
		t.Error("Expected an error")
	}

	if infos != nil {
		t.Errorf("Did not expect any data back but got %d rows", len(infos))
	}
}

func TestFileSystemFileService_AllQueryError(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size").
		WillReturnError(errors.New("random error"))

	infos, err := fs.All(ctx, 1, 3)
	if err == nil {
		t.Error("Expected an error")
	}

	if infos != nil {
		t.Errorf("Did not expect any data back but got %d rows", len(infos))
	}
}

func TestFileSystemFileService_FindById(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size").
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "file_path", "file_hash", "filename", "size"}).
				AddRow(int64(1), "/foo/bar.jpg", "ABCD1234", "bar.jpg", int64(1234567890)))

	info, err := fs.FindById(ctx, 1)
	if err != nil {
		t.Fatalf("Unable to find by id: %v", err)
	}

	if info.ID != 1 || info.FullPath != "/foo/bar.jpg" || info.FileHash != "ABCD1234" || info.Filename != "bar.jpg" || info.Size != 1234567890 {
		t.Errorf("Expected 1, /foo/bar.jpg ABCD1234 bar.jpg 1234567890 but got %d %s %s %s %d",
			info.ID, info.FullPath, info.FileHash, info.Filename, info.Size)
	}
}

func TestFileSystemFileService_FindByIdNotFound(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size FROM all_files WHERE id = (.*)").
		WillReturnError(pgx.ErrNoRows)

	_, err := fs.FindById(ctx, 1)
	if err != pgx.ErrNoRows {
		t.Fatalf("Should not have found any rows: %v", err)
	}
}

func TestFileSystemFileService_FindByHash(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size").
		WillReturnRows(pgxmock.NewRows([]string{"id", "file_path", "file_hash", "filename", "size"}).
			AddRow(int64(1), "/foo/bar.jpg", "ABCD1234", "bar.jpg", int64(35536)).
			AddRow(int64(2), "/bar/baz.jpg", "ABCD1234", "baz.jpg", int64(35536))).
		RowsWillBeClosed()

	files, err := fs.FindByHash(ctx, "ABCD1234")
	if err != nil {
		t.Fatalf("Error retrieve files by hash: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 results but got %d", len(files))
	}

	info := files[0]
	if info.ID != int64(1) || info.FullPath != "/foo/bar.jpg" || info.Filename != "bar.jpg" || info.Size != int64(35536) {
		t.Errorf("Expected file to be 1, /foo/bar.jpg, bar.jpg, 65535 but got %d %s %s %d",
			info.ID, info.FullPath, info.Filename, info.Size)
	}

	info = files[1]
	if info.ID != int64(2) || info.FullPath != "/bar/baz.jpg" || info.Filename != "baz.jpg" || info.Size != int64(35536) {
		t.Errorf("Expected file to be 2, /bar/baz.jpg, baz.jpg, 65535 but got %d %s %s %d",
			info.ID, info.FullPath, info.Filename, info.Size)
	}
}

func TestFileSystemFileService_FindByHashWithError(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size").
		WillReturnError(errors.New("random Database Error"))

	_, err := fs.FindByHash(ctx, "ABCD1234")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestFileSystemFileService_FindByHashScanError(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size").
		WillReturnRows(pgxmock.NewRows([]string{"id", "file_path", "file_hash", "filename", "size"}).
			AddRow("/foo/bar.jpg", int64(1), "ABCD1234", "bar.jpg", int64(35536)).
			AddRow("/bar/baz.jpg", int64(2), "ABCD1234", "baz.jpg", int64(35536))).
		RowsWillBeClosed()

	files, err := fs.FindByHash(ctx, "ABCD1234")
	if err == nil {
		t.Error("Expected an error")
	}

	if len(files) > 0 {
		t.Errorf("Expected 0 results but got %d", len(files))
	}
}

func TestFileSystemFileService_FindByHashNotFound(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size").
		WillReturnError(pgx.ErrNoRows)

	_, err := fs.FindByHash(ctx, "ABCD1234")
	if err != pgx.ErrNoRows {
		t.Error("Expected error")
	}
}

func TestFileSystemFileService_FindByName(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size").
		WillReturnRows(pgxmock.NewRows([]string{"id", "file_path", "file_hash", "filename", "size"}).
			AddRow(int64(1), "/foo/bar.jpg", "ABCD1234", "bar.jpg", int64(35536)).
			AddRow(int64(2), "/bar/bar.jpg", "ABCD1234", "bar.jpg", int64(35536))).
		RowsWillBeClosed()

	files, err := fs.FindByName(ctx, "bar.jpg")
	if err != nil {
		t.Fatalf("Error retrieve files by hash: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 results but got %d", len(files))
	}

	info := files[0]
	if info.ID != int64(1) || info.FullPath != "/foo/bar.jpg" || info.Filename != "bar.jpg" || info.Size != int64(35536) {
		t.Errorf("Expected file to be 1, /foo/bar.jpg, bar.jpg, 65535 but got %d %s %s %d",
			info.ID, info.FullPath, info.Filename, info.Size)
	}

	info = files[1]
	if info.ID != int64(2) || info.FullPath != "/bar/bar.jpg" || info.Filename != "bar.jpg" || info.Size != int64(35536) {
		t.Errorf("Expected file to be 2, /bar/bar.jpg, bar.jpg, 65535 but got %d %s %s %d",
			info.ID, info.FullPath, info.Filename, info.Size)
	}
}

func TestFileSystemFileService_FindByNameWithError(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size").
		WillReturnError(errors.New("random Database Error"))

	_, err := fs.FindByName(ctx, "ABCD1234")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestFileSystemFileService_FindByNameScanError(t *testing.T) {
	mockDB, fs, ctx := fileServiceTestSetup()

	mockDB.Conn.ExpectQuery("SELECT id, full_path, file_hash, filename, size").
		WillReturnRows(pgxmock.NewRows([]string{"id", "file_path", "file_hash", "filename", "size"}).
			AddRow("/foo/bar.jpg", int64(1), "ABCD1234", "bar.jpg", int64(35536)).
			AddRow("/bar/bar.jpg", int64(2), "ABCD1234", "bar.jpg", int64(35536))).
		RowsWillBeClosed()

	files, err := fs.FindByName(ctx, "bar.jpg")
	if err == nil {
		t.Error("Expected an error")
	}

	if len(files) > 0 {
		t.Errorf("Expected 0 results but got %d", len(files))
	}
}
