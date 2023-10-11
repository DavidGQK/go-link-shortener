package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

var ErrConflict = errors.New(`already exists`)

type Database struct {
	DBConnData string
	DB         *sql.DB
	Mode       int
	Links      map[string]string
}

func NewDB(dbConnData string, mode int) (*Database, error) {
	db, err := sql.Open("pgx", dbConnData)
	if err != nil {
		return nil, err
	}

	newDB := &Database{
		DBConnData: dbConnData,
		DB:         db,
		Mode:       mode,
		Links:      make(map[string]string),
	}

	return newDB, nil
}

func (db *Database) Restore() error {
	err := db.CreateDBScheme()
	if err != nil {
		logger.Log.Error("db restoring error", zap.Error(err))
		return err
	}

	return nil
}

func (db *Database) Add(key, value string) error {
	id := uuid.NewString()
	rec := models.Record{
		UUID:        id,
		ShortURL:    key,
		OriginalURL: value,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.SaveRecord(ctx, &rec)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return ErrConflict
		}

		logger.Log.Error("error while writing data to db", zap.Error(err))
		return err
	}

	db.Links[key] = value
	return nil
}

func (db *Database) AddBatch(ctx context.Context, records []models.Record) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.SaveRecordsBatch(ctx, records)
	if err != nil {
		logger.Log.Error("error while writing data batch to db", zap.Error(err))
	}

	for _, rec := range records {
		db.Links[rec.ShortURL] = rec.OriginalURL
	}

	return nil
}

func (db *Database) Get(key string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	rec, err := db.FindRecord(ctx, key)
	if err != nil {
		logger.Log.Error("db search query error", zap.Error(err))
		return "", false
	}

	return rec.OriginalURL, true
}

func (db *Database) GetMode() int {
	return db.Mode
}

func (db *Database) GetByOriginURL(originURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	rec, err := db.FindRecordByOriginURL(ctx, originURL)
	if err != nil {
		return "", err
	}

	return rec.ShortURL, nil
}

func (db *Database) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.DB.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func (db *Database) FindRecord(ctx context.Context, value string) (models.Record, error) {
	row := db.DB.QueryRowContext(ctx,
		`SELECT uuid, short_url, origin_url FROM urls WHERE short_url=$1 LIMIT 1`, value)

	var rec models.Record
	err := row.Scan(&rec.UUID, &rec.ShortURL, &rec.OriginalURL)
	if err != nil {
		return rec, err
	}

	return rec, err
}

func (db *Database) FindRecordByOriginURL(ctx context.Context, value string) (models.Record, error) {
	row := db.DB.QueryRowContext(ctx,
		`SELECT uuid, short_url, origin_url FROM urls WHERE origin_url=$1 LIMIT 1`, value)

	var rec models.Record
	err := row.Scan(&rec.UUID, &rec.ShortURL, &rec.OriginalURL)
	if err != nil {
		return rec, err
	}

	return rec, err
}

func (db *Database) Close() error {
	return db.DB.Close()
}

func (db *Database) SaveRecord(ctx context.Context, rec *models.Record) error {
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

func (db *Database) SaveRecordsBatch(ctx context.Context, records []models.Record) error {
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
