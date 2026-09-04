package client

import (
	"context"
	"fmt"

	"github.com/manrs-tools/manrs-observatory-go/model"
)

func (c *Client) Login(ctx context.Context, email string, password string) (*model.APIKey, error) {
	body, err := c.marshal(map[string]*model.Login{"login": {Email: email, Password: password}})
	if err != nil {
		return nil, fmt.Errorf("marshal login request: %w", err)
	}

	var response struct {
		Key *model.APIKey `json:"key"`
	}

	if _, err = c.POST(ctx, "auth/login", body, &response); err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	return response.Key, nil
}

func (c *Client) Register(ctx context.Context, register *model.Register) error {
	body, err := c.marshal(register)
	if err != nil {
		return fmt.Errorf("marshal register: %w", err)
	}

	if _, err = c.POST(ctx, "auth/register", body, nil); err != nil {
		return fmt.Errorf("register: %w", err)
	}

	return nil
}

func (c *Client) CreatePasswordResetToken(ctx context.Context, email string) error {
	body, err := c.marshal(map[string]string{"email": email})
	if err != nil {
		return fmt.Errorf("marshal password reset: %w", err)
	}

	if _, err = c.POST(ctx, "auth/reset-password", body, nil); err != nil {
		return fmt.Errorf("create password reset token: %w", err)
	}

	return nil
}

func (c *Client) ResetPassword(ctx context.Context, token string, password string) error {
	body, err := c.marshal(map[string]string{"password": password}) //nolint:goconst // Request body field name, unrelated to the redaction key set in http.go.
	if err != nil {
		return fmt.Errorf("marshal reset password: %w", err)
	}

	if _, err = c.POST(ctx, "auth/reset-password/"+token, body, nil); err != nil {
		return fmt.Errorf("reset password: %w", err)
	}

	return nil
}

func (c *Client) VerifyEmail(ctx context.Context, token string) error {
	if _, err := c.POST(ctx, "auth/verify/"+token, nil, nil); err != nil {
		return fmt.Errorf("verify email: %w", err)
	}

	return nil
}
