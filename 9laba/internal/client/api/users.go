package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Client) ListUsers(ctx context.Context, f ListFilter) ([]User, error) {

	q := url.Values{}
	if f.Limit > 0 {
		q.Set("limit", strconv.Itoa(f.Limit))
	}
	if f.Offset > 0 {
		q.Set("offset", strconv.Itoa(f.Offset))
	}
	if f.Name != "" {
		q.Set("name", f.Name)
	}
	if f.MinAge != nil {
		q.Set("min_age", strconv.Itoa(*f.MinAge))
	}
	if f.MaxAge != nil {
		q.Set("max_age", strconv.Itoa(*f.MaxAge))
	}
	path := "/users"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	var users []User
	err := c.do(ctx, http.MethodGet, path, nil, &users)
	return users, err
}

func (c *Client) GetUser(ctx context.Context, id int) (User, error) {
	var u User
	err := c.do(ctx, http.MethodGet, fmt.Sprintf("/users/%d", id), nil, &u)
	return u, err
}

func (c *Client) CreateUser(ctx context.Context, in UserInput) (User, error) {
	var u User
	err := c.do(ctx, http.MethodPost, "/users", in, &u)
	return u, err
}

func (c *Client) UpdateUser(ctx context.Context, id int, in UserInput) (User, error) {
	var u User
	err := c.do(ctx, http.MethodPut, fmt.Sprintf("/users/%d", id), in, &u)
	return u, err
}

func (c *Client) DeleteUser(ctx context.Context, id int) error {
	return c.do(ctx, http.MethodDelete, fmt.Sprintf("/users/%d", id), nil, nil)
}
