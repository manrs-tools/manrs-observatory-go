package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) FindROAs(ctx context.Context, filter *model.ROAFilter) ([]*model.ROA, error) {
	query := structToQueryParams(filter)

	var response struct {
		ROAs []*model.ROA `json:"roas"`
	}

	if _, err := c.GET(ctx, "roas?"+query, &response); err != nil {
		return nil, fmt.Errorf("find roas: %w", err)
	}

	return response.ROAs, nil
}

func (c *Client) FindROAsPaged(ctx context.Context, filter *model.ROAFilter) *Iterator[*model.ROA] {
	return &Iterator[*model.ROA]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.ROA, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				ROAs          []*model.ROA `json:"roas"`
				NextPageToken string       `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "roas?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find roas: %w", err)
			}

			return resp.ROAs, resp.NextPageToken, nil
		},
	}
}

func (c *Client) GetROAStats(ctx context.Context, filter *model.ROAStatsFilter) (*model.ROAStatsDetail, error) {
	query := structToQueryParams(filter)

	var response struct {
		Summary *model.ROAStatsDetail `json:"summary"`
	}

	if _, err := c.GET(ctx, "stats/roas?"+query, &response); err != nil {
		return nil, fmt.Errorf("get roa stats: %w", err)
	}

	return response.Summary, nil
}

// FindROAPrefixRecords returns per-prefix RPKI validation records, grouped by IP prefix.
func (c *Client) FindROAPrefixRecords(ctx context.Context, filter *model.ROARecordsFilter) ([]*model.ROAStatsRecord, error) {
	query := structToQueryParams(filter)

	var response struct {
		Records []*model.ROAStatsRecord `json:"records"`
	}

	if _, err := c.GET(ctx, "stats/roas/records?"+query, &response); err != nil {
		return nil, fmt.Errorf("find roa prefix records: %w", err)
	}

	return response.Records, nil
}

// FindROAPrefixRecordsPaged returns a paged iterator of per-prefix RPKI validation records, grouped by IP prefix.
func (c *Client) FindROAPrefixRecordsPaged(ctx context.Context, filter *model.ROARecordsFilter) *Iterator[*model.ROAStatsRecord] {
	return &Iterator[*model.ROAStatsRecord]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.ROAStatsRecord, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				Records       []*model.ROAStatsRecord `json:"records"`
				NextPageToken string                  `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "stats/roas/records?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find roa prefix records: %w", err)
			}

			return resp.Records, resp.NextPageToken, nil
		},
	}
}

// FindROAAnnouncementRecords returns per-announcement (prefix+ASN pair) daily RPKI validation outcomes.
func (c *Client) FindROAAnnouncementRecords(ctx context.Context, filter *model.ROAStatFilter) ([]*model.ROAStat, error) {
	query := structToQueryParams(filter)

	var response struct {
		Stats []*model.ROAStat `json:"stats"`
	}

	if _, err := c.GET(ctx, "roa-stats?"+query, &response); err != nil {
		return nil, fmt.Errorf("find roa announcement records: %w", err)
	}

	return response.Stats, nil
}

// FindROAAnnouncementRecordsPaged returns a paged iterator of per-announcement (prefix+ASN pair) daily RPKI validation outcomes.
func (c *Client) FindROAAnnouncementRecordsPaged(ctx context.Context, filter *model.ROAStatFilter) *Iterator[*model.ROAStat] {
	return &Iterator[*model.ROAStat]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.ROAStat, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				Stats         []*model.ROAStat `json:"stats"`
				NextPageToken string           `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "roa-stats?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find roa announcement records: %w", err)
			}

			return resp.Stats, resp.NextPageToken, nil
		},
	}
}
