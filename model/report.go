// Code generated from internal models; DO NOT EDIT.
package model

import (
	"net/netip"
	"time"
)

// ScoreWithSeverity represents a single score metric with its severity classification.
type ScoreWithSeverity struct {
	Value    any    `json:"value"`
	Severity string `json:"severity"`
}

// CategoryScores represents the five MANRS readiness score categories.
type CategoryScores struct {
	Filtering              *ScoreWithSeverity `json:"filtering"`
	Coordination           *ScoreWithSeverity `json:"coordination"`
	AntiSpoofing           *ScoreWithSeverity `json:"antiSpoofing"`
	RoutingInformationIRR  *ScoreWithSeverity `json:"routingInformationIRR"`
	RoutingInformationRPKI *ScoreWithSeverity `json:"routingInformationRPKI"`
}

// ASNActionScores holds the normalized 0.0-1.0 per-MANRS-action readiness scores.
type ASNActionScores struct {
	M1         float64  `json:"m1"`
	M1C        float64  `json:"m1c"`
	M2Grip     float64  `json:"m2GRIP"`
	M2CGrip    float64  `json:"m2cGRIP"`
	M3         float64  `json:"m3"`
	M3C        float64  `json:"m3c"`
	M3Admin    float64  `json:"m3Admin"`
	M3CAdmin   float64  `json:"m3cAdmin"`
	M4         float64  `json:"m4"`
	M4C        float64  `json:"m4c"`
	M4Admin    float64  `json:"m4Admin"`
	M4CAdmin   float64  `json:"m4cAdmin"`
	M5         *float64 `json:"m5"`
	M7Irr      *float64 `json:"m7IRR"`
	M7Rpki     *float64 `json:"m7RPKI"`
	M7RpkiN    *float64 `json:"m7RPKIN"`
	M7Conflict *float64 `json:"m7Conflict"`
	M8         *float64 `json:"m8"`
	M9         float64  `json:"m9"`
}

// IncidentCounts holds the per-ASN routing-incident counts for a report period (the conformance
// incident breakdown, culprit vs. accomplice).
type IncidentCounts struct {
	RouteLeaks                     uint16 `json:"routeLeaks"`
	RouteLeaksAccomplice           uint16 `json:"routeLeaksAccomplice"`
	RouteMisoriginations           uint16 `json:"routeMisoriginations"`
	RouteMisoriginationsAccomplice uint16 `json:"routeMisoriginationsAccomplice"`
	Total                          uint16 `json:"total"`
}

// ASNScoreSummary represents a single ASN's score data in the /scores listing response.
type ASNScoreSummary struct {
	ASN             uint32           `json:"asn"`
	ASNHolder       string           `json:"asnHolder"`
	ParticipantName string           `json:"participantName"`
	Scores          *CategoryScores  `json:"scores"`
	Actions         *ASNActionScores `json:"actions"`
	Incidents       *IncidentCounts  `json:"incidents"`
}

// ASNScoreDetail represents the detailed score data for a single ASN, including metrics.
type ASNScoreDetail struct {
	Actions          *ASNActionScores  `json:"actions,omitempty"`
	ContactInfoValid *bool             `json:"contactInfoValid"`
	Incidents        *IncidentCounts   `json:"incidents"`
	Metrics          *ASNMetrics       `json:"metrics,omitempty"`
	Routing          *ASNRoutingDetail `json:"routing,omitempty"`
	ROVSamples       []ASNROVSample    `json:"rovSamples,omitempty"`
	Scores           *CategoryScores   `json:"scores"`
	Trend            *ASNTrendData     `json:"trend,omitempty"`
}

// ASNMetrics represents traffic/routing statistics for an ASN.
type ASNMetrics struct {
	TotalIRRRoutes          uint16 `json:"totalIRRRoutes"`
	RouteLeaks              uint16 `json:"routeLeaks"`
	BogonAnnouncements      uint16 `json:"bogonAnnouncements"`
	AdminBogonAnnouncements uint16 `json:"adminBogonAnnouncements"`
}

// BogonDetail is a flattened, per-ASN view of a bogon the participant announced (culprit) or propagated
// (accomplice). It deliberately differs from model.Bogon, which stores one event's full culprit/accomplice
// graph (plus source date, IP version, etc.) and has no notion of the viewing ASN's role, none of which the
// detail view needs. Exactly one of ASN / IPPrefix is set, matching the bogon's event type.
type BogonDetail struct {
	ASN      uint32       `json:"asn,omitempty"`
	IPPrefix netip.Prefix `json:"ipPrefix,omitzero"`
	Class    string       `json:"class"`
	Role     string       `json:"role"`
}

