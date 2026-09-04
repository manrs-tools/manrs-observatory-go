// Code generated from internal models; DO NOT EDIT.
package model

import (
	"net/netip"
	"time"
)

// ROAStat records the RPKI validation outcome for a single BGP route announcement on a given day.
type ROAStat struct {
	ASN          uint32       `json:"asn"`
	CoveringROAs []string     `json:"coveringROAs"`
	IPPrefix     netip.Prefix `json:"ipPrefix"`
	IPVersion    uint8        `json:"ipVersion"`
	SourceDate   time.Time    `json:"sourceDate"`
	State        string       `json:"state"`
}

type ROAStatFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	ASNs []uint32 `json:"asns" query:"asns,omitempty"`

	IPPrefix  netip.Prefix `json:"ipPrefix,omitzero" query:"ipPrefix,omitempty"`
	IPVersion uint8        `json:"ipVersion" query:"ipVersion,omitempty"`

	SourceDate     time.Time `json:"sourceDate" query:"sourceDate,omitempty"`
	SourceDateFrom time.Time `json:"sourceDateFrom,omitzero" query:"sourceDateFrom,omitempty"`
	SourceDateTo   time.Time `json:"sourceDateTo,omitzero" query:"sourceDateTo,omitempty"`
	State          string    `json:"state" query:"state,omitempty"`
}
