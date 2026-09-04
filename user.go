package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) FindSessionUser(ctx context.Context) (*model.User, error) {
	var response struct {
		User *model.User `json:"user"`
	}

	if _, err := c.GET(ctx, "user", &response); err != nil {
		return nil, fmt.Errorf("find session user: %w", err)
	}

	return response.User, nil
}

func (c *Client) FindUser(ctx context.Context, userID string) (*model.User, error) {
	var response struct {
		User *model.User `json:"user"`
	}

	if _, err := c.GET(ctx, "users/"+userID, &response); err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return response.User, nil
}

func (c *Client) FindUsers(ctx context.Context, filter *model.UserFilter) ([]*model.User, error) {
	query := structToQueryParams(filter)

	var response struct {
		Users []*model.User `json:"users"`
	}

	if _, err := c.GET(ctx, "users?"+query, &response); err != nil {
		return nil, fmt.Errorf("find users: %w", err)
	}

	return response.Users, nil
}

func (c *Client) FindUsersPaged(ctx context.Context, filter *model.UserFilter) *Iterator[*model.User] {
	return &Iterator[*model.User]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.User, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				Users         []*model.User `json:"users"`
				NextPageToken string        `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "users?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find users: %w", err)
			}

			return resp.Users, resp.NextPageToken, nil
		},
	}
}

func (c *Client) UpdateUser(ctx context.Context, userID string, update *model.UserUpdate) (*model.User, error) {
	body, err := c.marshal(map[string]*model.UserUpdate{"user": update})
	if err != nil {
		return nil, fmt.Errorf("marshal user update: %w", err)
	}

	var response struct {
		User *model.User `json:"user"`
	}

	if _, err = c.PATCH(ctx, "users/"+userID, body, &response); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return response.User, nil
}
