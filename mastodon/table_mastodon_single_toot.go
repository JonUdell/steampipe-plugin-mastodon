package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonSingleToot() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_single_toot",
		Description: "Fetch a single toot by ID.",
		List: &plugin.ListConfig{
			Hydrate:    getSingleToot,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: commonAccountColumns(tootColumns()),
	}
}

func getSingleToot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	
	// Get the ID from the query
	id := d.EqualsQualString("id")
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	
	logger.Debug("mastodon_single_toot.getSingleToot", "toot_id", id)
	
	// Connect to Mastodon client
	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_single_toot.getSingleToot", "connect_error", err)
		return nil, err
	}
	
	// Fetch the specific toot
	status, err := client.GetStatus(ctx, mastodon.ID(id))
	if err != nil {
		logger.Error("mastodon_single_toot.getSingleToot", "query_error", err)
		return nil, err
	}
	
	// Stream the result like other tables do
	d.StreamListItem(ctx, status)
	
	return nil, nil
}
