package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) CreateTicket(ctx context.Context, ticket *model.Ticket) (*model.Ticket, error) {
	body, err := c.marshal(map[string]*model.Ticket{"ticket": ticket})
	if err != nil {
		return nil, fmt.Errorf("marshal ticket: %w", err)
	}

	var response struct {
		Ticket *model.Ticket `json:"ticket"`
	}

	if _, err = c.POST(ctx, "tickets", body, &response); err != nil {
		return nil, fmt.Errorf("create ticket: %w", err)
	}

	return response.Ticket, nil
}

func (c *Client) FindTicket(ctx context.Context, ticketID string) (*model.Ticket, error) {
	var response struct {
		Ticket *model.Ticket `json:"ticket"`
	}

	if _, err := c.GET(ctx, "tickets/"+ticketID, &response); err != nil {
		return nil, fmt.Errorf("find ticket: %w", err)
	}

	return response.Ticket, nil
}

func (c *Client) FindTickets(ctx context.Context, filter *model.TicketFilter) ([]*model.Ticket, error) {
	query := structToQueryParams(filter)

	var response struct {
		Tickets []*model.Ticket `json:"tickets"`
	}

	if _, err := c.GET(ctx, "tickets?"+query, &response); err != nil {
		return nil, fmt.Errorf("find tickets: %w", err)
	}

	return response.Tickets, nil
}

func (c *Client) FindTicketsPaged(ctx context.Context, filter *model.TicketFilter) *Iterator[*model.Ticket] {
	return &Iterator[*model.Ticket]{
		ctx: ctx,
		fetchPage: func(ctx context.Context, pageToken string) ([]*model.Ticket, string, error) {
			q := structToQueryParams(filter)
			q = withPageToken(q, pageToken)

			var resp struct {
				Tickets       []*model.Ticket `json:"tickets"`
				NextPageToken string          `json:"nextPageToken"`
			}

			if _, err := c.GET(ctx, "tickets?"+q, &resp); err != nil {
				return nil, "", fmt.Errorf("find tickets: %w", err)
			}

			return resp.Tickets, resp.NextPageToken, nil
		},
	}
}

func (c *Client) UpdateTicket(ctx context.Context, ticketID string, update *model.TicketUpdate) (*model.Ticket, error) {
	body, err := c.marshal(map[string]*model.TicketUpdate{"ticket": update})
	if err != nil {
		return nil, fmt.Errorf("marshal ticket update: %w", err)
	}

	var response struct {
		Ticket *model.Ticket `json:"ticket"`
	}

	if _, err = c.PATCH(ctx, "tickets/"+ticketID, body, &response); err != nil {
		return nil, fmt.Errorf("update ticket: %w", err)
	}

	return response.Ticket, nil
}

// DownloadTicketAttachment returns the bytes of an attachment. The endpoint answers with a
// redirect to storage rather than the file itself, so the object is fetched from the target the
// redirect names.
func (c *Client) DownloadTicketAttachment(ctx context.Context, ticketID string, attachmentID string) ([]byte, error) {
	data, err := c.download(ctx, "tickets/"+ticketID+"/attachments/"+attachmentID)
	if err != nil {
		return nil, fmt.Errorf("download ticket attachment: %w", err)
	}

	return data, nil
}

// CreateTicketReply posts a reply to a ticket. The endpoint answers with no content, so the
// stored reply is not returned; read it back with FindTicket.
func (c *Client) CreateTicketReply(ctx context.Context, reply *model.TicketReply) error {
	body, err := c.marshal(map[string]*model.TicketReply{"reply": reply})
	if err != nil {
		return fmt.Errorf("marshal ticket reply: %w", err)
	}

	if _, err = c.POST(ctx, "ticket-replies", body, nil); err != nil {
		return fmt.Errorf("create ticket reply: %w", err)
	}

	return nil
}
