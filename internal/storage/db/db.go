package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

type Database struct {
	dbConnData string
	DB         *sql.DB
	mode       int
	links      map[string]string
}

func NewDB(dbConnData string, mode int) (*Database, error) {
	db, err := sql.Open("pgx", dbConnData)
	if err != nil {
		return nil, err
	}

	newDB := &Database{
		dbConnData: dbConnData,
		DB:         db,
		mode:       mode,
		links:      make(map[string]string),
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

func (db *Database) Add(key, value, cookie string) error {
	id := uuid.NewString()
	rec := models.Record{
		UUID:        id,
		ShortURL:    key,
		OriginalURL: value,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := db.FindUserByCookie(ctx, cookie)
	if err != nil {
		return err
	}

	err = db.SaveRecord(ctx, &rec, user.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return models.ErrConflict
		}

		logger.Log.Error("error while writing data to db", zap.Error(err))
		return err
	}

	db.links[key] = value
	return nil
}

func (db *Database) AddBatch(ctx context.Context, records []models.Record) error {
	err := db.SaveRecordsBatch(ctx, records)
	if err != nil {
		logger.Log.Error("error while writing data batch to db", zap.Error(err))
	}

	for _, rec := range records {
		db.links[rec.ShortURL] = rec.OriginalURL
	}

	return nil
}

func (db *Database) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	rec, err := db.FindRecord(ctx, key)
	if err != nil {
		return "", fmt.Errorf("URL with the key \"%s\" is missing", key)
	}

	if rec.DeletedFlag {
		return "", models.ErrDeleted
	}

	return rec.OriginalURL, nil
}

func (db *Database) GetMode() int {
	return db.mode
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
		`SELECT uuid, short_url, origin_url, is_deleted FROM urls WHERE short_url=$1 LIMIT 1`, value)

	var rec models.Record
	err := row.Scan(&rec.UUID, &rec.ShortURL, &rec.OriginalURL, &rec.DeletedFlag)
	if err != nil {
		return rec, err
	}

	return rec, err
}

func (db *Database) FindRecordByOriginURL(ctx context.Context, value string) (models.Record, error) {
	row := db.DB.QueryRowContext(ctx,
		`SELECT uuid, short_url, origin_url, is_deleted FROM urls WHERE origin_url=$1 LIMIT 1`, value)

	var rec models.Record
	err := row.Scan(&rec.UUID, &rec.ShortURL, &rec.OriginalURL, &rec.DeletedFlag)
	if err != nil {
		return rec, err
	}

	return rec, err
}

func (db *Database) Close() error {
	return db.DB.Close()
}

func (db *Database) SaveRecord(ctx context.Context, rec *models.Record, userID int) error {
	_, err := db.DB.ExecContext(ctx,
		`INSERT INTO urls(uuid, short_url, origin_url, user_id) VALUES($1, $2, $3, $4)`,
		rec.UUID, rec.ShortURL, rec.OriginalURL, userID)
	return err
}

func (db *Database) CreateDBScheme() error {
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()

	_, err := db.DB.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS users(
												"id" SERIAL PRIMARY KEY,
												"cookie" VARCHAR)`)
	if err != nil {
		return err
	}

	_, err = db.DB.ExecContext(ctx,
		`CREATE UNIQUE INDEX IF NOT EXISTS cookie_idx on users(cookie)`)
	if err != nil {
		return err
	}

	_, err = db.DB.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS urls(
												"uuid" VARCHAR,
												"short_url" VARCHAR,
												"origin_url" VARCHAR,
												"user_id" INTEGER,
												"is_deleted" BOOLEAN DEFAULT false)`)
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

func (db *Database) CloseStorage() error {
	return db.DB.Close()
}

func (db *Database) FindRecordsByUserID(ctx context.Context, userID int) (records []models.Record, err error) {
	rows, err := db.DB.QueryContext(ctx,
		"SELECT uuid, short_url, origin_url, is_deleted FROM urls WHERE user_id=$1", userID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rec models.Record
		err = rows.Scan(&rec.UUID, &rec.ShortURL, &rec.OriginalURL, &rec.DeletedFlag)
		if err != nil {
			return
		}

		records = append(records, rec)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

func (db *Database) FindUserByCookie(ctx context.Context, cookie string) (*models.User, error) {
	row := db.DB.QueryRowContext(ctx,
		"SELECT id, cookie FROM users WHERE cookie=$1 LIMIT 1", cookie)

	var user models.User
	err := row.Scan(&user.UserID, &user.Cookie)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (db *Database) GetUserRecords(ctx context.Context, cookie string) ([]models.Record, error) {
	user, err := db.FindUserByCookie(ctx, cookie)
	if err != nil {
		return nil, err
	}

	records, err := db.FindRecordsByUserID(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (db *Database) FindUserByID(ctx context.Context, userID int) (*models.User, error) {
	row := db.DB.QueryRowContext(ctx,
		"SELECT id, cookie FROM users WHERE id=$1 LIMIT 1", userID)

	var user models.User
	err := row.Scan(&user.UserID, &user.Cookie)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (db *Database) CreateUser(ctx context.Context) (*models.User, error) {
	_, err := db.DB.ExecContext(ctx, `INSERT INTO users DEFAULT VALUES`)
	if err != nil {
		return nil, err
	}

	row := db.DB.QueryRowContext(ctx,
		"SELECT id FROM users ORDER BY id DESC LIMIT 1")
	var user models.User
	err = row.Scan(&user.UserID)
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (db *Database) UpdateUser(ctx context.Context, id int, cookie string) error {
	_, err := db.DB.ExecContext(ctx, `UPDATE users SET cookie=$1 WHERE id=$2`, cookie, id)
	return err
}
