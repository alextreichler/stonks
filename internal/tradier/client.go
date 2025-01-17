package tradier

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"log"
	"net/http"
	"time"
)

const apiURL = "https://api.tradier.com/v1"

type Client struct {
	apiKey  string
	http    *http.Client
	Program *tea.Program
}

func NewClient(apiKey string) *Client {
	if apiKey == "" {
		log.Fatal("API key is required. Set TRADIER_API_KEY environment variable.")
	}

	return &Client{
		apiKey: apiKey,
		http: &http.Client{
			Timeout: 9 * time.Second,
		},
	}
}

func (c *Client) createStreamingSession() (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/markets/events/session", apiURL), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Stream struct {
			SessionID string `json:"sessionid"`
			URL       string `json:"url"`
		} `json:"stream"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return response.Stream.SessionID, nil
}
