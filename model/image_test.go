package model

import (
	"context"
	"testing"
	"time"

	"github.com/darcinc/Simple/data"
)

type mockMetadataServer struct {
	results     []data.Metadata
	single      *data.Metadata
	returnError error
}

func (mms mockMetadataServer) Find(_ context.Context, _ data.MetadataQuery) ([]data.Metadata, error) {
	return mms.results, mms.returnError
}

func (mms mockMetadataServer) FindById(_ context.Context, _ int64) (*data.Metadata, error) {
	return mms.single, mms.returnError
}

func (mms mockMetadataServer) FindByTags(_ context.Context, _ []string) ([]data.Metadata, error) {
	return mms.results, mms.returnError
}

func (mms mockMetadataServer) FindByDateRange(_ context.Context, _, _ time.Time) ([]data.Metadata, error) {
	return mms.results, mms.returnError
}

func (mms mockMetadataServer) FindByLocation(_ context.Context, _ string) ([]data.Metadata, error) {
	return mms.results, mms.returnError
}

func (mms mockMetadataServer) FindByMimeType(_ context.Context, _ []string) ([]data.Metadata, error) {
	return mms.results, mms.returnError
}

func (mms mockMetadataServer) Create(_ context.Context, _ data.Metadata) (data.Metadata, error) {
	return *mms.single, mms.returnError
}

func (mms mockMetadataServer) Save(_ context.Context, _ data.Metadata) error {
	return mms.returnError
}

func TestNewImageRepository(t *testing.T) {
	var ir = NewImageRepository(mockMetadataServer{})
	dir, ok := ir.(dataImageRepository)
	if !ok {
		t.Fatalf("Expected to cast to data image repository")
	}

	if dir.metadataServer == nil {
		t.Errorf("Expected data image repository to have a metadata server")
	}
}

func TestDataImageRepository_Find(t *testing.T) {
	ir := NewImageRepository(mockMetadataServer{
		results: []data.Metadata{
			{ID: 1, Date: time.Now(), Location: "home", Tags: []string{"baz", "qux"}, Data: []data.Encoding{}},
			{ID: 2, Date: time.Now(), Location: "work", Tags: []string{"foo", "bar"}, Data: []data.Encoding{}},
			{ID: 3, Date: time.Now(), Location: "home", Tags: []string{"foo", "bar"}, Data: []data.Encoding{}},
		},
	})
	qp := QueryParameters{}

	images, err := ir.Find(context.Background(), qp)
	if err != nil {
		t.Fatalf("Find method returned an error: %v", err)
	}

	if len(images) != 3 {
		t.Errorf("Expected 3 rows but got %d", len(images))
	}

}
