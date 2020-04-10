package client

import (
	"./resource"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserClient struct {
	Host       string
	EntryPoint string
	Client     *http.Client
}

func NewUserClient(host string, timeout int) *UserClient {
	return &UserClient{
		Host:       host,
		EntryPoint: "/v1/user",
		Client: &http.Client{
			Timeout: time.Millisecond * time.Duration(timeout)},
	}
}

func (c *UserClient) PostRequest(requestResource *resource.LoginRequestResource) (*http.Response, error) {
	url := c.Host + c.EntryPoint

	body, err := json.Marshal(requestResource)
	if err != nil {
		return nil, fmt.Errorf("can not create login request. [%w]", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("can not create login request. [%w]", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client hass error. [%w]", err)
	}

	return resp, nil
}
