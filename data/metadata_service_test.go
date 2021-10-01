package data

import (
	"errors"
	"github.com/pashagolub/pgxmock"
	"log"
	"testing"
	"time"
)

func buildMetadataTestResults() *pgxmock.Rows {
	/*
		id  date        tags           location e.id e.runtime e.resolution e.mimetype e.hash   l.id l.did l.type l.path
		1   05/01/2020  {"foo", "bar"} home     7    0         1024x600 P   image/jpeg ABCD1234 100  7     file   /foo/bar.jpg
		1   05/01/2020  {"foo", "bar"} home     7    0         1024x600 P   image/jpeg ABCD1234 101  7     file   /baz/bar.jpg
		1   05/01/2020  {"foo", "bar"} home     8    0         4096x2160 P  image/tiff 1234ABCD 102  8     file   /archive/bar.tiff
		2   08/20/2020  {"baz"}        work     9    35s       1920x1080 P  video/mp4  78901234 103  9     file   /foo/some.mov
		2   08/20/2020  {"baz"}        work     10   25s       1920x1080 P  video/mp4  45678234 104  10    file   /archive/some.dv
		3   09/04/2020  {"qux"}        home     11   0         4096x1024 P  image/jpeg 2349848D 105  11    file   /foo/other.jpeg

		3 metadata records 1, 2, 3.
		metadata 1 would have 2 encodings, where encoding 1 had two locations.
		metadata 2 would have 2 encodings, with one location each
		metadata 3 would have 1 encoding in one location.

		SELECT metadata.id, date_captured, location, tags,
					encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
					locator.id, locator.source, locator.path

	*/

	rows := pgxmock.NewRows([]string{
		"id", "date_captured", "location", "tags", "encoding.id", "runtime", "resolution", "mime_type", "file_hash",
		"locator.id", "source", "path",
	})
	rows.AddRow(
		int64(1),
		time.Date(2021, 03, 20, 13, 24, 56, 0, time.Local),
		"home",
		[]string{"foo", "bar"},
		int64(10),
		time.Second*0,
		Resolution{
			Width:  1920,
			Height: 1080,
			Scan:   'P',
		},
		MimeJPEG,
		"ABCD1234",
		int64(100),
		"file",
		"/quz/baz.jpg",
	)
	rows.AddRow(
		int64(1),
		time.Date(2021, 03, 20, 13, 24, 56, 0, time.Local),
		"home",
		[]string{"foo", "bar"},
		int64(10),
		time.Second*0,
		Resolution{
			Width:  1920,
			Height: 1080,
			Scan:   'P',
		},
		MimeJPEG,
		"ABCD1234",
		int64(101),
		"file",
		"/quz/baz.jpg")
	rows.AddRow(
		int64(1),
		time.Date(2021, 03, 20, 13, 24, 56, 0, time.Local),
		"home",
		[]string{"foo", "bar"},
		int64(11),
		time.Second*0,
		Resolution{
			Width:  1024,
			Height: 600,
			Scan:   'P',
		},
		MimeTIFF,
		"ABCD1234",
		int64(103),
		"file",
		"/archive/baz.tiff")
	rows.AddRow(
		int64(2),
		time.Date(2021, 03, 20, 13, 24, 56, 0, time.Local),
		"work",
		[]string{"baz", "qux"},
		int64(12),
		time.Second*0,
		Resolution{
			Width:  1024,
			Height: 600,
			Scan:   'P',
		},
		MimeJPEG,
		"ABCD1234",
		int64(104),
		"file",
		"/foo/qux.jpg")

	return rows
}

func buildSingleResult() *pgxmock.Rows {
	rows := pgxmock.NewRows([]string{
		"id", "date_captured", "location", "tags", "encoding.id", "runtime", "resolution", "mime_type", "file_hash",
		"locator.id", "source", "path",
	})
	rows.AddRow(
		int64(1),
		time.Date(2021, 03, 20, 13, 24, 56, 0, time.Local),
		"home",
		[]string{"foo", "bar"},
		int64(10),
		time.Second*0,
		Resolution{
			Width:  1024,
			Height: 600,
			Scan:   'P',
		},
		MimeJPEG,
		"ABCD1234",
		int64(100),
		"file",
		"/quz/baz.jpg",
	)
	rows.AddRow(
		int64(1),
		time.Date(2021, 03, 20, 13, 24, 56, 0, time.Local),
		"home",
		[]string{"foo", "bar"},
		int64(10),
		time.Second*0,
		Resolution{
			Width:  1024,
			Height: 600,
			Scan:   'P',
		},
		MimeJPEG,
		"ABCD1234",
		int64(101),
		"file",
		"/quz/baz.jpg")
	rows.AddRow(
		int64(1),
		time.Date(2021, 03, 20, 13, 24, 56, 0, time.Local),
		"home",
		[]string{"foo", "bar"},
		int64(11),
		time.Second*0,
		Resolution{
			Width:  1024,
			Height: 600,
			Scan:   'P',
		},
		MimeTIFF,
		"ABCD1234",
		int64(103),
		"file",
		"/archive/baz.tiff")

	return rows
}

