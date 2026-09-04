// Code generated from internal models; DO NOT EDIT.
package model

type Login struct {
	Email    string `json:"email" yaml:"email" required:"true"`
	Password string `json:"password" yaml:"password" required:"true"`
}

type Register struct {
	Email          string `json:"email" required:"true"`
	FirstName      string `json:"firstName" required:"true"`
	LastName       string `json:"lastName" required:"true"`
	Password       string `json:"password" required:"true"`
	PrivacyConsent bool   `json:"privacyConsent" required:"true"`
	TermsConsent   bool   `json:"termsConsent" required:"true"`
}
