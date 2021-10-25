package model

type StableReference struct {
}

func NewStableReference(thing interface{}) StableReference {
	return StableReference{}
}

// StableReferenceRepository manage the stable references.  The Refer method
// will automatically create any stable references that are needed.
type StableReferenceRepository interface {
	Refer(things []interface{}) ([]StableReference, error)
}

// TODO: We need to new-up a new Stable Reference Repository
type dataStableReferenceRepository struct {
}

func (dsrr dataStableReferenceRepository) Refer(things []interface{}) ([]StableReference, error) {
	return nil, nil
}
