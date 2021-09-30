package data

import (
	"context"
	"github.com/pashagolub/pgxmock"
)

func createTestDBCaller() (*TestDBCaller, context.Context) {
	pgxIface, _ := pgxmock.NewPool()
	mockDB := &TestDBCaller{
		pgxIface,
	}

	return mockDB, context.Background()
}