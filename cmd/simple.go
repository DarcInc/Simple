package main

import (
	"context"
	"github.com/darcinc/Simple/data"
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
			log.Printf("Failed to acquire new database connection: %v", err)
			return nil, false
		}
		return data.NewDBCaller(conn), true
	})
}
