package client

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func TestRedactSecrets(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		excludes []string
	}{
		{
			name:     "redacts a nested password in a login body",
			input:    `{"login":{"email":"user@example.com","password":"hunter2"}}`,
			contains: []string{"user@example.com", "REDACTED"},
			excludes: []string{"hunter2"},
		},
		{
			name:     "redacts a nested key alongside unrelated fields",
			input:    `{"participant":{"id":"abc123"},"apiKey":"sk-live-secret"}`,
			contains: []string{"abc123", "REDACTED"},
			excludes: []string{"sk-live-secret"},
		},
		{
			name:     "redacts a secret nested inside an array",
			input:    `{"errors":[{"message":"invalid"},{"token":"eyJ.raw.jwt"}]}`,
			contains: []string{"invalid", "REDACTED"},
			excludes: []string{"eyJ.raw.jwt"},
		},
		{
			name:     "leaves non-sensitive bodies unchanged",
			input:    `{"name":"widget","count":3}`,
			contains: []string{"widget", "3"},
			excludes: []string{"REDACTED"},
		},
		{
			name:     "falls back to the raw bytes when the body isn't JSON",
			input:    `not json`,
			contains: []string{"not json"},
			excludes: []string{"REDACTED"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(redactSecrets([]byte(tt.input)))

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("redactSecrets(%q) = %q, want it to contain %q", tt.input, got, want)
				}
			}

			for _, forbidden := range tt.excludes {
				if strings.Contains(got, forbidden) {
					t.Errorf("redactSecrets(%q) = %q, want it to not contain %q", tt.input, got, forbidden)
				}
			}
		})
	}
}

// QueryInner must stay exported: structToQueryParams reads promoted fields through
// reflect.Value.Interface, which panics on an unexported embedded field.
type QueryInner struct {
	Alpha string    `json:"alpha" query:"alpha,omitempty"`
	Stamp time.Time `json:"stamp" query:"stamp,omitempty"`
	Tags  []string  `json:"tags" query:"tags,omitempty"`
}

type QueryOuter struct {
	QueryInner

	Beta int `json:"beta,omitempty" query:"beta"`
}

type QueryTags struct {
	Dashed   string `json:"dashed" query:"-"`
	Untagged string `json:"untagged"`
	Zeroable int    `json:"zeroable,omitzero" query:"zeroable"`
	Plain    bool   `json:"plain" query:"plain"`
	NilSlice []int  `json:"nilSlice" query:"nilSlice,omitempty"`
	Numbers  []int  `json:"numbers" query:"numbers,omitempty"`
}

