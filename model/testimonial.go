// Code generated from internal models; DO NOT EDIT.
package model

// Testimonial represents a testimonial provided by a participant.
type Testimonial struct {
	Attribution   string `json:"attribution,omitempty" yaml:"attribution,omitempty"`
	Body          string `json:"body,omitempty" yaml:"body,omitempty"`
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	ParticipantID string `json:"participantID,omitempty" yaml:"participantID,omitempty" readOnly:"true"`
	Showcase      bool   `json:"showcase" yaml:"showcase"`
}

// TestimonialFilter represents a filter passed to FindTestimonials().
type TestimonialFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	Name          string `json:"name" query:"name,omitempty"`
	ParticipantID string `json:"participantID" query:"participantID,omitempty"`
	Showcase      bool   `json:"showcase" query:"showcase,omitempty"`
}
