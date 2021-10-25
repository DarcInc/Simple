package data

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
)

// TestDBCaller a mocking object for database queries.
type TestDBCaller struct {
	Conn pgxmock.PgxPoolIface
}

func (tdbc *TestDBCaller) Query(ctx context.Context, query string, params ...interface{}) (pgx.Rows, error) {
	return tdbc.Conn.Query(ctx, query, params...)
}

func (tdbc *TestDBCaller) QueryRow(ctx context.Context, query string, params ...interface{}) pgx.Row {
	return tdbc.Conn.QueryRow(ctx, query, params...)
}
func (tdbc *TestDBCaller) Exec(ctx context.Context, query string, params ...interface{}) (pgconn.CommandTag, error) {
	return tdbc.Conn.Exec(ctx, query, params...)
}

func (tdbc *TestDBCaller) Begin(ctx context.Context) (DBCaller, error) {
	return tdbc, nil
}

func (tdbc *TestDBCaller) Release() {

}
