// Package connections persists the set of bank connection item_ids the server
// has seen, so list_connections/refresh_connections can enumerate them without
// a manually-maintained accounts.json. The Pluggy API has no "list all items"
// endpoint, so the server remembers each item_id the first time any tool
// resolves it successfully.
package connections

import (
	"encoding/json"
	"os"
	"slices"
	"time"
)

const itemsKey = "connections:items"

// Store is the encrypted persistence backing the registry. It is satisfied by
// sqlite.Store (the same store used for Pluggy credentials and tokens).
type Store interface {
	Get(key string) (string, bool, error)
	Set(key, value string, ttl time.Duration) error
}

// Registry is a store-backed set of connection item_ids. Safe for the single
// local writer the sqlite store already serializes.
type Registry struct {
	store Store
}

func NewRegistry(store Store) *Registry {
	return &Registry{store: store}
}

// Register records itemID if not already present. Empty ids are ignored and
// duplicates are a no-op, so callers can register on every successful item
// resolution without worrying about growth or ordering.
func (r *Registry) Register(itemID string) error {
	if itemID == "" {
		return nil
	}

	ids, err := r.List()
	if err != nil {
		return err
	}
	if slices.Contains(ids, itemID) {
		return nil
	}

	ids = append(ids, itemID)
	blob, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	// ttl 0 => never expires.
	return r.store.Set(itemsKey, string(blob), 0)
}

// List returns the registered item_ids. A missing key yields an empty slice
// (not an error): a fresh install simply has no connections yet.
func (r *Registry) List() ([]string, error) {
	raw, ok, err := r.store.Get(itemsKey)
	if err != nil {
		return nil, err
	}
	if !ok || raw == "" {
		return nil, nil
	}

	var ids []string
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// ImportLegacy back-compat: reads a legacy accounts.json (the old manual
// registry) if present and registers every item_id it lists. A missing file is
// a no-op; the file is never modified or removed.
func (r *Registry) ImportLegacy(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var entries []struct {
		ItemID string `json:"item_id"`
	}
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	for _, e := range entries {
		if err := r.Register(e.ItemID); err != nil {
			return err
		}
	}
	return nil
}
