package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) FindSentEmail(ctx context.Context, id string) (*model.Email, error) {
	var response struct {
		Email *model.Email `json:"email"`
	}

	if _, err := c.GET(ctx, "emails/sent/"+id, &response); err != nil {
		return nil, fmt.Errorf("find sent email: %w", err)
	}

	return response.Email, nil
}

func (c *Client) FindSentEmails(ctx context.Context, filter *model.EmailFilter) ([]*model.Email, error) {
	query := structToQueryParams(filter)

	var response struct {
		Emails []*model.Email `json:"emails"`
	}

	if _, err := c.GET(ctx, "emails/sent?"+query, &response); err != nil {
		return nil, fmt.Errorf("find sent emails: %w", err)
	}

	return response.Emails, nil
}

func (c *Client) FindSentEmailsPaged(ctx context.Context, filter *model.EmailFilter) *Iterator[*model.Email] {
	return &Iterator[*model.Email]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.Email, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				Emails        []*model.Email `json:"emails"`
				NextPageToken string         `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "emails/sent?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find sent emails: %w", err)
			}

			return resp.Emails, resp.NextPageToken, nil
		},
	}
}
