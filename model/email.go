// Code generated from internal models; DO NOT EDIT.
package model

import "time"

type Email struct {
	ID      string    `json:"id" yaml:"id" readOnly:"true"`
	Created time.Time `json:"created" yaml:"created" readOnly:"true"`

	Body      string `json:"body,omitempty" yaml:"body,omitempty"`
	Recipient string `json:"recipient" yaml:"recipient"`
	Subject   string `json:"subject" yaml:"subject"`
	UserID    string `json:"userID,omitempty" yaml:"userID,omitempty"`
}

type EmailFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	Recipient string `json:"recipient" query:"recipient,omitempty"`
	Subject   string `json:"subject" query:"subject,omitempty"`
	UserID    string `json:"userID" query:"userID,omitempty"`
}
