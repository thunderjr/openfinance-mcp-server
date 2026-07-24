package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/connections"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/secret"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/sqlite"
)

// fakeGateway stands in for the Pluggy network boundary (the only mock). Tests
// assert on the JSON the tools return, never on how many times this is called.
type fakeGateway struct {
	conns      map[string]pluggy.ConnectionInfo
	connErr    map[string]error
	refresh    map[string]pluggy.ConnectionInfo
	refreshErr map[string]error
}

func (f *fakeGateway) Connection(itemID string) (pluggy.ConnectionInfo, error) {
	if err := f.connErr[itemID]; err != nil {
		return pluggy.ConnectionInfo{}, err
	}
	info, ok := f.conns[itemID]
	if !ok {
		return pluggy.ConnectionInfo{}, fmt.Errorf("no connection for %s", itemID)
	}
	return info, nil
}

func (f *fakeGateway) RefreshItem(_ context.Context, itemID string) (pluggy.ConnectionInfo, error) {
	if err := f.refreshErr[itemID]; err != nil {
		return pluggy.ConnectionInfo{}, err
	}
	info, ok := f.refresh[itemID]
	if !ok {
		return pluggy.ConnectionInfo{}, fmt.Errorf("no refresh for %s", itemID)
	}
	return info, nil
}

// newRegistry builds a Registry over a real encrypted store (temp DB), so the
// tests exercise the same persistence the running server uses.
func newRegistry(t *testing.T) *connections.Registry {
	t.Helper()
	key := make([]byte, secret.KeySize)
	for i := range key {
		key[i] = byte(i)
	}
	s, err := sqlite.Open(filepath.Join(t.TempDir(), "test.db"), key)
	if err != nil {
		t.Fatalf("sqlite.Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return connections.NewRegistry(s)
}

func decodeResponse(t *testing.T, resp *mcp.ToolResponse, v any) {
	t.Helper()
	if resp == nil || len(resp.Content) == 0 || resp.Content[0].TextContent == nil {
		t.Fatalf("response has no text content: %+v", resp)
	}
	text := resp.Content[0].TextContent.Text
	if err := json.Unmarshal([]byte(text), v); err != nil {
		t.Fatalf("unmarshal %q: %v", text, err)
	}
}

func listConnections(t *testing.T, reg *connections.Registry, g connectionGateway) []Connection {
	t.Helper()
	resp, err := NewListConnectionsTool(reg, g).handleListConnections(ListConnectionsArgs{})
	if err != nil {
		t.Fatalf("handleListConnections: %v", err)
	}
	var rows []Connection
	decodeResponse(t, resp, &rows)
	return rows
}

func rowFor(rows []Connection, accountID string) (Connection, bool) {
	for _, r := range rows {
		if r.AccountID == accountID {
			return r, true
		}
	}
	return Connection{}, false
}

// AT1: nothing ever registered -> empty listing.
func TestListConnectionsEmptyWhenNothingRegistered(t *testing.T) {
	rows := listConnections(t, newRegistry(t), &fakeGateway{})
	if len(rows) != 0 {
		t.Fatalf("expected empty listing, got %v", rows)
	}
}

// AT2: once an item is registered (what every resolving tool does on success),
// it becomes listable with one row per account and live name/type, no file.
func TestConnectionListableAfterRegister(t *testing.T) {
	reg := newRegistry(t)
	g := &fakeGateway{conns: map[string]pluggy.ConnectionInfo{
		"item-1": {
			ItemID:        "item-1",
			ConnectorName: "MeuPluggy",
			Status:        "UPDATED",
			Accounts: []pluggy.ConnectionAccount{
				{ID: "acc-cc", Name: "Conta Corrente", Type: "BANK"},
				{ID: "acc-card", Name: "Cartão", Type: "CREDIT"},
			},
		},
	}}

	if err := reg.Register("item-1"); err != nil {
		t.Fatalf("Register: %v", err)
	}

	rows := listConnections(t, reg, g)
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows (one per account), got %v", rows)
	}

	cc, ok := rowFor(rows, "acc-cc")
	if !ok {
		t.Fatalf("checking account row missing: %v", rows)
	}
	if cc.ConnectionID != "item-1" || cc.AccountType != "BANK" || cc.Name != "MeuPluggy - Conta Corrente" {
		t.Fatalf("unexpected checking row: %+v", cc)
	}
	card, ok := rowFor(rows, "acc-card")
	if !ok {
		t.Fatalf("card account row missing: %v", rows)
	}
	if card.AccountType != "CREDIT" || card.Name != "MeuPluggy - Cartão" {
		t.Fatalf("unexpected card row: %+v", card)
	}
}

// AT3: the listing reflects live Pluggy data on each call, not a stored snapshot.
func TestListReflectsLiveData(t *testing.T) {
	reg := newRegistry(t)
	g := &fakeGateway{conns: map[string]pluggy.ConnectionInfo{
		"item-1": {
			ConnectorName: "MeuPluggy",
			Accounts:      []pluggy.ConnectionAccount{{ID: "acc-1", Name: "Conta", Type: "BANK"}},
		},
	}}
	if err := reg.Register("item-1"); err != nil {
		t.Fatalf("Register: %v", err)
	}

	if rows := listConnections(t, reg, g); len(rows) != 1 {
		t.Fatalf("expected 1 row initially, got %v", rows)
	}

	// Upstream changes: connector renamed, an account added.
	g.conns["item-1"] = pluggy.ConnectionInfo{
		ConnectorName: "Inter",
		Accounts: []pluggy.ConnectionAccount{
			{ID: "acc-1", Name: "Conta", Type: "BANK"},
			{ID: "acc-2", Name: "Poupança", Type: "BANK"},
		},
	}

	rows := listConnections(t, reg, g)
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows after upstream change, got %v", rows)
	}
	if r, ok := rowFor(rows, "acc-2"); !ok || r.Name != "Inter - Poupança" {
		t.Fatalf("live data not reflected: %v", rows)
	}
}

