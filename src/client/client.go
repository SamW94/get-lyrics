package client

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	baseURL = "https://api.genius.com"
)

func constructHTTPRequest(url, searchTerm string) (*http.Request, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error constructing new HTTP request to %v: %v", url, err)
	}

	clientAccessToken := os.Getenv("GENIUS_CLIENT_ACCESS_TOKEN")
	bearerString := fmt.Sprintf("Bearer %v", clientAccessToken)
	request.Header.Set("Authorization", bearerString)
	query := request.URL.Query()
	query.Add("q", searchTerm)
	request.URL.RawQuery = query.Encode()
	return request, nil
}

func (c *Client) makeHTTPRequest(request http.Request) (*http.Response, error) {
	response, err := c.httpClient.Do(&request)
	if err != nil {
		return nil, fmt.Errorf("Error completing HTTP request: %v", err)
	}
	return response, nil
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}
