package postgres

import (
	"context"
	"fmt"

	"github.com/LiliyaD/Reminder_telegram_bot/internal/journal"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

func connect() *pgxpool.Pool {
	ctx := context.Background()

	// connection string
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Password, DBname)

	config, err := pgxpool.ParseConfig(psqlConn)
	if err != nil {
		journal.LogFatal(errors.Wrap(err, "Parse Config: is failed"))
	}

	config.MaxConnIdleTime = MaxConnIdleTime
	config.MaxConnLifetime = MaxConnLifetime
	config.MinConns = MinConns
	config.MaxConns = MaxConns

	// fix for error: prepared statement "lrupsc_1_0" already exists (SQLSTATE 42P05)
	config.ConnConfig.PreferSimpleProtocol = true // disable prepared statements caching by default
	config.ConnConfig.RuntimeParams = map[string]string{"standard_conforming_strings": "on"}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		journal.LogFatal(errors.Wrap(err, "Connect to database is failed"))
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		journal.LogFatal(errors.Wrap(err, "Ping database error"))
	}

	return pool
}
