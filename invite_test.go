package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type stubRoundTripper struct {
	gotBody   []byte
	gotMethod string
	gotPath   string
	response  string
}

func (s *stubRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		s.gotBody, _ = io.ReadAll(req.Body)
	}

	s.gotMethod = req.Method
	s.gotPath = req.URL.Path

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(s.response)),
		Header:     http.Header{"Content-Type": []string{ContentTypeJSON}},
	}, nil
}

func TestRedeemInvite(t *testing.T) {
	stub := &stubRoundTripper{response: `{"user":{"id":"user-001","email":"alice@example.com"}}`}
	c := New("test-key", WithClient(&http.Client{Transport: stub}))

	user, err := c.RedeemInvite(context.Background(), "tok-abc", "s3cret-pass")
	if err != nil {
		t.Fatalf("RedeemInvite returned an error: %v", err)
	}

	if stub.gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", stub.gotMethod)
	}

	if want := "/api/v2/invites/tok-abc/redeem"; stub.gotPath != want {
		t.Errorf("path = %q, want %q", stub.gotPath, want)
	}

	var body map[string]string
	if err = json.Unmarshal(stub.gotBody, &body); err != nil {
		t.Fatalf("request body is not JSON: %v (%s)", err, stub.gotBody)
	}

	if body["password"] != "s3cret-pass" {
		t.Errorf("request body = %s, want the password under a %q key", stub.gotBody, "password")
	}

	if user == nil {
		t.Fatal("RedeemInvite returned a nil user; the response envelope should be unwrapped")
	}

	if user.ID != "user-001" {
		t.Errorf("user.ID = %q, want %q", user.ID, "user-001")
	}
}
