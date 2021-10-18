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
	Begin(ctx context.Context) (DBCaller, error)
}

type PGXDBCaller struct {
	pool *pgxpool.Conn
}

type PGXTxCaller struct {
	trans pgx.Tx
}

func NewDBCaller(pool *pgxpool.Conn) DBCaller {
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

func (p PGXDBCaller) Begin(ctx context.Context) (DBCaller, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return PGXTxCaller{
		trans: tx,
	}, nil
}

func (p PGXTxCaller) Query(ctx context.Context, query string, params ...interface{}) (pgx.Rows, error) {
	return p.trans.Query(ctx, query, params...)
}

func (p PGXTxCaller) QueryRow(ctx context.Context, query string, params ...interface{}) pgx.Row {
	return p.trans.QueryRow(ctx, query, params...)
}

func (p PGXTxCaller) Exec(ctx context.Context, query string, params ...interface{}) (pgconn.CommandTag, error) {
	return p.trans.Exec(ctx, query, params...)
}

func (p PGXTxCaller) Begin(ctx context.Context) (DBCaller, error) {
	return p, nil
}
