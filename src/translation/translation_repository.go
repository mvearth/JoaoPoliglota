package translation

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"joao_poliglota/config"
)

// Repository provides access to the translations store. It wraps a *sql.DB,
// which is a connection pool meant to be created once and shared.
type Repository struct {
	db *sql.DB
}

// NewRepository wires a Repository around an existing connection pool.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Connect opens a PostgreSQL connection pool using configuration from the
// environment. The returned *sql.DB is safe for concurrent use and should be
// closed by the caller on shutdown. It does not establish a connection until
// one is needed; use Ping to verify connectivity.
func Connect(cfg config.DB) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	return db, nil
}

// Insert stores a translation and reports whether it was persisted.
func (r *Repository) Insert(ctx context.Context, t Translation) (bool, error) {
	const query = `INSERT INTO translations (idiom, standard_key, translation)
		VALUES ($1, $2, $3) RETURNING translation_id`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, t.Idiom, t.StandardKey, t.Translation); err != nil {
		return false, err
	}
	return true, nil
}

// Get returns the cached translation for the given key and idiom. When no row
// matches it returns a zero-value Translation and a nil error, so callers
// should check whether StandardKey was populated.
func (r *Repository) Get(ctx context.Context, standardKey, idiom string) (Translation, error) {
	const query = `SELECT translation_id, idiom, standard_key, translation
		FROM translations WHERE idiom = $1 AND standard_key = $2`

	var t Translation
	err := r.db.QueryRowContext(ctx, query, idiom, standardKey).
		Scan(&t.ID, &t.Idiom, &t.StandardKey, &t.Translation)
	if err == sql.ErrNoRows {
		return Translation{}, nil
	}
	if err != nil {
		return Translation{}, err
	}
	return t, nil
}
