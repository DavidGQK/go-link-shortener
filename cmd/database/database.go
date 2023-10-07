package database

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type Database struct {
	DB *sql.DB
}

func New(dbConnData string) (*Database, error) {
	db, err := sql.Open("pgx", dbConnData)
	if err != nil {
		return nil, err
	}

	newDB := &Database{
		DB: db,
	}

	return newDB, nil
}

func (db *Database) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.DB.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}
