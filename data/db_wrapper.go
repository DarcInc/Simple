package data

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBCaller interface {
	Query(ctx context.Context, query string, params ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, params ...interface{}) pgx.Row
	Exec(ctx context.Context, query string, params ...interface{}) (pgconn.CommandTag, error)
}

type PGXDBCaller struct {
	pool *pgxpool.Pool
}

func NewDBCaller(pool *pgxpool.Pool) DBCaller {
	return PGXDBCaller{
		pool: pool,
	}
}

func (p PGXDBCaller) Query(ctx context.Context, query string, params ...interface{}) (pgx.Rows, error) {
	return p.pool.Query(ctx, query, params...)
}

func (p PGXDBCaller) QueryRow(ctx context.Context, query string, params ...interface{}) pgx.Row {
	return p.pool.QueryRow(ctx, query, params...)
}

func (p PGXDBCaller) Exec(ctx context.Context, query string, params ...interface{}) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, query, params...)
}