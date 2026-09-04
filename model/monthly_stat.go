// Code generated from internal models; DO NOT EDIT.
package model

import "time"

// ParticipantTypeCounts holds incident-related counts for one participant type within a single month.
type ParticipantTypeCounts struct {
	Count                     uint32 `json:"count"`
	ASNsWithIncidents         uint32 `json:"asnsWithIncidents"`
	ParticipantsWithIncidents uint32 `json:"participantsWithIncidents"`
}

// MonthlyStat is the merged Observatory-wide snapshot for a single month.
type MonthlyStat struct {
	Month             time.Time                        `json:"month"`
	Participants      map[string]ParticipantTypeCounts `json:"participants"`
	TotalASNs         uint32                           `json:"totalASNs"`
	ClaimedASNs       uint32                           `json:"claimedASNs"`
	ASNsWithIncidents uint32                           `json:"asnsWithIncidents"`
	UpdatedAt         time.Time                        `json:"updatedAt"`
}

// StatsOverview wraps a series of MonthlyStat rows plus live present-day totals.
type StatsOverview struct {
	Months []*MonthlyStat  `json:"months"`
	Live   StatsLiveCounts `json:"live"`
}

// StatsLiveCounts holds present-day Observatory totals fetched directly from
// the source tables, independent of the precomputed monthly stats rows.
type StatsLiveCounts struct {
	Participants uint32 `json:"participants"`
	Users        uint32 `json:"users"`
	TotalASNs    uint32 `json:"totalASNs"`
	ClaimedASNs  uint32 `json:"claimedASNs"`
}
