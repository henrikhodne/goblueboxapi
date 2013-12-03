package goblueboxapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// A client manages communication with the Blue Box API
type Client struct {
	client *http.Client

	BaseURL *url.URL

	// Credentials which is used for authentication during API requests
	CustomerId string
	ApiKey     string

	Blocks *BlocksService
}

func NewClient(customerId, apiKey string) *Client {
	c := &Client{
		client: http.DefaultClient,
		BaseURL: &url.URL{
			Scheme: "https",
			Host:   "boxpanel.bluebox.net",
		},
		CustomerId: customerId,
		ApiKey:     apiKey,
	}

	c.Blocks = &BlocksService{client: c}

	return c
}

func (c *Client) UrlForPath(path string) *url.URL {
	u, _ := url.Parse(fmt.Sprintf("api%s.json", path))
	return c.BaseURL.ResolveReference(u)
}

func (c *Client) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := c.UrlForPath(path)

	req, err := http.NewRequest(method, url.String(), body)

	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	req.SetBasicAuth(c.CustomerId, c.ApiKey)
	req.Header.Add("Accept", "application/json; charset=utf-8")

	return req, err
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if c := resp.StatusCode; c < 200 || c >= 300 {
		return errors.New(fmt.Sprintf("Expected status 2xx, got %d", c))
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return err
}
