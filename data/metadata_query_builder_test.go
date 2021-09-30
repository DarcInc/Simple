package data

import (
	"regexp"
	"strings"
	"testing"
)

var whitespace = regexp.MustCompilePOSIX("[\n\t ]+")

func TestMetadataQueryBuilder_FindById(t *testing.T) {
	qb := NewMetadataQueryBuilder()
	query := qb.FindById()

	target := `SELECT metadata.id, date_captured, location, tags, 
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata 
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		WHERE metadata.id = $1
		ORDER BY metadata.id, encoding.id, locator.id ASC`

	query = strings.Trim(whitespace.ReplaceAllString(query, " "), " ")
	target = strings.Trim(whitespace.ReplaceAllString(target, " "), " ")

	if target != query {
		t.Errorf("Expected %s but got %s", target, query)
	}
}

func TestMetadataQueryBuilder_AddTags(t *testing.T) {
	qb := NewMetadataQueryBuilder()
	query := qb.AddTags(3).String()

	target := `SELECT metadata.id, date_captured, location, tags, 
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata 
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		WHERE $1 = ANY(tags) AND $2 = ANY(tags) AND $3 = ANY(tags)
		ORDER BY metadata.id, encoding.id, locator.id ASC `

	query = strings.Trim(whitespace.ReplaceAllString(query, " "), " ")
	target = strings.Trim(whitespace.ReplaceAllString(target, " "), " ")

	if target != query {
		t.Errorf("Expected %s but got %s", target, query)
	}
}

func TestMetadataQueryBuilder_BetweenDates(t *testing.T) {
	qb := NewMetadataQueryBuilder()
	query := qb.BetweenDates().String()

	target := `SELECT metadata.id, date_captured, location, tags, 
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata 
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		WHERE date_captured BETWEEN $1 AND $2
		ORDER BY metadata.id, encoding.id, locator.id ASC `

	query = strings.Trim(whitespace.ReplaceAllString(query, " "), " ")
	target = strings.Trim(whitespace.ReplaceAllString(target, " "), " ")

	if target != query {
		t.Errorf("Expected %s but got %s", target, query)
	}
}

func TestMetadataQueryBuilder_AtLocation(t *testing.T) {
	qb := NewMetadataQueryBuilder()
	query := qb.AtLocation().String()

	target := `SELECT metadata.id, date_captured, location, tags, 
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata 
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		WHERE location = $1
		ORDER BY metadata.id, encoding.id, locator.id ASC `

	query = strings.Trim(whitespace.ReplaceAllString(query, " "), " ")
	target = strings.Trim(whitespace.ReplaceAllString(target, " "), " ")

	if target != query {
		t.Errorf("Expected %s but got %s", target, query)
	}
}

func TestMetadataQueryBuilder_BetweenDatesTagsLocation(t *testing.T) {
	qb := NewMetadataQueryBuilder()
	query := qb.BetweenDates().AddTags(3).AtLocation().String()

	target := `SELECT metadata.id, date_captured, location, tags, 
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata 
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		WHERE date_captured BETWEEN $1 AND $2 
          AND $3 = ANY(tags) AND $4 = ANY(tags) AND $5 = ANY(tags)
          AND location = $6
		ORDER BY metadata.id, encoding.id, locator.id ASC `

	query = strings.Trim(whitespace.ReplaceAllString(query, " "), " ")
	target = strings.Trim(whitespace.ReplaceAllString(target, " "), " ")

	if target != query {
		t.Errorf("Expected %s but got %s", target, query)
	}
}

func TestMetadataQueryBuilder_StringString(t *testing.T) {
	qb := NewMetadataQueryBuilder()
	query := qb.AtLocation().String()
	query = qb.String()

	target := `SELECT metadata.id, date_captured, location, tags, 
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata 
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		WHERE location = $1
		ORDER BY metadata.id, encoding.id, locator.id ASC `

	query = strings.Trim(whitespace.ReplaceAllString(query, " "), " ")
	target = strings.Trim(whitespace.ReplaceAllString(target, " "), " ")

	if target != query {
		t.Errorf("Expected %s but got %s", target, query)
	}
}

func TestMetadataQueryBuilder_NoMethodsCalled(t *testing.T) {
	qb := NewMetadataQueryBuilder()
	query := qb.String()

	target := `SELECT metadata.id, date_captured, location, tags, 
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata 
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		ORDER BY metadata.id, encoding.id, locator.id ASC `

	query = strings.Trim(whitespace.ReplaceAllString(query, " "), " ")
	target = strings.Trim(whitespace.ReplaceAllString(target, " "), " ")

	if target != query {
		t.Errorf("Expected %s but got %s", target, query)
	}
}

func TestMetadataQueryBuilder_AtLocations(t *testing.T) {
	qb := NewMetadataQueryBuilder()
	query := qb.AtLocations(3).String()

	target := `SELECT metadata.id, date_captured, location, tags, 
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata 
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		WHERE (location = $1 OR location = $2 OR location = $3)
		ORDER BY metadata.id, encoding.id, locator.id ASC `

	query = strings.Trim(whitespace.ReplaceAllString(query, " "), " ")
	target = strings.Trim(whitespace.ReplaceAllString(target, " "), " ")

	if target != query {
		t.Errorf("Expected %s but got %s", target, query)
	}
}

func TestMetadataQueryBuilder_AtTagsAndLocations(t *testing.T) {
	qb := NewMetadataQueryBuilder()
	query := qb.AddTags(2).AtLocations(3).String()

	target := `SELECT metadata.id, date_captured, location, tags, 
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata 
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		WHERE $1 = ANY(tags) AND $2 = ANY(tags) AND (location = $3 OR location = $4 OR location = $5)
		ORDER BY metadata.id, encoding.id, locator.id ASC `

	query = strings.Trim(whitespace.ReplaceAllString(query, " "), " ")
	target = strings.Trim(whitespace.ReplaceAllString(target, " "), " ")

	if target != query {
		t.Errorf("Expected %s but got %s", target, query)
	}
}

func TestMetadataQueryBuilder_AtLocationsAndTags(t *testing.T) {
	qb := NewMetadataQueryBuilder()
	query := qb.AtLocations(3).AddTags(2).String()

	target := `SELECT metadata.id, date_captured, location, tags, 
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata 
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		WHERE (location = $1 OR location = $2 OR location = $3) AND $4 = ANY(tags) AND $5 = ANY(tags) 
		ORDER BY metadata.id, encoding.id, locator.id ASC `

	query = strings.Trim(whitespace.ReplaceAllString(query, " "), " ")
	target = strings.Trim(whitespace.ReplaceAllString(target, " "), " ")

	if target != query {
		t.Errorf("Expected %s but got %s", target, query)
	}
}

func TestMetadataQueryBuilder_ByMimeTypes(t *testing.T) {
	qb := NewMetadataQueryBuilder()
	query := qb.ByMimeTypes(2).String()

	target := `SELECT metadata.id, date_captured, location, tags, 
			encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
			locator.id, locator.source, locator.path
 		FROM metadata 
			INNER JOIN encoding on metadata.id = encoding.metadata_id
    		INNER JOIN locator on encoding.id = locator.encoding_id
		WHERE (encoding.mime_type = $1 OR encoding.mime_type = $2) 
		ORDER BY metadata.id, encoding.id, locator.id ASC `

	query = strings.Trim(whitespace.ReplaceAllString(query, " "), " ")
	target = strings.Trim(whitespace.ReplaceAllString(target, " "), " ")

	if target != query {
		t.Errorf("Expected %s but got %s", target, query)
	}
}