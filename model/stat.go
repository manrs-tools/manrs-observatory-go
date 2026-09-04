// Code generated from internal models; DO NOT EDIT.
package model

import (
	"time"
)

// DatesAvailable holds the list of months that have report data, plus which are still incomplete.
type DatesAvailable struct {
	DataRange        []string `json:"dataRange"`
	IncompleteMonths []string `json:"incompleteMonths"`
}

// KeyFigure holds an aggregated score and severity breakdown for a single metric.
type KeyFigure struct {
	ID         string              `json:"id"`
	Value      float64             `json:"value"`
	Trend      *float64            `json:"trend,omitempty"`
	Severities []KeyFigureSeverity `json:"severities"`
}

// KeyFigureSeverity holds the count and percentage of ASNs in a given severity bucket.
type KeyFigureSeverity struct {
	ID         string  `json:"id"`
	Count      uint64  `json:"count"`
	Percentage float64 `json:"percentage"`
}

// DashboardStatistics holds aggregated report stats for the dashboard endpoint.
type DashboardStatistics struct {
	KeyFigures []KeyFigure        `json:"keyFigures"`
	Incidents  DashboardIncidents `json:"incidents"`
	Culprits   DashboardCulprits  `json:"culprits"`
	IRR        DashboardIRR       `json:"irr"`
	RPKI       DashboardRPKI      `json:"rpki"`
	ROV        DashboardROV       `json:"rov"`
}

type DashboardIncidents struct {
	RouteMisoriginations uint64 `json:"routeMisoriginations"`
	RouteLeaks           uint64 `json:"routeLeaks"`
	BogonAnnouncements   uint64 `json:"bogonAnnouncements"`
}

type DashboardCulprits struct {
	Total uint64 `json:"total"`
}

type DashboardIRR struct {
	Total        uint64 `json:"total"`
	Registered   uint64 `json:"registered"`
	Unregistered uint64 `json:"unregistered"`
}

type DashboardRPKI struct {
	Total   uint64 `json:"total"`
	Valid   uint64 `json:"valid"`
	Unknown uint64 `json:"unknown"`
	Invalid uint64 `json:"invalid"`
}

type DashboardROV struct {
	Valid float64 `json:"valid"`
}

// StatsHistory is the top-level response for the /stats/history endpoint.
type StatsHistory struct {
	Statistics []StatsHistoryStat     `json:"statistics"`
	History    []StatsHistoryMetric   `json:"history"`
	Coverage   []StatsHistoryCoverage `json:"coverage"`

	ParticipantsByType []StatsHistoryParticipantPoint `json:"participantsByType"`
}

// StatsHistoryCoverage sizes one month's AS population within the request's scope. Claims are read
// from the report snapshot rather than today's participant rows, so each month reflects who was
// claimed then. Callers wanting headline figures take the newest entry.
type StatsHistoryCoverage struct {
	MonthID           string `json:"monthID"`
	ASNs              uint64 `json:"asns"`
	ASNsWithIncidents uint64 `json:"asnsWithIncidents"`
	ClaimedASNs       uint64 `json:"claimedASNs"`
}

// StatsHistoryParticipantPoint is one month's approved-participant count for a single participant
// type within the request's scope. Unlike the global series in /stats/overview, which is a per-month
// snapshot taken at sync time, this is derived from current membership and current AS claims: a
// participant that has since left, or has since moved its claims, is counted per today's state in
// every month it was already approved for.
type StatsHistoryParticipantPoint struct {
	MonthID string `json:"monthID"`
	Type    string `json:"type"`
	Count   uint64 `json:"count"`
}

// StatsHistoryStat holds the grouped history for a set of related stats.
type StatsHistoryStat struct {
	ID        string              `json:"id"`
	Format    string              `json:"format"`
	History   []StatsHistoryPoint `json:"history"`
	Synthetic bool                `json:"synthetic"`
}

// StatsHistoryPoint is a single month's value for a metric.
type StatsHistoryPoint struct {
	ID      string  `json:"id"`
	MonthID string  `json:"monthID"`
	Value   float64 `json:"value"`
}

// StatsHistoryMetric holds the per-month severity breakdown for a single metric.
type StatsHistoryMetric struct {
	ID     string               `json:"id"`
	Public []StatsHistoryPublic `json:"public"`
}

// StatsHistoryPublic is a single severity bucket entry for a month.
type StatsHistoryPublic struct {
	MonthID    string  `json:"monthID"`
	ID         string  `json:"id"`
	Count      uint64  `json:"count"`
	Percentage float64 `json:"percentage"`
}

// ComparisonGeoFilter holds the geographic filter for a comparison selection.
type ComparisonGeoFilter struct {
	Selection []string `json:"selection"`
	Type      string   `json:"type"`
}

// ComparisonTimeRange specifies the time range for a comparison selection.
type ComparisonTimeRange struct {
	Start    string `json:"start"`
	Duration int    `json:"duration"`
}

// ComparisonSelection defines a single filter set for the comparison endpoint.
type ComparisonSelection struct {
	Geo       *ComparisonGeoFilter `json:"geo"`
	TimeRange ComparisonTimeRange  `json:"timeRange"`
}

// ComparisonScenario pairs a UUID with a selection for the comparison endpoint.
type ComparisonScenario struct {
	ID        string              `json:"id"`
	Selection ComparisonSelection `json:"selection"`
}

// ComparisonDataPoint is a single metric value within a comparison group.
type ComparisonDataPoint struct {
	Format      string   `json:"format"`
	GroupID     string   `json:"groupID"`
	Implemented bool     `json:"implemented"`
	MetricID    string   `json:"metricID"`
	MonthID     string   `json:"monthID"`
	ScenarioID  string   `json:"scenarioID"`
	Value       *float64 `json:"value"`
}

// ComparisonGroup holds a named set of data points for the comparison response.
type ComparisonGroup struct {
	Title      string                `json:"title"`
	DataPoints []ComparisonDataPoint `json:"dataPoints"`
}

type StatsHistoryFilter struct {
	ASNs   string `query:"asns,omitempty"`
	MANRS  string `query:"manrs,omitempty"`
	Region string `query:"region,omitempty"`
}

type ASDetailRow struct {
	ASN          uint32          `json:"asn"`
	Holder       string          `json:"holder"`
	Organisation string          `json:"organisation"`
	Economy      string          `json:"economy"`
	RIR          string          `json:"rir"`
	Region       string          `json:"region"`
	Scores       *CategoryScores `json:"scores"`
	Visible      bool            `json:"visible"`

	Claimed            bool `json:"claimed,omitempty"`
	HasOpenApplication bool `json:"hasOpenApplication,omitempty"`
}

type ASDetailsFilter struct {
	ASNs      []uint32  `query:"asns,omitempty"`
	Economies string    `query:"economies,omitempty"`
	MANRS     string    `query:"manrs,omitempty"`
	Month     time.Time `query:"month,omitempty"`
	Region    string    `query:"region,omitempty"`
	RIR       string    `query:"rir,omitempty"`
	Scope     string    `query:"scope,omitempty"`
	Severity  string    `query:"severity,omitempty"`
	SubRegion string    `query:"subRegion,omitempty"`
}
