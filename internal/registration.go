package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tomurashigaraki22/watchup-agent/config"
)

type RegistrationRequest struct {
	ProjectID        string `json:"project_id"`
	MasterAPIKey     string `json:"master_api_key"`
	ServerIdentifier string `json:"server_identifier"`
}

type RegistrationResponse struct {
	Success   bool   `json:"success"`
	ServerKey string `json:"server_key"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
}

type Registrar struct {
	apiEndpoint string
	httpClient  *http.Client
}

func NewRegistrar(apiEndpoint string) *Registrar {
	return &Registrar{
		apiEndpoint: apiEndpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (r *Registrar) PerformRegistration(cfg *config.Config, configPath string) error {
	fmt.Println("\n=== Watchup Agent Registration ===")
	fmt.Println("This agent needs to be registered with your Watchup project.")
	fmt.Println("Note: Free projects can only have ONE agent installed.\n")

	reader := bufio.NewReader(os.Stdin)

	// Get Project ID
	fmt.Print("Enter your Project ID: ")
	projectID, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read project ID: %w", err)
	}
	projectID = strings.TrimSpace(projectID)

	// Get Master API Key
	fmt.Print("Enter your Master API Key: ")
	masterAPIKey, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read master API key: %w", err)
	}
	masterAPIKey = strings.TrimSpace(masterAPIKey)

	// Get Server Identifier (optional, generate default)
	fmt.Print("Enter Server Identifier (press Enter for auto-generated): ")
	serverIdentifier, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read server identifier: %w", err)
	}
	serverIdentifier = strings.TrimSpace(serverIdentifier)
	
	if serverIdentifier == "" {
		hostname, _ := os.Hostname()
		if hostname == "" {
			hostname = "server"
		}
		serverIdentifier = fmt.Sprintf("%s-%d", hostname, time.Now().Unix())
	}

	fmt.Printf("\nRegistering agent with Watchup...\n")
	fmt.Printf("Project ID: %s\n", projectID)
	fmt.Printf("Server: %s\n", serverIdentifier)

	// Send registration request
	serverKey, err := r.register(projectID, masterAPIKey, serverIdentifier)
	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	// Update configuration
	cfg.ProjectID = projectID
	cfg.ServerKey = serverKey
	cfg.ServerIdentifier = serverIdentifier
	cfg.Registered = true

	// Save configuration
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println("\n✓ Registration successful!")
	fmt.Printf("Server Key: %s\n", serverKey)
	fmt.Println("Configuration saved. Agent is ready to start monitoring.\n")

	return nil
}

func (r *Registrar) register(projectID, masterAPIKey, serverIdentifier string) (string, error) {
	reqBody := RegistrationRequest{
		ProjectID:        projectID,
		MasterAPIKey:     masterAPIKey,
		ServerIdentifier: serverIdentifier,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.apiEndpoint + "/agent/register"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var regResp RegistrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !regResp.Success {
		if regResp.Error != "" {
			return "", fmt.Errorf("%s", regResp.Error)
		}
		return "", fmt.Errorf("registration failed: %s", regResp.Message)
	}

	if regResp.ServerKey == "" {
		return "", fmt.Errorf("server key not provided in response")
	}

	return regResp.ServerKey, nil
}

func (r *Registrar) CheckRegistration(cfg *config.Config) bool {
	return cfg.IsRegistered()
}
