package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) GetDashboardFilters(ctx context.Context) (*model.DatesAvailable, error) {
	var response struct {
		DatesAvailable *model.DatesAvailable `json:"datesAvailable"`
	}

	if _, err := c.GET(ctx, "filters", &response); err != nil {
		return nil, fmt.Errorf("get dashboard filters: %w", err)
	}

	return response.DatesAvailable, nil
}

func (c *Client) GetScoresSummary(ctx context.Context, filter *model.ScoresSummaryFilter) (*model.DashboardStatistics, error) {
	query := structToQueryParams(filter)

	var response struct {
		Scores *model.DashboardStatistics `json:"scores"`
	}

	if _, err := c.GET(ctx, "scores/summary?"+query, &response); err != nil {
		return nil, fmt.Errorf("get scores summary: %w", err)
	}

	return response.Scores, nil
}

func (c *Client) GetStatsComparison(ctx context.Context, selection model.ComparisonSelection, scenarios []model.ComparisonScenario) ([]model.ComparisonGroup, error) {
	body, err := c.marshal(map[string]any{"selection": selection, "scenarios": scenarios})
	if err != nil {
		return nil, fmt.Errorf("marshal stats comparison request: %w", err)
	}

	var response struct {
		Groups []model.ComparisonGroup `json:"groups"`
	}

	if _, err = c.POST(ctx, "stats/comparison", body, &response); err != nil {
		return nil, fmt.Errorf("get stats comparison: %w", err)
	}

	return response.Groups, nil
}

func (c *Client) GetStatsHistory(ctx context.Context, filter *model.StatsHistoryFilter) (*model.StatsHistory, error) {
	query := structToQueryParams(filter)

	var response struct {
		History *model.StatsHistory `json:"history"`
	}

	if _, err := c.GET(ctx, "stats/history?"+query, &response); err != nil {
		return nil, fmt.Errorf("get stats history: %w", err)
	}

	return response.History, nil
}

func (c *Client) GetStatsOverview(ctx context.Context) (*model.StatsOverview, error) {
	var response struct {
		Overview *model.StatsOverview `json:"overview"`
	}

	if _, err := c.GET(ctx, "stats/overview", &response); err != nil {
		return nil, fmt.Errorf("get stats overview: %w", err)
	}

	return response.Overview, nil
}

func (c *Client) GetConformanceScores(ctx context.Context, filter *model.ConformanceFilter) (map[string]any, error) {
	query := structToQueryParams(filter)

	var response map[string]any

	if _, err := c.GET(ctx, "scores/conformance?"+query, &response); err != nil {
		return nil, fmt.Errorf("get conformance scores: %w", err)
	}

	return response, nil
}
