package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type Database struct {
	DB *sql.DB
}

func NewDB(dbConnData string) (*Database, error) {
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

func (db *Database) FindRecordByOriginURL(ctx context.Context, value string) (Record, error) {
	row := db.DB.QueryRowContext(ctx,
		`SELECT uuid, short_url, origin_url FROM urls WHERE origin_url=$1 LIMIT 1`, value)

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

func (db *Database) CreateDBScheme() error {
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()

	_, err := db.DB.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS urls(
												"uuid" VARCHAR,
												"short_url" VARCHAR,
												"origin_url" VARCHAR)`)
	if err != nil {
		return err
	}

	_, err = db.DB.ExecContext(ctx,
		`CREATE UNIQUE INDEX IF NOT EXISTS origin_url_idx on urls(origin_url)`)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) SaveRecordsBatch(ctx context.Context, records []Record) error {
	tx, err := db.DB.Begin()
	if err != nil {
		rb := tx.Rollback()
		if rb != nil {
			return rb
		}
		return err
	}

	for _, rec := range records {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO urls(uuid, short_url, origin_url) VALUES($1, $2, $3)`,
			rec.UUID, rec.ShortURL, rec.OriginalURL)

		if err != nil {
			rb := tx.Rollback()
			if rb != nil {
				return rb
			}
			return err
		}
	}

	return tx.Commit()
}
