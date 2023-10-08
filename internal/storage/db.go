package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

type Database struct {
	DB *sql.DB
}

func NewDB(dbConnData string) (*Database, error) {
	db, err := sql.Open("pgx", dbConnData)
	if err != nil {
		fmt.Println("ERROR")
		return nil, err
	}

	newDB := &Database{
		DB: db,
	}
	//logger.Log.Info("connection to database was successful")

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

func (db *Database) FindRecord(ctx context.Context, value string) (Record, error) {
	row := db.DB.QueryRowContext(ctx,
		`SELECT uuid, short_url, origin_url FROM urls WHERE short_url=$1 LIMIT 1`, value)

	var rec Record
	err := row.Scan(&rec.UUID, &rec.ShortURL, &rec.OriginalURL)
	if err != nil {
		return rec, err
	}

	return rec, err
}

func (db *Database) Close() error {
	return db.DB.Close()
}

func (db *Database) SaveRecord(ctx context.Context, rec *Record) error {
	_, err := db.DB.ExecContext(ctx,
		`INSERT INTO urls(uuid, short_url, origin_url) VALUES($1, $2, $3)`,
		rec.UUID, rec.ShortURL, rec.OriginalURL)
	return err
}

func (db *Database) Restore(records []Record) error {
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()

	_, err := db.DB.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS urls(
												"uuid" VARCHAR,
												"short_url" VARCHAR,
												"origin_url" VARCHAR)`)

	row := db.DB.QueryRowContext(ctx, `SELECT count(*) as count from urls`)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	if count < 1 {
		tx, err := db.DB.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		for _, rec := range records {
			_, err := tx.ExecContext(ctx,
				`INSERT INTO urls(uuid, short_url, origin_url) VALUES($1, $2, $3)`,
				rec.UUID, rec.ShortURL, rec.OriginalURL)
			if err != nil {
				logger.Log.Error("db restore error", zap.Error(err))
			}
		}

		return tx.Commit()
	}

	return nil
}
