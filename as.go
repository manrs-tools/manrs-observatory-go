package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) FindAS(ctx context.Context, asn uint32) (*model.AS, error) {
	var response struct {
		AS *model.AS `json:"as"`
	}

	if _, err := c.GET(ctx, "ases/"+strconv.FormatUint(uint64(asn), 10), &response); err != nil {
		return nil, fmt.Errorf("find as: %w", err)
	}

	return response.AS, nil
}

func (c *Client) FindASes(ctx context.Context, filter *model.ASFilter) ([]*model.AS, error) {
	query := structToQueryParams(filter)

	var response struct {
		ASes []*model.AS `json:"ases"`
	}

	if _, err := c.GET(ctx, "ases?"+query, &response); err != nil {
		return nil, fmt.Errorf("find ases: %w", err)
	}

	return response.ASes, nil
}

func (c *Client) FindASsesPaged(ctx context.Context, filter *model.ASFilter) *Iterator[*model.AS] {
	return &Iterator[*model.AS]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.AS, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				ASes          []*model.AS `json:"ases"`
				NextPageToken string      `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "ases?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find ases: %w", err)
			}

			return resp.ASes, resp.NextPageToken, nil
		},
	}
}

func (c *Client) ListApprovedASNs(ctx context.Context) ([]int32, error) {
	var response struct {
		ASNs []int32 `json:"asns"`
	}

	if _, err := c.GET(ctx, "asns", &response); err != nil {
		return nil, fmt.Errorf("list approved asns: %w", err)
	}

	return response.ASNs, nil
}

func (c *Client) GetAllASNScores(ctx context.Context, filter *model.ScoresFilter) ([]*model.ASNScoreSummary, error) {
	query := structToQueryParams(filter)

	var response struct {
		Scores []*model.ASNScoreSummary `json:"scores"`
		Month  string                   `json:"month"`
	}

	if _, err := c.GET(ctx, "scores?"+query, &response); err != nil {
		return nil, fmt.Errorf("get scores: %w", err)
	}

	return response.Scores, nil
}

// GetASNScores returns an ASN's current MANRS readiness scores alongside the monthly history the
// endpoint reports with them.
func (c *Client) GetASNScores(ctx context.Context, filter *model.ScoresDetailsFilter) (*model.ASNScoreDetail, []*model.ASNScoreHistoryPoint, error) {
	query := structToQueryParams(filter)

	var response struct {
		Score   *model.ASNScoreDetail         `json:"score"`
		History []*model.ASNScoreHistoryPoint `json:"history"`
		Month   string                        `json:"month"`
	}

	if _, err := c.GET(ctx, "scores/details?"+query, &response); err != nil {
		return nil, nil, fmt.Errorf("get scores detailed: %w", err)
	}

	return response.Score, response.History, nil
}

func (c *Client) GetASNIncidents(ctx context.Context, filter *model.IncidentsFilter) ([]*model.IncidentRow, error) {
	query := structToQueryParams(filter)

	var response struct {
		Incidents []*model.IncidentRow `json:"incidents"`
		Month     string               `json:"month"`
	}

	if _, err := c.GET(ctx, "incidents?"+query, &response); err != nil {
		return nil, fmt.Errorf("get incidents: %w", err)
	}

	return response.Incidents, nil
}

// FindASDetails returns a page of per-AS detail rows and the total number of rows matching the
// filter, which is not bounded by the page the call returns.
func (c *Client) FindASDetails(ctx context.Context, filter *model.ASDetailsFilter) ([]*model.ASDetailRow, int64, error) {
	query := structToQueryParams(filter)

	var response struct {
		Count int64                `json:"count"`
		Data  []*model.ASDetailRow `json:"data"`
		Month string               `json:"month"`
	}

	if _, err := c.GET(ctx, "stats/details?"+query, &response); err != nil {
		return nil, 0, fmt.Errorf("find as details: %w", err)
	}

	return response.Data, response.Count, nil
}
