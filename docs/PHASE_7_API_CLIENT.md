# Phase 7 - Watchup API Client

## Status: ✓ Complete

## Overview

Handles all HTTPS communication with the Watchup backend API.

## File: transport/api_client.go

## API Endpoints

### 1. Agent Registration
```
POST /agent/register
```

**Request:**
```json
{
  "project_id": "proj_12345",
  "master_api_key": "ma_abcdef123",
  "server_identifier": "prod-api-1"
}
```

**Response:**
```json
{
  "success": true,
  "server_key": "srv_89sd0a",
  "message": "Agent registered successfully"
}
```

### 2. Send Metrics
```
POST /server/metrics
```

**Request:**
```json
{
  "server_key": "srv_89sd0a",
  "cpu": 45.2,
  "ram": 67.8,
  "timestamp": "2026-03-25T14:10:00Z"
}
```

**Frequency**: Every 5 seconds (every monitoring tick)

### 3. Send Alerts
```
POST /server/alerts
```

**Request:**
```json
{
  "server_key": "srv_89sd0a",
  "metric": "cpu",
  "usage": 87.2,
  "duration": 300,
  "top_processes": [...],
  "timestamp": "2026-03-25T14:11:00Z"
}
```

**Frequency**: When threshold violations occur

### 4. Fetch Configuration
```
GET /agent/config?server_key=srv_89sd0a
```

**Response:**
```json
{
  "server_key": "srv_89sd0a",
  "project_id": "proj_12345",
  "sampling_interval": 5,
  "alerts": {
    "cpu": {
      "threshold": 80,
      "duration": 300
    },
    "ram": {
      "threshold": 75,
      "duration": 600
    }
  }
}
```

**Frequency**: Every 60 seconds

## API Client Structure

```go
type APIClient struct {
    baseURL    string
    serverKey  string
    httpClient *http.Client
}
```

### Initialization

```go
apiClient := transport.NewAPIClient(cfg)
// Uses cfg.APIEndpoint and cfg.ServerKey
```

## Authentication

### Registration
- Uses `master_api_key` in request body
- One-time authentication

### Ongoing Operations
- Uses `server_key` in:
  - Request body (metrics, alerts)
  - Query parameter (config fetch)
  - Authorization header: `Bearer srv_89sd0a`

## HTTP Client Configuration

```go
httpClient: &http.Client{
    Timeout: 10 * time.Second,
}
```

- 10-second timeout for all requests
- Prevents hanging on network issues
- Automatic retry not implemented (fail fast)

## Error Handling

### Network Errors
```go
if err := apiClient.SendMetrics(cpu, ram); err != nil {
    fmt.Printf("Failed to send metrics: %v\n", err)
    // Continue monitoring (don't exit)
}
```

### HTTP Status Codes
- **200-299**: Success
- **400-499**: Client error (logged, continue)
- **500-599**: Server error (logged, continue)

### Registration Errors
```go
if err := registrar.PerformRegistration(cfg, configPath); err != nil {
    fmt.Printf("Registration failed: %v\n", err)
    os.Exit(1)  // Exit on registration failure
}
```

## Request/Response Flow

### Sending Metrics

```go
func (c *APIClient) SendMetrics(cpu, ram float64) error {
    payload := MetricsPayload{
        ServerKey: c.serverKey,
        CPU:       cpu,
        RAM:       ram,
        Timestamp: time.Now().Format(time.RFC3339),
    }
    
    return c.post("/server/metrics", payload)
}
```

### Sending Alerts

```go
func (c *APIClient) SendAlert(alert alerts.AlertPayload) error {
    return c.post("/server/alerts", alert)
}
```

### Fetching Config

```go
func (c *APIClient) FetchConfig() (*config.Config, error) {
    url := fmt.Sprintf("%s/agent/config?server_key=%s", 
        c.baseURL, c.serverKey)
    
    resp, err := c.httpClient.Get(url)
    // ... decode response
    
    return &cfg, nil
}
```

## Security Features

### HTTPS Only
- All communication over TLS
- No plaintext transmission
- Certificate validation enabled

### Authentication
- Server key required for all operations
- Master API key never stored locally
- Keys transmitted in headers/body (not URL for sensitive ops)

### Headers
```
Content-Type: application/json
Authorization: Bearer srv_89sd0a
```

## Performance

### Metrics Sending
- **Frequency**: Every 5 seconds
- **Payload Size**: ~100 bytes
- **Network Usage**: ~20 bytes/sec

### Alert Sending
- **Frequency**: On threshold violation
- **Payload Size**: ~500 bytes (with processes)
- **Network Usage**: Minimal (rare events)

### Config Fetching
- **Frequency**: Every 60 seconds
- **Payload Size**: ~300 bytes
- **Network Usage**: ~5 bytes/sec

**Total Network Usage**: < 30 bytes/sec (~2.5 KB/minute)

## Integration

### With Alert Manager
```go
alertManager := alerts.NewAlertManager(
    cfg.ServerKey,
    func(alert alerts.AlertPayload) {
        apiClient.SendAlert(alert)
    },
)
```

### With Scheduler
```go
monitoringTick := func() error {
    // ... collect metrics ...
    apiClient.SendMetrics(cpuUsage, ramUsage)
    return nil
}

configReload := func() error {
    newCfg, _ := apiClient.FetchConfig()
    // ... update config ...
    return nil
}
```

## Retry Logic

Currently: **No automatic retry**
- Fail fast approach
- Errors logged and monitoring continues
- Next tick will attempt again

Future enhancement:
- Exponential backoff
- Request queuing
- Offline buffering

## Testing

### Mock Server
```bash
# Start mock API server
go run test/mock_api.go

# Run agent with mock endpoint
export API_ENDPOINT=http://localhost:8080
go run cmd/agent/main.go
```

### Network Failure Simulation
```bash
# Block API endpoint
sudo iptables -A OUTPUT -d api.watchup.site -j DROP

# Agent continues monitoring, logs errors
```

## Next Phase

Phase 8 - Configuration System (Already Implemented)