func TestStructToQueryParams(t *testing.T) {
	stamp := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{name: "nil returns an empty query", input: nil, want: ""},
		{name: "nil pointer returns an empty query", input: (*QueryOuter)(nil), want: ""},
		{name: "non-struct returns an empty query", input: "not a struct", want: ""},
		{name: "zero struct omits every omitempty field", input: QueryInner{}, want: ""},
		{name: "pointer is dereferenced", input: &QueryInner{Alpha: "a"}, want: "alpha=a"},
		{name: "time formats as RFC3339", input: QueryInner{Stamp: stamp}, want: "stamp=2026-03-04T05%3A06%3A07Z"},
		{name: "slices join with commas", input: QueryTags{Numbers: []int{1, 2, 3}}, want: "numbers=1%2C2%2C3&plain=false"},
		{
			name:  "dashed and untagged fields are skipped",
			input: QueryTags{Dashed: "d", Untagged: "u"},
			want:  "plain=false",
		},
		{
			name:  "an omitzero json tag suppresses a field its query tag would emit",
			input: QueryTags{Zeroable: 4},
			want:  "plain=false&zeroable=4",
		},
		{
			name:  "embedded fields encode as if declared on the outer struct",
			input: QueryOuter{QueryInner: QueryInner{Alpha: "a", Tags: []string{"x", "y"}}, Beta: 7},
			want:  "alpha=a&beta=7&tags=x%2Cy",
		},
		{
			name:  "embedded values are escaped exactly once",
			input: QueryOuter{QueryInner: QueryInner{Alpha: "a b&c=d", Stamp: stamp}},
			want:  "alpha=a+b%26c%3Dd&stamp=2026-03-04T05%3A06%3A07Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := structToQueryParams(tt.input); got != tt.want {
				t.Errorf("structToQueryParams(%#v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestStructToQueryParamsEmbeddingIsTransparent asserts that an embedded field and the same field
// declared directly encode identically.
func TestStructToQueryParamsEmbeddingIsTransparent(t *testing.T) {
	inner := QueryInner{
		Alpha: "a b&c=d/e?f",
		Stamp: time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC),
		Tags:  []string{"x y", "z&w"},
	}

	direct := structToQueryParams(inner)
	embedded := structToQueryParams(QueryOuter{QueryInner: inner})

	if embedded != direct {
		t.Errorf("embedded encoding %q does not match direct encoding %q", embedded, direct)
	}
}

func TestRedactStorageURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "redacts every value of a signed storage URL",
			input: "https://bucket.s3.amazonaws.com/reports/2026.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256" +
				"&X-Amz-Credential=AKIAEXAMPLE%2F20260904%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Signature=deadbeef",
			want: "https://bucket.s3.amazonaws.com/reports/2026.csv?X-Amz-Algorithm=%2A%2A%2AREDACTED%2A%2A%2A" +
				"&X-Amz-Credential=%2A%2A%2AREDACTED%2A%2A%2A&X-Amz-Signature=%2A%2A%2AREDACTED%2A%2A%2A",
		},
		{
			name:  "leaves a URL without a query unchanged",
			input: "https://bucket.s3.amazonaws.com/reports/2026.csv",
			want:  "https://bucket.s3.amazonaws.com/reports/2026.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redactStorageURL(tt.input); got != tt.want {
				t.Errorf("redactStorageURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRedactURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "redacts an invite token in the path",
			input: EndpointURL + "/invites/secret-token/redeem",
			want:  EndpointURL + "/invites/***REDACTED***/redeem",
		},
		{
			name:  "redacts a password reset token in the path",
			input: EndpointURL + "/auth/reset-password/secret-token",
			want:  EndpointURL + "/auth/reset-password/***REDACTED***",
		},
		{
			name:  "redacts an email verification token in the path",
			input: EndpointURL + "/auth/verify/secret-token",
			want:  EndpointURL + "/auth/verify/***REDACTED***",
		},
		{
			name:  "redacts a token in the query string",
			input: EndpointURL + "/invites?token=secret-token",
			want:  EndpointURL + "/invites?token=%2A%2A%2AREDACTED%2A%2A%2A",
		},
		{
			name:  "leaves other URLs unchanged",
			input: EndpointURL + "/keys/abc123?limit=10",
			want:  EndpointURL + "/keys/abc123?limit=10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redactURL(tt.input); got != tt.want {
				t.Errorf("redactURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRedactSecretsIsValidJSONWhenInputIsJSON(t *testing.T) {
	input := `{"login":{"email":"user@example.com","password":"hunter2"}}`

	got := redactSecrets([]byte(input))

	var parsed map[string]any
	if err := json.Unmarshal(got, &parsed); err != nil {
		t.Fatalf("redacted output is not valid JSON: %v", err)
	}
}

func TestGenericResponseIncludesErrorLocations(t *testing.T) {
	body := `{"title":"Unprocessable Entity","status":422,"detail":"validation failed","instance":"/api/v2/favorites","type":"https://example.com/errors/422","errors":[{"location":"body.favorite.manrs","message":"expected value to be one of \"participants, non-participants\"","value":"nope"}]}`

	var resp GenericResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	got := resp.ToErrorString()
	for _, want := range []string{"validation failed", "422", "body.favorite.manrs", "non-participants"} {
		if !strings.Contains(got, want) {
			t.Errorf("error string %q missing %q", got, want)
		}
	}

	if resp.Instance != "/api/v2/favorites" {
		t.Errorf("instance = %q", resp.Instance)
	}

	if resp.Type != "https://example.com/errors/422" {
		t.Errorf("type = %q", resp.Type)
	}

	if resp.Errors[0].Value != "nope" {
		t.Errorf("value = %v", resp.Errors[0].Value)
	}
}

func TestResponseErrorUnwrapsWithErrorsAs(t *testing.T) {
	var err error = ResponseError{GenericResponse{
		Status: 422,
		Detail: "validation failed",
		Errors: []ErrorDetail{{Location: "body.favorite.manrs", Message: "bad value"}},
	}}

	var resp ResponseError
	if !errors.As(err, &resp) {
		t.Fatal("errors.As did not recover the ResponseError")
	}

	if resp.Errors[0].Location != "body.favorite.manrs" {
		t.Errorf("location = %q", resp.Errors[0].Location)
	}
}

// TestJSONExcludedFieldsStillEncodeAsQueryParams asserts that a filter field the API binds from
// the query string but excludes from JSON reaches the query string and stays out of a request body.
func TestJSONExcludedFieldsStillEncodeAsQueryParams(t *testing.T) {
	filter := model.FavoriteFilter{ID: "f1", Limit: 5}

	want := "id=f1&limit=5"
	if got := structToQueryParams(&filter); got != want {
		t.Errorf("structToQueryParams = %q, want %q", got, want)
	}

	body, err := New("test-key").marshal(map[string]model.FavoriteFilter{"filter": filter})
	if err != nil {
		t.Fatalf("marshal returned an error: %v", err)
	}

	if string(body) != `{"filter":{"limit":5}}` {
		t.Errorf("marshal = %s, want the id excluded", body)
	}
}
