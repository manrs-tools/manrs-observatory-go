package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) CreateAPIKey(ctx context.Context, apiKey *model.APIKey) (*model.APIKey, error) {
	body, err := c.marshal(map[string]*model.APIKey{"key": apiKey}) //nolint:goconst // Request body field name, unrelated to the redaction key set in http.go.
	if err != nil {
		return nil, fmt.Errorf("marshal api key: %w", err)
	}

	var response struct {
		Key *model.APIKey `json:"key"`
	}

	if _, err = c.POST(ctx, "keys", body, &response); err != nil {
		return nil, fmt.Errorf("create api key: %w", err)
	}

	return response.Key, nil
}

func (c *Client) DeleteAPIKey(ctx context.Context, keyID string) error {
	if _, err := c.DELETE(ctx, "keys/"+keyID, nil); err != nil {
		return fmt.Errorf("delete api key: %w", err)
	}

	return nil
}

func (c *Client) FindAPIKey(ctx context.Context, keyID string) (*model.APIKey, error) {
	var response struct {
		Key *model.APIKey `json:"key"`
	}

	if _, err := c.GET(ctx, "keys/"+keyID, &response); err != nil {
		return nil, fmt.Errorf("find api key: %w", err)
	}

	return response.Key, nil
}

func (c *Client) FindAPIKeys(ctx context.Context, filter *model.APIKeyFilter) ([]*model.APIKey, error) {
	query := structToQueryParams(filter)

	var response struct {
		Keys []*model.APIKey `json:"keys"`
	}

	if _, err := c.GET(ctx, "keys?"+query, &response); err != nil {
		return nil, fmt.Errorf("find api keys: %w", err)
	}

	return response.Keys, nil
}

func (c *Client) FindAPIKeysPaged(ctx context.Context, filter *model.APIKeyFilter) *Iterator[*model.APIKey] {
	return &Iterator[*model.APIKey]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.APIKey, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				Keys          []*model.APIKey `json:"keys"`
				NextPageToken string          `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "keys?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find api keys: %w", err)
			}

			return resp.Keys, resp.NextPageToken, nil
		},
	}
}

func (c *Client) UpdateAPIKey(ctx context.Context, keyID string, update *model.APIKeyUpdate) (*model.APIKey, error) {
	body, err := c.marshal(map[string]*model.APIKeyUpdate{"key": update})
	if err != nil {
		return nil, fmt.Errorf("marshal api key update: %w", err)
	}

	var response struct {
		Key *model.APIKey `json:"key"`
	}

	if _, err = c.PATCH(ctx, "keys/"+keyID, body, &response); err != nil {
		return nil, fmt.Errorf("update api key: %w", err)
	}

	return response.Key, nil
}
