// Code generated from internal models; DO NOT EDIT.
package model

import "time"

type Favorite struct {
	Created time.Time `json:"created" readOnly:"true"`
	ID      string    `json:"id" readOnly:"true"`

	MANRS     string `json:"manrs,omitempty"`
	Name      string `json:"name"`
	Selection string `json:"selection"`
	UserID    string `json:"userID" readOnly:"true"`
}

type FavoriteFilter struct {
	Limit int `json:"limit,omitempty" query:"limit"`

	ID string `json:"-" query:"id,omitempty"`
}

type FavoriteUpdate struct {
	MANRS     *string `json:"manrs,omitempty"`
	Name      *string `json:"name,omitempty"`
	Selection *string `json:"selection,omitempty"`
}