func TestNewMetadataServer(t *testing.T) {
	tdc := &TestDBCaller{}

	ms := NewMetadataServer(tdc)
	if ms == nil {
		t.Fatal("Returned Metadata server was nil")
	}

	if someServer, ok := ms.(dbMetadataServer); !ok {
		t.Fatal("Unable to cast server to dbMetadataServer")
	} else {
		if someServer.db != tdc {
			t.Error("Expected server caller to be the test caller")
		}
	}
}

func TestDbMetadataServer_FindById(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery(`SELECT metadata\.id, date_captured, location, tags, 
			encoding\.id, encoding\.runtime, encoding\.resolution, encoding\.mime_type, encoding\.file_hash,
			locator\.id, locator\.source, locator\.path
 		FROM metadata 
			INNER JOIN encoding on metadata\.id = encoding\.metadata_id
    		INNER JOIN locator on encoding\.id = locator\.encoding_id
		WHERE metadata.id = \$1
		ORDER BY metadata\.id, encoding\.id, locator\.id ASC`).WillReturnRows(buildSingleResult())

	ms := NewMetadataServer(caller)
	metadata, err := ms.FindById(ctx, 1)
	if err != nil {
		t.Fatalf("Unexpected error when retrieving by ID: %v", err)
	}

	if metadata == nil {
		t.Error("We expected to find a result")
	}

	if metadata.ID != 1 {
		t.Errorf("Expected id = 1 but got %d", metadata.ID)
	}

}

func TestDbMetadataServer_Find(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery(`SELECT metadata\.id, date_captured, location, tags, 
			encoding\.id, encoding\.runtime, encoding\.resolution, encoding\.mime_type, encoding\.file_hash,
			locator\.id, locator\.source, locator\.path
 		FROM metadata 
			INNER JOIN encoding on metadata\.id = encoding\.metadata_id
    		INNER JOIN locator on encoding\.id = locator\.encoding_id
		WHERE \$1 = ANY\(tags\) AND \$2 = ANY\(tags\)
   			AND date_captured BETWEEN \$3 AND \$4
   			AND location = \$5
 		ORDER BY metadata\.id, encoding\.id, locator\.id ASC`).WillReturnRows(buildMetadataTestResults())

	ms := NewMetadataServer(caller)
	query := MetadataQuery{
		Tags:      []string{"foo", "bar"},
		StartDate: time.Now(),
		EndDate:   time.Now(),
		LocatedAt: []string{"home"},
	}
	metadata, err := ms.Find(ctx, query)
	if err != nil {
		t.Fatalf("Error returned from find: %v", err)
	}

	if len(metadata) != 2 {
		t.Errorf("Expected 2 result but got %d", len(metadata))
	}

	if len(metadata[0].Data) != 2 {
		t.Fatalf("Expected first result to have 2 encodings but got: %d", len(metadata[0].Data))
	}

	if metadata[0].Data[0].MimeType != MimeJPEG || metadata[0].Data[1].MimeType != MimeTIFF {
		t.Errorf("Expected a JPEG and a TIFF for the first metadata but got %s and %s",
			metadata[0].Data[0].MimeType, metadata[0].Data[1].MimeType)
	}

	res := metadata[0].Data[0].Resolution
	if res.Width != 1920 || res.Height != 1080 || res.Scan != 'P' {
		t.Errorf("Expected a 1920x1080 P but got a %dx%d %c", res.Width, res.Height, res.Scan)
	}

	res = metadata[0].Data[1].Resolution
	if res.Width != 1024 || res.Height != 600 || res.Scan != 'P' {
		t.Errorf("Expected a 1024x600 P but got a %dx%d %c", res.Width, res.Height, res.Scan)
	}

	if len(metadata[0].Data[0].Locator) != 2 {
		t.Errorf("Expected first result, first encoding to have 2 locators but got: %d", len(metadata[0].Data[0].Locator))
	}

	if len(metadata[0].Data[1].Locator) != 1 {
		t.Errorf("Expected first result, second encoding to have 1 locators but got: %d", len(metadata[0].Data[1].Locator))
	}

	if len(metadata[1].Data) != 1 {
		t.Fatalf("Expected second result to have 1 encodings but got: %d", len(metadata[1].Data))
	}

	if len(metadata[0].Data[0].Locator) != 2 {
		t.Errorf("Expected second result, first encoding to have 1 locator but got: %d", len(metadata[1].Data[0].Locator))
	}
}

