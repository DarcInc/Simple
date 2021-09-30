package data

import (
	"context"
	"errors"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

const (
	EnvDBURI = "DB_URI"
)

var (
	ErrDatabaseURINotSet = errors.New("database URI is not set")
)


func registerResolutionType(_ context.Context, conn *pgx.Conn) error {
	var oid uint32
	row := conn.QueryRow(context.Background(), "select 'resolution'::regtype::oid")
	if err := row.Scan(&oid); err != nil {
		log.Printf("Failed to scan oid: %v", err)
		return err
	}

	// Create the custom type
	ctype, err := pgtype.NewCompositeType("resolution", []pgtype.CompositeTypeField{
		{"width", pgtype.Int4OID},
		{"height", pgtype.Int4OID},
		{"scan", pgtype.BPCharOID},
	}, conn.ConnInfo())
	if err != nil {
		log.Printf("Failed to register new type: %v", err)
		return err
	}

	// Register the custom type with our connection.
	conn.ConnInfo().RegisterDataType(pgtype.DataType{
		Value: ctype,
		Name:  ctype.TypeName(),
		OID:   oid,
	})

	return nil
}

func DBConnect(ctx context.Context) (DBCaller, error) {
	DBURI := os.Getenv(EnvDBURI)
	if DBURI == "" {
		return nil, ErrDatabaseURINotSet
	}

	config := pgxpool.Config{
		AfterConnect: registerResolutionType,
	}

	pool, err := pgxpool.ConnectConfig(ctx, &config)
	if err != nil {
		return nil, err
	}

	return NewDBCaller(pool), nil
}