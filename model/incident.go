// Code generated from internal models; DO NOT EDIT.
package model

import (
	"net/netip"
	"time"
)

// IncidentRow is one routing-security incident as served by GET /incidents, attributed to a
// single ASN. An incident affecting several ASNs yields one row per ASN.
type IncidentRow struct {
	ASN        uint32       `json:"asn"`
	ID         string       `json:"id"`
	MetricType string       `json:"metricType"`
	IsCulprit  bool         `json:"isCulprit"`
	StartTime  time.Time    `json:"startTime"`
	EndTime    time.Time    `json:"endTime"`
	Duration   float64      `json:"duration"`
	IPPrefix   netip.Prefix `json:"ipPrefix"`
	IPVersion  uint8        `json:"ipVersion,omitempty"`
	Paths      []string     `json:"paths,omitempty"`
	Type       string       `json:"type,omitempty"`
	Cleared    bool         `json:"cleared,omitempty"`
}

type IncidentSummaryRow struct {
	ASN       uint32 `json:"asn"`
	ASNHolder string `json:"asnHolder"`
	M1        uint16 `json:"m1"`
	M1C       uint16 `json:"m1c"`
	M2        uint16 `json:"m2"`
	M2C       uint16 `json:"m2c"`
}
