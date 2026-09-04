// Code generated from internal models; DO NOT EDIT.
package model

import "time"

type IncidentFeedback struct {
	ID      string    `json:"id,omitzero" readOnly:"true"`
	Created time.Time `json:"created" readOnly:"true"`

	IncidentID    string `json:"incidentID"`
	ReportMonth   string `json:"reportMonth"`
	ASN           uint32 `json:"asn"`
	ParticipantID string `json:"participantID,omitempty"`
	FeedbackType  string `json:"feedbackType"`
	Comment       string `json:"comment,omitempty"`

	Cleared bool `json:"cleared" readOnly:"true"`

	AuditorID       string    `json:"auditorID,omitempty" readOnly:"true"`
	AuditorDecision string    `json:"auditorDecision,omitempty" readOnly:"true"`
	AuditorComment  string    `json:"auditorComment,omitempty" readOnly:"true"`
	AuditorReviewed time.Time `json:"auditorReviewed,omitzero" readOnly:"true"`
}

type IncidentFeedbackFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	IncidentID    string `json:"incidentID" query:"incidentID,omitempty"`
	ReportMonth   string `json:"reportMonth" query:"reportMonth,omitempty"`
	ASN           uint32 `json:"asn" query:"asn,omitempty"`
	ParticipantID string `json:"participantID" query:"participantID,omitempty"`
}
