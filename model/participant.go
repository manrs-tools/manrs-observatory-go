// Code generated from internal models; DO NOT EDIT.
package model

import "time"

// ExternalID is a participant's identifier in an external registry. Provider names
// the registry/namespace; ID is that registry's identifier: numeric for the current
// providers, kept as a string to accommodate future registries with handle-style IDs.
// A participant may hold many: an IXP running multiple exchange points has one IX-F /
// PeeringDB ix pair per exchange.
type ExternalID struct {
	Provider string `json:"provider" yaml:"provider"`
	ID       string `json:"id" yaml:"id"`
}

// Participant represents a registered participant in the system.
type Participant struct {
	ID       string    `json:"id,omitempty" yaml:"id,omitempty" readOnly:"true"`
	Approved time.Time `json:"approved,omitzero" yaml:"approved,omitempty" readOnly:"true"`
	Created  time.Time `json:"created,omitzero" yaml:"created,omitempty" required:"false" readOnly:"true"`
	Deleted  bool      `json:"deleted,omitzero" yaml:"deleted,omitempty"`
	Updated  time.Time `json:"updated,omitzero" yaml:"updated,omitempty" required:"false" readOnly:"true"`

	Abbreviation  string       `json:"abbreviation,omitempty" yaml:"abbreviation,omitempty"`
	ASNs          []uint32     `json:"asns,omitempty" yaml:"asns,omitempty" readOnly:"true"`
	AccessASNs    []uint32     `json:"accessASNs,omitempty" yaml:"accessASNs,omitempty"`
	Economies     []string     `json:"economies,omitempty" yaml:"economies,omitempty"`
	RIRs          []string     `json:"rirs,omitempty" yaml:"rirs,omitempty"`
	HideScores    bool         `json:"hideScores,omitzero" yaml:"hideScores,omitempty"`
	ExternalIDs   []ExternalID `json:"externalIDs,omitempty" yaml:"externalIDs,omitempty"`
	LegacyID      int64        `json:"legacyID,omitempty" yaml:"legacyID,omitempty" readOnly:"true"`
	LogoURL       string       `json:"logoURL,omitempty" yaml:"logoURL,omitempty"`
	Name          string       `json:"name" yaml:"name" required:"true"`
	PrimaryUserID string       `json:"primaryUserID,omitempty" yaml:"primaryUserID,omitempty"`
	Regions       []string     `json:"regions,omitempty" yaml:"regions,omitempty"`
	Status        string       `json:"status,omitempty" yaml:"status,omitempty"`
	Type          string       `json:"type" yaml:"type" required:"true"`
	URL           string       `json:"url,omitempty" yaml:"url,omitempty"`
}

// ParticipantFilter represents a filter passed to FindParticipants().
type ParticipantFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	IDs           []string `json:"ids" query:"ids,omitempty"`
	LegacyID      int64    `json:"legacyID" query:"legacyID,omitempty"`
	Abbreviation  string   `json:"abbreviation" query:"abbreviation,omitempty"`
	ASN           uint32   `json:"asn" query:"asn,omitempty"`
	Name          string   `json:"name" query:"name,omitempty"`
	PrimaryUserID string   `json:"primaryUserID" query:"primaryUserID,omitempty"`
	Region        string   `json:"region" query:"region,omitempty"`
	Status        string   `json:"status" query:"status,omitempty"`
	Type          string   `json:"type" query:"type,omitempty"`
	Deleted       bool     `json:"deleted" query:"deleted,omitempty"`
	MemberOf      bool     `json:"memberOf" query:"memberOf,omitempty"`
}

// ParticipantBadgeFilter represents query parameters for the GET /participants/{id}/badge endpoint.
type ParticipantBadgeFilter struct {
	Layout   string `query:"layout,omitempty"`
	Style    string `query:"style,omitempty"`
	Download bool   `query:"download,omitempty"`
}
