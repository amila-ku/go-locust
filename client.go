// Package locust implements a locust client using native Go data structures
// Locust API https://docs.locust.io/en/stable/api.html#
// Locust endpoints https://github.com/locustio/locust/blob/master/locust/web.py
package locust

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// defines the timeout value of httpclient as 2 seconds
const (
	defaultTimeout = 2 * time.Second
)

// Client defines a structure for a locust client. 
// Client contains URL of the locust endpoint and a httpclient
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
	Error       string `json:"error"`
	Method      string `json:"method"`
	Name        string `json:"name"`
	Occurrences int    `json:"occurrences"`
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

// GenerateLoad starts locust swarming or modifes if the load generation has already started,
// hatch rate and number of users to simulate are inputs.
func (c *Client) GenerateLoad(users int, hatchrate float64) (*SwarmResponse, error) {
	s := SwarmResponse{}
	u, err := c.BaseURL.Parse("/swarm")
	if err != nil {
		return nil, err
	}
	// sets payload for post as hatch rate and user count
	data := url.Values{"user_count": {strconv.Itoa(users)}, "hatch_rate": {strconv.FormatFloat(hatchrate, 'E', -1, 32)}}

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
func (c *Client) StopLoad() (*SwarmResponse, error) {
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
func (c *Client) isReady() error {
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

// Stats gets the execution metrics from locust
// this provides error ratio, current users hatchd, state and many other defined in
// StatsResponse structure
func (c *Client) Stats() (*StatsResponse, error) {
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

func (c *Client) getCurrentRps() (float64, error) {
	s, err := c.Stats()
	if err != nil {
		return s.TotalRps, err
	}
	return s.TotalRps, err
}

//Swarm handles load test when given maximum requests rate and ramp up time.
func (c *Client) Swarm(rps float64, duration string) (*SwarmResponse, error) {
	// check if the locust is ready
	if err := c.isReady(); err != nil {
		return &SwarmResponse{ Message: "Locust not ready", Success: false}, err
	}

	userCount := 1
	hatchRate := 0.1
	targetRps := rps
	loadtestDuration, err := time.ParseDuration(duration)

	if err != nil {
		loadtestDuration, _ = time.ParseDuration("1h")
	}

	// generate load with 1 user to identify rps, this will be used to calculate rps later.
	c.GenerateLoad(userCount, hatchRate)

	/// get rps for 1 user
	initrps, err := c.getCurrentRps()
	if err != nil {
		return &SwarmResponse{ Message: "Locust failed to generate load", Success: false}, err
	}
	currentRps := initrps

	// timed wait and start ramping up load untill expected rps

	for currentRps < targetRps {
		userTarget := calculateUsersTarget(targetRps,currentRps,userCount)
		if userCount < userTarget {
			userCount = +5
		} else {
			userCount = +1
		}

		c.GenerateLoad(userCount, hatchRate)

		// get rps for current execution
		r, err := c.getCurrentRps()
		if err != nil {
			return &SwarmResponse{ Message: "Failed to get current RPS from Locust for initial attempt", Success: false}, err
		}

		for r < initrps*float64(userCount)/2 {
			// get rps for current number of usrs, sleep for two seconds if not expected rps is achive
			r, err = c.getCurrentRps()
			if err != nil {
				return &SwarmResponse{ Message: "Failed to get current RPS from Locust", Success: false}, err
			}
			time.Sleep(2 * time.Second)
		}

		// if time exceeds stop load test
		fmt.Println(loadtestDuration)

	}

	return &SwarmResponse{ Message: "Started generating load", Success: true}, nil

}

// calculate users required for rps
func calculateUsersTarget(targetrps float64, currentrps float64, currentusers int) int {
	return int(targetrps / (currentrps / float64(currentusers)))
}
