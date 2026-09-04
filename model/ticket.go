// Code generated from internal models; DO NOT EDIT.
package model

import (
	"time"
)

type Ticket struct {
	ID      string    `json:"id,omitzero" yaml:"id,omitempty" readOnly:"true"`
	Created time.Time `json:"created" yaml:"created" readOnly:"true"`

	AssignedTo    string              `json:"assignedTo,omitempty" yaml:"assignedTo,omitempty" required:"false"`
	Emails        []string            `json:"emails,omitempty" yaml:"emails,omitempty"`
	LastReplyAt   time.Time           `json:"lastReplyAt" yaml:"lastReplyAt" readOnly:"true"`
	Priority      string              `json:"priority" yaml:"priority" required:"true"`
	Queue         string              `json:"queue" yaml:"queue"`
	Replies       []*TicketReply      `json:"replies,omitempty" yaml:"replies,omitempty" readOnly:"true"`
	Attachments   []*TicketAttachment `json:"attachments,omitempty" yaml:"attachments,omitempty" readOnly:"true"`
	Status        string              `json:"status" yaml:"status" required:"false"`
	Subject       string              `json:"subject" yaml:"subject"`
	Type          string              `json:"type" yaml:"type" required:"false"`
	UserID        string              `json:"userID,omitempty" yaml:"userID,omitempty" required:"false"`
	ParticipantID string              `json:"participantID,omitempty" yaml:"participantID,omitempty" required:"false"`
	ASNs          []uint32            `json:"asns,omitempty" yaml:"asns,omitempty"`
	UserName      string              `json:"userName,omitempty" yaml:"userName,omitempty" readOnly:"true"`
}

type TicketReply struct {
	Created  time.Time `json:"created" yaml:"created" readOnly:"true"`
	ID       string    `json:"id,omitzero" yaml:"id,omitempty" readOnly:"true"`
	TicketID string    `json:"ticketID" yaml:"ticketID" required:"false"`
	UserID   string    `json:"userID,omitempty" yaml:"userID,omitempty" required:"false"`

	Attachments []*TicketAttachment `json:"attachments,omitempty" yaml:"attachments,omitempty" readOnly:"true"`
	Email       string              `json:"email,omitempty" yaml:"email,omitempty"`
	Mailed      bool                `json:"mailed" yaml:"mailed" readOnly:"true"`
	Message     string              `json:"message" yaml:"message"`
	Staff       bool                `json:"staff" yaml:"staff" readOnly:"true"`
	UserName    string              `json:"userName,omitempty" yaml:"userName,omitempty" readOnly:"true"`
}

type TicketAttachment struct {
	ID       string    `json:"id,omitzero" yaml:"id,omitempty" readOnly:"true"`
	Created  time.Time `json:"created" yaml:"created" readOnly:"true"`
	TicketID string    `json:"ticketID" yaml:"ticketID"`
	ReplyID  string    `json:"replyID,omitempty" yaml:"replyID,omitempty"`
	UserID   string    `json:"userID,omitempty" yaml:"userID,omitempty"`

	Filename    string `json:"filename" yaml:"filename"`
	ContentType string `json:"contentType" yaml:"contentType"`
	Size        int64  `json:"size" yaml:"size"`

	Inline bool   `json:"inline,omitempty" yaml:"inline,omitempty"`
	URL    string `json:"url,omitempty" yaml:"url,omitempty" readOnly:"true"`
}

type TicketFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	UserID string `json:"userID" query:"userID,omitempty"`

	ParticipantID string   `json:"participantID" query:"participantID,omitempty"`
	ASNs          []uint32 `json:"asns" query:"asns,omitempty"`
	AssignedTo    string   `json:"assignedTo" query:"assignedTo,omitempty"`
	Type          string   `json:"type" query:"type,omitempty"`
	Queue         string   `json:"queue" query:"queue,omitempty"`
	Status        []string `json:"status" query:"status,omitempty"`
	StatusNot     []string `json:"statusNot" query:"statusNot,omitempty"`
}

type TicketUpdate struct {
	Queue    *string `json:"queue,omitempty"`
	Type     *string `json:"type,omitempty"`
	Priority *string `json:"priority,omitempty"`
	Subject  *string `json:"subject,omitempty"`
	Status   *string `json:"status,omitempty"`

	UserID        *string   `json:"userID,omitempty"`
	ParticipantID *string   `json:"participantID,omitempty"`
	AssignedTo    *string   `json:"assignedTo,omitempty"`
	Emails        *[]string `json:"emails,omitempty"`
	ASNs          *[]uint32 `json:"asns,omitempty"`
	Reply         *string   `json:"reply,omitempty"`
}
