package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
)

type MimeType string

const (
	MimeBinary MimeType = "application/octet-stream"
	MimeBMP             = "image/bmp"
	MimeGIF             = "image/gif"
	MimeJPEG            = "image/jpeg"
	MimeMPEG            = "audio/mpeg"
	MimeMP4             = "video/mp4"
	MimeOGG             = "video/ogg"
	MimePNG             = "image/png"
	MimeSVG             = "image/svg+xml"
	MimeFlash           = "application/x-shockwave-flash"
	MimeTIFF            = "image/tiff"
	MimeWEBM            = "video/webm"
	MimeWEBP            = "image/webp"
	Mime3GPP            = "video/3gpp"
	Mime3GPP2           = "video/3gpp2"
)

type Resolution struct {
	// The width in units for the media.
	Width int
	// The height in units for the media
	Height int
	// The scan, for example, P or I for progressive or interlaced.
	Scan rune
}

// Metadata is the core information about a given stream of media.
// Each stream may have multiple copies at each of the given locations.
// For example, we have a JPEG file on the filesystem.  The same
// JPEG exists in another directory and at a bucket in AWS.
// If I have a thumbnail of a picture, a JPEG picture, and a TIFF
// encoded version, those are three separate metadata records.
// One for the JPEG picture, one for the TIFF, and one for the
// thumbnail.  If the JPEG version is stored in two place, I still
// have 3 metadata instances, but the JPEG one will have two
// locations.
type Metadata struct {
	ID       int64
	Date     time.Time
	Tags     []string
	Location string
	Data     []Encoding
}
/*
SELECT metadata.id, date_captured, location, tags,
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path

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

 */

type Encoding struct {
	ID         int64
	Metadata   Metadata // -> metadata id on the table.
	Runtime    time.Duration // file specific
	Locator    []Locator  // file specific - a copy of the same data stream
	Resolution Resolution // file specific
	MimeType   string     // file specific
	Hash       string     // file specific
}

// MetadataQuery represents the parameters passed to the Find
// method.
type MetadataQuery struct {
	Tags      []string
	StartDate time.Time
	EndDate   time.Time
	LocatedAt []string
	MimeType  []string
}

// MetadataServer is an interface to a repository of stored Metadata
// information.  The service provides an interface to search the
// repository for matching metadata.
type MetadataServer interface {
	Find(ctx context.Context, query MetadataQuery) ([]Metadata, error)
	FindById(ctx context.Context, id int64) (*Metadata, error)
	FindByTags(ctx context.Context, tags []string) ([]Metadata, error)
	FindByDateRange(ctx context.Context, start, end time.Time) ([]Metadata, error)
	// TODO add FindByMimeType
	FindByLocation(ctx context.Context, location string) ([]Metadata, error)

	// Create stores new metadata and will assign a new ID to that
	// metadata every time.
	Create(ctx context.Context, metadata Metadata) (Metadata, error)
	// Save updates existing metadata and will not assign a new id.
	Save(ctx context.Context, metadata Metadata) error
}

// NewMetadataServer returns an instance of the MetadataServer, in this
// case returning a database based metadata server.
func NewMetadataServer(db DBCaller) MetadataServer {
	return dbMetadataServer{
		db: db,
	}
}

type dbMetadataServer struct {
	db DBCaller
}
/*
SELECT metadata.id, date_captured, location, tags,
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		WHERE (location = $1 OR location = $2 OR location = $3)
		ORDER BY metadata.id, encoding.id, locator.id ASC
 */
