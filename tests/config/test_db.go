//go:build integration
// +build integration

package config

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/edbmanniwood/pgxpoolmock"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type TestDB struct {
	Pool pgxpoolmock.PgxPool
}

func Connect() *TestDB {
	cfg := FromEnv()

	ctx := context.Background()

	// connection string
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.DbPort, cfg.User, cfg.Password, cfg.DBname)

	config, err := pgxpool.ParseConfig(psqlConn)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Parse Config: is failed"))
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Connect to database is failed"))
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		log.Fatal(errors.Wrap(err, "Ping database error"))
	}

	return &TestDB{Pool: pool}
}

func (db *TestDB) Truncate() {
	ctx := context.Background()

	query := "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name != 'goose_db_version'"
	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	tables := make([]string, 2)

	i := 0
	for rows.Next() {
		err := rows.Scan(&tables[i])
		if err != nil {
			log.Fatal(err)
		}
		i++
	}

	if len(tables) == 0 {
		log.Fatal("No tables, need to run migration")
	}

	query = fmt.Sprintf("Truncate table %s", strings.Join(tables, ","))
	if _, err := db.Pool.Exec(ctx, query); err != nil {
		log.Fatal(err)
	}
}
