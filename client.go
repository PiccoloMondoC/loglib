package logclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PiccoloMondoC/sky-common/logtypes"
)

type Client struct {
	HttpClient *http.Client
	BaseURL    string // Base URL of LogAggregator service
}

func NewClient(baseURL string) *Client {
	return &Client{
		HttpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		BaseURL: baseURL,
	}
}

func (c *Client) AggregateLogs(ctx context.Context, logEntry logtypes.LogEntry) error {
	// Prepare the request body
	body, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	// Prepare the request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/logs", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Log this error in your service's logs
		return fmt.Errorf("log aggregator returned status: %d", resp.StatusCode)
	}

	return nil
}
