package pg

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/katyafirstova/auth_service/internal/client/db"
)

type pgClient struct {
	masterDBC db.DB
}

func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	return pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}

func (p pgClient) DB() db.DB {
	return p.masterDBC
}

func (p pgClient) Close() error {
	if p.masterDBC != nil {
		p.masterDBC.Close()
	}

	return nil
}
