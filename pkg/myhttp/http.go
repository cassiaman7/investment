package myhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	HTTPDefaultTimeout = 60 * time.Second
)

type Client struct {
	body    []byte
	timeout time.Duration
	header  map[string]string
	cookies []*http.Cookie
}

type HTTPFunc func(c *Client)

func HTTPTimeout(timeout time.Duration) HTTPFunc {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func HTTPHeader(m map[string]string) HTTPFunc {
	return func(c *Client) {
		c.header = m
	}
}

func HTTPBody(body []byte) HTTPFunc {
	return func(c *Client) {
		c.body = body
	}
}

func HTTPCookies(cookies []*http.Cookie) HTTPFunc {
	return func(c *Client) {
		c.cookies = cookies
	}
}

func NewClient(fn ...HTTPFunc) *Client {
	c := &Client{
		timeout: HTTPDefaultTimeout,
		header:  make(map[string]string),
	}
	for _, f := range fn {
		f(c)
	}

	return c
}

func (c *Client) HTTPGet(url string, rs interface{}) (err error) {
	return c.req(http.MethodGet, url, rs)
}

func (c *Client) HTTPPost(url string, rs interface{}) (err error) {
	return c.req(http.MethodPost, url, rs)
}

func (c *Client) req(method, rawURL string, rs interface{}) (err error) {
	req, err := http.NewRequest(method, rawURL, bytes.NewBuffer(c.body))
	if err != nil {
		return
	}
	client := &http.Client{
		Timeout: c.timeout,
	}

	for k, v := range c.header {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("return code: %d", resp.StatusCode)
	}
	rt, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(rt, rs)
}
