package locust

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	locusturl = "http://localhost:8089"

	locustStatsResponce = `{
		"current_response_time_percentile_50": 11, 
		"current_response_time_percentile_95": 22, 
		"errors": [], 
		"fail_ratio": 0.31311475409836065, 
		"state": "running", 
		"stats": [], 
		"total_rps": 9.9, 
		"user_count": 5
	}`

	locustTestStoppedResponce = `{
		"message": "Test stopped", 
		"success": true
	}`

	locustTestStartedResponce = `{
		"message": "Swarming started", 
		"success": true
	}`
)

func TestNewClientURLSetting(t *testing.T) {
	c, err := New(locusturl)
	assert.Nil(t, err)
	url := c.BaseURL.String()
	assert.Equal(t, locusturl, url)
}
