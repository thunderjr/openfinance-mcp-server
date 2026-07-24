package connections

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/secret"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/sqlite"
)

// openStore opens a real encrypted store at path with a fixed test key, so
// tests exercise the same persistence the app uses (mirrors store_test.go).
func openStore(t *testing.T, path string) *sqlite.Store {
	t.Helper()
	key := make([]byte, secret.KeySize)
	for i := range key {
		key[i] = byte(i)
	}
	s, err := sqlite.Open(path, key)
	if err != nil {
		t.Fatalf("sqlite.Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func newRegistry(t *testing.T) *Registry {
	t.Helper()
	return NewRegistry(openStore(t, filepath.Join(t.TempDir(), "test.db")))
}

func contains(ids []string, want string) bool {
	for _, id := range ids {
		if id == want {
			return true
		}
	}
	return false
}

func TestListEmptyWhenNothingRegistered(t *testing.T) {
	r := newRegistry(t)
	ids, err := r.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected empty, got %v", ids)
	}
}

func TestRegisterThenList(t *testing.T) {
	r := newRegistry(t)
	if err := r.Register("item-a"); err != nil {
		t.Fatalf("Register: %v", err)
	}
	ids, err := r.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if !contains(ids, "item-a") {
		t.Fatalf("expected item-a in %v", ids)
	}
}

func TestRegisterDedupes(t *testing.T) {
	r := newRegistry(t)
	_ = r.Register("item-a")
	_ = r.Register("item-a")
	ids, _ := r.List()
	if len(ids) != 1 {
		t.Fatalf("expected 1 entry after duplicate register, got %v", ids)
	}
}

func TestRegisterMultiplePreservesAll(t *testing.T) {
	r := newRegistry(t)
	_ = r.Register("item-a")
	_ = r.Register("item-b")
	ids, _ := r.List()
	if !contains(ids, "item-a") || !contains(ids, "item-b") {
		t.Fatalf("expected both items, got %v", ids)
	}
}

func TestRegisterEmptyIsNoop(t *testing.T) {
	r := newRegistry(t)
	if err := r.Register(""); err != nil {
		t.Fatalf("Register empty: %v", err)
	}
	ids, _ := r.List()
	if len(ids) != 0 {
		t.Fatalf("empty id should not be registered, got %v", ids)
	}
}

func TestImportLegacyRegistersItems(t *testing.T) {
	r := newRegistry(t)
	path := filepath.Join(t.TempDir(), "accounts.json")
	body := `[{"name":"Inter","item_id":"item-a","account_id":"acc-1"},
	          {"name":"Inter","item_id":"item-a","account_id":"acc-2"},
	          {"name":"MeuPluggy","item_id":"item-b","account_id":"acc-3"}]`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write legacy file: %v", err)
	}

	if err := r.ImportLegacy(path); err != nil {
		t.Fatalf("ImportLegacy: %v", err)
	}

	ids, _ := r.List()
	if len(ids) != 2 || !contains(ids, "item-a") || !contains(ids, "item-b") {
		t.Fatalf("expected item-a and item-b once each, got %v", ids)
	}
}

func TestImportLegacyMissingFileIsNoop(t *testing.T) {
	r := newRegistry(t)
	if err := r.ImportLegacy(filepath.Join(t.TempDir(), "does-not-exist.json")); err != nil {
		t.Fatalf("missing file should be a no-op, got %v", err)
	}
	if ids, _ := r.List(); len(ids) != 0 {
		t.Fatalf("expected empty, got %v", ids)
	}
}

func TestPersistsAcrossReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	r1 := NewRegistry(openStore(t, path))
	if err := r1.Register("item-a"); err != nil {
		t.Fatalf("Register: %v", err)
	}

	// Reopen the store at the same path with the same key: the registration
	// must survive, proving it lives in the encrypted store, not memory.
	r2 := NewRegistry(openStore(t, path))
	ids, err := r2.List()
	if err != nil {
		t.Fatalf("List after reopen: %v", err)
	}
	if !contains(ids, "item-a") {
		t.Fatalf("registration did not persist, got %v", ids)
	}
}
