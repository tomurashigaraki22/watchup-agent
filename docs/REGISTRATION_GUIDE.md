# Agent Registration Guide

## Why Manual Registration is Required

The Watchup Server Agent requires interactive input during first-time registration. Since systemd services run in the background without a terminal, they cannot accept user input. Therefore, you must run the agent manually first to complete registration.

---

## Registration Process

### Step 1: Stop the Service

If the service is running and failing due to missing registration:

```bash
sudo systemctl stop watchup-agent
```

### Step 2: Run Agent Manually

```bash
sudo /usr/local/bin/watchup-agent
```

### Step 3: Enter Credentials

You'll see prompts like this:

```
Watchup Server Agent started

⚠️  Agent is not registered.

=== Watchup Agent Registration ===
This agent needs to be registered with your Watchup project.
Note: Free projects can only have ONE agent installed.

Enter your Project ID: 
```

**Enter your Project ID** (e.g., `proj_12345`)

```
Enter your Master API Key: 
```

**Enter your Master API Key** (e.g., `ma_abcdef123456`)

```
Enter Server Identifier (press Enter for auto-generated): 
```

**Enter a server name** (e.g., `production-api-1`) or press Enter for auto-generated

### Step 4: Verify Registration

After entering credentials, you should see:

```
Registering agent with Watchup...
Project ID: proj_12345
Server: production-api-1

✓ Registration successful!
Server Key: srv_89sd0a7f3b2c1d4e
Configuration saved. Agent is ready to start monitoring.

✓ Agent registered successfully
Project ID: proj_12345
Server: production-api-1
Sampling interval: 5s
Alert Thresholds - CPU: 80%, RAM: 75%

Starting monitoring loop...
[14:23:15] CPU: 12.3%, RAM: 45.6% | CPU: 12.3% (0/60 samples), RAM: 45.6% (0/120 samples), Process: 0.0% (0/24 samples)
```

### Step 5: Stop Manual Run

Press `Ctrl+C` to stop the manual run:

```
^C
Shutdown signal received
Scheduler stopped
Watchup Server Agent stopped
```

### Step 6: Start as Service

Now start the agent as a systemd service:

```bash
sudo systemctl start watchup-agent
```

### Step 7: Verify Service

Check that the service is running:

```bash
sudo systemctl status watchup-agent
```

Expected output:
```
● watchup-agent.service - Watchup Server Monitoring Agent
   Loaded: loaded (/etc/systemd/system/watchup-agent.service; enabled)
   Active: active (running) since Tue 2026-03-25 14:23:00 UTC; 5s ago
 Main PID: 12345 (watchup-agent)
   CGroup: /system.slice/watchup-agent.service
           └─12345 /usr/local/bin/watchup-agent
```

### Step 8: View Logs

Monitor the agent logs:

```bash
sudo journalctl -u watchup-agent -f
```

You should see:
```
Mar 25 14:23:00 server watchup-agent[12345]: Watchup Server Agent started
Mar 25 14:23:00 server watchup-agent[12345]: ✓ Agent registered successfully
Mar 25 14:23:00 server watchup-agent[12345]: Project ID: proj_12345
Mar 25 14:23:00 server watchup-agent[12345]: Server: production-api-1
Mar 25 14:23:00 server watchup-agent[12345]: Starting monitoring loop...
Mar 25 14:23:05 server watchup-agent[12345]: [14:23:05] CPU: 12.3%, RAM: 45.6%
```

---

## Where to Get Credentials

### Project ID

1. Log in to Watchup dashboard: https://watchup.space
2. Go to **Project Settings**
3. Copy your **Project ID** (format: `proj_xxxxx`)

### Master API Key

1. In Watchup dashboard, go to **Settings** → **API Keys**
2. Copy your **Master API Key** (format: `ma_xxxxx`)
3. **Important**: This key is only used for registration and is not stored by the agent

### Server Identifier

- Choose any friendly name for your server
- Examples: `production-api`, `staging-db`, `web-server-1`
- Must be alphanumeric with hyphens/underscores
- 3-50 characters long

---

## Troubleshooting

### Error: "failed to read project ID: EOF"

**Cause**: The agent is running as a systemd service and cannot read input.

**Solution**: Stop the service and run manually:
```bash
sudo systemctl stop watchup-agent
sudo /usr/local/bin/watchup-agent
```

### Error: "This project already has an agent installed"

**Cause**: Free projects are limited to one agent.

**Solutions**:
1. Use the existing agent
2. Upgrade to a paid plan
3. Use a different project

### Error: "Invalid master API key"

**Cause**: The API key is incorrect or has been revoked.

**Solution**: 
1. Verify the key in Watchup dashboard
2. Generate a new key if needed
3. Try registration again

### Error: "Project with ID 'proj_12345' does not exist"

**Cause**: The project ID is incorrect.

**Solution**:
1. Verify the project ID in Watchup dashboard
2. Ensure you're using the correct project
3. Try registration again

### Service Keeps Restarting

**Cause**: Agent is not registered and exits immediately.

**Solution**: Complete manual registration first:
```bash
sudo systemctl stop watchup-agent
sudo /usr/local/bin/watchup-agent
# Complete registration
# Press Ctrl+C
sudo systemctl start watchup-agent
```

---

## Configuration File

After successful registration, the configuration is saved to `/etc/watchup/config.yaml`:

```yaml
server_key: "srv_89sd0a7f3b2c1d4e"
project_id: "proj_12345"
server_identifier: "production-api-1"
sampling_interval: 5
api_endpoint: "https://watchup.space"
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

**Important**: 
- File permissions are set to `600` (owner read/write only)
- The `server_key` is used for all future API communication
- The `master_api_key` is NOT stored in this file

---

## Re-registration

If you need to re-register the agent:

1. Stop the service:
   ```bash
   sudo systemctl stop watchup-agent
   ```

2. Edit the config file:
   ```bash
   sudo nano /etc/watchup/config.yaml
   ```

3. Set `registered: false` and clear the `server_key`:
   ```yaml
   server_key: ""
   registered: false
   ```

4. Run manual registration again:
   ```bash
   sudo /usr/local/bin/watchup-agent
   ```

---

## Security Notes

- The **Master API Key** is only used during registration
- It is NOT stored in the configuration file
- The **Server Key** is generated during registration and stored locally
- The Server Key is used for all ongoing API communication
- Configuration file has restricted permissions (600)

---

## Next Steps

After successful registration:

1. ✅ Agent is monitoring your server
2. ✅ Metrics are being sent every 5 seconds
3. ✅ Alerts will trigger on threshold violations
4. ✅ View data in Watchup dashboard

**Dashboard**: https://watchup.space/dashboard

---

## Support

If you encounter issues during registration:

- **Documentation**: See [README.md](../README.md)
- **API Docs**: See [API_SPECIFICATION.md](API_SPECIFICATION.md)
- **GitHub Issues**: https://github.com/tomurashigaraki22/watchup-agent/issues
- **Email**: support@watchup.space
