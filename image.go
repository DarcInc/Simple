package simple

import (
	"Simple/data"
	"context"
	"net/url"
	"time"
)



// Source indicates where an image exists.  For example,
// the samge image may be both a TIFF stored in an AWS
// bucket and a JPEG on the local filesystem.
type Source struct {
	// Location indicates where the image resides.  For example, a
	// filesystem image would file:///some/path/to/image.jpg and
	// an API might be https://foo.bar.com/images?id=1234&format=jpeg
	Location url.URL
	// Resolution gives the width and the height of the encoded image
	Resolution data.Resolution
	// Encoding is the format in which the image is encoded as a
	// mime type.  For example image/tiff, if it's a tiff.
	Encoding data.MimeType
}

// Image represents an image which may be encoded into various
// formats, or sizes, but is a version of the same photograph.
// For example, I have a TIFF on the local filesystem.  I have
// another copy as a JPEG in AWS.  And I have a thumbnail in an API, and
// they all reference the same image.
type Image struct {
	id int64
	// Date is the time the original image was captured or created.
	// This is not the same as the file time.  It is
	// possible that the other representations were created
	// later, but this is the best estimate of the date
	// and time the image was first captured.
	Date time.Time
	// Subjects are the focus of the image.  Since an image may
	// capture multiple objects of interest or note at the
	// same time, the subjects are kept as a slice.
	Subjects []string
	// Location indicates where the image was taken, in a
	// human read-able context.  It could be an address or
	// place name like 'home'.  Features such as the
	// latitude or longitude the image was taken would be
	// in an extended image data section.
	Location string
	// Description is the human readable information about
	// the image.
	Description string
	// Sources are the various resource that can contain the
	// image data.  For example, I can take an image and
	// store it in its original JPEG format but also save a
	// copy of the same image as an uncompressed TIFF for
	// archival reasons.  It is still the same image.
	Sources []Source
}

// QueryParameters are the values that are passed to the find
// method to find all the different versions of the same
// image.
type QueryParameters struct {
	// FromDate indicates the starting point in time for the
	// search.
	FromDate time.Time
	// ToDate indicates the ending point in time for the search.
	ToDate time.Time
	// Subjects are the tags that must match.
	Subjects []string
	// Locations represents a list of OR'd together locations.
	Locations []string
}

// ImageRepository allows the user to query for images and open
// streams to the image data.
type ImageRepository interface {
	// Find queries for the image information, given the query
	// parameters.
	Find(ctx context.Context, qp QueryParameters) ([]Image, error)
}

type dataImageRepository struct {
	metadataServer data.MetadataServer
}

func NewImageRepository(server data.MetadataServer) ImageRepository {
	return dataImageRepository{
		metadataServer: server,
	}
}

func (dir dataImageRepository) Find(ctx context.Context, qp QueryParameters) ([]Image, error) {
	mq := data.MetadataQuery{}
	metadata, err := dir.metadataServer.Find(ctx, mq)
	if err != nil {
		return nil, err
	}

	result := make([]Image, len(metadata))

	for i := range metadata {
		result[i] = Image{
			id:       metadata[i].ID,
			Subjects: metadata[i].Tags,
			Date:     metadata[i].Date,
			Location: metadata[i].Location,
		}

		result[i].Sources = make([]Source, len(metadata[i].Locator))
		for l := range metadata[i].Locator {
			//result[i].Sources[l] = Source{Name: "foo"}
		}
	}

	return result, nil
}
