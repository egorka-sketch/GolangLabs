package api

import (
	"context"
	"net/http"
)

func (c *Client) Register(ctx context.Context, login, password string) (Account, error) {
	var acc Account
	err := c.do(ctx, http.MethodPost, "/auth/register", Credentials{Login: login, Password: password}, &acc)
	return acc, err
}

func (c *Client) Login(ctx context.Context, login, password string) error {
	var resp struct {
		Token string `json:"token"`
	}
	if err := c.do(ctx, http.MethodPost, "/auth/login", Credentials{Login: login, Password: password}, &resp); err != nil {
		return err
	}
	c.token = resp.Token
	return nil
}

func (c *Client) Logout(ctx context.Context) error {
	err := c.do(ctx, http.MethodPost, "/auth/logout", nil, nil)
	c.token = ""
	return err
}

func (c *Client) LoggedIn() bool {
	return c.token != ""
}
