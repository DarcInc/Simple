package data

import (
	"context"
	"errors"
	"time"
)

// PublishedEntity is a data entity that is published outside the system.  The id is an internal
// identifier for the database.  The RelatedId is the id in the foreign table.  The Type is the
// type of entity represented.  Created indicates the date time the entity was created, and along
// with the RelatedID forms the basis of a portable identifier.  Finally, there's the
// PublishedIdentifier which in its string form.
type PublishedEntity struct {
	id                  int64
	RelatedId           int64
	Type                string
	Created             time.Time
	PublishedIdentifier string
}

// LookupResult is a returned lookup from the batch find operation.  SearchFor contains the published
// identifier that was searched.  Found is the value found in the database.  OfType indicates the type
// of data found.  WithError indicates an error finding that value.
type LookupResult struct {
	SearchedFor string
	Found       interface{}
	OfType      string
	WithError   error
}

type PublishedEntityService interface {
	// Lookup returns the underlying data member, its type, and any error.  For now handling not found
	// as an error.
	Lookup(ctx context.Context, publishedId string) (LookupResult, error)
	// LookupAll returns the found entities for a given list of published ids.
	LookupAll(ctx context.Context, publishedIds []string) ([]LookupResult, error)
	// Create takes an entity such as a "metadata" and an id and returns a PublishedEntity that references
	// that underlying datum.
	Create(ctx context.Context, entityName string, id int64) (PublishedEntity, error)
	// CreateAll takes a list of entities and their ids.  It creates the PublishedEntity elements and
	// returns those elements and possibly an error.
	CreateAll(ctx context.Context, entityNames []string, ids []int64) ([]PublishedEntity, error)
	// Find a specific published information for a given entity.
	Find(ctx context.Context, entityType string, id int64) (PublishedEntity, error)
	// FindAll returns all the published entities for a given set of database entities.
	FindAll(ctx context.Context, entityType []string, ids []int64) ([]PublishedEntity, error)
	// FindOrCreate attempts to find a the published entity information for a given database entity.  If
	// there is no published entity information, it creates it.
	FindOrCreate(ctx context.Context, entityName string, id int64) (PublishedEntity, error)
	// FindOrCreateAll takes a given collection for entities, attempt to find or create the published
	// entity information for each one.
	FindOrCreateAll(ctx context.Context, entityNames []string, ids []int64) ([]PublishedEntity, error)
}

type dbPublishedEntityServer struct {
	db DBCaller
}

func NewPublishedEntityServer(db DBCaller) PublishedEntityService {
	return dbPublishedEntityServer{db: db}
}

func (pes dbPublishedEntityServer) Lookup(ctx context.Context, publishedId string) (LookupResult, error) {
	return LookupResult{}, errors.New("not implemented")
}

func (pes dbPublishedEntityServer) LookupAll(ctx context.Context, publishedIds []string) ([]LookupResult, error) {
	return nil, errors.New("not implemented")
}

func (pes dbPublishedEntityServer) Create(ctx context.Context, entityName string, id int64) (PublishedEntity, error) {
	return PublishedEntity{}, errors.New("not implemented")
}

func (pes dbPublishedEntityServer) CreateAll(ctx context.Context, entityNames []string, ids []int64) ([]PublishedEntity, error) {
	return nil, errors.New("not implemented")
}

func (pes dbPublishedEntityServer) Find(ctx context.Context, entityType string, id int64) (PublishedEntity, error) {
	return PublishedEntity{}, errors.New("not implemented")
}

func (pes dbPublishedEntityServer) FindAll(ctx context.Context, entityType []string, ids []int64) ([]PublishedEntity, error) {
	return nil, errors.New("not implemented")
}

func (pes dbPublishedEntityServer) FindOrCreate(ctx context.Context, entityName string, id int64) (PublishedEntity, error) {
	return PublishedEntity{}, errors.New("not implemented")
}

func (pes dbPublishedEntityServer) FindOrCreateAll(ctx context.Context, entityNames []string, ids []int64) ([]PublishedEntity, error) {
	return nil, errors.New("not implemented")
}
