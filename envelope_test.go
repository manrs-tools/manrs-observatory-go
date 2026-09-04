package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

// The API nests each of these payloads under a named key. A client decoding the payload at the
// body root instead gets a zero value and a nil error, so every one of these asserts on a field
// carried inside the envelope.
func TestResponseEnvelopesAreUnwrapped(t *testing.T) {
	t.Run("GetStatsComparison", func(t *testing.T) {
		stub := &stubRoundTripper{response: `{"groups":[{"title":"Filtering","dataPoints":[]}]}`}
		c := New("test-key", WithClient(&http.Client{Transport: stub}))

		groups, err := c.GetStatsComparison(context.Background(), model.ComparisonSelection{}, nil)
		if err != nil {
			t.Fatalf("GetStatsComparison returned an error: %v", err)
		}

		if len(groups) != 1 || groups[0].Title != "Filtering" {
			t.Fatalf("groups = %+v, want one group titled Filtering", groups)
		}
	})

	t.Run("GetStatsHistory", func(t *testing.T) {
		stub := &stubRoundTripper{response: `{"history":{"statistics":[{}],"coverage":[]}}`}
		c := New("test-key", WithClient(&http.Client{Transport: stub}))

		history, err := c.GetStatsHistory(context.Background(), &model.StatsHistoryFilter{})
		if err != nil {
			t.Fatalf("GetStatsHistory returned an error: %v", err)
		}

		if history == nil || len(history.Statistics) != 1 {
			t.Fatalf("history = %+v, want one statistics entry", history)
		}
	})

	t.Run("GetStatsOverview", func(t *testing.T) {
		stub := &stubRoundTripper{response: `{"overview":{"months":[],"live":{"users":7}}}`}
		c := New("test-key", WithClient(&http.Client{Transport: stub}))

		overview, err := c.GetStatsOverview(context.Background())
		if err != nil {
			t.Fatalf("GetStatsOverview returned an error: %v", err)
		}

		if overview == nil || overview.Live.Users != 7 {
			t.Fatalf("overview = %+v, want live.users 7", overview)
		}
	})

	t.Run("GetROAStats", func(t *testing.T) {
		stub := &stubRoundTripper{response: `{"summary":{"total":5,"valid":4,"dateFetched":"2026-08-25"}}`}
		c := New("test-key", WithClient(&http.Client{Transport: stub}))

		summary, err := c.GetROAStats(context.Background(), &model.ROAStatsFilter{})
		if err != nil {
			t.Fatalf("GetROAStats returned an error: %v", err)
		}

		if summary == nil || summary.Total != 5 {
			t.Fatalf("summary = %+v, want total 5", summary)
		}
	})
}
