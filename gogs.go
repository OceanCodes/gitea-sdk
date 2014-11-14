// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gogs

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Version() string {
	return "0.0.1"
}

// Client represents a Gogs API client.
type Client struct {
	url         string
	accessToken string
	client      *http.Client
}

// NewClient initializes and returns a API client.
func NewClient(url, token string) *Client {
	return &Client{
		url:         strings.TrimSuffix(url, "/"),
		accessToken: token,
		client:      &http.Client{},
	}
}

func (c *Client) getResponse(method, path string, header http.Header, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, c.url+"/api/v1"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.accessToken)
	for k, v := range header {
		req.Header[k] = v
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, errors.New("page not found")
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		errMap := make(map[string]interface{})
		if err = json.Unmarshal(data, &errMap); err != nil {
			return nil, err
		}
		return nil, errors.New(errMap["message"].(string))
	}

	return data, nil
}

func (c *Client) getParsedResponse(method, path string, header http.Header, body io.Reader, obj interface{}) error {
	data, err := c.getResponse(method, path, header, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}
