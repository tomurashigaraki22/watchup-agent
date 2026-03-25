# Scheduler Fix - Monitoring Interval Separation

## Problem
The monitoring ticker and config reload ticker were not properly separated, causing monitoring to only happen every 60 seconds instead of every 5 seconds.

## Solution

### 1. Enhanced Scheduler Logging
Added detailed logging to show:
- Actual intervals being used for monitoring and config reload
- Tick counts for both monitoring and config checks
- Initial tick runs immediately on startup

### 2. Improved Startup Messages
The agent now displays comprehensive configuration on startup:

```
✓ Instance lock acquired (PID: 12345)
✓ Agent registered successfully
Project ID: proj_abc123
Server: my-server-1
API Endpoint: https://watchup.space

--- Monitoring Configuration ---
Sampling interval: 5s (metrics collected every 5 seconds)
Config reload: every 60s

--- Alert Thresholds ---
CPU: 80% sustained for 300s (60 samples)
RAM: 75% sustained for 600s (120 samples)
Process CPU: 60% sustained for 120s (24 samples)

Starting monitoring loop (interval: 5s)...
Scheduler started - Monitoring: every 5s, Config reload: every 1m0s
```

### 3. Prefixed Log Messages
All log messages now have clear prefixes:
- `[MONITOR]` - Monitoring tick events (every 60s summary)
- `[CONFIG]` - Configuration reload events (every 60s)

Example output:
```
[MONITOR] [23:34:05] CPU: 2.1%, RAM: 34.0% | CPU: 2.1% (0/60 samples), RAM: 34.0% (0/120 samples)
[CONFIG] Checking for configuration updates...
[CONFIG] Warning: API returned sampling_interval=60s, keeping local value=5s
[CONFIG] Alert thresholds updated - CPU: 80%, RAM: 75%, Process: 60%
[CONFIG] Resource check: CPU: 2.1%, RAM: 34.0%
```

### 4. Understanding Sample Counts

The sample counts show spike detection progress:
- `CPU: 2.1% (0/60 samples)` means:
  - Current CPU usage is 2.1%
  - It's been above the threshold for 0 out of 60 required samples
  - Alert triggers only after 60 consecutive samples above threshold
  
Example progression:
- `CPU: 85.0% (1/60 samples)` - Just exceeded threshold
- `CPU: 87.0% (30/60 samples)` - Halfway to alert
- `CPU: 88.0% (60/60 samples)` - Alert triggered!
- `CPU: 70.0% (0/60 samples)` - Dropped below, counter reset

### 5. Config Reload Protection

The config reload now:
- Only updates alert thresholds from API
- Preserves local `sampling_interval` setting
- Warns if API tries to change sampling interval
- Immediately checks resource usage after config update
- Never overwrites server_key, project_id, or other local settings

## Expected Behavior

### Monitoring Ticks (every 5 seconds)
- Collect CPU, RAM, and process metrics
- Check for threshold violations
- Send metrics to API
- Print status summary every 12 ticks (60 seconds)
- Sample counters increment when usage exceeds thresholds

### Config Reload (every 60 seconds)
- Fetch latest alert thresholds from API
- Preserve local sampling_interval setting
- Update spike detector thresholds dynamically
- Immediately check resource usage
- No agent restart required

### Sample Counter Behavior
- Increments when metric exceeds threshold
- Resets to 0 when metric drops below threshold
- Alert triggers when counter reaches required samples
- Example: CPU threshold 80% for 300s at 5s interval = 60 samples needed

## Verification

After deploying this fix, you should see:
1. Startup message showing "Monitoring: every 5s"
2. `[MONITOR]` logs appearing every 60 seconds (12 ticks × 5s)
3. `[CONFIG]` logs appearing every 60 seconds
4. Both events may coincide occasionally but run independently

## Files Changed
- `internal/scheduler.go` - Enhanced logging and tick counting
- `cmd/agent/main.go` - Improved startup messages and log prefixes
