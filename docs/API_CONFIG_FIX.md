# API Configuration Format Fix

## Problem

The API client wasn't using the correct endpoint and response format, causing:
1. Wrong sample counts in logs
2. Thresholds not updating properly
3. Durations being misinterpreted

## Correct API Endpoint

**Endpoint:** `GET /agent/config`  
**Authentication:** Bearer token in Authorization header  
**Header:** `Authorization: Bearer {server_key}`

## API Response Format

The API returns:
```json
{
  "success": true,
  "updated_at": "2026-03-25T23:38:50",
  "config": {
    "server_key": "srv_xxx",
    "project_id": "proj_xxx",
    "server_identifier": "my-server",
    "sampling_interval": 5,
    "api_endpoint": "https://watchup.space",
    "alerts": {
      "cpu": {
        "threshold": 80,
        "duration": 300,
        "enabled": true
      },
      "ram": {
        "threshold": 75,
        "duration": 600,
        "enabled": true
      },
      "process_cpu": {
        "threshold": 60,
        "duration": 120,
        "enabled": true
      }
    },
    "features": {
      "disk_monitoring": false,
      "network_monitoring": false,
      "custom_metrics": false
    }
  }
}
```

## Agent Config Format

The agent uses:
```yaml
server_key: "srv_xxx"
project_id: "proj_xxx"
server_identifier: "my-server"
sampling_interval: 5
api_endpoint: "https://watchup.space"
registered: true

alerts:
  cpu:
    threshold: 80
    duration: 300    # seconds
  ram:
    threshold: 75
    duration: 600    # seconds
  process_cpu:
    threshold: 60
    duration: 120    # seconds
```

## Solution

### 1. Correct API Endpoint
Changed from:
- ❌ `GET /agent/config?server_key={key}` (query parameter)
- ✅ `GET /agent/config` with `Authorization: Bearer {key}` header

### 2. API Response Struct
Created `APIConfigResponse` struct matching the actual API response:
- Includes `success` field
- Includes `updated_at` timestamp
- Properly nested `config` object with `alerts` structure
- Includes `features` for future extensibility

### 3. Format Conversion
The `FetchConfig()` method now:
- Uses correct endpoint without query parameters
- Sends Bearer token in Authorization header
- Decodes the full API response correctly
- Validates `success` field
- Maps all fields properly to agent config

### 4. Dynamic Threshold Updates
Added `UpdateThresholds()` method to `SpikeDetector`:
- Updates thresholds without restarting
- Recalculates required sample counts
- Resets violation counters to prevent false alerts

### 5. Config Reload Improvements
The config reload now:
- Preserves local `sampling_interval` (uses local config.yaml value)
- Updates only alert thresholds from API
- Updates spike detector dynamically
- Shows detailed threshold changes in logs

## Sample Count Calculation

With typical API config:
- CPU: 80% for 300s at 5s interval = **60 samples** needed
- RAM: 75% for 600s at 5s interval = **120 samples** needed
- Process: 60% for 120s at 5s interval = **24 samples** needed

Logs will show:
```
[CONFIG] Alert thresholds updated:
  CPU: 80% for 300s (60 samples)
  RAM: 75% for 600s (120 samples)
  Process: 60% for 120s (24 samples)
```

## Expected Behavior

### On Config Reload
```
[CONFIG] Checking for configuration updates...
[CONFIG] Alert thresholds updated:
  CPU: 80% for 300s (60 samples)
  RAM: 75% for 600s (120 samples)
  Process: 60% for 120s (24 samples)
[CONFIG] Resource check: CPU: 2.1%, RAM: 34.0%
```

### In Monitoring Logs
```
[MONITOR] [23:42:35] CPU: 2.0%, RAM: 34.2% | CPU: 2.0% (0/60 samples), RAM: 34.2% (0/120 samples), Process: 4.1% (0/24 samples)
```

The sample counts now correctly reflect the API's duration settings.

## Files Changed
- `transport/api_client.go` - Updated endpoint, added proper APIConfigResponse struct, fixed authentication
- `detectors/spike_detector.go` - Added UpdateThresholds() method
- `cmd/agent/main.go` - Updated config reload to use UpdateThresholds()
