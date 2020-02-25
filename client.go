// Package locust implements a locust client using native Go data structures
package locust

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// defines the timeout value as 2 seconds
const (
	defaultTimeout = 2 * time.Second
)

// Client defines a structure for a locust client.
type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
}

// SwarmResponse defines the structure of a response for locust swarm endpoint
// Locust returns a json with status of the action and a message.
// start and stop locust swarming both returns a message with same structure.
type SwarmResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// StatsResponse defines the structure for the response  from locust stats endpoint.
type StatsResponse struct {
	CurrentResponseTimePercentile50 float64 `json:"current_response_time_percentile_50"`
	CurrentResponseTimePercentile95 float64 `json:"current_response_time_percentile_95"`
	Errors                          []Error `json:"errors"`
	FailRatio                       float64 `json:"fail_ratio"`
	State                           string  `json:"state"`
	Statistics                      []Stat  `json:"stats"`
	TotalRps                        float64 `json:"total_rps"`
	UserCount                       int     `json:"user_count"`
}

// Stat defines locust stats structure from locust stats endpoint, this is part of StatsResponse
type Stat struct {
	AvgContentLength   float64 `json:"avg_content_length"`
	AvgResponseTime    float64 `json:"avg_response_time"`
	CurrentRps         float64 `json:"current_rps"`
	MaxResponseTime    float64 `json:"max_response_time"`
	MedianResponseTime float64 `json:"median_response_time"`
	Method             string  `json:"method"`
	MinResponseTime    float64 `json:"min_response_time"`
	Name               string  `json:"name"`
	NumFailures        int     `json:"num_failures"`
	NumRequests        int     `json:"num_requests"`
}

// Error defines structure of error records in locust stats from locust stats endpoint,
// this is part of StatsResponse
type Error struct {
	Error      string `json:"error"`
	Method     string `json:"method"`
	Name       string `json:"name"`
	Occurences int    `json:"occurences"`
}

// StartLoad starts locust swarming or modifes if the load generation has already started,
// hatch rate and number of users to simulate are inputs.
func (c *Client) startLoad(users int, hatchrate int) (*SwarmResponse, error) {
	s := SwarmResponse{}
	u, err := c.BaseURL.Parse("/swarm")
	if err != nil {
		return nil, err
	}
	// sets payload for post as hatch rate and user count
	data := url.Values{"locust_count": {strconv.Itoa(users)}, "hatch_rate": {strconv.Itoa(hatchrate)}}

	resp, err := c.httpClient.PostForm(u.String(), data)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 || s.Success != true {
		return nil, err
	}

	return &s, nil
}

// StopLoad stops an existing locust execution
func (c *Client) stopLoad() (*SwarmResponse, error) {
	s := SwarmResponse{}
	u, err := c.BaseURL.Parse("/stop")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 || s.Success != true {
		return nil, err
	}

	return &s, nil
}

// isReady probes reset endpoint of locust to check if the service is ready
func (c *Client) isReady()  error {
	u, err := c.BaseURL.Parse("/stats/reset")
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return err
	}

	respdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && string(respdata) == "ok" {
		return err
	}

	return nil
}

// GetStatus gets the execution metrics from locust
// this provides error ratio, current users hatchd, state and many other defined in
// StatsResponse structure
func (c *Client) getStatus() (*StatsResponse, error) {
	s := StatsResponse{}
	u, err := c.BaseURL.Parse("/stats/requests")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, err
	}

	return &s, nil
}

// New initiantes a new client to control locust, url of the locust endpoint is required as a paramenter
func New(endpoint string) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	} else if u.Scheme == "" {
		return nil, fmt.Errorf("invalid url, protocol scheme is empty")
	} else if u.Host == "" {
		return nil, fmt.Errorf("invalid url, host is empty")
	}

	client := Client{
		BaseURL: u,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	return &client, err
}