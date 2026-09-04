// Code generated from internal models; DO NOT EDIT.
package model

// Declaration represents a participant's declaration for specific actions.
type Declaration = map[string]ActionDeclaration

// ActionDeclaration represents a single action in a participant's declaration.
//
// Status is the canonical conformance state. Conformant is retained for
// backward compatibility (it mirrors Status == "conformant") so existing API
// consumers and stored JSON blobs continue to work; records written before
// Status existed carry only Conformant, and EffectiveStatus derives the state
// from it.
type ActionDeclaration struct {
	Comments    string `json:"comments" yaml:"comments"`
	Conformant  bool   `json:"conformant" yaml:"conformant"`
	Description string `json:"description" yaml:"description"`
	Status      string `json:"status,omitempty" yaml:"status,omitempty"`
}
