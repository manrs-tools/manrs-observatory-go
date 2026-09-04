package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) CreateIncidentFeedback(ctx context.Context, feedback *model.IncidentFeedback) (*model.IncidentFeedback, error) {
	body, err := c.marshal(map[string]*model.IncidentFeedback{"feedback": feedback})
	if err != nil {
		return nil, fmt.Errorf("marshal incident feedback: %w", err)
	}

	var response struct {
		Feedback *model.IncidentFeedback `json:"feedback"`
	}

	if _, err = c.POST(ctx, "incident-feedback", body, &response); err != nil {
		return nil, fmt.Errorf("create incident feedback: %w", err)
	}

	return response.Feedback, nil
}

func (c *Client) DeleteIncidentFeedback(ctx context.Context, id string) error {
	if _, err := c.DELETE(ctx, "incident-feedback/"+id, nil); err != nil {
		return fmt.Errorf("delete incident feedback: %w", err)
	}

	return nil
}

func (c *Client) FindIncidentFeedback(ctx context.Context, id string) (*model.IncidentFeedback, error) {
	var response struct {
		Feedback *model.IncidentFeedback `json:"feedback"`
	}

	if _, err := c.GET(ctx, "incident-feedback/"+id, &response); err != nil {
		return nil, fmt.Errorf("find incident feedback: %w", err)
	}

	return response.Feedback, nil
}

func (c *Client) FindIncidentFeedbacks(ctx context.Context, filter *model.IncidentFeedbackFilter) ([]*model.IncidentFeedback, error) {
	query := structToQueryParams(filter)

	var response struct {
		Feedbacks []*model.IncidentFeedback `json:"feedbacks"`
	}

	if _, err := c.GET(ctx, "incident-feedback?"+query, &response); err != nil {
		return nil, fmt.Errorf("find incident feedbacks: %w", err)
	}

	return response.Feedbacks, nil
}

func (c *Client) FindIncidentFeedbacksPaged(ctx context.Context, filter *model.IncidentFeedbackFilter) *Iterator[*model.IncidentFeedback] {
	return &Iterator[*model.IncidentFeedback]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.IncidentFeedback, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				Feedbacks     []*model.IncidentFeedback `json:"feedbacks"`
				NextPageToken string                    `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "incident-feedback?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find incident feedbacks: %w", err)
			}

			return resp.Feedbacks, resp.NextPageToken, nil
		},
	}
}

func (c *Client) GetIncidentSummary(ctx context.Context, filter *model.IncidentSummaryFilter) ([]*model.IncidentSummaryRow, error) {
	query := structToQueryParams(filter)

	var response struct {
		ReportMonth string                      `json:"reportMonth"`
		Incidents   []*model.IncidentSummaryRow `json:"incidents"`
	}

	if _, err := c.GET(ctx, "incidents/summary?"+query, &response); err != nil {
		return nil, fmt.Errorf("get incident summary: %w", err)
	}

	return response.Incidents, nil
}
