package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tomurashigaraki22/watchup-agent/alerts"
	"github.com/tomurashigaraki22/watchup-agent/config"
)

type MetricsPayload struct {
	ServerKey string  `json:"server_key"`
	CPU       float64 `json:"cpu"`
	RAM       float64 `json:"ram"`
	Timestamp string  `json:"timestamp"`
}

type APIClient struct {
	baseURL    string
	serverKey  string
	httpClient *http.Client
}

func NewAPIClient(cfg *config.Config) *APIClient {
	return &APIClient{
		baseURL:   cfg.APIEndpoint,
		serverKey: cfg.ServerKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *APIClient) SendMetrics(cpu, ram float64) error {
	payload := MetricsPayload{
		ServerKey: c.serverKey,
		CPU:       cpu,
		RAM:       ram,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return c.post("/server/metrics", payload)
}

func (c *APIClient) SendAlert(alert alerts.AlertPayload) error {
	return c.post("/server/alerts", alert)
}

func (c *APIClient) FetchConfig() (*config.Config, error) {
	url := fmt.Sprintf("%s/agent/config?server_key=%s", c.baseURL, c.serverKey)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add Authorization header
	req.Header.Set("Authorization", "Bearer "+c.serverKey)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cfg config.Config
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &cfg, nil
}

func (c *APIClient) post(endpoint string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.serverKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
