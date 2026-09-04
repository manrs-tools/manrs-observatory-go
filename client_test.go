package client

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestCheckRetry(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		noValue bool // Context carries no method at all, simulating a caller that forgot withMethod.
		resp    *http.Response
		err     error
		want    bool
	}{
		{
			name:   "retries a GET on a 503",
			method: http.MethodGet,
			resp:   &http.Response{StatusCode: http.StatusServiceUnavailable},
			want:   true,
		},
		{
			name:   "retries a GET on a transport error",
			method: http.MethodGet,
			err:    errors.New("connection reset by peer"),
			want:   true,
		},
		{
			name:   "does not retry a GET on a 200",
			method: http.MethodGet,
			resp:   &http.Response{StatusCode: http.StatusOK},
			want:   false,
		},
		{
			name:   "does not retry a POST on a 503",
			method: http.MethodPost,
			resp:   &http.Response{StatusCode: http.StatusServiceUnavailable},
			want:   false,
		},
		{
			name:   "does not retry a POST on a transport error",
			method: http.MethodPost,
			err:    errors.New("connection reset by peer"),
			want:   false,
		},
		{
			name:   "does not retry a PATCH on a transport error",
			method: http.MethodPatch,
			err:    errors.New("connection reset by peer"),
			want:   false,
		},
		{
			name:   "does not retry a DELETE on a transport error",
			method: http.MethodDelete,
			err:    errors.New("connection reset by peer"),
			want:   false,
		},
		{
			name:    "fails closed when the context carries no method",
			noValue: true,
			err:     errors.New("connection reset by peer"),
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if !tt.noValue {
				ctx = withMethod(ctx, tt.method)
			}

			got, _ := checkRetry(ctx, tt.resp, tt.err)
			if got != tt.want {
				t.Errorf("checkRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}