// AT4: a registered item Pluggy can no longer resolve is skipped, not fatal.
func TestBrokenItemSkippedNotFatal(t *testing.T) {
	reg := newRegistry(t)
	g := &fakeGateway{
		conns: map[string]pluggy.ConnectionInfo{
			"good": {ConnectorName: "Inter", Accounts: []pluggy.ConnectionAccount{{ID: "acc-g", Name: "Conta", Type: "BANK"}}},
		},
		connErr: map[string]error{"bad": errors.New("item not found")},
	}
	if err := reg.Register("good"); err != nil {
		t.Fatalf("Register good: %v", err)
	}
	if err := reg.Register("bad"); err != nil {
		t.Fatalf("Register bad: %v", err)
	}

	rows := listConnections(t, reg, g)
	if len(rows) != 1 {
		t.Fatalf("expected only the healthy item, got %v", rows)
	}
	if rows[0].ConnectionID != "good" {
		t.Fatalf("expected 'good' connection, got %+v", rows[0])
	}
}

// AT5: refresh acts on every registered item and reports per-item status;
// registering the same item twice yields a single entry (idempotent).
func TestRefreshRefreshesRegisteredItemsIdempotently(t *testing.T) {
	reg := newRegistry(t)
	g := &fakeGateway{refresh: map[string]pluggy.ConnectionInfo{
		"item-1": {Status: "UPDATED", ConnectorName: "Inter"},
		"item-2": {Status: "UPDATED", ConnectorName: "MeuPluggy"},
	}}
	for _, id := range []string{"item-1", "item-1", "item-2"} {
		if err := reg.Register(id); err != nil {
			t.Fatalf("Register %s: %v", id, err)
		}
	}

	resp, err := NewRefreshConnectionsTool(reg, g).handleRefreshConnections(RefreshConnectionsArgs{})
	if err != nil {
		t.Fatalf("handleRefreshConnections: %v", err)
	}
	var results []refreshResult
	decodeResponse(t, resp, &results)

	byItem := map[string]int{}
	for _, r := range results {
		byItem[r.ItemID]++
		if r.Status != "UPDATED" {
			t.Fatalf("item %s not refreshed: %+v", r.ItemID, r)
		}
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 distinct results, got %v", results)
	}
	if byItem["item-1"] != 1 {
		t.Fatalf("idempotency broken: item-1 appears %d times", byItem["item-1"])
	}
}
