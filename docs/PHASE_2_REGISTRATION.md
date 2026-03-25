# Phase 2 - Project Registration System

## Status: ✓ Complete

## Overview

Implements the one-agent-per-free-project enforcement through a registration system.

## Implementation Details

### Files Modified/Created

1. **config/config.go**
   - Added `ProjectID`, `ServerIdentifier`, and `Registered` fields
   - Added `Save()` method to persist configuration
   - Added `IsRegistered()` method to check registration status

2. **internal/registration.go** (NEW)
   - `Registrar` struct for handling registration
   - Interactive CLI prompts for user input
   - Registration API communication
   - Configuration persistence

3. **cmd/agent/main.go**
   - Added registration check on startup
   - Blocks agent execution if registration fails
   - Interactive registration flow

## Registration Flow

```
Agent Start
    │
    ▼
Check if registered (config.yaml)
    │
    ├─ YES ──> Continue to monitoring
    │
    └─ NO ──> Prompt for credentials
              │
              ├─ Project ID
              ├─ Master API Key
              └─ Server Identifier (optional)
              │
              ▼
         POST /agent/register
              │
              ├─ Success ──> Save server_key ──> Continue
              │
              └─ Failure ──> Exit (enforces free limit)
```

## API Endpoint

**POST** `https://api.watchup.site/agent/register`

### Request
```json
{
  "project_id": "proj_12345",
  "master_api_key": "ma_abcdef123",
  "server_identifier": "prod-api-1"
}
```

### Response (Success)
```json
{
  "success": true,
  "server_key": "srv_89sd0a",
  "message": "Agent registered successfully"
}
```

### Response (Failure - Already Registered)
```json
{
  "success": false,
  "error": "This project already has an agent installed. Free projects are limited to one agent.",
  "message": "Registration failed"
}
```

## Configuration File

After successful registration, `/etc/watchup/config.yaml` contains:

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

## Free vs Paid Users

### Free Users
- Limited to **1 agent per project**
- Backend checks `agent_installed` flag
- Registration fails if agent already exists
- Agent refuses to start without valid registration

### Paid Users
- Multiple agents allowed per project
- Each agent gets unique `server_key`
- Backend allows multiple registrations

## Security Features

- Master API key required for registration (not stored locally)
- Server key used for ongoing authentication
- Configuration file permissions: `0600` (owner read/write only)
- HTTPS-only communication

## Testing

### Test Registration Flow

```bash
# Remove existing config
rm config.yaml

# Run agent
go run cmd/agent/main.go config.yaml
```

Expected prompts:
```
Watchup Server Agent started

⚠️  Agent is not registered.

=== Watchup Agent Registration ===
This agent needs to be registered with your Watchup project.
Note: Free projects can only have ONE agent installed.

Enter your Project ID: proj_12345
Enter your Master API Key: ma_abcdef123
Enter Server Identifier (press Enter for auto-generated): prod-api-1

Registering agent with Watchup...
Project ID: proj_12345
Server: prod-api-1

✓ Registration successful!
Server Key: srv_89sd0a
Configuration saved. Agent is ready to start monitoring.
```

### Test Already Registered

```bash
# Run agent again with same config
go run cmd/agent/main.go config.yaml
```

Expected output:
```
Watchup Server Agent started

✓ Agent registered successfully
Project ID: proj_12345
Server: prod-api-1
Sampling interval: 5s
Alert Thresholds - CPU: 80%, RAM: 75%

Starting monitoring loop...
```

## Error Handling

1. **Invalid Credentials**: Agent exits with error message
2. **Network Failure**: Agent exits with connection error
3. **Already Registered**: Agent exits with clear message about free limit
4. **Missing Config**: Agent creates default and prompts for registration

## Next Phase

Phase 3 - System Metrics Collection (Already Implemented)