func TestDbMetadataServer_FindMultipleLocations(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery(`SELECT metadata\.id, date_captured, location, tags, 
			encoding\.id, encoding\.runtime, encoding\.resolution, encoding\.mime_type, encoding\.file_hash,
			locator\.id, locator\.source, locator\.path
 		FROM metadata 
			INNER JOIN encoding on metadata\.id = encoding\.metadata_id
    		INNER JOIN locator on encoding\.id = locator\.encoding_id
 		WHERE \$1 = ANY\(tags\) AND \$2 = ANY\(tags\)
   			AND date_captured BETWEEN \$3 AND \$4
   			AND \(location = \$5 OR location = \$6\) 
 		ORDER BY metadata\.id, encoding\.id, locator\.id ASC`).WillReturnRows(buildMetadataTestResults())

	ms := NewMetadataServer(caller)
	query := MetadataQuery{
		Tags:      []string{"foo", "bar"},
		StartDate: time.Now(),
		EndDate:   time.Now(),
		LocatedAt: []string{"home", "work"},
	}
	metadata, err := ms.Find(ctx, query)
	if err != nil {
		t.Fatalf("Error returned from find: %v", err)
	}

	if len(metadata) != 2 {
		t.Errorf("Expected 2 result but got %d", len(metadata))
	}
}

func TestDbMetadataServer_FindEmptyQueryParameters(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery(`SELECT metadata\.id, date_captured, location, tags, 
			encoding\.id, encoding\.runtime, encoding\.resolution, encoding\.mime_type, encoding\.file_hash,
			locator\.id, locator\.source, locator\.path
 		FROM metadata 
			INNER JOIN encoding on metadata\.id = encoding\.metadata_id
    		INNER JOIN locator on encoding\.id = locator\.encoding_id
	 	ORDER BY metadata\.id, encoding\.id, locator\.id ASC`).WillReturnRows(buildMetadataTestResults())

	ms := NewMetadataServer(caller)
	query := MetadataQuery{}
	metadata, err := ms.Find(ctx, query)
	if err != nil {
		log.Fatalf("Error returned from find: %v", err)
	}

	if len(metadata) != 2 {
		t.Errorf("Expected 2 result but got %d", len(metadata))
	}
}

func TestDbMetadataServer_FindByTags(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery(`SELECT metadata\.id, date_captured, location, tags, 
			encoding\.id, encoding\.runtime, encoding\.resolution, encoding\.mime_type, encoding\.file_hash,
			locator\.id, locator\.source, locator\.path
 		FROM metadata 
			INNER JOIN encoding on metadata\.id = encoding\.metadata_id
    		INNER JOIN locator on encoding\.id = locator\.encoding_id
 		WHERE \$1 = ANY\(tags\) AND \$2 = ANY\(tags\)
 		ORDER BY metadata\.id, encoding\.id, locator\.id ASC`).WillReturnRows(buildMetadataTestResults())

	ms := NewMetadataServer(caller)
	metadata, err := ms.FindByTags(ctx, []string{"foo", "bar"})

	if err != nil {
		t.Fatalf("Error returned from find by tags: %v", err)
	}

	if len(metadata) != 2 {
		t.Errorf("Expected 2 metadata records but got: %d", len(metadata))
	}
}

func TestDbMetadataServer_FindMimeType(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery(`SELECT metadata\.id, date_captured, location, tags, 
			encoding\.id, encoding\.runtime, encoding\.resolution, encoding\.mime_type, encoding\.file_hash,
			locator\.id, locator\.source, locator\.path
 		FROM metadata 
			INNER JOIN encoding on metadata\.id = encoding\.metadata_id
    		INNER JOIN locator on encoding\.id = locator\.encoding_id
 		WHERE \(encoding\.mime_type = \$1 OR encoding\.mime_type = \$2\)
 		ORDER BY metadata\.id, encoding\.id, locator\.id ASC`).WillReturnRows(buildMetadataTestResults())

	ms := NewMetadataServer(caller)
	metadata, err := ms.FindByMimeType(ctx, []string{MimeJPEG, MimeTIFF})

	if err != nil {
		t.Fatalf("Error returned from find by tags: %v", err)
	}

	if len(metadata) != 2 {
		t.Errorf("Expected 2 metadata records but got: %d", len(metadata))
	}
}

