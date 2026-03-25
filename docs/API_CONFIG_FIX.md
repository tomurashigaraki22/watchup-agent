# API Configuration Format Fix

## Problem

The API was returning configuration in a different format than the agent expected, causing:
1. Wrong sample counts in logs
2. Thresholds not updating properly
3. Durations being misinterpreted

## API Response Format

The API returns:
```json
{
  "config": {
    "agent_id": "69b0da1c-14d2-431d-b68f-6016340ba61f",
    "sampling_interval": 5,
    "threshold_cpu": 80,
    "threshold_ram": 38,
    "threshold_process": 60,
    "duration_cpu": 10,      // seconds
    "duration_ram": 5,       // seconds
    "duration_process": 10   // seconds
  }
}
```

## Agent Config Format

The agent uses:
```yaml
sampling_interval: 5
alerts:
  cpu:
    threshold: 80
    duration: 10    # seconds
  ram:
    threshold: 38
    duration: 5     # seconds
  process_cpu:
    threshold: 60
    duration: 10    # seconds
```

## Solution

### 1. API Response Struct
Created `APIConfigResponse` struct to properly decode the API's nested JSON format.

### 2. Format Conversion
The `FetchConfig()` method now:
- Decodes the API response correctly
- Converts to agent's internal config format
- Maps all fields properly

### 3. Dynamic Threshold Updates
Added `UpdateThresholds()` method to `SpikeDetector`:
- Updates thresholds without restarting
- Recalculates required sample counts
- Resets violation counters to prevent false alerts

### 4. Config Reload Improvements
The config reload now:
- Preserves local `sampling_interval` (uses local config.yaml value)
- Updates only alert thresholds from API
- Updates spike detector dynamically
- Shows detailed threshold changes in logs

## Sample Count Calculation

With the API config example:
- CPU: 80% for 10s at 5s interval = **2 samples** needed
- RAM: 38% for 5s at 5s interval = **1 sample** needed
- Process: 60% for 10s at 5s interval = **2 samples** needed

Logs will now show:
```
[CONFIG] Alert thresholds updated:
  CPU: 80% for 10s (2 samples)
  RAM: 38% for 5s (1 sample)
  Process: 60% for 10s (2 samples)
```

## Expected Behavior

### On Config Reload
```
[CONFIG] Checking for configuration updates...
[CONFIG] Alert thresholds updated:
  CPU: 80% for 10s (2 samples)
  RAM: 38% for 5s (1 sample)
  Process: 60% for 10s (2 samples)
[CONFIG] Resource check: CPU: 2.1%, RAM: 34.0%
```

### In Monitoring Logs
```
[MONITOR] [23:42:35] CPU: 2.0%, RAM: 34.2% | CPU: 2.0% (0/2 samples), RAM: 34.2% (0/1 samples), Process: 4.1% (0/2 samples)
```

The sample counts now correctly reflect the API's duration settings.

## Files Changed
- `transport/api_client.go` - Added APIConfigResponse struct and proper decoding
- `detectors/spike_detector.go` - Added UpdateThresholds() method
- `cmd/agent/main.go` - Updated config reload to use UpdateThresholds()