// ConflictingIRRPrefix holds an IRR-registered prefix that is not announced in BGP and has a non-valid
// RPKI state, along with the IRR source(s) and the RPKI state that caused the conflict.
type ConflictingIRRPrefix struct {
	IRRSources []string     `json:"irrSources"`
	IPPrefix   netip.Prefix `json:"ipPrefix"`
	RPKIState  string       `json:"rpkiState"`
}

// CoveringROA describes a single ROA that covers an RPKI-invalid route announcement.
type CoveringROA struct {
	ASN       uint32       `json:"asn"`
	MaxLength int          `json:"maxLength"`
	IPPrefix  netip.Prefix `json:"ipPrefix"`
}

// InvalidRPKIPrefix holds an RPKI-invalid announced prefix and the ROA(s) causing invalidation.
type InvalidRPKIPrefix struct {
	CoveringROAs []CoveringROA `json:"coveringRoas"`
	IPPrefix     netip.Prefix  `json:"ipPrefix"`
	State        string        `json:"state"`
}

// ASNRoutingDetail holds per-ASN IRR and RPKI routing coverage stats.
type ASNRoutingDetail struct {
	TotalIRRRoutes        uint16                 `json:"totalIRRRoutes"`
	RegisteredIRRRoutes   uint16                 `json:"registeredIRRRoutes"`
	UnregisteredIRRRoutes uint16                 `json:"unregisteredIRRRoutes"`
	UnregisteredPrefixes  []string               `json:"unregisteredPrefixes,omitempty"`
	ConflictingPrefixes   []ConflictingIRRPrefix `json:"conflictingPrefixes,omitempty"`
	RegisteredRPKIRoutes  uint16                 `json:"registeredRPKIRoutes"`
	ValidRPKIRoutes       uint16                 `json:"validRPKIRoutes"`
	UnknownRPKIRoutes     uint16                 `json:"unknownRPKIRoutes"`
	InvalidRPKIRoutes     uint16                 `json:"invalidRPKIRoutes"`
	InvalidPrefixes       []InvalidRPKIPrefix    `json:"invalidPrefixes,omitempty"`
	UnknownPrefixes       []string               `json:"unknownPrefixes,omitempty"`
	Bogons                []BogonDetail          `json:"bogons,omitempty"`
}

// ASNROVSample is a per-economy ROV deployment sample for an ASN.
type ASNROVSample struct {
	EconomyCode        string  `json:"economyCode"`
	ValidRate          float32 `json:"validRate"`
	Samples            int     `json:"samples"`
	TotalIRRRoutes     uint16  `json:"totalIRRRoutes"`
	RouteLeaks         uint16  `json:"routeLeaks"`
	BogonAnnouncements uint16  `json:"bogonAnnouncements"`
}

// ScoresFilter represents query parameters for the GET /scores endpoint.
type ScoresFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	Month          time.Time `query:"month,omitempty"`
	ASNs           []uint32  `query:"asns"`
	Economies      []string  `query:"economies"`
	Region         string    `query:"region,omitempty"`
	ParticipantIDs []string  `query:"participantIDs"`
}

// ScoresDetailsFilter represents query parameters for the GET /scores/details endpoint.
type ScoresDetailsFilter struct {
	ASN   uint32    `query:"asn" required:"true"`
	Month time.Time `query:"month,omitempty"`
}

// ScoresSummaryFilter represents query parameters for the GET /scores/summary endpoint.
type ScoresSummaryFilter struct {
	ASNs      []uint32  `query:"asns"`
	Economies []string  `query:"economies"`
	MANRS     string    `query:"manrs,omitempty"`
	Month     time.Time `query:"month,omitempty"`
	Region    string    `query:"region,omitempty"`
}

// IncidentsFilter represents query parameters for the GET /incidents endpoint.
type IncidentsFilter struct {
	ASNs  []uint32  `query:"asns"`
	Month time.Time `query:"month,omitempty"`
	IDs   []string  `query:"ids,omitempty"`
}

// IncidentSummaryFilter represents query parameters for the GET /incidents/summary endpoint.
type IncidentSummaryFilter struct {
	ReportMonth time.Time `query:"reportMonth,omitempty"`
	ASNs        []uint32  `query:"asns"`
}
