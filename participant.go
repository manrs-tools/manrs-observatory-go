package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) FindParticipant(ctx context.Context, participantID string) (*model.Participant, error) {
	var response struct {
		Participant *model.Participant `json:"participant"`
	}

	if _, err := c.GET(ctx, "participants/"+participantID, &response); err != nil {
		return nil, fmt.Errorf("find participant: %w", err)
	}

	return response.Participant, nil
}

func (c *Client) FindParticipants(ctx context.Context, filter *model.ParticipantFilter) ([]*model.Participant, error) {
	query := structToQueryParams(filter)

	var response struct {
		Participants []*model.Participant `json:"participants"`
	}

	if _, err := c.GET(ctx, "participants?"+query, &response); err != nil {
		return nil, fmt.Errorf("find participants: %w", err)
	}

	return response.Participants, nil
}

func (c *Client) FindParticipantsPaged(ctx context.Context, filter *model.ParticipantFilter) *Iterator[*model.Participant] {
	return &Iterator[*model.Participant]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.Participant, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				Participants  []*model.Participant `json:"participants"`
				NextPageToken string               `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "participants?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find participants: %w", err)
			}

			return resp.Participants, resp.NextPageToken, nil
		},
	}
}

// GetParticipantBadge returns the participant's "MANRS Participant" badge as an SVG document.
func (c *Client) GetParticipantBadge(ctx context.Context, participantID string, filter *model.ParticipantBadgeFilter) ([]byte, error) {
	query := structToQueryParams(filter)

	badge, err := c.GET(ctx, "participants/"+participantID+"/badge?"+query, nil)
	if err != nil {
		return nil, fmt.Errorf("get participant badge: %w", err)
	}

	return badge, nil
}
