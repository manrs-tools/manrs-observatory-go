package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) CreateFavorite(ctx context.Context, favorite *model.Favorite) (*model.Favorite, error) {
	body, err := c.marshal(map[string]*model.Favorite{"favorite": favorite}) //nolint:goconst // Request body field name, unrelated to the identically spelled response envelope keys.
	if err != nil {
		return nil, fmt.Errorf("marshal favorite: %w", err)
	}

	var response struct {
		Favorite *model.Favorite `json:"favorite"`
	}

	if _, err = c.POST(ctx, "favorites", body, &response); err != nil {
		return nil, fmt.Errorf("create favorite: %w", err)
	}

	return response.Favorite, nil
}

func (c *Client) DeleteFavorite(ctx context.Context, favoriteID string) error {
	if _, err := c.DELETE(ctx, "favorites/"+favoriteID, nil); err != nil {
		return fmt.Errorf("delete favorite: %w", err)
	}

	return nil
}

func (c *Client) FindFavorite(ctx context.Context, favoriteID string) (*model.Favorite, error) {
	var response struct {
		Favorite *model.Favorite `json:"favorite"`
	}

	if _, err := c.GET(ctx, "favorites/"+favoriteID, &response); err != nil {
		return nil, fmt.Errorf("find favorite: %w", err)
	}

	return response.Favorite, nil
}

func (c *Client) FindFavorites(ctx context.Context, filter *model.FavoriteFilter) ([]*model.Favorite, error) {
	query := structToQueryParams(filter)

	var response struct {
		Favorites []*model.Favorite `json:"favorites"`
	}

	if _, err := c.GET(ctx, "favorites?"+query, &response); err != nil {
		return nil, fmt.Errorf("find favorites: %w", err)
	}

	return response.Favorites, nil
}

func (c *Client) FindFavoritesPaged(ctx context.Context, filter *model.FavoriteFilter) *Iterator[*model.Favorite] {
	return &Iterator[*model.Favorite]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.Favorite, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				Favorites     []*model.Favorite `json:"favorites"`
				NextPageToken string            `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "favorites?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find favorites: %w", err)
			}

			return resp.Favorites, resp.NextPageToken, nil
		},
	}
}

func (c *Client) UpdateFavorite(ctx context.Context, favoriteID string, update *model.FavoriteUpdate) (*model.Favorite, error) {
	body, err := c.marshal(map[string]*model.FavoriteUpdate{"update": update})
	if err != nil {
		return nil, fmt.Errorf("marshal favorite update: %w", err)
	}

	var response struct {
		Favorite *model.Favorite `json:"favorite"`
	}

	if _, err = c.PATCH(ctx, "favorites/"+favoriteID, body, &response); err != nil {
		return nil, fmt.Errorf("update favorite: %w", err)
	}

	return response.Favorite, nil
}
