// Package sqlite provides a local, encrypted key/value store backed by a
// single SQLite file (pure-Go driver, no CGO). It replaces the previous Redis
// cache: secret values are sealed with AES-256-GCM before they touch disk and
// entries support optional TTL expiry.
package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/secret"
)

const schema = `
CREATE TABLE IF NOT EXISTS secrets (
	key        TEXT PRIMARY KEY,
	value      BLOB NOT NULL,
	expires_at INTEGER,
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL
);`

// Store is an encrypted key/value store. It is safe for concurrent use.
type Store struct {
	db  *sql.DB
	key []byte
	now func() time.Time
}

// Open opens (creating if needed) the SQLite database at path and prepares the
// schema. key is the AES-256 master key used to seal/open stored values.
func Open(path string, key []byte) (*Store, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(on)", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlite: open: %w", err)
	}

	// Serialize access: a single local writer avoids SQLITE_BUSY and keeps the
	// store simple. Reads are fast and infrequent for this workload.
	db.SetMaxOpenConns(1)

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite: migrate: %w", err)
	}

	return &Store{db: db, key: key, now: time.Now}, nil
}

// Get returns the decrypted value for key. The boolean is false when the key is
// absent or expired (expired entries are deleted lazily on read).
func (s *Store) Get(key string) (string, bool, error) {
	var (
		blob      []byte
		expiresAt sql.NullInt64
	)
	err := s.db.QueryRow(`SELECT value, expires_at FROM secrets WHERE key = ?`, key).Scan(&blob, &expiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("sqlite: get %q: %w", key, err)
	}

	if expiresAt.Valid && s.now().Unix() >= expiresAt.Int64 {
		if err := s.Delete(key); err != nil {
			return "", false, err
		}
		return "", false, nil
	}

	plaintext, err := secret.Open(s.key, blob)
	if err != nil {
		return "", false, fmt.Errorf("sqlite: decrypt %q: %w", key, err)
	}
	return string(plaintext), true, nil
}

// Set stores value (encrypted) under key. A ttl <= 0 means the entry never
// expires.
func (s *Store) Set(key, value string, ttl time.Duration) error {
	blob, err := secret.Seal(s.key, []byte(value))
	if err != nil {
		return fmt.Errorf("sqlite: encrypt %q: %w", key, err)
	}

	now := s.now().Unix()
	var expiresAt sql.NullInt64
	if ttl > 0 {
		expiresAt = sql.NullInt64{Int64: s.now().Add(ttl).Unix(), Valid: true}
	}

	_, err = s.db.Exec(`
		INSERT INTO secrets (key, value, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			expires_at = excluded.expires_at,
			updated_at = excluded.updated_at`,
		key, blob, expiresAt, now, now)
	if err != nil {
		return fmt.Errorf("sqlite: set %q: %w", key, err)
	}
	return nil
}

// Delete removes key if present.
func (s *Store) Delete(key string) error {
	if _, err := s.db.Exec(`DELETE FROM secrets WHERE key = ?`, key); err != nil {
		return fmt.Errorf("sqlite: delete %q: %w", key, err)
	}
	return nil
}

// Close closes the underlying database.
func (s *Store) Close() error {
	return s.db.Close()
}