func (dms dbMetadataServer) processRows(rows pgx.Rows) ([]Metadata, error) {
	var result []Metadata
	var currentID int64 = -1
	var currentEncodingID int64 =-1

	var currentMetadata Metadata
	var currentEncoding Encoding
	hasResults := false
	for rows.Next() {
		hasResults = true

		var locatorID, encodingID, ID int64
		var source, location, path, fileHash, mimeType string
		var date time.Time
		var runtime time.Duration
		var foundTags []string
		var resolution Resolution

		err := rows.Scan(&ID, &date, &location, &foundTags,
			&encodingID, &runtime, &resolution, &mimeType, &fileHash,
			&locatorID, &source, &path)
		if err != nil {
			return nil, err
		}

		switch {
		case currentID < 0:
			currentID = ID
			currentMetadata = Metadata{
				ID: ID,
				Date: date,
				Location: location,
				Tags: foundTags,
			}
			currentEncoding = Encoding{
				ID: encodingID,
				Runtime: runtime,
				Resolution: resolution,
				MimeType: mimeType,
				Hash: fileHash,
			}
			currentEncodingID = encodingID
		case currentID != ID:
			currentMetadata.Data = append(currentMetadata.Data, currentEncoding)
			result = append(result, currentMetadata)
			currentMetadata = Metadata{
				ID: ID,
				Date: date,
				Location: location,
				Tags: foundTags,
			}
			currentID = ID
			currentEncoding = Encoding{
				ID: encodingID,
				Runtime: runtime,
				Resolution: resolution,
				MimeType: mimeType,
				Hash: fileHash,
			}
			currentEncodingID = encodingID
		case currentEncodingID != encodingID:
			currentMetadata.Data = append(currentMetadata.Data, currentEncoding)
			currentEncoding = Encoding{
				ID: encodingID,
				Runtime: runtime,
				Resolution: resolution,
				MimeType: mimeType,
				Hash: fileHash,
			}
			currentEncodingID = encodingID
		}


		currentEncoding.Locator = append(currentEncoding.Locator, &fileSystemLocator{
			path,
		})
	}

	if hasResults {
		currentMetadata.Data = append(currentMetadata.Data, currentEncoding)
		result = append(result, currentMetadata)
	}

	return result, nil
}

// Find searches for matching Metadata, given the query parameters.
func (dms dbMetadataServer) Find(ctx context.Context, query MetadataQuery) ([]Metadata, error) {
	builder := NewMetadataQueryBuilder()
	var args []interface{}

	if len(query.Tags) > 0 {
		builder = builder.AddTags(len(query.Tags))
		for _, t := range query.Tags {
			args = append(args, t)
		}
	}

	if !query.StartDate.IsZero() || !query.EndDate.IsZero() {
		builder = builder.BetweenDates()
		args = append(args, query.StartDate, query.EndDate)
	}

	if len(query.LocatedAt) == 1 {
		builder = builder.AtLocation()
		args = append(args, query.LocatedAt)
	}

	if len(query.LocatedAt) > 1 {
		builder = builder.AtLocations(len(query.LocatedAt))
		for _, l := range query.LocatedAt {
			args = append(args, l)
		}
	}

	rows, err := dms.db.Query(ctx, builder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return dms.processRows(rows)
}

func (dms dbMetadataServer) FindById(ctx context.Context, id int64) (*Metadata, error) {
	qb := NewMetadataQueryBuilder()

	rows, err := dms.db.Query(ctx, qb.FindById(), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metadata, err := dms.processRows(rows)
	if err != nil {
		return nil, err
	}

	if len(metadata) == 0 {
		return nil, nil
	}

	return &metadata[0], nil
}

func (dms dbMetadataServer) FindByTags(ctx context.Context, tags []string) ([]Metadata, error) {
	query := MetadataQuery{
		Tags: tags,
	}

	return dms.Find(ctx, query)
}

func (dms dbMetadataServer) FindByDateRange(ctx context.Context, start, end time.Time) ([]Metadata, error) {
	query := MetadataQuery{
		StartDate: start,
		EndDate:   end,
	}

	return dms.Find(ctx, query)
}

func (dms dbMetadataServer) FindByLocation(ctx context.Context, location string) ([]Metadata, error) {
	query := MetadataQuery{
		LocatedAt: []string{location},
	}

	return dms.Find(ctx, query)
}

func (dms dbMetadataServer) Create(_ context.Context, _ Metadata) (Metadata, error) {
	return Metadata{}, errors.New("not implemented")
}

func (dms dbMetadataServer) Save(_ context.Context, _ Metadata) error {
	return errors.New("not implemented")
}
