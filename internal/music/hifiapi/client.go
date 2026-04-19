package hifiapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HifiClient interface {
	Search(params SearchParams) (*SearchResponse, error)
}

type defaultHifiClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) HifiClient {
	return &defaultHifiClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *defaultHifiClient) Search(params SearchParams) (*SearchResponse, error) {
	query := url.Values{}

	if params.Track != "" {
		query.Set("s", params.Track)
	}

	requestURL := fmt.Sprintf("%s/search?%s", c.BaseURL, query.Encode())

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("Request URL: %s\n", requestURL)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp SearchResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &searchResp, nil
}
