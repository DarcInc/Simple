package data

import (
	"github.com/pashagolub/pgxmock"
	"reflect"
	"testing"
)

func TestNewPublishedEntityServer(t *testing.T) {
	dbcCaller := &TestDBCaller{}
	pes := NewPublishedEntityServer(dbcCaller)

	if pes == nil {
		t.Fatalf("Expected a published entity server but got nil")
	}

	if someServer, ok := pes.(dbPublishedEntityServer); ok {
		if someServer.db == nil {
			t.Error("Expected the published entity server to have a database")
		}
	} else {
		t.Errorf("Expected a dbPublishedEntityServer but got %v", reflect.TypeOf(someServer))
	}
}

// What do we want from our generated identifier:
// 1) Stable - given a set of inputs - always the same output.  (No GUID).
// 2) Unique for the range of inputs (int64).  Low or no collisions.
// 3) Not guessable or incremented.  If I create something I shouldn't be able to
//    guess the identifier of the next thing.
// 4) Short - so it can be typed.  (Throws out SHA of number).
// 5) Independent of the number of digits/size of string passed int.
//
// 1) Use the low bytes of time in millis.   low 3 bytes
// 2) The identity of the thing referenced.  low 7 bytes
//                                           -----------
//                                              10 bytes
// So pass in a random seed.  Then xor the 10 bytes with the random seed.
//
// Algorithm: Take the low four bytes of time in micros.  Take the id of the thing.
// Take a random seed of 10 bytes.  Xor them.
// Encode the 10 bytes using upper and lower case 16 letters of the alphabet.
// That's 5 bits per glyph, 80 bits total, or 16 glyphs per id.  Split into 2 groups
// of 8. aDkfbzQr-Tmqldfd
//
func TestDbPublishedEntityServer_Create(t *testing.T) {
	caller, ctx := createTestDBCaller()
	caller.Conn.ExpectQuery("SELECT id FROM metadata WHERE id = 1234").
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1234)))
	caller.Conn.ExpectQuery(`INSERT INTO published_entity 
		(identifier, referenced_entity, referenced_id)
		VALUES `).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1)))

	pes := NewPublishedEntityServer(caller)
	pe, err := pes.Create(ctx, "metadata", 1234)
	if err != nil {
		t.Fatalf("Unexpected error when calling PublishedEntityServer#Create: %v", err)
	}

	if pe.PublishedIdentifier == "" {
		t.Error("Expected published entity identifier to be populated")
	}

	if pe.id == 0 {
		t.Error("Expected id to be populated")
	}

	if pe.Type != "metadata" {
		t.Errorf("Expected entity key to be 'metadata' but got %s", pe.Type)
	}

	if pe.RelatedId != 1234 {
		t.Errorf("Expected entity id to be 1234 but got %d", pe.RelatedId)
	}
}

func TestDbPublishedEntityServer_CreateAll(t *testing.T) {
	t.Error("Not implemented")
}

func TestDbPublishedEntityServer_Find(t *testing.T) {
	t.Error("Not implemented")
}

func TestDbPublishedEntityServer_FindAll(t *testing.T) {
	t.Error("Not implemented")
}

func TestDbPublishedEntityServer_FindOrCreate(t *testing.T) {
	t.Error("Not implemented")
}

func TestDbPublishedEntityServer_FindOrCreateAll(t *testing.T) {
	t.Error("Not implemented")
}

func TestDbPublishedEntityServer_Lookup(t *testing.T) {
	t.Error("Not implemented")
}

func TestDbPublishedEntityServer_LookupAll(t *testing.T) {
	t.Error("Not implemented")
}
