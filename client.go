package pool

import (
	"io"
	"net/http"
	"net/url"
)

// ClientInterface is a interface build around the http.client
type ClientInterface interface {
	Do(request *http.Request) (resp ResponseInterface, err error)
	Get(url string) (resp ResponseInterface, err error)
	Post(url string, bodyType string, body io.Reader) (resp ResponseInterface, err error)
	PostForm(url string, data url.Values) (resp ResponseInterface, err error)
	Head(url string) (resp ResponseInterface, err error)
}

// Client is a wrapper around the http.Client,but instead or returning the http.Response it will
// return a ResponseInterface. This is used by the pool and can be overwritten to use you own
// client. Fo example a client that returns a wrapped response or a pre configured client.
type Client struct {
	c *http.Client
}
// @inheritdoc
func NewClient() *Client {
	return &Client{c: &http.Client{}}
}
// @inheritdoc
func (c *Client) Do(request *http.Request) (resp ResponseInterface, err error) {
	if resp, err := c.c.Do(request); err != nil {
		return nil, err
	} else {
		return ResponseInterface(resp), nil
	}
}
// @inheritdoc
func (c *Client) Get(url string) (resp ResponseInterface, err error) {
	if resp, err := c.c.Get(url); err != nil {
		return nil, err
	} else {
		return ResponseInterface(resp), nil
	}
}
// @inheritdoc
func (c *Client) Post(url string, bodyType string, body io.Reader) (resp ResponseInterface, err error) {
	if resp, err := c.c.Post(url, bodyType, body); err != nil {
		return nil, err
	} else {
		return ResponseInterface(resp), nil
	}
}
// @inheritdoc
func (c *Client) PostForm(url string, data url.Values) (resp ResponseInterface, err error) {
	if resp, err := c.c.PostForm(url, data); err != nil {
		return nil, err
	} else {
		return ResponseInterface(resp), nil
	}
}
// @inheritdoc
func (c *Client) Head(url string) (resp ResponseInterface, err error) {
	if resp, err := c.c.Head(url); err != nil {
		return nil, err
	} else {
		return ResponseInterface(resp), nil
	}
}
