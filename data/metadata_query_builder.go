package data

import (
	"fmt"
	"strings"
)

// TODO add encoding id for control-break on new encodings
const (
	queryBase = `SELECT metadata.id, date_captured, location, tags,
		encoding.id, encoding.runtime, encoding.resolution, encoding.mime_type, encoding.file_hash,
		locator.id, locator.source, locator.path
	FROM metadata
         INNER JOIN encoding on metadata.id = encoding.metadata_id
         INNER JOIN locator on encoding.id = locator.encoding_id`
	orderClause = `ORDER BY metadata.id, encoding.id, locator.id ASC`
)

// TODO add mimetype query

// MetadataQueryBuilder builds a metadata query in a fluent manner.
type MetadataQueryBuilder struct {
	b strings.Builder
	idx int
}

// NewMetadataQueryBuilder returns a query builder.
func NewMetadataQueryBuilder() *MetadataQueryBuilder {
	return &MetadataQueryBuilder{
		b: strings.Builder{},
		idx: 1,
	}
}

func (qb *MetadataQueryBuilder) FindById() string {
	qb.addFrontMatter()
	qb.b.WriteString(fmt.Sprintf("metadata.id = $%d %s", qb.idx, orderClause))
	return qb.b.String()
}

func (qb *MetadataQueryBuilder) addFrontMatter() *MetadataQueryBuilder {
	if qb.b.Len() == 0 {
		_, _ = qb.b.WriteString(fmt.Sprintf("%s WHERE ", queryBase))
	} else {
		_, _ = qb.b.WriteString(" AND ")
	}
	return qb
}

// AddTags adds placeholders for tags to pass to the query.
func (qb *MetadataQueryBuilder) AddTags(ntags int) *MetadataQueryBuilder {
	qb.addFrontMatter()

	for i := qb.idx; i < qb.idx + ntags; i++ {
		qb.b.WriteString(fmt.Sprintf("$%d = ANY(tags) ", i))
		if i < (qb.idx + ntags - 1) {
			qb.b.WriteString("AND ")
		}
	}

	qb.idx = qb.idx + ntags

	return qb
}

// BetweenDates adds a date range clause to the query.
func (qb *MetadataQueryBuilder) BetweenDates() *MetadataQueryBuilder {
	qb.addFrontMatter()

	qb.b.WriteString(fmt.Sprintf(`date_captured BETWEEN $%d AND $%d `, qb.idx, qb.idx + 1))
	qb.idx = qb.idx + 2

	return qb
}

// AtLocation adds a location clause to the query.
func (qb *MetadataQueryBuilder) AtLocation() *MetadataQueryBuilder {
	qb.addFrontMatter()

	qb.b.WriteString(fmt.Sprintf(`location = $%d `, qb.idx))
	qb.idx = qb.idx + 1

	return qb
}

func (qb *MetadataQueryBuilder) AtLocations(nlocs int) *MetadataQueryBuilder {
	qb.addFrontMatter()

	qb.b.WriteString("(")
	for i := qb.idx; i < qb.idx + nlocs; i++ {
		qb.b.WriteString(fmt.Sprintf("location = $%d", i))
		if i < (qb.idx + nlocs - 1) {
			qb.b.WriteString(" OR ")
		}
	}
	qb.idx = qb.idx + nlocs

	qb.b.WriteString(") ")

	return qb
}

func (qb *MetadataQueryBuilder) ByMimeTypes(ntypes int) *MetadataQueryBuilder {
	qb.addFrontMatter()
	qb.b.WriteString("(")
	for i := qb.idx; i < qb.idx + ntypes; i++ {
		qb.b.WriteString(fmt.Sprintf("encoding.mime_type = $%d", i))
		if i < (qb.idx + ntypes - 1) {
			qb.b.WriteString(" OR ")
		}
	}
	qb.idx = qb.idx + ntypes

	qb.b.WriteString(") ")
	return qb
}

// String creates the final query string and resets the builder.
func (qb *MetadataQueryBuilder) String() string {
	if qb.b.Len() == 0 {
		return fmt.Sprintf("%s %s", queryBase, orderClause)
	}
	result := fmt.Sprintf("%s %s", qb.b.String(), orderClause)
	return result
}
