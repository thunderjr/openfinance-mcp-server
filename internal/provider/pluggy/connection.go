package pluggy

import "context"

// ConnectionAccount is an exported, caller-facing view of an account under a
// connection. It intentionally exposes only what the connection-listing tools
// need, keeping the internal account struct unexported.
type ConnectionAccount struct {
	ID   string
	Name string
	Type string
}

// ConnectionInfo is a resolved connection: the connector's display name, its
// current sync status, and its accounts. It is the single value that hides the
// multi-call GetItem+GetAccounts orchestration behind one result (a Facade over
// the Pluggy item subsystem), so callers depend on this instead of the raw,
// unexported item/account types.
type ConnectionInfo struct {
	ItemID        string
	ConnectorName string
	Status        string
	Accounts      []ConnectionAccount
}

// Connection resolves a single connection: it fetches the item (for connector
// name and status) and its accounts, in one call.
func (c *Client) Connection(itemID string) (ConnectionInfo, error) {
	item, err := c.GetItem(itemID)
	if err != nil {
		return ConnectionInfo{}, err
	}

	accounts, err := c.GetAccounts(itemID)
	if err != nil {
		return ConnectionInfo{}, err
	}

	info := ConnectionInfo{
		ItemID:        itemID,
		ConnectorName: item.Connector.Name,
		Status:        item.Status,
	}
	for _, a := range accounts.Results {
		info.Accounts = append(info.Accounts, ConnectionAccount{
			ID:   a.ID,
			Name: a.Name,
			Type: a.Type,
		})
	}
	return info, nil
}

// RefreshItem triggers a sync, waits for it to finish, and returns the item's
// resulting status. It does not fetch accounts, keeping refresh's cost the same
// as before (update + wait + get item).
func (c *Client) RefreshItem(ctx context.Context, itemID string) (ConnectionInfo, error) {
	if _, err := c.UpdateItem(itemID); err != nil {
		return ConnectionInfo{}, err
	}
	if err := c.WaitUpdated(ctx, itemID); err != nil {
		return ConnectionInfo{}, err
	}

	item, err := c.GetItem(itemID)
	if err != nil {
		return ConnectionInfo{}, err
	}
	return ConnectionInfo{
		ItemID:        itemID,
		ConnectorName: item.Connector.Name,
		Status:        item.Status,
	}, nil
}
