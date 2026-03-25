# Phase 8 - Configuration System

## Status: ✓ Complete

## Overview

Manages agent configuration with support for local files and dynamic API updates.

## File: config/config.go

## Configuration Structure

```go
type Config struct {
    ServerKey        string       `yaml:"server_key"`
    ProjectID        string       `yaml:"project_id"`
    ServerIdentifier string       `yaml:"server_identifier"`
    SamplingInterval int          `yaml:"sampling_interval"`
    Alerts           AlertsConfig `yaml:"alerts"`
    APIEndpoint      string       `yaml:"api_endpoint"`
    Registered       bool         `yaml:"registered"`
}

type AlertsConfig struct {
    CPU        AlertConfig `yaml:"cpu"`
    RAM        AlertConfig `yaml:"ram"`
    ProcessCPU AlertConfig `yaml:"process_cpu"`
}

type AlertConfig struct {
    Threshold int `yaml:"threshold"`
    Duration  int `yaml:"duration"`
}
```

## Configuration File

**Location**: `/etc/watchup/config.yaml`

**Example**:
```yaml
server_key: "srv_89sd0a"
project_id: "proj_12345"
server_identifier: "prod-api-1"
sampling_interval: 5
api_endpoint: "https://api.watchup.site"
registered: true

alerts:
  cpu:
    threshold: 80
    duration: 300
  ram:
    threshold: 75
    duration: 600
  process_cpu:
    threshold: 60
    duration: 120
```

## Configuration Loading

### On Startup

```go
cfg, err := config.Load("/etc/watchup/config.yaml")
if err != nil {
    // Handle error
}
```

### Load Behavior

1. **File exists**: Parse YAML
2. **File missing**: Use defaults
3. **Parse error**: Exit with error

### Default Configuration

```go
var defaultConfig = Config{
    ServerKey:        "",
    ProjectID:        "",
    ServerIdentifier: "",
    SamplingInterval: 5,
    APIEndpoint:      "https://api.watchup.site",
    Registered:       false,
    Alerts: AlertsConfig{
        CPU: AlertConfig{
            Threshold: 80,
            Duration:  300,
        },
        RAM: AlertConfig{
            Threshold: 75,
            Duration:  600,
        },
        ProcessCPU: AlertConfig{
            Threshold: 60,
            Duration:  120,
        },
    },
}
```

## Configuration Saving

```go
func (c *Config) Save(path string) error {
    data, err := yaml.Marshal(c)
    if err != nil {
        return err
    }
    
    // File permissions: 0600 (owner read/write only)
    return os.WriteFile(path, data, 0600)
}
```

**Used during**:
- Initial registration
- Manual configuration updates

## Dynamic Configuration Updates

### Fetch from API

```go
configReload := func() error {
    newCfg, err := apiClient.FetchConfig()
    if err != nil {
        return err
    }
    
    // Update in-memory configuration
    *cfg = *newCfg
    
    return nil
}
```

### Update Frequency

- **Interval**: Every 60 seconds
- **Endpoint**: `GET /agent/config?server_key=xxx`
- **No Restart Required**: Changes applied immediately

### Updatable Fields

- `sampling_interval`
- `alerts.cpu.threshold`
- `alerts.cpu.duration`
- `alerts.ram.threshold`
- `alerts.ram.duration`
- `alerts.process_cpu.threshold`
- `alerts.process_cpu.duration`

### Non-Updatable Fields

- `server_key` (security)
- `project_id` (identity)
- `server_identifier` (identity)
- `registered` (state)

## Configuration Validation

### Registration Check

```go
func (c *Config) IsRegistered() bool {
    return c.Registered && 
           c.ServerKey != "" && 
           c.ProjectID != ""
}
```

### Sampling Duration

```go
func (c *Config) GetSamplingDuration() time.Duration {
    return time.Duration(c.SamplingInterval) * time.Second
}
```

## Configuration Hierarchy

1. **Command-line argument**: `go run cmd/agent/main.go /path/to/config.yaml`
2. **Default path**: `/etc/watchup/config.yaml`
3. **Built-in defaults**: If file doesn't exist

## File Permissions

- **Owner**: root (when installed via systemd)
- **Permissions**: `0600` (read/write owner only)
- **Reason**: Protects `server_key` from unauthorized access

## Configuration Scenarios

### First Run (Unregistered)

```yaml
server_key: ""
project_id: ""
server_identifier: ""
registered: false
# ... defaults ...
```

Agent prompts for registration.

### After Registration

```yaml
server_key: "srv_89sd0a"
project_id: "proj_12345"
server_identifier: "prod-api-1"
registered: true
# ... defaults ...
```

Agent starts monitoring immediately.

### Custom Thresholds

```yaml
# ... identity fields ...
alerts:
  cpu:
    threshold: 90  # More lenient
    duration: 180  # Shorter duration
  ram:
    threshold: 85
    duration: 300
```

## Environment Variables

Currently not supported. Future enhancement:

```bash
export WATCHUP_API_ENDPOINT=https://api.staging.watchup.site
export WATCHUP_SERVER_KEY=srv_test_key
```

## Configuration Updates Without Restart

### Scenario

1. User updates thresholds in dashboard
2. Backend stores new configuration
3. Agent fetches config every 60s
4. New thresholds applied immediately
5. Spike detector uses new values

### Example Timeline

```
00:00 - Agent starts (CPU threshold: 80%)
05:00 - User changes threshold to 90% in dashboard
05:30 - Agent fetches new config
05:30 - CPU threshold now 90%
05:30 - Spike detector counters reset
```

## Error Handling

### Load Errors

```go
cfg, err := config.Load(path)
if err != nil {
    fmt.Printf("Failed to load config: %v\n", err)
    os.Exit(1)
}
```

### Save Errors

```go
if err := cfg.Save(path); err != nil {
    fmt.Printf("Failed to save config: %v\n", err)
    // Continue with in-memory config
}
```

### API Fetch Errors

```go
newCfg, err := apiClient.FetchConfig()
if err != nil {
    fmt.Printf("Failed to fetch config: %v\n", err)
    // Continue with current config
    return nil
}
```

## Next Phase

Phase 9 - Installation System (Already Implemented)
