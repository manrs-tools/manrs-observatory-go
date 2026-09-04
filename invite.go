package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) RedeemInvite(ctx context.Context, token string, password string) (*model.User, error) {
	body, err := c.marshal(map[string]string{"password": password}) //nolint:goconst // Request body field name, unrelated to the redaction key set in http.go.
	if err != nil {
		return nil, fmt.Errorf("marshal redeem invite: %w", err)
	}

	var response struct {
		User *model.User `json:"user"`
	}

	if _, err = c.POST(ctx, "invites/"+token+"/redeem", body, &response); err != nil {
		return nil, fmt.Errorf("redeem invite: %w", err)
	}

	return response.User, nil
}
