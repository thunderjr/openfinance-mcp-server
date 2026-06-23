package sqlite

import (
	"bytes"
	"path/filepath"
	"testing"
	"time"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/secret"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	key := make([]byte, secret.KeySize)
	for i := range key {
		key[i] = byte(i)
	}
	s, err := Open(filepath.Join(t.TempDir(), "test.db"), key)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestSetGetRoundTrip(t *testing.T) {
	s := newTestStore(t)
	if err := s.Set("k", "v", 0); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, ok, err := s.Get("k")
	if err != nil || !ok {
		t.Fatalf("Get: ok=%v err=%v", ok, err)
	}
	if got != "v" {
		t.Fatalf("got %q want %q", got, "v")
	}
}

func TestGetMissing(t *testing.T) {
	s := newTestStore(t)
	if _, ok, err := s.Get("nope"); ok || err != nil {
		t.Fatalf("expected miss, got ok=%v err=%v", ok, err)
	}
}

func TestSetOverwrites(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("k", "one", 0)
	_ = s.Set("k", "two", 0)
	got, _, _ := s.Get("k")
	if got != "two" {
		t.Fatalf("got %q want %q", got, "two")
	}
}

func TestTTLExpiry(t *testing.T) {
	s := newTestStore(t)
	base := time.Unix(1_700_000_000, 0)
	s.now = func() time.Time { return base }

	if err := s.Set("k", "v", time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// Just before expiry.
	s.now = func() time.Time { return base.Add(59 * time.Second) }
	if _, ok, _ := s.Get("k"); !ok {
		t.Fatal("entry expired early")
	}

	// At/after expiry -> gone and lazily deleted.
	s.now = func() time.Time { return base.Add(time.Minute) }
	if _, ok, _ := s.Get("k"); ok {
		t.Fatal("entry should be expired")
	}
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM secrets WHERE key='k'`).Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 0 {
		t.Fatal("expired entry not lazily deleted")
	}
}

func TestDelete(t *testing.T) {
	s := newTestStore(t)
	_ = s.Set("k", "v", 0)
	if err := s.Delete("k"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, ok, _ := s.Get("k"); ok {
		t.Fatal("entry not deleted")
	}
}

func TestValuesEncryptedAtRest(t *testing.T) {
	s := newTestStore(t)
	plaintext := "plaintext-secret"
	if err := s.Set("k", plaintext, 0); err != nil {
		t.Fatalf("Set: %v", err)
	}

	var raw []byte
	if err := s.db.QueryRow(`SELECT value FROM secrets WHERE key='k'`).Scan(&raw); err != nil {
		t.Fatalf("raw read: %v", err)
	}
	if bytes.Contains(raw, []byte(plaintext)) {
		t.Fatal("value stored in plaintext on disk")
	}
}

func TestWrongKeyCannotDecrypt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	good := make([]byte, secret.KeySize)
	s, err := Open(path, good)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := s.Set("k", "v", 0); err != nil {
		t.Fatalf("Set: %v", err)
	}
	s.Close()

	bad := make([]byte, secret.KeySize)
	bad[0] = 1
	s2, err := Open(path, bad)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer s2.Close()
	if _, _, err := s2.Get("k"); err == nil {
		t.Fatal("expected decrypt error with wrong key")
	}
}
