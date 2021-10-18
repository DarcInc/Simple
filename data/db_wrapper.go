package data

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type DBCaller interface {
	Query(ctx context.Context, query string, params ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, params ...interface{}) pgx.Row
	Exec(ctx context.Context, query string, params ...interface{}) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (DBCaller, error)
	Release()
}

type PGXDBCaller struct {
	conn *pgxpool.Conn
}

type PGXTxCaller struct {
	trans pgx.Tx
}

func NewDBCaller(conn *pgxpool.Conn) DBCaller {
	return PGXDBCaller{
		conn: conn,
	}
}

func (p PGXDBCaller) Query(ctx context.Context, query string, params ...interface{}) (pgx.Rows, error) {
	return p.conn.Query(ctx, query, params...)
}

func (p PGXDBCaller) QueryRow(ctx context.Context, query string, params ...interface{}) pgx.Row {
	return p.conn.QueryRow(ctx, query, params...)
}

func (p PGXDBCaller) Exec(ctx context.Context, query string, params ...interface{}) (pgconn.CommandTag, error) {
	return p.conn.Exec(ctx, query, params...)
}

func (p PGXDBCaller) Release() {
	p.conn.Release()
}

func (p PGXDBCaller) Begin(ctx context.Context) (DBCaller, error) {
	tx, err := p.conn.Begin(ctx)
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

func (p PGXTxCaller) Release() {
	log.Printf("Warning - calling release on a transaction")
}
