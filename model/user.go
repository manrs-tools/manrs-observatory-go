// Code generated from internal models; DO NOT EDIT.
package model

import "time"

// User represents a user in the system.
type User struct {
	ID      string    `json:"id,omitempty" yaml:"id,omitempty" readOnly:"true"`
	Created time.Time `json:"created,omitzero" yaml:"created,omitempty" required:"false" readOnly:"true"`
	Updated time.Time `json:"updated,omitzero" yaml:"updated,omitempty" required:"false" readOnly:"true"`

	Email                  string   `json:"email" yaml:"email" required:"true"`
	FirstName              string   `json:"firstName" yaml:"firstName"`
	LastName               string   `json:"lastName" yaml:"lastName"`
	MailConformanceEnabled bool     `json:"mailConformanceEnabled" yaml:"mailConformanceEnabled"`
	MailListEnabled        bool     `json:"mailListEnabled" yaml:"mailListEnabled"`
	ParticipantIDs         []string `json:"participantIDs,omitempty" yaml:"participantIDs,omitempty"`
	Password               string   `json:"password" yaml:"password" required:"true"`
	Role                   string   `json:"role" yaml:"role" required:"true"`
	Status                 string   `json:"status" yaml:"status" required:"true"`
}

// UserFilter represents a filter passed to FindUsers().
type UserFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	Email                  string `json:"email" query:"email,omitempty"`
	FirstName              string `json:"firstName" query:"firstName,omitempty"`
	LastName               string `json:"lastName" query:"lastName,omitempty"`
	MailConformanceEnabled bool   `json:"mailConformanceEnabled" query:"mailConformanceEnabled,omitempty"`
	MailListEnabled        bool   `json:"mailListEnabled" query:"mailListEnabled,omitempty"`
	ParticipantID          string `json:"participantID" query:"participantID,omitempty"`
	Role                   string `json:"role" query:"role,omitempty"`
	Status                 string `json:"status" query:"status,omitempty"`
}

// UserUpdate represents a set of fields to be updated via UpdateUser().
type UserUpdate struct {
	Email                  *string   `json:"email,omitempty"`
	FirstName              *string   `json:"firstName,omitempty" query:"firstName,omitempty"`
	LastName               *string   `json:"lastName,omitempty" query:"lastName,omitempty"`
	MailConformanceEnabled *bool     `json:"mailConformanceEnabled,omitempty"`
	MailListEnabled        *bool     `json:"mailListEnabled,omitempty"`
	ParticipantIDs         *[]string `json:"participantIDs,omitempty"`
	Password               *string   `json:"password,omitempty"`
	Role                   *string   `json:"role,omitempty"`
	Status                 *string   `json:"status,omitempty"`
}
