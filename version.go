package client

import (
	"context"
	"fmt"
)

func (c *Client) GetVersion(ctx context.Context) (string, error) {
	var response struct {
		Version string `json:"version"`
	}

	if _, err := c.GET(ctx, "version", &response); err != nil {
		return "", fmt.Errorf("get version: %w", err)
	}

	return response.Version, nil
}
