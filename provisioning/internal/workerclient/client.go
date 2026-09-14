package workerclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/BARGHEST-ngo/MESH/provisioning/internal/state"
	"github.com/BARGHEST-ngo/MESH/provisioning/internal/workerapi"
)

// Client implements api.ContainerService by delegating to the internal
// worker process over HTTP instead of talking to Docker directly.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Start(d state.Deployment, token string) error {
	bodyStruct := workerapi.HandleStartRequest{
		Deployment: d,
		Token:      token,
	}
	body, err := json.Marshal(bodyStruct)
	if err != nil {
		return fmt.Errorf("failed to marshal deployment: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/containers", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	return c.do(req)
}

func (c *Client) Stop(slug string) error {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/containers/"+url.PathEscape(slug), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	return c.do(req)
}

func (c *Client) do(req *http.Request) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to reach worker: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("worker returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}
