// Code generated from internal models; DO NOT EDIT.
package model

type TrendDirection string

type TrendData struct {
	Value    float64        `json:"value"`
	Trend    TrendDirection `json:"trend"`
	Severity string         `json:"severity"`
	Previous *float64       `json:"previous,omitempty"`
	Change   *float64       `json:"change,omitempty"`
	ChangeP  *float64       `json:"changePercent,omitempty"`
}

// ASNScoreHistoryPoint is a single month of MANRS readiness score averages for one ASN.
type ASNScoreHistoryPoint struct {
	Month                  string   `json:"month"`
	Filtering              *float64 `json:"filtering,omitempty"`
	AntiSpoofing           *float64 `json:"antiSpoofing,omitempty"`
	Coordination           *float64 `json:"coordination,omitempty"`
	RoutingInformationIRR  *float64 `json:"routingInformationIRR,omitempty"`
	RoutingInformationRPKI *float64 `json:"routingInformationRPKI,omitempty"`
}
