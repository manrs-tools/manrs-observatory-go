// Code generated from internal models; DO NOT EDIT.
package model

import "time"

type ASNTrendData struct {
	Month                  time.Time  `json:"month"`
	Filtering              *TrendData `json:"filtering,omitempty"`
	FilteringGRIP          *TrendData `json:"filteringGRIP,omitempty"`
	AntiSpoofing           *TrendData `json:"antiSpoofing,omitempty"`
	Coordination           *TrendData `json:"coordination,omitempty"`
	RoutingInformationIRR  *TrendData `json:"routingInformationIRR,omitempty"`
	RoutingInformationRPKI *TrendData `json:"routingInformationRPKI,omitempty"`
}

// ConformanceFilter represents query parameters for the GET /scores/conformance endpoint.
type ConformanceFilter struct {
	Type string `query:"type,omitempty"`
}