func TestDbMetadataServer_FindByTagsNoData(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery("SELECT").WillReturnRows(pgxmock.NewRows([]string{
		"id", "date_captured", "location", "tags", "encoding.id", "runtime", "resolution", "mime_type", "file_hash",
		"locator.id", "source", "path",
	}))

	ms := NewMetadataServer(caller)
	metadata, err := ms.FindByTags(ctx, []string{"foo", "bar"})

	if err != nil {
		t.Fatalf("Error returned from find by tags: %v", err)
	}

	if len(metadata) != 0 {
		t.Errorf("Expected no metadata records but got: %d", len(metadata))
	}
}

func TestDbMetadataServer_FindByTagsDatabaseError(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery("SELECT").WillReturnError(errors.New("random database error"))

	ms := NewMetadataServer(caller)
	metadata, err := ms.FindByTags(ctx, []string{"foo", "bar"})

	if err == nil {
		t.Error("Expected database error")
	}

	if len(metadata) != 0 {
		t.Errorf("Expected no metadata but got %d rows", len(metadata))
	}
}

func TestDbMetadataServer_FindByDateRange(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery(`SELECT metadata\.id, date_captured, location, tags, 
			encoding\.id, encoding\.runtime, encoding\.resolution, encoding\.mime_type, encoding\.file_hash,
			locator\.id, locator\.source, locator\.path
 		FROM metadata 
			INNER JOIN encoding on metadata\.id = encoding\.metadata_id
    		INNER JOIN locator on encoding\.id = locator\.encoding_id
 		WHERE date_captured BETWEEN \$1 AND \$2
 		ORDER BY metadata\.id, encoding\.id, locator\.id ASC`).WillReturnRows(buildMetadataTestResults())

	ms := NewMetadataServer(caller)
	startDate, err := time.Parse("2006-01-02T15:04:05", "2021-02-20T00:00:00")
	if err != nil {
		t.Fatalf("Failed to parse date: %v", err)
	}

	endDate, err := time.Parse("2006-01-02T15:04:05", "2021-03-21T00:00:00")
	if err != nil {
		t.Fatalf("Failed to parse date: %v", err)
	}

	metadata, err := ms.FindByDateRange(ctx, startDate, endDate)

	if err != nil {
		t.Fatalf("Error returned from find by tags: %v", err)
	}

	if len(metadata) != 2 {
		t.Errorf("Expected 2 metadata records but got: %d", len(metadata))
	}
}

func TestDbMetadataServer_FindByDateRangeDatabaseError(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery("SELECT").WillReturnError(errors.New("random database error"))

	ms := NewMetadataServer(caller)
	startDate, _ := time.Parse(time.RFC3339, "2021-02-20T00:00:00Z-4:00")
	endDate, _ := time.Parse(time.RFC3339, "2021-03-21T00:00:00Z-4:00")
	metadata, err := ms.FindByDateRange(ctx, startDate, endDate)

	if err == nil {
		t.Error("Expected database error")
	}

	if len(metadata) != 0 {
		t.Errorf("Expected no metadata but got %d rows", len(metadata))
	}
}

func TestDbMetadataServer_FindByLocation(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery(`SELECT metadata\.id, date_captured, location, tags, 
			encoding\.id, encoding\.runtime, encoding\.resolution, encoding\.mime_type, encoding\.file_hash,
			locator\.id, locator\.source, locator\.path
 		FROM metadata 
			INNER JOIN encoding on metadata\.id = encoding\.metadata_id
    		INNER JOIN locator on encoding\.id = locator\.encoding_id
 		WHERE location = \$1
 		ORDER BY metadata\.id, encoding\.id, locator\.id ASC`).WillReturnRows(buildMetadataTestResults())

	ms := NewMetadataServer(caller)
	metadata, err := ms.FindByLocation(ctx, "home")

	if err != nil {
		t.Fatalf("Error returned from find by tags: %v", err)
	}

	if len(metadata) != 2 {
		t.Errorf("Expected 2 metadata records but got: %d", len(metadata))
	}
}

func TestDbMetadataServer_FindByLocationDatabaseError(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery("SELECT").WillReturnError(errors.New("random database error"))

	ms := NewMetadataServer(caller)
	metadata, err := ms.FindByLocation(ctx, "home")

	if err == nil {
		t.Error("Expected database error")
	}

	if len(metadata) != 0 {
		t.Errorf("Expected no metadata but got %d rows", len(metadata))
	}
}
