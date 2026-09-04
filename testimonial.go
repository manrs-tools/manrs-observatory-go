package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) FindTestimonials(ctx context.Context, filter *model.TestimonialFilter) ([]*model.Testimonial, error) {
	query := structToQueryParams(filter)

	var response struct {
		Testimonials []*model.Testimonial `json:"testimonials"`
	}

	if _, err := c.GET(ctx, "testimonials?"+query, &response); err != nil {
		return nil, fmt.Errorf("find testimonials: %w", err)
	}

	return response.Testimonials, nil
}

func (c *Client) FindTestimonialsPaged(ctx context.Context, filter *model.TestimonialFilter) *Iterator[*model.Testimonial] {
	return &Iterator[*model.Testimonial]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.Testimonial, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				Testimonials  []*model.Testimonial `json:"testimonials"`
				NextPageToken string               `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "testimonials?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find testimonials: %w", err)
			}

			return resp.Testimonials, resp.NextPageToken, nil
		},
	}
}
