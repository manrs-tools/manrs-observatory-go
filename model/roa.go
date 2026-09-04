// Code generated from internal models; DO NOT EDIT.
package model

import (
	"net/netip"
	"time"
)

// ROA represents a route origin authorization from an RPKI.
type ROA struct {
	ASN       uint32       `json:"asn" yaml:"asn"`
	IPPrefix  netip.Prefix `json:"ipPrefix" yaml:"ipPrefix"`
	IPVersion uint8        `json:"ipVersion" yaml:"ipVersion"`
	MaxLength int          `json:"maxLength" yaml:"maxLength"`
	NotAfter  time.Time    `json:"notAfter" yaml:"notAfter"`
	NotBefore time.Time    `json:"notBefore" yaml:"notBefore"`
	Source    string       `json:"source" yaml:"source"`

	Valid bool `json:"valid" yaml:"valid"`
}

type ROAFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	ASNs []uint32 `json:"asns" query:"asns,omitempty"`

	IPPrefix   netip.Prefix `json:"ipPrefix,omitzero" query:"ipPrefix,omitempty"`
	IPVersion  uint8        `json:"ipVersion" query:"ipVersion,omitempty"`
	MaxLength  int          `json:"maxLength" query:"maxLength,omitempty"`
	NotAfter   time.Time    `json:"notAfter,omitzero" query:"notAfter,omitempty"`
	NotBefore  time.Time    `json:"notBefore,omitzero" query:"notBefore,omitempty"`
	Source     string       `json:"source" query:"source,omitempty"`
	SourceDate time.Time    `json:"sourceDate,omitzero" query:"sourceDate,omitempty"`
	Valid      bool         `json:"valid" query:"valid,omitempty"`
}

type ROAStatsFilter struct {
	ASN     uint32    `query:"asn,omitempty"`
	Economy string    `query:"economy,omitempty"`
	Region  string    `query:"region,omitempty"`
	Month   time.Time `query:"month,omitempty"`
}

type ROARecordsFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	ASN     uint32    `query:"asn,omitempty"`
	Economy string    `query:"economy,omitempty"`
	Region  string    `query:"region,omitempty"`
	Status  string    `query:"status,omitempty"`
	Month   time.Time `query:"month,omitempty"`
}

// ROAStats holds RPKI validation counts, combined across address families.
type ROAStats struct {
	Total   int `json:"total"`
	Invalid int `json:"invalid"`
	Unknown int `json:"unknown"`
	Valid   int `json:"valid"`
}

// ROAVersionStats is the RPKI validation breakdown for a single IP version. It lets
// callers recover the IPv4/IPv6 split without the response hardcoding it as top-level
// fields.
type ROAVersionStats struct {
	ROAStats

	IPVersion uint8 `json:"ipVersion"`
}

// ROAStatsDetail is the aggregate ROA stats response for a scope (ASN, economy, or
// region). Counts are combined across address families; ByVersion carries the optional
// per-version breakdown, and ByEconomy the per-economy breakdown (region scope only).
type ROAStatsDetail struct {
	ROAStats

	DateFetched string            `json:"dateFetched,omitempty"`
	Entity      ROAStatsEntity    `json:"entity,omitzero"`
	ByVersion   []ROAVersionStats `json:"byVersion,omitempty"`
	ByEconomy   []ROAEconomyStats `json:"byEconomy,omitempty"`
}

// ROAStatsEntity describes the entity that a ROA stats response was filtered by.
type ROAStatsEntity struct {
	Label string `json:"label"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ROAStatsRecord is a single prefix entry in a ROA stats response.
type ROAStatsRecord struct {
	ASN       uint32       `json:"asn"`
	ASNName   string       `json:"asnName"`
	Economy   string       `json:"economy"`
	IPPrefix  netip.Prefix `json:"ipPrefix"`
	IPVersion uint8        `json:"ipVersion"`
	Status    string       `json:"status"`
}

// ROAEconomyStats is one economy's aggregate RPKI validity within a region: the
// per-economy breakdown carried in ROAStatsDetail.ByEconomy for region scope. It
// exposes no ASN-level data.
type ROAEconomyStats struct {
	ROAStats

	Economy   string            `json:"economy"`
	Label     string            `json:"label"`
	ByVersion []ROAVersionStats `json:"byVersion,omitempty"`
}
