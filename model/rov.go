// Code generated from internal models; DO NOT EDIT.
package model

import "time"

type ROV struct {
	EconomyCode string    `json:"economyCode"`
	Filtered    int       `json:"filtered"`
	FilterRate  float32   `json:"filterRate"`
	Seen        int       `json:"seen"`
	SourceDate  time.Time `json:"sourceDate" yaml:"sourceDate"`
}
