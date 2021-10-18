package main

import (
	"context"
	"github.com/darcinc/Simple/data"
	"github.com/darcinc/Simple/model"
	"github.com/darcinc/Simple/reflex"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"time"
)

func main() {
	DBURI := os.Getenv("DB_URI")

	poolConfig, err := pgxpool.ParseConfig(DBURI)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		os.Exit(1)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		os.Exit(1)

	}

	r := reflex.GlobalReflex()

	r.Register("timeout", 15*time.Second)
	r.Register("caller", func(dm reflex.Reflex) (interface{}, bool) {
		timeout, ok := dm.MustGet("timeout").(time.Duration)
		if !ok {
			log.Printf("Warning - timeout was not set to an instance of duration, using 15 seconds")
			timeout = 15 * time.Second
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		conn, err := pool.Acquire(ctx)
		if err != nil {
			log.Printf("Error - Failed to acquire new database connection: %v", err)
			return nil, false
		}
		return data.NewDBCaller(conn), true
	})
	r.Register("fileService", func(dm reflex.Reflex) (interface{}, bool) {
		caller, ok := dm.MustGet("caller").(data.DBCaller)
		if !ok {
			log.Printf("Error - Database connection is not a DBCaller")
			return nil, false
		}

		return data.NewFileService(caller), true
	})
	r.Register("metadataService", func(dm reflex.Reflex) (interface{}, bool) {
		caller, ok := dm.MustGet("caller").(data.DBCaller)
		if !ok {
			log.Printf("Error - Database connection is not a DBCaller")
			return nil, false
		}

		return data.NewMetadataServer(caller), true
	})

	r.Register("imageRepository", func(dm reflex.Reflex) (interface{}, bool) {
		metadataService, ok := dm.MustGet("metadataService").(data.MetadataServer)
		if !ok {
			log.Printf("Error - Metadata service is not a data.MetadataServer")
			return nil, false
		}

		return model.NewImageRepository(metadataService), true
	})
}
