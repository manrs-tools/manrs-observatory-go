// Code generated from internal models; DO NOT EDIT.
package model

type AS struct {
	ID           uint32 `json:"id"`
	Economy      string `json:"economy"`
	Name         string `json:"name"`
	Notes        string `json:"notes,omitempty"`
	Organization string `json:"organization"`
	Region       string `json:"region"`
	RIR          string `json:"rir"`
	VerifiedBy   string `json:"verifiedBy,omitempty"`
	Visible      bool   `json:"visible"`

	ClaimeeIDNetOps string `json:"claimeeIDNetOps,omitempty"`
	ClaimeeIDIXP    string `json:"claimeeIDIXP,omitempty"`
	ClaimeeIDCDN    string `json:"claimeeIDCDN,omitempty"`
	ClaimeeIDVendor string `json:"claimeeIDVendor,omitempty"`
}

type ASFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	Economies []string `json:"economies" query:"economies,omitempty"`
	RIRs      []string `json:"rirs" query:"rirs,omitempty"`
}
